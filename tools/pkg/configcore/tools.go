package configcore

import "time"

const defaultEnvFile = ".env"

type ServerConfig struct {
	ServerHost       string        `envconfig:"SERVER_HOST"`
	ServerPort       string        `envconfig:"SERVER_PORT"`
	ServerType       string        `envconfig:"SERVER_TYPE" default:"grpc"`
	HttpReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"10s"`
	HttpWriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"10s"`
	HttpIdleTimeout  time.Duration `envconfig:"HTTP_IDLE_TIMEOUT" default:"1200s"`
}

type Observer struct {
	ServiceName    string        `envconfig:"SERVICE_NAME"`
	ServiceVersion string        `envconfig:"SERVICE_VERSION" default:"v1"`
	TraceTimeout   time.Duration `envconfig:"TRACE_TIMEOUT" default:"1s"`
	MetricsTimeout time.Duration `envconfig:"METRICS_TIMEOUT" default:"3s"`
}
