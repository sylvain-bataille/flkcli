package flickrclient

import (
	"gopkg.in/masci/flickr.v3"
	"gopkg.in/masci/flickr.v3/photos"
)

type PhotoClient struct {
	flickrClient *flickr.FlickrClient
}

type PhotoInfoRetriever interface {
	GetPhotoInfo(photoId string) (photos.PhotoInfo, error)
}

func GetPhotoClient() (*PhotoClient, error) {
	c, err := GetFlickrClient()
	if err != nil {
		return nil, err
	}
	return &PhotoClient{
		flickrClient: c,
	}, nil
}

func (p *PhotoClient) GetPhotoInfo(photoId string) (photos.PhotoInfo, error) {
	resp, err := photos.GetInfo(p.flickrClient, photoId, "")
	if err != nil {
		return photos.PhotoInfo{}, err
	}
	return resp.Photo, nil
}
