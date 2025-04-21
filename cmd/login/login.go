package login

import (
	"flkcli/cmd"
	"flkcli/config"
	"flkcli/flickrclient"
	"fmt"

	"github.com/spf13/cobra"
)

func login(c config.ConfigManager, lh flickrclient.LoginHandler) {
	token, secret, err := lh.Login()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	//Print that the login was successful
	fmt.Println("Login successful")
	c.SetTokenConfig(token, secret)
}

// LoginCmd represents the login command
// The login command is used to authenticate the user with the Flickr API
// The login command uses the flickr package to authenticate the user
// The login command uses the config package to save the OAuth token and secret
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Interactive login",
	Long:  `Interactive login to authenticate the user with the Flickr API.`,
	Run: func(command *cobra.Command, args []string) {
		c, err := config.GetUserProfileConfig()
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		lh, err := flickrclient.GetLoginHandler()
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		login(c, lh)
	},
}

func init() {
	cmd.RootCmd.AddCommand(LoginCmd)
}
