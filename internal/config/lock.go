package config

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/shin1x1/docker-base-chaser/internal/handler"
	"gopkg.in/yaml.v2"
)

type Lock struct {
	LastUpdatedAt *time.Time `yaml:"last_updated_at"`
	Images        []Image    `yaml:"images"`
}

func CreateLockPath(path string) string {
	return strings.Replace(path, ".yaml", ".lock", 1)
}

func LoadLock(path string) (*Lock, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	l := Lock{}
	if err := yaml.Unmarshal(data, &l); err != nil {
		return nil, err
	}

	return &l, nil
}

func NewLockWithTargets(targets *handler.Targets, now *time.Time) *Lock {
	l := Lock{
		LastUpdatedAt: now,
	}

	var images []Image
	for _, target := range *targets {
		image := Image{
			Template: Template{
				Source:      target.Template.Source,
				Destination: target.Template.Destination,
			},
			Command: target.Command,
			Base: ImageBase{
				Provider: target.Provider,
				Image:    target.Image,
				Tags:     nil,
			},
		}

		var tags []Tag
		for _, tag := range target.Tags {
			t := Tag{
				Pattern:       tag.Pattern,
				Version:       tag.Version,
				Resolved:      tag.Tag,
				Os:            tag.Os,
				Architecture:  tag.Architecture,
				LastUpdatedAt: tag.LastUpdated,
			}

			tags = append(tags, t)
		}
		image.Base.Tags = tags

		images = append(images, image)
	}
	l.Images = images

	return &l
}

func (l *Lock) CreateTargets() *handler.Targets {
	targets := handler.Targets{}

	for _, lImage := range l.Images {
		t := handler.Target{
			Provider: lImage.Base.Provider,
			Image:    lImage.Base.Image,
			Template: &handler.Template{
				Source:      lImage.Template.Source,
				Destination: lImage.Template.Destination,
			},
			Command: lImage.Command,
			Tags:    make([]*handler.TargetTag, 0),
		}

		for _, lTag := range lImage.Base.Tags {
			tag := handler.NewTargetTag(lTag.Pattern, lTag.Version, lTag.Resolved, lTag.Os, lTag.Architecture, lTag.LastUpdatedAt)
			t.Tags = append(t.Tags, tag)
		}

		targets = append(targets, &t)
	}

	return &targets
}

func (l *Lock) Save(path string) error {
	b, err := yaml.Marshal(l)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, 0666)
}
