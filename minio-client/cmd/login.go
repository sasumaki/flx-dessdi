package cmd

import (
	"aiga/api"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/theherk/viper"
	"golang.org/x/crypto/ssh/terminal"
)

// Login allows users to log in using an API token.
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Enter url and authentication for your bucket.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		valid, err := AuthValid()
		if valid == true {
			fmt.Println("Saved credentials work")
			return
		}

		reader := bufio.NewReader(os.Stdin)

		ogEndpoint := viper.GetString("endpoint")
		if ogEndpoint != "" {
			fmt.Print("Enter bucket endpoint(" + ogEndpoint + "): ")

		}
		fmt.Print("Enter bucket endpoint: ")
		endpoint, _ := reader.ReadString('\n')
		if endpoint == "" {
			endpoint = ogEndpoint
		}

		viper.Set("endpoint", endpoint)
		fmt.Print("Enter bucket access key ID: ")
		accessKeyID, _ := reader.ReadString('\n')

		fmt.Print("Enter bucket secret key: ")
		byteSecretKey, err := terminal.ReadPassword(int(syscall.Stdin))

		secretKey := string(byteSecretKey)

		endpoint, accessKeyID, secretKey = strings.TrimSpace(endpoint), strings.TrimSpace(accessKeyID), strings.TrimSpace(secretKey)
		fmt.Println("")
		viper.Set("endpoint", endpoint)
		viper.Set("accesskeyid", accessKeyID)
		viper.Set("secretKey", secretKey)

		valid, err = AuthValid()
		if valid == true {
			fmt.Println("Logged in.")
		} else {
			er(err)
		}
		err = viper.WriteConfig()
		if err != nil {
			er("Could not write config: " + err.Error())
		}

		return
	},
}

// AuthValid checks if you can access your bucket.
func AuthValid() (bool, error) {
	endpoint := viper.GetString("endpoint")
	accesskeyid := viper.GetString("accesskeyid")
	secretKey := viper.GetString("secretKey")

	if viper.GetString("endpoint") != "" && viper.GetString("accesskeyid") != "" && viper.GetString("secretkey") != "" {
		ok, err := api.TestConnection(endpoint, accesskeyid, secretKey)
		if err != nil {
			return false, err
		}
		return ok, nil
	}
	return false, errors.New("Credentials empty")
}
