package config

import (
	"github.com/shin1x1/docker-base-chaser/internal/handler"
	"reflect"
	"testing"
	"time"
)

func TestConfig_LoadTargets(t *testing.T) {
	c, err := LoadConfig("tests/base-chaser.yaml")
	if err != nil {
		t.Error(err)
	}

	targets := c.LoadTargets()

	t.Run("test", func(t *testing.T) {
		lastUpdated, _ := time.Parse(time.RFC3339, "2021-07-09T19:34:56.308213Z")

		want := []*handler.TargetTag{
			handler.NewTargetTag(
				"^7\\-cli\\-buster$",
				"7",
				"7-cli-buster",
				"linux",
				"amd64",
				&lastUpdated,
			),
			handler.NewTargetTag(
				"^8\\-cli\\-buster$",
				"8",
				"", // unresolved
				"linux",
				"amd64",
				nil, // unresolved
			),
		}
		got := (*targets)[0].Tags
		if !reflect.DeepEqual(got, want) {
			t.Errorf("LoadTargets() = %v, want %v", got, want)
		}
	})
}
