package config

import (
	"io/ioutil"

	"github.com/shin1x1/docker-base-chaser/internal/handler"
	"gopkg.in/yaml.v2"
)

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c Config) CreateTargets(conf *Config) *handler.Targets {
	targets := handler.Targets{}

	for _, cImage := range conf.Images {
		t := handler.Target{
			Provider: cImage.Base.Provider,
			Image:    cImage.Base.Image,
			Template: &handler.Template{
				Source:      cImage.Template.Source,
				Destination: cImage.Template.Destination,
			},
			Command: cImage.Command,
			Tags:    make([]*handler.TargetTag, 0),
		}

		for _, cTag := range cImage.Base.Tags {
			tag := handler.NewTargetTag(cTag.Pattern, cTag.Version, "", cTag.Os, cTag.Architecture, nil)
			t.Tags = append(t.Tags, tag)
		}

		targets = append(targets, &t)
	}

	return &targets
}
