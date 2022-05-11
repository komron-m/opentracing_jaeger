package main

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
)

// in production these variables should be resolved from .env or some kind of configuration file
const (
	serviceName  = "a_creator"
	samplerType  = "const"
	samplerParam = 1.0
)

func NewJaegerOpentracingTracer() (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  samplerType,
			Param: samplerParam,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	return cfg.NewTracer(config.Logger(jaeger.StdLogger))
}
