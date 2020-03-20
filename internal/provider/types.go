package provider

import (
	"time"
)

type Provider interface {
	Fetch() ([]*Image, error)
	HasNext() bool
}

type Image struct {
	Image        string
	Tag          string
	Os           string
	Architecture string
	LastUpdated  time.Time
}

func NewImage(image string, tag string, os string, architecture string, lastUpdated time.Time) *Image {
	return &Image{Image: image, Tag: tag, Os: os, Architecture: architecture, LastUpdated: lastUpdated}
}
