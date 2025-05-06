package flickrclient

import (
	"fmt"
	"strings"

	"gopkg.in/masci/flickr.v3"
	"gopkg.in/masci/flickr.v3/photosets"
)

type SetsClient struct {
	flickrClient *flickr.FlickrClient
}

type SetsLister interface {
	ListUserSets(userid string) (total int, photosetitems []photosets.Photoset, error error)
	GetSetByName(setName, userid string) (string, error)
	GetPhotosInSet(setId, userId string) (photosetitems []photosets.Photo, error error)
}

type SetsEditor interface {
	AddToPhotoSet(photoId, photoSetId string) error
	CreateSet(title, description, primaryPhotoId string) (string, error)
	OrderSet(photosetId string, photoIds []string) error
}

type SetsManager interface {
	SetsLister
	SetsEditor
}

func GetSetsManager() (SetsManager, error) {
	c, err := GetFlickrClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get flickr client: %w", err)
	}

	return &SetsClient{
		flickrClient: c,
	}, nil
}

// ListUserSets returns the list of sets for a user
func (c SetsClient) ListUserSets(userid string) (total int, photosetitems []photosets.Photoset, error error) {
	// Get the list of sets
	response, err := photosets.GetList(c.flickrClient, true, userid, 0)

	if err != nil {
		return 0, nil, fmt.Errorf("failed to get list of sets: %w", err)
	}

	photoSetItems := response.Photosets.Items

	for response.Photosets.Page < response.Photosets.Pages {
		response, err = photosets.GetList(c.flickrClient, true, userid, response.Photosets.Page+1)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to get list of sets: %w", err)
		}
		photoSetItems = append(photoSetItems, response.Photosets.Items...)
	}

	return response.Photosets.Total, photoSetItems, nil
}

func (c SetsClient) AddToPhotoSet(photoId, photoSetId string) error {
	// Get the list of sets
	response, err := photosets.AddPhoto(c.flickrClient, photoSetId, photoId)
	if err != nil {
		return fmt.Errorf("failed to add picture to set: %w", err)
	}

	if response.Status != "ok" {
		return fmt.Errorf("failed to add picture to set: %s", response.Extra)
	}

	return nil
}

func (c SetsClient) GetSetByName(setName, userid string) (string, error) {
	_, sets, err := c.ListUserSets(userid)
	if err != nil {
		return "", fmt.Errorf("failed to get list of sets: %w", err)
	}
	for _, set := range sets {
		if strings.EqualFold(strings.ToLower(set.Title), strings.ToLower(setName)) {
			return set.Id, nil
		}
	}
	return "", nil
}

func (c SetsClient) GetPhotosInSet(setId, userId string) (photosetitems []photosets.Photo, error error) {
	response, err := photosets.GetPhotos(c.flickrClient, true, setId, userId, 0)

	if err != nil {
		return nil, fmt.Errorf("failed to get list of sets: %w", err)
	}

	photoItems := response.Photoset.Photos

	for response.Photoset.Page < response.Photoset.Pages {
		response, err = photosets.GetPhotos(c.flickrClient, true, setId, userId, response.Photoset.Page+1)
		if err != nil {
			return nil, fmt.Errorf("failed to get list of sets: %w", err)
		}
		photoItems = append(photoItems, response.Photoset.Photos...)
	}
	return photoItems, nil
}

func (c SetsClient) CreateSet(title, description, primaryPhotoId string) (string, error) {
	response, err := photosets.Create(c.flickrClient, title, description, primaryPhotoId)
	if err != nil {
		return "", fmt.Errorf("failed to create set: %w", err)
	}

	if response.Status != "ok" {
		return "", fmt.Errorf("failed to create set: %s", response.Extra)
	}

	return response.Set.Id, nil
}

func (c SetsClient) OrderSet(photosetId string, photoIds []string) error {
	if len(photoIds) < 1 {
		return fmt.Errorf("at least 1 id must be specified to order a set")
	}
	_, err := photosets.ReorderPhotos(c.flickrClient, photosetId, photoIds[0], photoIds)
	if err != nil {
		return fmt.Errorf("failed to create set: %w", err)
	}
	return nil
}
