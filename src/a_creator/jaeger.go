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
	samplerParam, _ := strconv.ParseFloat(os.Getenv("JAEGER_SAMPLER_PARAM"), 64)
	reporterLogSpans, _ := strconv.ParseBool(os.Getenv("JAEGER_LOGS_ENABLED"))

	cfg := &config.Configuration{
		ServiceName: os.Getenv("JAEGER_SERVICE_NAME"),
		Sampler: &config.SamplerConfig{
			Type:  os.Getenv("JAEGER_SAMPLER_TYPE"),
			Param: samplerParam,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: reporterLogSpans,
		},
	}

	return cfg.NewTracer(config.Logger(jaeger.StdLogger))
}
