package handler

import (
	"regexp"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/shin1x1/docker-base-chaser/internal/provider"
)

type FetchedImages []*FetchedImage

type FetchedImage struct {
	Image       string
	Tag         string
	LastUpdated time.Time
}

// Target Models
type Targets []*Target

func (t *Targets) AreUpdated() bool {
	for _, target := range *t {
		for _, tag := range target.Tags {
			if tag.IsUpdated() {
				return true
			}
		}
	}

	return false
}

type Target struct {
	Provider string
	Image    string
	Template *Template
	Command  string
	Tags     []*TargetTag
}

func (t *Target) Done() bool {
	for _, tag := range t.Tags {
		if !tag.IsFetched() {
			return false
		}
	}

	return true
}

type Template struct {
	Source      string
	Destination string
}

type TargetTag struct {
	Tag         string
	Pattern     string
	Version     string
	LastUpdated *time.Time
	Fetched     *time.Time
	Updated     bool

	re   *regexp.Regexp
	cons *semver.Constraints
}

func NewTargetTag(pattern, version, resolved string, lastUpdated *time.Time) *TargetTag {
	t := TargetTag{
		Tag:         resolved,
		Pattern:     pattern,
		Version:     version,
		LastUpdated: lastUpdated,
		Updated:     false,
	}
	t.re = regexp.MustCompile(t.Pattern)

	var err error
	if t.cons, err = semver.NewConstraint(t.Version); err != nil {
		panic(err)
	}

	return &t
}

func (t *TargetTag) MatchPattern(v string) bool {
	return t.re.MatchString(v)
}

func (t *TargetTag) CheckVersion(v *semver.Version) bool {
	return t.cons.Check(v)
}

func (t *TargetTag) Before(image *provider.Image) bool {
	if t.LastUpdated == nil {
		return true
	}
	return image.LastUpdated.After(*t.LastUpdated)
}

func (t *TargetTag) Update(image *provider.Image) {
	t.Tag = image.Tag
	t.LastUpdated = &image.LastUpdated
	t.Fetched = &image.LastUpdated
	t.Updated = true
}

func (t *TargetTag) IsFetched() bool {
	return t.Fetched != nil
}

func (t *TargetTag) IsUpdated() bool {
	return t.Updated
}
