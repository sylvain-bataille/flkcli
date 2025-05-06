package organize

import (
	"flkcli/cmd"
	"flkcli/flickrclient"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/masci/flickr.v3/photos"
)

func runSortSetCmd(sm flickrclient.SetsManager, pr flickrclient.PhotoInfoRetriever, setName string, sortOption string, descending bool) error {
	set, err := sm.GetSetByName(setName, "")
	if err != nil {
		return fmt.Errorf("cannot resolve photoset id of %s: %e", setName, err)
	}
	photoset, err := sm.GetPhotosInSet(set, "")
	if err != nil {
		return fmt.Errorf("error getting photos in set %s: %e", setName, err)
	}

	photos := make([]photos.PhotoInfo, 0)
	for _, photo := range photoset {
		// Sort the photos based on some criteria
		// For example, sort by title
		fmt.Printf("Photo ID: %s, Title: %s\n", photo.Id, photo.Title)
		//TODO replace by loading bar, optmize photo info loading
		photoInfo, err := pr.GetPhotoInfo(photo.Id)
		if err != nil {
			return fmt.Errorf("error getting photo info of %s: %e", photo.Title, err)
		}
		photos = append(photos, photoInfo)
	}
	photos = sortPhotos(photos, sortOption, descending)
	photosOrder := []string{}
	for _, photo := range photos {
		photosOrder = append(photosOrder, photo.Id)
	}
	sm.OrderSet(set, photosOrder)
	return nil
}

func sortPhotos(photosList []photos.PhotoInfo, o string, descending bool) []photos.PhotoInfo {
	fmt.Println("Sorting photos...")
	switch o {
	case "date":
		fmt.Println("Sorting by date...")
		slices.SortFunc(photosList,
			func(a, b photos.PhotoInfo) int {
				return comparePhotosByDate(a, b)
			})
	default:
		fmt.Println("Sorting by title...")
		slices.SortFunc(photosList,
			func(a, b photos.PhotoInfo) int {
				return strings.Compare(a.Title, b.Title)
			})
	}

	if descending {
		fmt.Println("Reversing order...")
		slices.Reverse(photosList)
	}
	return photosList
}

func comparePhotosByDate(photo1, photo2 photos.PhotoInfo) int {
	start, _ := time.Parse(time.DateTime, photo1.Dates.Taken)
	end, _ := time.Parse(time.DateTime, photo2.Dates.Taken)
	diff := start.Sub(end)
	return int(diff.Milliseconds())
}

var SortSetCmd = &cobra.Command{
	Use:   "sort-set",
	Short: "Sort the photos in a specified set",
	Long:  `Sort the photos in a specified set`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Sorting photos in the specified set...")
		if len(args) == 0 {
			cmd.Println("Please provide a set name.")
			return
		}
		setName := args[0]
		sortOrder, err := cmd.Flags().GetString("sort")
		if err != nil {
			cmd.Println("error getting sort order:", err)
			return
		}
		descending, err := cmd.Flags().GetBool("descending")
		if err != nil {
			cmd.Println("error getting descending flag:", err)
			return
		}
		setsManager, err := flickrclient.GetSetsManager()
		if err != nil {
			cmd.Println("error getting sets manager:", err)
			return
		}
		photoRetreiver, err := flickrclient.GetPhotoClient()
		if err != nil {
			fmt.Println("error getting photo client:", err)
			return
		}
		err = runSortSetCmd(setsManager, photoRetreiver, setName, sortOrder, descending)
		if err != nil {
			fmt.Println("error sorting photos:", err)
			return
		}
		cmd.Println("Photos sorted successfully.")
	},
}

var OrganizeCmd = &cobra.Command{
	Use:   "organize",
	Short: "Organize flickr resources",
	Long: `Organize flickr resources
	sort-set: Sort the photos in a specified set`,
}

func init() {
	SortSetCmd.Flags().StringP("sort", "s", "title", "Sort order: date or title")
	SortSetCmd.Flags().BoolP("descending", "d", false, "Sort in descending order")
	OrganizeCmd.AddCommand(SortSetCmd)
	cmd.RootCmd.AddCommand(OrganizeCmd)
}
