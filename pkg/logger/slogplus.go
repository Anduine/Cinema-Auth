package logger

import (
	"context"
	"encoding/json"
	"io"
	stdLog "log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type PlusHandler struct {
	slog.Handler
	l        *stdLog.Logger
	attrs    []slog.Attr
	location *time.Location
}

func SetupGlobalLogger(env string, levelStr string, utcTimeZone int) {
	level := parseLevel(levelStr)
	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler

	switch env {
	case "prod":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "local":
		handler = NewPlusHandler(os.Stdout, opts, utcTimeZone)
	default:
		handler = NewPlusHandler(os.Stdout, opts, utcTimeZone)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func NewPlusHandler(out io.Writer, opts *slog.HandlerOptions, utcTimeZone int) *PlusHandler {
	loc := time.FixedZone("UTC", utcTimeZone*60*60)

	return &PlusHandler{
		Handler:  slog.NewJSONHandler(out, opts),
		l:        stdLog.New(out, "", 0),
		location: loc,
	}
}

func (h *PlusHandler) Handle(_ context.Context, rec slog.Record) error {
	level := rec.Level.String() + ":"

	switch rec.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]any, rec.NumAttrs())

	rec.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	t := rec.Time
	if h.location != nil {
		t = t.In(h.location)
	}

	timeStr := t.Format("15:04:05.000")
	msg := color.CyanString(rec.Message)

	h.l.Println(
		timeStr,
		level,
		msg,
		color.WhiteString(string(b)),
	)

	return nil
}

func (h *PlusHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PlusHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

func (h *PlusHandler) WithGroup(name string) slog.Handler {
	return &PlusHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}

func parseLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
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
