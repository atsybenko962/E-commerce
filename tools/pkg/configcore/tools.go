package configcore

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"time"
)

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

// Load метод для чтения конфига из окружения или .env файла
func Load(cfg interface{}, envNamespace string) error {
	return LoadWithEnv(cfg, envNamespace, "")
}

// LoadWithEnv метод для чтения конфига из окружения или .env файла
func LoadWithEnv(cfg interface{}, envNamespace, envFile string) error {
	if envFile == "" {
		envFile = defaultEnvFile
	}

	// Load environment variables from the .env file
	if err := godotenv.Load(envFile); err != nil {
		log.Println("config file is not exists")
	}

	// Parse environment variables into the Config struct
	if err := envconfig.Process(envNamespace, cfg); err != nil {
		log.Fatalf("config not loaded: %s", err)
		return nil
	}

	// Return the loaded configuration
	return nil

}
