package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	ztrace "github.com/hduhelp/go-zero/core/trace"
	"github.com/hduhelp/go-zero/rest/chain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TestOtelHandler(t *testing.T) {
	ztrace.StartAgent(ztrace.Config{
		Name:     "go-zero-test",
		Endpoint: "http://localhost:14268/api/traces",
		Batcher:  "jaeger",
		Sampler:  1.0,
	})
	defer ztrace.StopAgent()

	for _, test := range []string{"", "bar"} {
		t.Run(test, func(t *testing.T) {
			h := chain.New(TracingHandler("foo", test)).Then(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
					spanCtx := trace.SpanContextFromContext(ctx)
					assert.True(t, spanCtx.IsValid())
				}))
			ts := httptest.NewServer(h)
			defer ts.Close()

			client := ts.Client()
			err := func(ctx context.Context) error {
				ctx, span := otel.Tracer("httptrace/client").Start(ctx, "test")
				defer span.End()

				req, _ := http.NewRequest("GET", ts.URL, nil)
				otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

				res, err := client.Do(req)
				assert.Nil(t, err)
				return res.Body.Close()
			}(context.Background())

			assert.Nil(t, err)
		})
	}
}

func TestDontTracingSpanName(t *testing.T) {
	ztrace.StartAgent(ztrace.Config{
		Name:     "go-zero-test",
		Endpoint: "http://localhost:14268/api/traces",
		Batcher:  "jaeger",
		Sampler:  1.0,
	})

	DontTraceSpan("bar")

	for _, test := range []string{"", "bar", "foo"} {
		t.Run(test, func(t *testing.T) {
			h := chain.New(TracingHandler("foo", test)).Then(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					spanCtx := trace.SpanContextFromContext(r.Context())
					if test == "bar" {
						assert.False(t, spanCtx.IsValid())
						return
					}

					assert.True(t, spanCtx.IsValid())
				}))
			ts := httptest.NewServer(h)
			defer ts.Close()

			client := ts.Client()
			err := func(ctx context.Context) error {
				ctx, span := otel.Tracer("httptrace/client").Start(ctx, "test")
				defer span.End()

				req, _ := http.NewRequest("GET", ts.URL, nil)
				otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

				res, err := client.Do(req)
				assert.Nil(t, err)
				return res.Body.Close()
			}(context.Background())

			assert.Nil(t, err)
		})
	}
}
