package flickrclient

import (
	"fmt"

	"gopkg.in/masci/flickr.v3"
	"gopkg.in/masci/flickr.v3/people"
)

type PeopleClient struct {
	flickrClient *flickr.FlickrClient
}

type IdResolver interface {
	ResolveId(username string) (id string, error error)
}

type PeopleManager interface {
	IdResolver
}

func GetPeopleManager() (PeopleManager, error) {
	c, err := GetFlickrClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get flickr client: %w", err)
	}
	return &PeopleClient{
		flickrClient: c,
	}, nil
}

func (p PeopleClient) ResolveId(username string) (id string, error error) {
	user_response, err := people.FindByUsername(p.flickrClient, username)
	if err != nil {
		return "", fmt.Errorf("failed to resolve username: %w", err)
	}
	return user_response.User.Id, nil
}
