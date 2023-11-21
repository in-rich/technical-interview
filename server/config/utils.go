package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type EnvLoader map[string][]byte

func loadEnv(files EnvLoader, out interface{}) error {
	if ENV == "" {
		ENV = DevENV
	}

	for env, file := range files {
		if env != DefaultENV && env != ENV {
			continue
		}

		if err := yaml.Unmarshal([]byte(os.ExpandEnv(string(file))), out); err != nil {
			return err
		}
	}

	return nil
}
