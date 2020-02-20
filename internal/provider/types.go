package provider

import (
	"time"
)

type Provider interface {
	Fetch() ([]*Image, error)
	HasNext() bool
}

type Image struct {
	Image       string
	Tag         string
	LastUpdated time.Time
}

func NewImage(image string, tag string, lastUpdated time.Time) *Image {
	return &Image{Image: image, Tag: tag, LastUpdated: lastUpdated}
}
