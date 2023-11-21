package config

import (
	_ "embed"
	"log"
)

//go:embed app.yml
var appFile []byte

//go:embed app-dev.yml
var appDevFile []byte

//go:embed app-prod.yml
var appProdFile []byte

type appConfig struct {
	Name      string `yaml:"name"`
	Port      int    `yaml:"port"`
	ProjectID string `yaml:"project_id"`
}

var App *appConfig

func init() {
	cfg := new(appConfig)
	if err := loadEnv(EnvLoader{DefaultENV: appFile, ProdENV: appProdFile, DevENV: appDevFile}, cfg); err != nil {
		log.Fatalf("error loading app configuration: %v\n", err)
	}

	App = cfg
}
