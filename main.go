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
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"

	// "github.com/PuerkitoBio/goquery"
	// "github.com/spf13/cobra"
	// "github.com/spf13/viper"

	"fmt"
	// "log"
	// "net/http"
	// "net/url"
	"strings"

	// "time"

	fyne1 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	// "github.com/PuerkitoBio/goquery"
)

func CreateStorage() {
	// Get the root URI
	rootURI := fyne1.CurrentApp().Storage().RootURI()

	// Define the file URI
	fileURI := storage.NewFileURI(rootURI.Path() + "/usrData.json")

	// Check if the file exists
	exists, err := storage.Exists(fileURI)
	if err != nil {
		fmt.Println("Error checking file existence:", err)
		return
	}

	// If the file does not exist, create a new file
	if !exists {
		writer, err := storage.Writer(fileURI)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer writer.Close()

		// Write initial content to the file
		_, err = writer.Write([]byte("Initial content"))
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		fmt.Println("File created successfully")
	} else {
		fmt.Println("File already exists")
	}
}

type UserData struct {
	Address       string  `json:"address"`
	Username      string  `json:"username"`
	Password      string  `json:"password"`
	UserAgent     string  `json:"user_agent"`
	KeepAlive     float64 `json:"keep_alive"`
	RetryInterval float64 `json:"retry_interval"`
	AutoUpdates   bool    `json:"auto_updates"`
	Port          string  `json:"port"`
}

const dataFile = "usrData.json"

func saveData(data UserData) error {
	// Convert struct to JSON
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Get the root URI
	rootURI := fyne1.CurrentApp().Storage().RootURI()

	// Define the file URI
	fileURI := storage.NewFileURI(rootURI.Path() + "/" + dataFile)

	// Write data to the file
	writer, err := storage.Writer(fileURI)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = writer.Write(content)
	return err
}

func loadData() (UserData, error) {
	var data UserData

	// Get the root URI
	rootURI := fyne1.CurrentApp().Storage().RootURI()

	// Define the file URI
	fileURI := storage.NewFileURI(rootURI.Path() + "/" + dataFile)

	// Check if the file exists
	exists, err := storage.Exists(fileURI)
	if err != nil {
		return data, err
	}

	// If the file does not exist, return empty data
	if !exists {
		return data, nil
	}

	// Read data from the file
	reader, err := storage.Reader(fileURI)
	if err != nil {
		return data, err
	}
	defer reader.Close()

	// Decode JSON data
	err = json.NewDecoder(reader).Decode(&data)
	return data, err
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("OpenFortiAP Authenticator")

	// Load existing data
	userData, err := loadData()
	if err != nil {
		fmt.Println("Error loading data:", err)
	}

	// Address
	addressEntry := widget.NewEntry()
	addressEntry.SetPlaceHolder("172.18.10.10")
	addressEntry.SetText(userData.Address)
	addressEntry.OnChanged = func(content string) {
		userData.Address = content
	}

	//Port
	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("1000")
	portEntry.SetText(userData.Port)
	portEntry.OnChanged = func(content string) {
		userData.Port = content
	}

	// Credentials
	usernameEntry := widget.NewEntry()
	usernameEntry.SetText(userData.Username)
	usernameEntry.OnChanged = func(content string) {
		userData.Username = content
	}

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetText(userData.Password)
	passwordEntry.OnChanged = func(content string) {
		userData.Password = content
	}

	// User Agent
	userAgentEntry := widget.NewEntry()
	userAgentEntry.SetText("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:129.0) Gecko/20100101 Firefox/129.0")
	userData.UserAgent = userAgentEntry.Text
	userAgentEntry.OnChanged = func(content string) {
		userData.UserAgent = content
	}

	// Keep alive interval
	keepAliveSlider := widget.NewSlider(0, 180)
	keepAliveSlider.SetValue(userData.KeepAlive)
	keepAliveSlider.OnChanged = func(value float64) {
		userData.KeepAlive = value
	}

	// Retry interval
	retrySlider := widget.NewSlider(0, 60)
	retrySlider.SetValue(userData.RetryInterval)
	retrySlider.OnChanged = func(value float64) {
		userData.RetryInterval = value
	}

	// Automatic updates
	autoUpdatesCheck := widget.NewCheck("Automatic updates", func(checked bool) {
		userData.AutoUpdates = checked
	})
	autoUpdatesCheck.SetChecked(userData.AutoUpdates)

	// Save button
	saveButton := widget.NewButton("Save", func() {
		if userData.Address == "" || userData.Username == "" || userData.Password == "" || userData.Port == "" {
			dialog.ShowError(fmt.Errorf("please fill in all fields"), myWindow)
		}
		err := saveData(userData)
		if err != nil {
			fmt.Println("Error saving data:", err)
			dialog.ShowError(err, myWindow)
		}
		dialog.ShowInformation("Data saved", "Your data has been saved successfully", myWindow)
	})

	// Dynamic text widget
	dynamicText := widget.NewLabel("")

	forceReLoginButton := widget.NewButton("Force Re-Login", func() {
		//Progress bar
		progressBar := widget.NewProgressBar()
		if userData.Address == "" || userData.Username == "" || userData.Password == "" || userData.Port == "" {
			dialog.ShowError(fmt.Errorf("please fill in all fields and click save"), myWindow)
			progressBar.SetValue(0)
		}
		session, err := Magic(userData.Address, userData.Port)
		if err != nil {
			log.Error("Unable to get Magic ID :", err)
		}
		dynamicText.SetText(session)
		Login(userData.Address, userData.Username, userData.Password, session, userData.Port)
		progressBar.SetValue(1)
	})
	logoutButton := widget.NewButton("Logout", func() {
		if userData.Address == "" || userData.Port == "" {
			dialog.ShowError(fmt.Errorf("please fill in all fields and click save"), myWindow)
		}
		err := Logout(userData.Address, userData.Port, userData.UserAgent)
		if err != nil {
			log.Error("Unable to logout:", err)
		}
		dynamicText.SetText("")
	}) // New logout button

	// IP Addresses
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	var text map[string]string
	// var err error
	// var ips string
	var richText *widget.RichText

	for range time.Tick(5 * time.Second) {
		<-ticker.C
		text, err = detect_interfaces()
		if err != nil {
			// fmt.Println("Error detecting interfaces:", err)
			log.Debug("Error detecting interfaces:", err)
		} else {
			// fmt.Println("Detected interfaces:", text)
			log.Debug("Detected interfaces:", text)
			var builder strings.Builder
			for key, value := range text {
				builder.WriteString(key)
				builder.WriteString(" : ")
				builder.WriteString(value)
				builder.WriteString("\n")
			}

			ips := builder.String()
			// Create the RichText widget
			richText = widget.NewRichText(&widget.TextSegment{
				Text:  ips,
				Style: widget.RichTextStyle{ColorName: "Blue"},
			})

		}
	}

	// Layout
	content := container.NewVBox(
		widget.NewLabel("Assigned IP Address:"),
		richText,
		widget.NewLabel("Address"),
		addressEntry,
		widget.NewLabel("Username"),
		usernameEntry,
		widget.NewLabel("Password"),
		passwordEntry,
		widget.NewLabel("User Agent"),
		userAgentEntry,
		widget.NewLabel("Keep alive interval (seconds):"),
		keepAliveSlider,
		widget.NewLabel("Retry interval on fail (seconds):"),
		retrySlider,
		autoUpdatesCheck,
		// saveButton,
		container.NewHBox(saveButton, forceReLoginButton, logoutButton), // Buttons in horizontal layout
		widget.NewLabel("Magic ID: "+dynamicText.Text),
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne1.NewSize(300, 200))
	myWindow.ShowAndRun()
}

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
