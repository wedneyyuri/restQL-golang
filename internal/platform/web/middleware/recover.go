package middleware

import (
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/http"
	"runtime/debug"
)

type Recover struct {
	log *logger.Logger
}

func NewRecover(log *logger.Logger) Middleware {
	return Recover{log: log}
}

func (r Recover) Apply(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if reason := recover(); reason != nil {
				err := errors.Errorf("reason : %v", reason)
				r.log.Error("application recovered from panic", err, "stack", string(debug.Stack()))

				ctx.SetStatusCode(http.StatusInternalServerError)
			}
		}()

		h(ctx)
	}
}
