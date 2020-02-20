package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Count   int      `json:"count"`
	Next    string   `json:"next"`
	Results []result `json:"results"`
}

type result struct {
	Name        string    `json:"name"`
	LastUpdated time.Time `json:"last_updated"`
}

type DockerHub struct {
	name   string
	next   string
	logger log.Logger
}

func NewDockerHub(name string) *DockerHub {
	logger := log.Logger{}
	logger.SetOutput(os.Stdout)
	return &DockerHub{name: name, next: "", logger: logger}
}

const apiUrl = "https://registry.hub.docker.com/v2/repositories/%s/tags/?page_size=100"

func (p *DockerHub) Fetch() ([]*Image, error) {
	resp, err := p.readImages()
	if err != nil {
		return nil, err
	}

	p.next = resp.Next

	return p.parse(resp)
}

func (p *DockerHub) HasNext() bool {
	return p.next != ""
}

func (p *DockerHub) readImages() (*Response, error) {
	url := fmt.Sprintf(apiUrl, p.name)
	if p.next != "" {
		url = p.next
	}

	p.logger.Println("fetch:" + url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse Response
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse, nil
}

func (p *DockerHub) parse(resp *Response) ([]*Image, error) {
	images := make([]*Image, 0)
	for _, i := range resp.Results {
		img := NewImage("php", i.Name, i.LastUpdated)

		images = append(images, img)
	}

	return images, nil
}
