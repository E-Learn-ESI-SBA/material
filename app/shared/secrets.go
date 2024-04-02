package shared

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"log"
)

func GetSecrets() *koanf.Koanf {
	k := koanf.New("/")
	err := k.Load(file.Provider("config/secrets/env.yaml"), yaml.Parser())
	if err != nil {
		log.Fatal("Env file not found")

	}
	return k
}
