package upload

import (
	"flkcli/flickrclient"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"gopkg.in/masci/flickr.v3/photosets"
)

func TestSyncModeWithoutSetShouldFail(t *testing.T) {
	args := []string{}
	syncMode := true
	overwrite := false
	delete := false
	albumName, err := validateArgs(args, syncMode, overwrite, delete)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if albumName != "" {
		t.Errorf("Expected empty album name, got '%s'", albumName)
	}
}

func TestOverwriteAndDeleteWithoutSyncShouldFail(t *testing.T) {
	args := []string{"albumName"}
	syncMode := false
	overwrite := true
	delete := true
	albumName, err := validateArgs(args, syncMode, overwrite, delete)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if albumName != "" {
		t.Errorf("Expected empty album name, got '%s'", albumName)
	}
}

func TestOverWriteOrDeleteWithSyncShouldPass(t *testing.T) {
	args := []string{"album1"}
	syncMode := false
	overwrite := true
	delete := false
	albumName, err := validateArgs(args, syncMode, overwrite, delete)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if albumName != "" {
		t.Errorf("Expected empty album name, got '%s'", albumName)
	}
}

func TestAddToSetsWithSyncShouldFail(t *testing.T) {
	args := []string{}
	syncMode := true
	overwrite := false
	delete := false
	albumName, err := validateArgs(args, syncMode, overwrite, delete)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if albumName != "" {
		t.Errorf("Expected empty album name, got '%s'", albumName)
	}
}

func TestAddToSetsWithoutSyncShouldPass(t *testing.T) {
	args := []string{}
	syncMode := false
	overwrite := false
	delete := false
	albumName, err := validateArgs(args, syncMode, overwrite, delete)
	if err != nil {
		t.Errorf("Expected nil error, got '%s'", err)
	}
	if albumName != "" {
		t.Errorf("Expected empty album name, got '%s'", albumName)
	}
}

func TestValidateArgsWithAlbumName(t *testing.T) {
	args := []string{"album1"}
	syncMode := false
	overwrite := false
	delete := false
	albumName, err := validateArgs(args, syncMode, overwrite, delete)
	if err != nil {
		t.Errorf("Expected nil error, got '%s'", err)
	}
	if albumName != "album1" {
		t.Errorf("Expected album name 'album1', got '%s'", albumName)
	}
}

func TestValidateArgsWithAlbumNameAndSync(t *testing.T) {
	args := []string{"album1"}
	syncMode := true
	overwrite := false
	delete := false
	albumName, err := validateArgs(args, syncMode, overwrite, delete)
	if err != nil {
		t.Errorf("Expected nil error, got '%s'", err)
	}
	if albumName != "album1" {
		t.Errorf("Expected album name 'album1', got '%s'", albumName)
	}
}

type SetsManagerMock struct {
	sets map[string][]string
}

func (s *SetsManagerMock) OrderSet(photosetId string, photoIds []string) error {
	fmt.Println("Not implemented")
	return nil
}

func (s *SetsManagerMock) ListUserSets(userid string) (total int, photosetitems []photosets.Photoset, error error) {
	sets := []photosets.Photoset{}
	for setName, _ := range s.sets {
		sets = append(sets, photosets.Photoset{
			Id:    setName,
			Title: setName,
		})
	}
	return len(sets), sets, nil
}
func (s *SetsManagerMock) GetSetByName(setName, userid string) (string, error) {
	if _, ok := s.sets[setName]; ok {
		return setName, nil
	}
	return "", nil
}
func (s *SetsManagerMock) GetPhotosInSet(setId, userId string) (photosetitems []photosets.Photo, error error) {
	photos := []photosets.Photo{}
	if set, ok := s.sets[setId]; ok {
		for _, photoId := range set {
			extension := filepath.Ext(photoId)
			fileName := strings.TrimSuffix(photoId, extension)
			photos = append(photos, photosets.Photo{
				Id:    photoId,
				Title: fileName,
			})
		}
	}
	return photos, nil
}
func (s *SetsManagerMock) AddToPhotoSet(photoId, photoSetId string) error {
	if _, ok := s.sets[photoSetId]; !ok {
		s.sets[photoSetId] = []string{}
	}
	s.sets[photoSetId] = append(s.sets[photoSetId], photoId)
	fmt.Printf("Added photo %s to set %s\n", photoId, photoSetId)
	return nil
}
func (s *SetsManagerMock) CreateSet(title, description, primaryPhotoId string) (string, error) {
	if _, ok := s.sets[title]; !ok {
		s.sets[title] = []string{}
	}
	s.sets[title] = append(s.sets[title], primaryPhotoId)
	return title, nil
}

type DirectoryAccessorMock struct {
	GetSrcPhotosCallCount bool
	Directories           []fs.DirEntry
}

func (d *DirectoryAccessorMock) GetSrcPhotos() ([]fs.DirEntry, error) {
	// Mock implementation
	d.GetSrcPhotosCallCount = true
	return d.Directories, nil
}

type FlickrClientUploadHandlerFactoryMock struct {
	UploadHandlerMock *UploadHandlerMock
}

func (f FlickrClientUploadHandlerFactoryMock) GetUploadHandler() (flickrclient.UploadHandler, error) {
	return f.UploadHandlerMock, nil
}

type UploadHandlerMock struct {
	UploadedPhotos []string
	mu             sync.Mutex
	sm             *SetsManagerMock
}

func (u *UploadHandlerMock) UploadPhoto(isPublic bool, isFamily bool, isFriend bool, path string) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.UploadedPhotos = append(u.UploadedPhotos, path)
	return "", nil
}
func (u *UploadHandlerMock) DeletePhoto(photoId string) error {
	// Remove the photo from all sets
	for setName, photos := range u.sm.sets {
		for i, photo := range photos {
			if photo == photoId {
				u.sm.sets[setName] = append(photos[:i], photos[i+1:]...)
				break
			}
		}
	}
	return nil
}

type DireentryMock struct {
	DirName string
}

func (d *DireentryMock) Name() string {
	return d.DirName
}
func (d *DireentryMock) IsDir() bool {
	return false
}
func (d *DireentryMock) Type() fs.FileMode {
	return 0
}
func (d *DireentryMock) Info() (fs.FileInfo, error) {
	return nil, nil
}

func TestRunCmdWithoutPhotoToUpload(t *testing.T) {
	da := &DirectoryAccessorMock{}
	factory := FlickrClientUploadHandlerFactoryMock{}
	sm := &SetsManagerMock{}
	params := syncParams{
		syncMode:  false,
		overwrite: false,
		delete:    false,
		setName:   "",
	}
	runSyncCmd(params, sm, da, factory)
	if !da.GetSrcPhotosCallCount {
		t.Errorf("Expected GetSrcPhotos to be called")
	}
}
func TestRunCmdWithPhotoToUpload(t *testing.T) {
	da := &DirectoryAccessorMock{
		Directories: []fs.DirEntry{
			&DireentryMock{DirName: "photo1.jpg"},
			&DireentryMock{DirName: "photo2.jpg"},
		},
	}
	factory := FlickrClientUploadHandlerFactoryMock{
		UploadHandlerMock: &UploadHandlerMock{},
	}
	sm := &SetsManagerMock{
		sets: map[string][]string{},
	}
	params := syncParams{
		syncMode:  false,
		overwrite: false,
		delete:    false,
		setName:   "",
	}
	runSyncCmd(params, sm, da, factory)
	// Verify that the 2 photos were uploaded
	if len(factory.UploadHandlerMock.UploadedPhotos) != 2 {
		t.Errorf("Expected 2 photos to be uploaded, got %d", len(factory.UploadHandlerMock.UploadedPhotos))
	}
	photosUploaded := factory.UploadHandlerMock.UploadedPhotos
	photosToUpload := []string{"photo1.jpg", "photo2.jpg"}

	// Check if photosUploaded contains the expected photos
	for _, photo := range photosToUpload {
		if !slices.Contains(photosUploaded, photo) {
			t.Errorf("Expected photo %s to be uploaded, but it wasn't", photo)
		}
	}
}

func TestRunCmdWithSetName(t *testing.T) {
	da := &DirectoryAccessorMock{
		Directories: []fs.DirEntry{
			&DireentryMock{DirName: "photo1.jpg"},
			&DireentryMock{DirName: "photo2.jpg"},
		},
	}
	factory := FlickrClientUploadHandlerFactoryMock{
		UploadHandlerMock: &UploadHandlerMock{},
	}
	sm := &SetsManagerMock{
		sets: map[string][]string{},
	}
	params := syncParams{
		syncMode:  false,
		overwrite: false,
		delete:    false,
		setName:   "album1",
	}
	runSyncCmd(params, sm, da, factory)
	if len(factory.UploadHandlerMock.UploadedPhotos) != 2 {
		t.Errorf("Expected 2 photos to be uploaded, got %d", len(factory.UploadHandlerMock.UploadedPhotos))
	}
	// Verify if the set is created
	if _, ok := sm.sets["album1"]; !ok {
		t.Errorf("Expected set 'album1' to be created, but it wasn't")
	}

	// Verify if the set contains the uploaded photos
	if len(sm.sets["album1"]) != 2 {
		t.Errorf("Expected set 'album1' to contain 2 photos, got %d", len(sm.sets["album1"]))
	}
}

func TestCmdWithSetnameExistingSetAndNoSync(t *testing.T) {
	da := &DirectoryAccessorMock{
		Directories: []fs.DirEntry{
			&DireentryMock{DirName: "photo1.jpg"},
			&DireentryMock{DirName: "photo2.jpg"},
		},
	}
	factory := FlickrClientUploadHandlerFactoryMock{
		UploadHandlerMock: &UploadHandlerMock{},
	}
	sm := &SetsManagerMock{
		sets: map[string][]string{
			"album1": {"photo1.jpg"},
		},
	}
	params := syncParams{
		syncMode:  false,
		overwrite: false,
		delete:    false,
		setName:   "album1",
	}
	runSyncCmd(params, sm, da, factory)
	if len(factory.UploadHandlerMock.UploadedPhotos) != 2 {
		t.Errorf("Expected 2 photos to be uploaded, got %d", len(factory.UploadHandlerMock.UploadedPhotos))
	}
	// verify if the set contains 3 pictures
	if len(sm.sets["album1"]) != 3 {
		t.Errorf("Expected set 'album1' to contain 3 photos, got %d", len(sm.sets["album1"]))
	}
}

func TestRunCmdWithSetNameAndSync(t *testing.T) {
	da := &DirectoryAccessorMock{
		Directories: []fs.DirEntry{
			&DireentryMock{DirName: "photo1.jpg"},
			&DireentryMock{DirName: "photo2.jpg"},
		},
	}
	factory := FlickrClientUploadHandlerFactoryMock{
		UploadHandlerMock: &UploadHandlerMock{},
	}
	sm := &SetsManagerMock{
		sets: map[string][]string{
			"album1": {"photo1.jpg"},
		},
	}
	params := syncParams{
		syncMode:  true,
		overwrite: false,
		delete:    false,
		setName:   "album1",
	}
	runSyncCmd(params, sm, da, factory)
	if len(factory.UploadHandlerMock.UploadedPhotos) != 1 {
		t.Errorf("Expected 1 photo to be uploaded, got %d", len(factory.UploadHandlerMock.UploadedPhotos))
	}
	// verify if the set contains 2
	if len(sm.sets["album1"]) != 2 {
		t.Errorf("Expected set 'album1' to contain 2 photos, got %d", len(sm.sets["album1"]))
	}
}

func TestRunCmdWithSetNameAndSyncOverwrite(t *testing.T) {
	da := &DirectoryAccessorMock{
		Directories: []fs.DirEntry{
			&DireentryMock{DirName: "photo1.jpg"},
			&DireentryMock{DirName: "photo2.jpg"},
		},
	}
	sm := &SetsManagerMock{
		sets: map[string][]string{
			"album1": {"photo1.jpg"},
		},
	}
	factory := FlickrClientUploadHandlerFactoryMock{
		UploadHandlerMock: &UploadHandlerMock{
			sm: sm,
		},
	}
	params := syncParams{
		syncMode:  true,
		overwrite: true,
		delete:    false,
		setName:   "album1",
	}
	runSyncCmd(params, sm, da, factory)
	if len(factory.UploadHandlerMock.UploadedPhotos) != 2 {
		t.Errorf("Expected 1 photo to be uploaded, got %d", len(factory.UploadHandlerMock.UploadedPhotos))
	}
	// verify if the set contains 2
	if len(sm.sets["album1"]) != 2 {
		t.Errorf("Expected set 'album1' to contain 2 photos, got %d", len(sm.sets["album1"]))
	}
}

func TestRunCMDDelete(t *testing.T) {
	da := &DirectoryAccessorMock{
		Directories: []fs.DirEntry{
			&DireentryMock{DirName: "photo3.jpg"},
		},
	}
	sm := &SetsManagerMock{
		sets: map[string][]string{
			"album1": {"photo1.jpg", "photo2.jpg"},
		},
	}
	factory := FlickrClientUploadHandlerFactoryMock{
		UploadHandlerMock: &UploadHandlerMock{
			sm: sm,
		},
	}
	params := syncParams{
		syncMode:  true,
		overwrite: false,
		delete:    true,
		setName:   "album1",
	}
	runSyncCmd(params, sm, da, factory)
	if len(factory.UploadHandlerMock.UploadedPhotos) != 1 {
		t.Errorf("Expected 1 photos to be uploaded, got %d", len(factory.UploadHandlerMock.UploadedPhotos))
	}
	if len(sm.sets["album1"]) != 1 {
		t.Errorf("Expected set 'album1' to contain 1 photos, got %d", len(sm.sets["album1"]))
	}
}
