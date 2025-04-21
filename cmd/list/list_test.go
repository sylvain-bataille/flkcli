package list

import (
	"testing"

	"gopkg.in/masci/flickr.v3/photosets"
)

type fakeResolver struct{}

func (f *fakeResolver) ResolveId(username string) (id string, error error) {
	return "12345", nil
}

type fakeLister struct {
	userId string
}

func (f *fakeLister) ListUserSets(userid string) (total int, photosetitems []photosets.Photoset, error error) {
	photoSetItems := []photosets.Photoset{
		{
			Id: "1",
		},
		{
			Id: "2",
		},
		{
			Id: "3",
		},
	}
	f.userId = userid

	return 3, photoSetItems, nil
}

func (f *fakeLister) GetSetByName(setName, userid string) (string, error) {
	return "12345", nil
}

func (f *fakeLister) GetPhotosInSet(setId, userId string) (photosetitems []photosets.Photo, error error) {
	photoSetItems := []photosets.Photo{
		{
			Id: "1",
		},
		{
			Id: "2",
		},
		{
			Id: "3",
		},
	}
	return photoSetItems, nil
}

func TestGetUserIDFromEmptyArgs(t *testing.T) {
	args := []string{""}
	id, err := getUserIDFromArgs(args, &fakeResolver{}, false)
	if id != "" {
		t.Errorf("Expected empty id, got %s", id)
	}
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestGetUserIDFromArgs(t *testing.T) {
	args := []string{"testuser"}
	id, err := getUserIDFromArgs(args, &fakeResolver{}, false)
	if id != "testuser" {
		t.Errorf("Expected id testuser, got %s", id)
	}
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestGetUserIDFromArgsEmpty(t *testing.T) {
	args := []string{}
	id, err := getUserIDFromArgs(args, &fakeResolver{}, false)
	if id != "" {
		t.Errorf("Expected empty id, got %s", id)
	}
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestGetUserIDFromArgsWithUsername(t *testing.T) {
	args := []string{"testuser"}
	id, err := getUserIDFromArgs(args, &fakeResolver{}, true)
	if id != "12345" {
		t.Errorf("Expected id 12345, got %s", id)
	}
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestGetUserIDFromArgsWithError(t *testing.T) {
	args := []string{"testuser", "extra"}
	_, err := getUserIDFromArgs(args, &fakeResolver{}, false)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestListSets(t *testing.T) {
	f := &fakeLister{}
	listSets("user", f)
	if f.userId != "user" {
		t.Errorf("Expected user id user, got %s", f.userId)
	}
}

func TestRunListSetsCmd(t *testing.T) {
	f := &fakeLister{}
	idResolver := &fakeResolver{}
	args := []string{"testuser"}
	runListSetsCmd(idResolver, f, args, true)
	if f.userId != "12345" {
		t.Errorf("Expected user id 12345, got %s", f.userId)
	}
}

func TestRunListSetsCmdWithError(t *testing.T) {
	f := &fakeLister{}
	idResolver := &fakeResolver{}
	args := []string{"testuser", "extra"}
	runListSetsCmd(idResolver, f, args, true)
	if f.userId != "" {
		t.Errorf("Expected empty user id, got %s", f.userId)
	}
}
