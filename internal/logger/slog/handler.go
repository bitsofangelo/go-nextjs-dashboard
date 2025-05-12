package slog

import (
	"context"
	"log/slog"

	"go-nextjs-dashboard/internal/http"
)

type ctxHandler struct {
	base slog.Handler
}

func (h *ctxHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.base.Enabled(ctx, level)
}

func (h *ctxHandler) Handle(ctx context.Context, r slog.Record) error {
	// Pull req_id (or any other attrs) from context
	if reqID, ok := http.ReqID(ctx); ok {
		r.Add("req_id", slog.StringValue(reqID))
	}

	return h.base.Handle(ctx, r)
}

func (h *ctxHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ctxHandler{base: h.base.WithAttrs(attrs)}
}

func (h *ctxHandler) WithGroup(name string) slog.Handler {
	return &ctxHandler{base: h.base.WithGroup(name)}
}
