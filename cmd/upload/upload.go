package upload

import (
	"flkcli/cmd"
	"flkcli/flickrclient"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"gopkg.in/masci/flickr.v3/photosets"
)

type LoadingBar struct {
	// The total number of items to load
	Total int
	// The number of items loaded
	Loaded int
	// The last item loaded
	LastLoaded string
}

type PhotoUploadInfo struct {
	FilePath        string
	Overwrite       bool
	PreviousPhotoId string
}

type syncParams struct {
	syncMode  bool
	overwrite bool
	delete    bool
	setName   string
}

type PhotosAccessor interface {
	GetSrcPhotos() ([]fs.DirEntry, error)
}

type UploadHandlerFactory interface {
	GetUploadHandler() (flickrclient.UploadHandler, error)
}

type FlickrClientUploadHandlerFactory struct{}

func (f FlickrClientUploadHandlerFactory) GetUploadHandler() (flickrclient.UploadHandler, error) {
	uh, err := flickrclient.GetUploadHandler()
	if err != nil {
		return nil, fmt.Errorf("failed to get upload handler: %w", err)
	}
	return uh, nil
}

type LocalDirectoryAccessor struct{}

func (lda LocalDirectoryAccessor) GetSrcPhotos() ([]fs.DirEntry, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	return files, nil
}

func runSyncCmd(params syncParams, sm flickrclient.SetsManager, pa PhotosAccessor, uhf UploadHandlerFactory) {
	setId := ""
	photosInSet := []photosets.Photo{}
	files, err := pa.GetSrcPhotos()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	if params.setName != "" {
		setId, err = sm.GetSetByName(params.setName, "")
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		if setId != "" {
			photosInSet, _ = sm.GetPhotosInSet(setId, "")
		}
	}

	var photosToUpload []PhotoUploadInfo

	// Identify files to upload
	for _, file := range files {
		if !isPhoto(file) {
			fmt.Printf("\nIgnoring: %s\n", file.Name())
			continue
		}

		if params.syncMode {
			if result, id := isPhotoInSet(file, photosInSet); result {
				if params.overwrite {
					photosToUpload = append(photosToUpload, PhotoUploadInfo{FilePath: file.Name(), Overwrite: true, PreviousPhotoId: id})
				}
				continue
			}
		}
		photosToUpload = append(photosToUpload, PhotoUploadInfo{FilePath: file.Name(), Overwrite: false, PreviousPhotoId: ""})
	}

	if len(photosToUpload) == 0 {
		fmt.Println("No photos to upload")
		return
	}
	bar := LoadingBar{Total: len(photosToUpload), Loaded: 0}
	bar.printLoadingBar()

	var mutex sync.Mutex
	var wg sync.WaitGroup
	maxGoroutines := 5
	guard := make(chan struct{}, maxGoroutines)
	wg.Add(len(photosToUpload))
	for _, pui := range photosToUpload {
		guard <- struct{}{}
		go func(pui PhotoUploadInfo) {
			defer wg.Done()
			defer func() { <-guard }()

			uh, err := uhf.GetUploadHandler()
			if err != nil {
				fmt.Printf("Error: %s", err)
				return
			}
			id, err := uploadPhoto(uh, pui)
			if err != nil {
				fmt.Printf("Error: %s", err)
				return
			}

			mutex.Lock()
			defer mutex.Unlock()
			// TODO handle multiple sets in non sync mode
			// TODO error if syncmode and multiple sets
			// TODO param to create set if not exist
			if setId == "" {
				setId, err = sm.CreateSet(params.setName, "", id)
				if err != nil {
					fmt.Printf("Error: %s", err)
				}
			} else {
				err = sm.AddToPhotoSet(id, setId)
				if err != nil {
					fmt.Printf("Error: %s", err)
					return
				}
			}
			bar.Loaded++
			bar.LastLoaded = pui.FilePath
			fmt.Printf("\nUploaded: %s\n", pui.FilePath)
			bar.printLoadingBar()
		}(pui)
	}
	wg.Wait()
	fmt.Println("\nAll photos uploaded")

	var photosToDelete []string

	if params.delete {
		for _, photo := range photosInSet {
			found := false
			for _, file := range files {
				extension := filepath.Ext(file.Name())
				fileName := strings.TrimSuffix(file.Name(), extension)
				if strings.EqualFold(strings.ToLower(photo.Title), strings.ToLower(fileName)) {
					found = true
					break
				}
			}
			if !found {
				photosToDelete = append(photosToDelete, photo.Id)
			}
		}
		if len(photosToDelete) > 0 {
			fmt.Printf("\nDeleting %d photos...\n", len(photosToDelete))
			for _, photoId := range photosToDelete {
				uh, err := uhf.GetUploadHandler()
				if err != nil {
					fmt.Printf("Error: %s", err)
					return
				}
				err = uh.DeletePhoto(photoId)
				if err != nil {
					fmt.Printf("Error: %s", err)
				}
			}

		}
	}
}

func validateArgs(args []string, syncMode bool, overwrite bool, delete bool) (string, error) {
	albumName := ""
	if (len(args) == 0 || args[0] == "") && syncMode {
		return "", fmt.Errorf("album name is required in sync mode")
	}
	if (overwrite || delete) && !syncMode {
		return "", fmt.Errorf("overwrite and delete flags can only be used in sync mode")
	}
	if len(args) > 0 {
		albumName = args[0]
	}
	return albumName, nil
}

var UploadCmd = &cobra.Command{
	Use:   "upload [albumName]",
	Short: "Upload photos to flickr",
	Long:  `Upload all the photos in the specified directory to flickr`,
	Run: func(command *cobra.Command, args []string) {
		syncMode, _ := command.Flags().GetBool("sync")
		overwrite, _ := command.Flags().GetBool("overwrite")
		delete, _ := command.Flags().GetBool("delete")
		setName, err := validateArgs(args, syncMode, overwrite, delete)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		sm, err := flickrclient.GetSetsManager()
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		runSyncCmd(syncParams{
			syncMode:  syncMode,
			overwrite: overwrite,
			delete:    delete,
			setName:   setName,
		}, sm, LocalDirectoryAccessor{}, FlickrClientUploadHandlerFactory{})
	},
}

func isPhotoInSet(file fs.DirEntry, photosInSet []photosets.Photo) (bool, string) {
	for _, photo := range photosInSet {
		extension := filepath.Ext(file.Name())
		fileName := strings.TrimSuffix(file.Name(), extension)
		if strings.EqualFold(strings.ToLower(photo.Title), strings.ToLower(fileName)) {
			// Print that the file already exists in the set
			return true, photo.Id
		}
	}
	return false, ""
}

func init() {
	// Initialize empty string array
	UploadCmd.PersistentFlags().BoolP("sync", "s", false, "Sync mode")
	UploadCmd.PersistentFlags().BoolP("overwrite", "o", false, "Overwrite already existing photos")
	UploadCmd.PersistentFlags().BoolP("delete", "d", false, "Delete distant photos from flickr that are not in the local directory")

	cmd.RootCmd.AddCommand(UploadCmd)
}

func uploadPhoto(uh flickrclient.UploadHandler, pui PhotoUploadInfo) (string, error) {
	id, err := uh.UploadPhoto(false, false, false, pui.FilePath)
	if err != nil {
		return "", fmt.Errorf("failed to upload photo: %w", err)
	}
	if pui.Overwrite {
		uh.DeletePhoto(pui.PreviousPhotoId)
	}
	return id, nil
}

func isPhoto(file os.DirEntry) bool {
	// Check if the file has a valid photo extension
	extensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, ext := range extensions {
		if strings.HasSuffix(strings.ToLower(file.Name()), ext) {
			return true
		}
	}
	return false
}

func (bar LoadingBar) printLoadingBar() {
	progress := float64(bar.Loaded) / float64(bar.Total)
	barLength := 30
	filledLength := int(progress * float64(barLength))

	barstr := strings.Repeat("█", filledLength) + strings.Repeat(" ", barLength-filledLength)

	// Clear console content
	fmt.Print("\033[H\033[2J")
	// Print total number of files being uploaded
	fmt.Printf("Uploading %d files...\n", bar.Total)
	fmt.Printf("\r[%s] %d/%d", barstr, bar.Loaded, bar.Total)
	// Print last successfully uploaded file if not empty
	if bar.LastLoaded != "" {
		fmt.Printf("\nLast uploaded: %s", bar.LastLoaded)
	}
}
