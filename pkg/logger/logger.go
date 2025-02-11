package logger

import (
	"context"
	"fmt"
	"go/token"
	"log/slog"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ssa"
)

const (
	valueKey = "value"
	posKey   = "pos"
)

func NewLogger(level slog.Level, pass *analysis.Pass) *slog.Logger {
	return slog.New(&ValueHandler{
		handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}),
		pass:    pass,
	})
}

func Attr(value ssa.Value) slog.Attr {
	return slog.Group(valueKey, slog.Any(posKey, value.Pos()), slog.String("name", value.Name()), slog.String("str", value.String()))
}

type ValueHandler struct {
	handler slog.Handler
	pass    *analysis.Pass
}

func (h *ValueHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *ValueHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.handler.WithAttrs(attrs)
}

func (h *ValueHandler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}

func (h *ValueHandler) Handle(ctx context.Context, r slog.Record) error {
	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	f := func(attr slog.Attr) bool {
		if attr.Key != valueKey {
			newRecord.AddAttrs(attr)
			return true
		}
		groups := attr.Value.Group()
		newGroups := make([]any, 0, len(groups))
		for _, subAttr := range groups {
			p, ok := getPos(subAttr)
			if !ok {
				newGroups = append(newGroups, subAttr)
				continue
			}
			posStr := getPosStr(h.pass, p)
			newGroups = append(newGroups, slog.String(posKey, posStr))
		}
		newRecord.AddAttrs(slog.Group(valueKey, newGroups...))
		return true
	}
	r.Attrs(f)
	return h.handler.Handle(ctx, newRecord)
}

func getPos(attr slog.Attr) (token.Pos, bool) {
	if attr.Key != posKey {
		return 0, false
	}
	value, ok := attr.Value.Any().(token.Pos)
	if !ok {
		return 0, false
	}
	return value, true
}

func getPosStr(pass *analysis.Pass, pos token.Pos) string {
	position := pass.Fset.Position(pos)
	return fmt.Sprintf("%s:%d:%d", getFileName(position.Filename), position.Line, position.Column)
}

func getFileName(path string) string {
	if !strings.HasPrefix(path, "/var/folders") {
		return path
	}
	elems := strings.Split(path, "/")
	return elems[len(elems)-1]
}

func ConvertLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
