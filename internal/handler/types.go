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
			if tag.IsExecuted() {
				return true
			}
		}
	}

	return false
}

// Merge create new targets.
// New targets is based on ts, and override the tag that ts1 has it.
func (ts Targets) Merge(ts1 *Targets) *Targets {
	targets := Targets{}

	for _, t := range ts {
		for _, t1 := range *ts1 {
			if t.Image != t1.Image {
				continue
			}

			for k, tg := range t.Tags {
				for _, tg1 := range t1.Tags {
					if tg.Pattern == tg1.Pattern {
						t.Tags[k] = tg1
						break
					}
				}
			}
		}

		targets = append(targets, t)
	}

	return &targets
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
		if tag.IsNotMatched() {
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
	Tag          string
	Pattern      string
	Version      string
	Os           string
	Architecture string
	LastUpdated  *time.Time
	mode         TargetTagMode

	re   *regexp.Regexp
	cons *semver.Constraints
}

type TargetTagMode int

const (
	notMatched TargetTagMode = iota
	matched
	updated
	executed
	notExecuted
)

func NewTargetTag(pattern, version, resolved, os, architecture string, lastUpdated *time.Time) *TargetTag {
	t := TargetTag{
		Tag:          resolved,
		Pattern:      pattern,
		Version:      version,
		Os:           os,
		Architecture: architecture,
		LastUpdated:  lastUpdated,
		mode:         notMatched,
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

func (t *TargetTag) CanMatch() bool {
	return t.mode == notMatched
}

func (t *TargetTag) CanExecute() bool {
	return t.mode == updated
}

func (t *TargetTag) ShouldUpdate(image *provider.Image) bool {
	if t.LastUpdated == nil {
		return true
	}
	return image.LastUpdated.After(*t.LastUpdated)
}

func (t *TargetTag) Matched() {
	t.mode = matched
}

func (t *TargetTag) NotExecuted() {
	t.mode = notExecuted
}

func (t *TargetTag) Executed() {
	t.mode = executed
}

func (t *TargetTag) Update(image *provider.Image) {
	t.Tag = image.Tag
	t.LastUpdated = &image.LastUpdated
	t.mode = updated
}

func (t *TargetTag) IsNotMatched() bool {
	return t.mode == notMatched
}

func (t *TargetTag) IsExecuted() bool {
	return t.mode == executed
}
