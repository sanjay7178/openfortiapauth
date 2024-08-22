package main

import (
	// "crypto/md5"
	// "encoding/hex"
	// "fmt"

	// // "log"
	// "net/http"
	// "net/url"
	// "os"
	// "strings"

	// // "github.com/pingcap/log"
	// log "github.com/sirupsen/logrus"

	// "github.com/PuerkitoBio/goquery"
	// "github.com/spf13/cobra"
	// "github.com/spf13/viper"

	fyne1 "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

// func magic(fgtIP string) (string, error) {
// 	// The URL to send the GET request to
// 	// url := "http://172.18.10.10:1000/login?"
// 	url := "http://"+fgtIP+":1000/login?"


// 	// Send the GET request
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		log.Fatalf("Failed to make GET request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Check if the request was successful
// 	if resp.StatusCode != http.StatusOK {
// 		log.Fatalf("Received non-200 response: %d %s", resp.StatusCode, resp.Status)
// 	}

// 	// Parse the HTML response
// 	doc, err := goquery.NewDocumentFromReader(resp.Body)
// 	if err != nil {
// 		log.Fatalf("Failed to parse HTML: %v", err)
// 	}

// 	// Find the hidden input element with name="magic" and extract its value
// 	magicValue, exists := doc.Find(`input[name="magic"]`).Attr("value")
// 	if !exists {
// 		log.Fatalf("Magic value not found")
// 		return "", fmt.Errorf("Magic value not found")
// 	}
// 	return magicValue, nil

// }

// func login(fgtIP, uname, passw, sessionID string) bool {
// 	postData := url.Values{}
// 	magic ,err :=  magic(fgtIP)
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	sessionID =  magic
// 	log.Debug(magic)
// 	postData.Set("magic", sessionID)
// 	postData.Set("username", uname)
// 	postData.Set("password", passw)

// 	client := &http.Client{}
// 	req, err := http.NewRequest("POST", "http://"+fgtIP+":1000/fgtauth", strings.NewReader(postData.Encode()))
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return false
// 	}
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error sending request:", err)
// 		return false
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode == http.StatusOK {
// 		fmt.Println("Successfully authenticated....!")
// 		return true
// 	} else {
// 		fmt.Println("Failed to authenticate.")
// 		return false
// 	}
// }

// func main() {

// 	go fyne() 

// 	var username, password, fgtIP, sessionID string

// 	rootCmd := &cobra.Command{
// 		Use:   "login",
// 		Short: "Login CLI application",
// 		Long:  `A CLI tool for logging into a specified gateway using provided credentials.`,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			viper.BindPFlag("username", cmd.Flags().Lookup("username"))
// 			viper.BindPFlag("password", cmd.Flags().Lookup("password"))
// 			viper.BindPFlag("fgtIP", cmd.Flags().Lookup("fgtIP"))
// 			// viper.BindPFlag("sessionID", cmd.Flags().Lookup("sessionID"))

// 			username = viper.GetString("username")
// 			password = viper.GetString("password")
// 			fgtIP = viper.GetString("fgtIP")
// 			// sessionID = viper.GetString("sessionID")

// 			if username == "" || password == "" || fgtIP == "" {
// 				fmt.Println("Username, password, FortiGate IP  are required.")
// 				os.Exit(1)
// 			}

// 			hash := md5.Sum([]byte(password))
// 			hashStr := hex.EncodeToString(hash[:])

// 			fmt.Printf("Checking for user: %s, password: %s\n", username, hashStr)

// 			if login(fgtIP, username, password, sessionID) {
// 				fmt.Println("Login successful.....!\n")
// 			} else {
// 				fmt.Println("Login Error. Try again......!\n")
// 			}
// 		},
// 	}

// 	rootCmd.Flags().StringP("username", "u", "", "Username for login")
// 	rootCmd.Flags().StringP("password", "p", "", "Password for login")
// 	rootCmd.Flags().StringP("fgtIP", "i", "", "FortiGate IP address")
// 	rootCmd.Flags().StringP("sessionID", "s", "", "Session ID")

// 	rootCmd.SetHelpCommand(&cobra.Command{
// 		Use:   "help",
// 		Short: "Help command for the login CLI application",
// 		Long:  "Provides detailed descriptions of the available commands and flags.",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			fmt.Println("Usage: login [flags]\n")
// 			fmt.Println("Flags:")
// 			fmt.Println("  -u, --username string   Username for login")
// 			fmt.Println("  -p, --password string   Password for login")
// 			fmt.Println("  -i, --fgtIP string      FortiGate IP address")
// 			fmt.Println("  -s, --sessionID string  Session ID")
// 			fmt.Println("\nHelp command for the login CLI application provides detailed descriptions of the available commands and flags.")
// 		},
// 	})

// 	if err := rootCmd.Execute(); err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// }

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("OpenFortiAP Authenticator")

	// Address
	addressEntry := widget.NewEntry()
	addressEntry.SetPlaceHolder("https://hfw.vitap.ac.in:8090")

	// Credentials
	usernameEntry := widget.NewEntry()
	usernameEntry.SetText("21MIC7178")
	passwordEntry := widget.NewPasswordEntry()

	// User Agent
	userAgentEntry := widget.NewEntry()
	userAgentEntry.SetText("OpenXGAuthenticator GUI (Via libopenvpn v0.1.3; F")

	// Keep alive interval
	keepAliveSlider := widget.NewSlider(0, 180)
	keepAliveSlider.SetValue(90)

	// Retry interval
	retrySlider := widget.NewSlider(0, 60)
	retrySlider.SetValue(5)

	// Automatic updates
	autoUpdatesCheck := widget.NewCheck("Automatic updates", nil)
	autoUpdatesCheck.SetChecked(true)

	// Buttons
	saveButton := widget.NewButton("Save", nil)
	forceReLoginButton := widget.NewButton("Force Re-Login", nil)
	logoutButton := widget.NewButton("Logout", nil)  // New logout button

	// Status
	statusLabel := widget.NewLabel("Connection Status: Logged in (71s keepalive)")

	// Layout
	content := container.NewVBox(
		widget.NewLabel("Assigned IP Address:"),

		widget.NewLabel("Address"),
		addressEntry,
		widget.NewLabel("Credentials"),
		usernameEntry,
		passwordEntry,
		widget.NewLabel("User Agent"),
		userAgentEntry,
		widget.NewLabel("Keep alive interval (seconds):"),
		keepAliveSlider,
		widget.NewLabel("Retry interval on fail (seconds):"),
		retrySlider,
		autoUpdatesCheck,
		container.NewHBox(saveButton, forceReLoginButton, logoutButton),  // Buttons in horizontal layout
		widget.NewLabel("Status"),
		statusLabel,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne1.NewSize(300, 200))
	myWindow.ShowAndRun()
}