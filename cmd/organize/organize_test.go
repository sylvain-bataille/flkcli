package organize

import (
	"fmt"
	"testing"

	"gopkg.in/masci/flickr.v3/photos"
	"gopkg.in/masci/flickr.v3/photosets"
)

type SetsManagerMock struct {
	Photos     []photosets.Photo
	DatesTaken map[string]string // "2006-01-02 15:04:05"
	Name       string
	Id         string
}

type PhotoRetreiverMock struct {
	SetsManager SetsManagerMock
}

func (pr PhotoRetreiverMock) GetPhotoInfo(photoId string) (photos.PhotoInfo, error) {
	pi := photos.PhotoInfo{}
	for _, p := range pr.SetsManager.Photos {
		if p.Id != photoId {
			continue
		}
		pi.Id = p.Id
		pi.Title = p.Title
		pi.Dates.Taken = pr.SetsManager.DatesTaken[p.Id]
		return pi, nil
	}
	return pi, fmt.Errorf("Not found")
}

func TestComparePhotoByDate(t *testing.T) {
	// Create two mock PhotoInfo objects with different dates
	photo1 := photos.PhotoInfo{
		Id: "1",
	}
	photo1.Dates.Taken = "2023-10-01 12:00:00"
	photo2 := photos.PhotoInfo{
		Id: "2",
	}
	photo2.Dates.Taken = "2023-10-02 12:00:00"
	// Compare the two photos
	result := comparePhotosByDate(photo1, photo2)

	// Check if the result is as expected
	if result >= 0 {
		t.Errorf("Expected photo1 to be less than photo2, but got %d", result)
	}

}

func (s SetsManagerMock) AddToPhotoSet(photoId, photoSetId string) error {
	fmt.Println("Not implemented")
	return nil
}

func (s SetsManagerMock) CreateSet(title, description, primaryPhotoId string) (string, error) {
	fmt.Println("Not implemented")
	return "", nil
}

func (s SetsManagerMock) GetPhotosInSet(setId, userId string) (photosetitems []photosets.Photo, error error) {
	return s.Photos, nil
}

func (s SetsManagerMock) GetSetByName(setName, userid string) (string, error) {
	return s.Id, nil
}

func (s SetsManagerMock) ListUserSets(userid string) (total int, photosetitems []photosets.Photoset, error error) {
	return 0, nil, nil
}

func (s *SetsManagerMock) OrderSet(photosetId string, photoIds []string) error {
	newPhotos := []photosets.Photo{}
	for _, id := range photoIds {
		for _, p := range s.Photos {
			if p.Id == id {
				newPhotos = append(newPhotos, p)
			}
		}
	}
	s.Photos = newPhotos
	return nil
}

func TestSortSetByNameAscend(t *testing.T) {
	set := SetsManagerMock{
		Name: "myset",
		Id:   "111",
		Photos: []photosets.Photo{
			{
				Id:    "2",
				Title: "B-photo",
			},
			{
				Id:    "1",
				Title: "A-photo",
			},
			{
				Id:    "3",
				Title: "C-photo",
			},
		},
		DatesTaken: map[string]string{
			"1": "2006-01-02 15:04:05",
			"2": "2006-01-02 15:04:05",
			"3": "2006-01-02 15:04:05",
		},
	}
	pr := PhotoRetreiverMock{
		SetsManager: set,
	}
	runSortSetCmd(&set, pr, "myset", "title", false)

	if set.Photos[0].Id != "1" {
		t.Errorf("Expected photo 1 to be first")
	}
	if set.Photos[1].Id != "2" {
		t.Errorf("Expected photo 2 to be second")
	}
	if set.Photos[2].Id != "3" {
		t.Errorf("Expected photo 3 to be last")
	}
}
