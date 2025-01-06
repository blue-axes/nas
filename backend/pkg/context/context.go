package context

import "context"

type (
	Context struct {
		context.Context
		TraceID string
	}
	Option func(c *Context)
)

func WithTraceID(traceID string) Option {
	return func(c *Context) {
		c.TraceID = traceID
	}
}

func WithCtx(ctx context.Context) Option {
	return func(c *Context) {
		c.Context = ctx
	}
}

func New(options ...Option) *Context {
	ctx := &Context{
		Context: context.Background(),
	}
	for _, opt := range options {
		opt(ctx)
	}
	if ctx.Context == nil {
		ctx.Context = context.Background()
	}

	return ctx
}
