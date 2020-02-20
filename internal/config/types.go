package config

import "time"

type Config struct {
	Images []Image `yaml:"images"`
}

type Image struct {
	Template Template  `yaml:"template"`
	Command  string    `yaml:"command"`
	Base     ImageBase `yaml:"base"`
}

type Template struct {
	Source      string `yaml:"src"`
	Destination string `yaml:"dest"`
}

type ImageBase struct {
	Provider string `yaml:"provider"`
	Image    string `yaml:"image"`
	Tags     []Tag  `yaml:"tags"`
}

type Tag struct {
	Pattern       string     `yaml:"pattern"`
	Version       string     `yaml:"version"`
	Resolved      string     `yaml:"resolved"`
	LastUpdatedAt *time.Time `yaml:"last_updated_at"`
}
