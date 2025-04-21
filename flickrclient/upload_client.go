package flickrclient

import (
	"fmt"

	"gopkg.in/masci/flickr.v3"
	"gopkg.in/masci/flickr.v3/photos"
)

type UploadHandler interface {
	UploadPhoto(isPublic bool, isFamily bool, isFriend bool, path string) (string, error)
	DeletePhoto(photoId string) error
}

type UploadClient struct {
	flickrClient *flickr.FlickrClient
}

func (uc *UploadClient) UploadPhoto(isPublic bool, isFamily bool, isFriend bool, path string) (string, error) {
	params := flickr.NewUploadParams()
	// Restrict the photo to private by default
	params.IsPublic = isPublic
	params.IsFamily = isFamily
	params.IsFriend = isFriend
	resp, err := flickr.UploadFile(uc.flickrClient, path, params)
	if err != nil {
		return "", fmt.Errorf("failed to upload photo: %w", err)
	}
	if resp.Status != "ok" {
		return "", fmt.Errorf("upload failed: %s", resp.ErrorMsg())
	}
	return resp.ID, nil
}

func (uc *UploadClient) DeletePhoto(photoId string) error {
	_, err := photos.Delete(uc.flickrClient, photoId)
	if err != nil {
		return fmt.Errorf("failed to delete photo: %w", err)
	}
	return nil
}

func GetUploadHandler() (UploadHandler, error) {
	// Get the flickr client
	client, err := GetFlickrClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get flickr client: %w", err)
	}
	return &UploadClient{
		flickrClient: client,
	}, nil
}
