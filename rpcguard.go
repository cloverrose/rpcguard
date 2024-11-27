package rpcguard

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"log"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const doc = "rpcguard checks if connect endpoint is implemented properly."

// Analyzer checks if connect endpoint is implemented properly.
var Analyzer = &analysis.Analyzer{
	Name:     "rpcguard",
	Doc:      doc,
	Run:      run,
	Requires: []*analysis.Analyzer{},
	Flags:    *flag.NewFlagSet("rpcguard", flag.ExitOnError),
}

func init() {
	Analyzer.Flags.BoolVar(&Verbose, "Verbose", Verbose, "verbose logging")
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		fileName := pass.Fset.Position(file.Pos()).Filename
		if skipFile(fileName) {
			continue
		}
		aliases, err := getAliasData(pass, file)
		if err != nil {
			// skip file that imports package more than once.
			return nil, nil
		}
		// Process each function declaration
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			// Check if the function has the signature we're looking for
			if !isTargetGRPCMethod(funcDecl, aliases) {
				continue
			}

			if Verbose {
				log.Printf("%s %s is RPC method\n", fileName, funcDecl.Name.Name)
			}

			validationFound := false
			// Walk through the function body to find protovalidate.Validate usage
			ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
				if isValidateMethodCall(n, aliases) {
					validationFound = true
					return false
				}
				return true
			})

			if !validationFound {
				pass.Reportf(funcDecl.Pos(), "gRPC endpoint %s does not call protovalidate.Validate", funcDecl.Name.Name)
			}
		}

	}
	return nil, nil
}

func skipFile(fileName string) bool {
	if strings.HasSuffix(fileName, "_test.go") {
		return true
	}
	if strings.HasSuffix(fileName, ".connect.go") {
		return true
	}
	return false
}

type aliasData struct {
	context       string
	connect       string
	protovalidate string
}

func (a *aliasData) String() string {
	return fmt.Sprintf("aliasData{%s %s %s}", a.context, a.connect, a.protovalidate)
}

func getAliasData(pass *analysis.Pass, file *ast.File) (*aliasData, error) {
	contextAlias, err := getImportAlias(pass, file, "context", "context")
	if err != nil {
		return nil, err
	}
	connectAlias, err := getImportAlias(pass, file, "connectrpc.com/connect", "connect")
	if err != nil {
		return nil, err
	}
	protovalidateAlias, err := getImportAlias(pass, file, "github.com/bufbuild/protovalidate-go", "protovalidate")
	if err != nil {
		return nil, err
	}
	return &aliasData{
		context:       contextAlias,
		connect:       connectAlias,
		protovalidate: protovalidateAlias,
	}, nil
}

// Without alias
// getImportAlias(x, "connectrpc.com/connect", "connect") returns "connect"
// With alias e.g. cc "connectrpc.com/connect"
// getImportAlias(x, "connectrpc.com/connect", "connect") returns "cc"
func getImportAlias(pass *analysis.Pass, file *ast.File, targetImportPath string, defaultRef string) (string, error) {
	var alias string
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if path != targetImportPath {
			continue
		}
		if imp.Name == nil {
			if alias != "" {
				pass.Reportf(imp.Path.Pos(), "duplicated import %s", imp.Path.Value)
				return "", errors.New("duplicated import")
			}
			alias = defaultRef
		} else {
			if alias != "" {
				pass.Reportf(imp.Path.Pos(), "duplicated import %s", imp.Path.Value)
				return "", errors.New("duplicated import")
			}
			alias = imp.Name.Name
		}
	}
	return alias, nil
}

// Method(ctx context.Context, req *connect.Request[T]) (*connect.Response[S], error)
func isTargetGRPCMethod(funcDecl *ast.FuncDecl, aliases *aliasData) bool {
	if funcDecl.Type == nil || funcDecl.Type.Params == nil || funcDecl.Type.Results == nil {
		return false
	}

	// Check parameters (ctx context.Context, req *connect.Request[T])
	if len(funcDecl.Type.Params.List) != 2 {
		return false
	}

	// Check results: (*connect.Response[S], error)
	if len(funcDecl.Type.Results.List) != 2 {
		return false
	}

	// First parameter check: context.Context
	if !isContextType(funcDecl.Type.Params.List[0].Type, aliases, "Context") {
		return false
	}

	// Second parameter check: *connect.Request[T]
	if !isConnectType(funcDecl.Type.Params.List[1].Type, aliases, "Request") {
		return false
	}

	// First result check: *connect.Response[S]
	if !isConnectType(funcDecl.Type.Results.List[0].Type, aliases, "Response") {
		return false
	}

	return true
}

// context.Context
func isContextType(expr ast.Expr, aliases *aliasData, typeName string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	x, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	if aliases.context == "" {
		return false // not imported
	}
	if x.Name == aliases.context && sel.Sel.Name == typeName {
		return true
	}

	return false
}

// *connect.Request[T] or *connect.Response[S]
func isConnectType(expr ast.Expr, aliases *aliasData, typeName string) bool {
	star, ok := expr.(*ast.StarExpr)
	if !ok {
		return false
	}

	index, ok := star.X.(*ast.IndexExpr)
	if !ok {
		return false
	}

	sel, ok := index.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	x, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	if aliases.connect == "" {
		return false // not imported
	}
	if x.Name == aliases.connect && sel.Sel.Name == typeName {
		return true
	}

	return false
}

func isValidateMethodCall(n ast.Node, aliases *aliasData) bool {
	callExpr, ok := n.(*ast.CallExpr)
	if !ok {
		return false
	}

	sel, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	x, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	if aliases.protovalidate == "" {
		return false // not imported
	}
	if x.Name == aliases.protovalidate && sel.Sel.Name == "Validate" {
		return true
	}

	return false
}
