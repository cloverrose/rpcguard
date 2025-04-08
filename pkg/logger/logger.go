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

// Config is logging configuration.
type Config struct {
	Level  string // "debug", "info", "warn", "error"
	File   string // file name to write log
	Format string // "json" or "text"
}

func SetDefault(cfg Config, pass *analysis.Pass) (closer func() error, err error) {
	opts := &slog.HandlerOptions{
		Level: convertLogLevel(cfg.Level),
	}

	// default closer
	closer = func() error { return nil }

	// configure writer
	writer := os.Stdout
	if cfg.File != "" {
		file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		writer = file
		closer = file.Close
	}

	// configure handler
	var handler slog.Handler
	switch strings.ToLower(cfg.Format) {
	case "text":
		handler = slog.NewTextHandler(writer, opts)
	default:
		handler = slog.NewJSONHandler(writer, opts)
	}

	// set default logger
	slog.SetDefault(slog.New(&ValueHandler{
		handler: handler,
		pass:    pass,
	}))
	return closer, nil
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

func convertLogLevel(level string) slog.Level {
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
