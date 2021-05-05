package config

import (
	"os"

	"github.com/joho/godotenv"
)

type GoDotEnvProvider struct {
	configFolder string
	logger       logger
}

type logger interface {
	Log(args ...interface{})
	Logf(format string, a ...interface{})
}

func NewGoDotEnvProvider(l logger, configFolder string) *GoDotEnvProvider {
	provider := &GoDotEnvProvider{
		configFolder: configFolder,
		logger:       l,
	}

	provider.readConfig(configFolder)

	return provider
}

// readConfig(logger Logger) loads the environment variables from .env file
// Priority Order is Environment Variable > .env.X file > .env file
// if there is a need to overwrite any of the environment variable present in the ./env
// then it can be done by creating .env.local file
// or by specifying the file prefix in environment variable GOFR_ENV.
func (g *GoDotEnvProvider) readConfig(confLocation string) {
	defaultFile := confLocation + "/.env"

	env := os.Getenv("GOFR_ENV")
	if env == "" {
		env = "local"
	}

	overrideFile := confLocation + "/." + env + ".env"

	err := godotenv.Load(overrideFile)
	if err == nil {
		g.logger.Log("Loaded config from file: ", overrideFile)
	}

	err = godotenv.Load(defaultFile)
	if err == nil {
		g.logger.Log("Loaded config from file: ", defaultFile)
	}
}

func (g *GoDotEnvProvider) Get(key string) string {
	return os.Getenv(key)
}

func (g *GoDotEnvProvider) GetOrDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	return defaultValue
}
