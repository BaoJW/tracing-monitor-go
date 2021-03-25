package tracer

import (
	"io"
	"log"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"

	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func NewZipkinTracer() {

}

func NewJaegerTracer(serviceName string, withEnv bool, options ...jaegercfg.Option) (opentracing.Tracer, io.Closer, error) {
	if withEnv {
		return newJaegerTracerFromEnv()
	}

	return newDefaultJaegerTracer(serviceName, options...)
}

func newJaegerTracerFromEnv() (opentracing.Tracer, io.Closer, error) {
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		// parsing errors might happen here, such as when we get a string where we expect a number
		// todo 替换日志库
		log.Printf("Could not parse Jaeger env vars: %s", err.Error())
		return nil, nil, err
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		// todo 替换日志库
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return nil, nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, nil
}

func newDefaultJaegerTracer(serviceName string, options ...jaegercfg.Option) (opentracing.Tracer, io.Closer, error) {
	// Recommended configuration for production.
	cfg := jaegercfg.Configuration{}

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	// todo 用自己的prometheus和zap去替换
	options = append(options, jaegercfg.Logger(jLogger), jaegercfg.Metrics(jMetricsFactory))
	closer, err := cfg.InitGlobalTracer(serviceName, options...)
	if err != nil {
		// todo 替换日志库
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return nil, nil, err
	}

	return opentracing.GlobalTracer(), closer, nil
}
