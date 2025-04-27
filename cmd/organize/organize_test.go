package organize

import (
	"testing"

	"gopkg.in/masci/flickr.v3/photos"
)

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
