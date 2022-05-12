package main

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"os"
	"strconv"
)

func NewJaegerOpentracingTracer() (opentracing.Tracer, io.Closer, error) {
	sp, _ := strconv.ParseFloat(os.Getenv("JAEGER_SAMPLER_PARAM"), 64)
	repLogSpans, _ := strconv.ParseBool(os.Getenv("JAEGER_LOGS_ENABLED"))

	cfg := &config.Configuration{
		ServiceName: os.Getenv("JAEGER_SERVICE_NAME"),
		Sampler: &config.SamplerConfig{
			Type:  os.Getenv("JAEGER_SAMPLER_TYPE"),
			Param: sp,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: repLogSpans,
		},
	}

	return cfg.NewTracer(config.Logger(jaeger.StdLogger))
}
