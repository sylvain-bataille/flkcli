package list

import (
	"flkcli/cmd"
	"flkcli/flickrclient"
	"fmt"

	"github.com/spf13/cobra"
)

// GetUserFromArgs returns the user id from the arguments and resolves it if necessary
func getUserIDFromArgs(args []string, res flickrclient.IdResolver, asUsername bool) (userid string, err error) {
	if len(args) == 0 {
		return "", nil
	}
	if len(args) > 1 {
		return "", fmt.Errorf("too many arguments")
	}

	userid = args[0]

	if !asUsername {
		return userid, nil
	}

	userid, err = res.ResolveId(userid)
	if err != nil {
		return "", fmt.Errorf("failed to resolve user id: %w", err)
	}
	return userid, nil
}

func listSets(uid string, sl flickrclient.SetsLister) error {
	// Get the list of sets
	total, photoSetItems, err := sl.ListUserSets(uid)
	if err != nil {
		return fmt.Errorf("failed to get users sets for id %s: %w", uid, err)
	}

	fmt.Printf("Total sets: %d\n", total)
	for _, item := range photoSetItems {
		fmt.Printf("- %s\n", item.Title)
	}
	return nil
}

func runListSetsCmd(idResolver flickrclient.IdResolver, setsLister flickrclient.SetsLister, args []string, asusername bool) {
	userid, err := getUserIDFromArgs(args, idResolver, asusername)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	err = listSets(userid, setsLister)
	if err != nil {
		fmt.Printf("cannot list user sets: %s", err)
		return
	}
}

// The list command is used to list the photo sets of a user
var listSetsCmd = &cobra.Command{
	Use:   "sets [userid]",
	Short: "List photo sets",
	Long: `List photo sets of a user
		If no user is specified, the currents user sets are listed`,
	Args: cobra.MaximumNArgs(1),
	Run: func(command *cobra.Command, args []string) {
		asUsernameFlag, err := command.Flags().GetBool("as-username")
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		idResolver, err := flickrclient.GetPeopleManager()
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		setsLister, err := flickrclient.GetSetsManager()
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		runListSetsCmd(idResolver, setsLister, args, asUsernameFlag)
	},
}

var SetCmd = &cobra.Command{
	Use:   "list",
	Short: "List flickr resources",
	Long: `List flickr resources
	You must specify the resource type to list
	Available resource types are: sets`,
}

func init() {
	listSetsCmd.Flags().Bool("as-username", false, "Treat the user id as a username")
	SetCmd.AddCommand(listSetsCmd)
	cmd.RootCmd.AddCommand(SetCmd)
}
