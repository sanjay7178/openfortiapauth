package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func login(uname, passw string) bool {
	url1 := "http://www.google.co.in"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url1, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.99 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	defer resp.Body.Close()

	magic := resp.Request.URL.RawQuery

	payload := url.Values{}
	payload.Set("4Tredir", "http://google.com/")
	payload.Set("magic", magic)
	payload.Set("username", uname)
	payload.Set("password", passw)

	url2 := "https://gateway.iitk.ac.in:1003"

	req, _ = http.NewRequest("POST", url2, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.99 Safari/537.36")
	req.PostForm = payload

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) == "Failed" {
		return false
	} else {
		fmt.Println("Successfully authenticated....!")
		return true
	}
}

func main() {
	var username, password string

	rootCmd := &cobra.Command{
		Use:   "login",
		Short: "Login CLI application",
		Long:  `A CLI tool for logging into a specified gateway using provided credentials.`,
		Run: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("username", cmd.Flags().Lookup("username"))
			viper.BindPFlag("password", cmd.Flags().Lookup("password"))

			username = viper.GetString("username")
			password = viper.GetString("password")

			if username == "" || password == "" {
				fmt.Println("Username and password are required.")
				os.Exit(1)
			}

			hash := md5.Sum([]byte(password))
			hashStr := hex.EncodeToString(hash[:])

			fmt.Printf("Checking for user: %s, password: %s\n", username, hashStr)

			resp, err := http.Head("http://www.google.co.in")
			if err == nil {
				fmt.Printf("Already connected, pinged %s status: %d\n", resp.Request.URL, resp.StatusCode)
			} else {
				if login(username, password) {
					fmt.Println("Login successful.....!\n")
				} else {
					fmt.Println("Login Error. Try again......!\n")
				}
			}
		},
	}

	rootCmd.Flags().StringP("username", "u", "", "Username for login")
	rootCmd.Flags().StringP("password", "p", "", "Password for login")

	rootCmd.SetHelpCommand(&cobra.Command{
		Use:   "help",
		Short: "Help command for the login CLI application",
		Long:  "Provides detailed descriptions of the available commands and flags.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Usage: login [flags]\n")
			fmt.Println("Flags:")
			fmt.Println("  -u, --username string   Username for login")
			fmt.Println("  -p, --password string   Password for login")
			fmt.Println("\nHelp command for the login CLI application provides detailed descriptions of the available commands and flags.")
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
