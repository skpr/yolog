package yolog

import (
	"context"
	"io"
	"log/slog"
	"time"
)

// NewNoMsgJSONHandler creates a new slog Handler that outputs JSON logs without the built-in message key.
func NewNoMsgJSONHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	var o slog.HandlerOptions

	if opts != nil {
		o = *opts // copy so we don't mutate caller's options
	}

	prev := o.ReplaceAttr

	o.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		// Drop the built-in message key. We don't use it in our logs.
		if a.Key == slog.MessageKey {
			return slog.Attr{}
		}

		if a.Key == slog.TimeKey {
			a.Value = slog.StringValue(a.Value.Time().Format(time.DateTime))
		}

		if prev != nil {
			return prev(groups, a)
		}

		return a
	}

	return NoMsgHandler{h: slog.NewJSONHandler(w, &o)}
}

// NoMsgHandler is a slog Handler that omits the built-in message key.
type NoMsgHandler struct{ h slog.Handler }

// Enabled reports whether the handler is enabled for the given level.
func (n NoMsgHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return n.h.Enabled(ctx, level)
}

// Handle processes a log record.
func (n NoMsgHandler) Handle(ctx context.Context, r slog.Record) error { return n.h.Handle(ctx, r) }

// WithAttrs returns a new handler with the given attributes.
func (n NoMsgHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NoMsgHandler{h: n.h.WithAttrs(attrs)}
}

// WithGroup returns a new handler with the given group.
func (n NoMsgHandler) WithGroup(name string) slog.Handler {
	return NoMsgHandler{h: n.h.WithGroup(name)}
}
