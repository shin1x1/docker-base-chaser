package handler

import (
	"testing"
	"time"

	"github.com/shin1x1/docker-base-chaser/internal/provider"
)

func TestRootHandler_updateTags(t *testing.T) {
	now := time.Date(2020, 2, 1, 12, 34, 56, 0, time.Local)
	past := now.Add(time.Duration(1) * time.Second * -1)

	images := []*provider.Image{
		{
			Image:       "php",
			Tag:         "7.4.2",
			LastUpdated: now,
		},
		{
			Image:       "php",
			Tag:         "7.3.1",
			LastUpdated: now,
		},
		{
			Image:       "php",
			Tag:         "7.4",
			LastUpdated: now,
		},
		{
			Image:       "php",
			Tag:         "7.3",
			LastUpdated: now,
		},
	}
	target := &Target{
		Image: "php",
		Tags: []*TargetTag{
			NewTargetTag("^7.3$", "7.3", "7.3", &now),
			NewTargetTag("^7.3.[0-9]+$", "^7.3.0", "7.3.0", &past),
			NewTargetTag("^7.4$", "7.4", "7.4", &past),
			NewTargetTag("^7.4.[0-9]+$", "^7.4.0", "", nil),
			NewTargetTag("^7.5.[0-9]+$", "^7.5.0", "", nil),
		},
	}

	r := NewRootHandler(false)
	if err := r.updateTags(images, target); err != nil {
		t.Errorf("%+v", err)
	}

	expected := []struct {
		tag         string
		pattern     string
		version     string
		lastUpdated *time.Time
		mode        TargetTagMode
	}{
		{
			tag:         "7.3",
			pattern:     "^7.3$",
			version:     "7.3",
			lastUpdated: &now,
			mode:        notExecuted,
		},
		{
			tag:         "7.3.1",
			pattern:     "^7.3.[0-9]+$",
			version:     "^7.3.0",
			lastUpdated: &now,
			mode:        updated,
		},
		{
			tag:         "7.4",
			pattern:     "^7.4$",
			version:     "7.4",
			lastUpdated: &now,
			mode:        updated,
		},
		{
			tag:         "7.4.2",
			pattern:     "^7.4.[0-9]+$",
			version:     "^7.4.0",
			lastUpdated: &now,
			mode:        updated,
		},
		{
			tag:         "",
			pattern:     "^7.5.[0-9]+$",
			version:     "^7.5.0",
			lastUpdated: nil,
			mode:        notMatched,
		},
	}

	for i, exp := range expected {
		tag := target.Tags[i]

		if exp.tag != tag.Tag {
			t.Errorf("%d: expected: %s actual:%s", i, exp.tag, tag.Tag)
		}
		if exp.pattern != tag.Pattern {
			t.Errorf("%d: expected: %s actual:%s", i, exp.pattern, tag.Pattern)
		}
		if exp.version != tag.Version {
			t.Errorf("%d: expected: %s actual:%s", i, exp.version, tag.Version)
		}
		if exp.lastUpdated == nil {
			if tag.LastUpdated != nil {
				t.Errorf("%d: lastUpdated should be nil: %+v", i, *tag.LastUpdated)
			}
		} else {
			if *exp.lastUpdated != *tag.LastUpdated {
				t.Errorf("%d: expected: %+v actual:%+v", i, *exp.lastUpdated, *tag.LastUpdated)
			}
		}
		if exp.mode != tag.mode {
			t.Errorf("expected: %d actual:%d", exp.mode, tag.mode)
		}
	}
}
