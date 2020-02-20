package provider

const NameDockerHub = "dockerhub"

func New(p, name string) Provider {
	switch p {
	case NameDockerHub:
		return NewDockerHub(name)
	}

	return nil
}
