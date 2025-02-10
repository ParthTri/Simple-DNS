package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ReadConfig(configs []string) (Configs, error) {
	var config Configs

	for _, path := range configs {
		b := Domain{}

		content, err := os.ReadFile(path)
		if err != nil {
			return config, err
		}

		err = yaml.Unmarshal(content, &b)
		if err != nil {
			return config, err
		}

		config.Domains = append(config.Domains, b)
	}

	return config, nil
}


func ReadConfigDir(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	var configFiles []string = []string{};
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		name := strings.Split(file.Name(), ".")
		suffix := name[len(name)-1]
		if suffix == "yaml" || suffix == "yml" {
			configFiles = append(configFiles, path + "/" + file.Name())
		}
	}

	return configFiles, err
}
