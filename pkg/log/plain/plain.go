/*
Package plain implements a plain handler, similar to the default one, but that
also handles the custom levels sets in ../log.go
*/
package plain

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"

	"github.com/nobe4/gh-ln/pkg/log"
)

const buflen = 1024

type Handler struct {
	opts log.Options
	mu   *sync.Mutex
	out  io.Writer

	indent int
}

func New(out io.Writer, o log.Options) *Handler {
	h := &Handler{
		out:    out,
		opts:   o,
		mu:     &sync.Mutex{},
		indent: 0,
	}

	return h
}

func (h *Handler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.opts.Level.Level()
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	level := ""

	switch r.Level {
	case log.LevelDebug:
		level = "[D]"
	case log.LevelWarn:
		level = "[W]"
	case log.LevelError:
		level = "[E]"
	case log.LevelNotice:
		level = "[N]"
	case log.LevelInfo:
		level = "[I]"

	case log.LevelGroup:
		level = "[G]"
	case log.LevelGroupEnd:
		h.indent = 0
		level = "[/G]"
		r.Message = "\n"
	}

	buf := make([]byte, 0, buflen)

	buf = fmt.Appendf(buf,
		"%*s%s %s %s\n",
		h.indent,
		"",
		level,
		r.Message,
		h.formatAttrs(r),
	)

	if r.Level == log.LevelGroup {
		h.indent = 2
	}

	return h.write(buf)
}

func (h *Handler) WithAttrs(_ []slog.Attr) slog.Handler {
	// TODO: implement?
	return h
}

func (h *Handler) WithGroup(_ string) slog.Handler {
	// TODO: implement?
	return h
}

func (h *Handler) write(p []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := h.out.Write(p)

	return fmt.Errorf("%w: %w", log.ErrCannotWrite, err)
}

func (*Handler) formatAttrs(r slog.Record) string {
	attrs := []string{}

	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s=%s", a.Key, a.Value))

		return true
	})

	if len(attrs) > 0 {
		return strings.Join(attrs, " ")
	}

	return ""
}
