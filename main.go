package main

import (
    "encoding/json"
    "fmt"
    fyne1 "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/storage"
    "fyne.io/fyne/v2/widget"
    log "github.com/sirupsen/logrus"
    "strings"
    // "time"
)

func CreateStorage() {
    // Get the root URI
    rootURI := fyne1.CurrentApp().Storage().RootURI()

    // Define the file URI
    fileURI := storage.NewFileURI(rootURI.Path() + "/usrData.json")

    // Check if the file exists
    exists, err := storage.Exists(fileURI)
    if err != nil {
        //fmt.Println("Error checking file existence:", err)
        return
    }

    // If the file does not exist, create a new file
    if !exists {
        writer, err := storage.Writer(fileURI)
        if err != nil {
            //fmt.Println("Error creating file:", err)
            return
        }
        defer writer.Close()

        // Write initial content to the file
        _, err = writer.Write([]byte("Initial content"))
        if err != nil {
            //fmt.Println("Error writing to file:", err)
            return
        }
        //fmt.Println("File created successfully")
    } else {
        //fmt.Println("File already exists")
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
        //fmt.Println("Error loading data:", err)
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

    progressBar := widget.NewProgressBarInfinite()
    progressBar.Resize(fyne1.NewSize(1, 1)) // Set the height to 20

	

    // Save button
    saveButton := widget.NewButton("Save", func() {
        progressBar.Start()
        defer progressBar.Stop()

        if userData.Address == "" || userData.Username == "" || userData.Password == "" || userData.Port == "" {
            dialog.ShowError(fmt.Errorf("please fill in all fields"), myWindow)
            return
        }
        err := saveData(userData)
        if err != nil {
            //fmt.Println("Error saving data:", err)
            dialog.ShowError(err, myWindow)
            return
        }
        dialog.ShowInformation("Data saved", "Your data has been saved successfully", myWindow)
    })

    // Dynamic text widget
    dynamicText := widget.NewLabel("")

    forceReLoginButton := widget.NewButton("Force Re-Login", func() {
        progressBar.Start()
        defer progressBar.Stop()

        if userData.Address == "" || userData.Username == "" || userData.Password == "" || userData.Port == "" {
            dialog.ShowError(fmt.Errorf("please fill in all fields and click save"), myWindow)
            return
        }
        session, err := Magic(userData.Address, userData.Port)
        if err != nil {
            log.Error("Unable to get Magic ID :", err)
            dialog.ShowError(err, myWindow)
            return
        }
        dynamicText.SetText(session)
        err = Login(userData.Address, userData.Username, userData.Password, session, userData.Port)
        if err != nil {
            log.Error("Unable to login:", err)
            dialog.ShowError(err, myWindow)
            return
        }
        dialog.ShowInformation("Re-Login", "Re-Login successful", myWindow)
    })

    logoutButton := widget.NewButton("Logout", func() {
        progressBar.Start()
        defer progressBar.Stop()

        if userData.Address == "" || userData.Port == "" {
            dialog.ShowError(fmt.Errorf("please fill in all fields and click save"), myWindow)
            return
        }
        err := Logout(userData.Address, userData.Port, userData.UserAgent)
        if err != nil {
            log.Error("Unable to logout:", err)
            dialog.ShowError(err, myWindow)
            return
        }
        dynamicText.SetText("")
        dialog.ShowInformation("Logout", "Logout successful", myWindow)
    })

    var text map[string]string
    var err1 error
    text, err1 = detect_interfaces()

    var builder strings.Builder
    for key, value := range text {
        builder.WriteString(key)
        builder.WriteString(" : ")
        builder.WriteString(value)
        builder.WriteString("\n")
    }
    var ips string
    if err1 != nil {
        log.Error("Error getting IP address:", err1)
        ips = "Error getting Assigned IP addresses"
    } else {
        ips = builder.String()
    }
    richText := widget.NewRichText(&widget.TextSegment{
        Text:  "Assigned IP Address: \n" + ips,
        Style: widget.RichTextStyle{TextStyle: fyne1.TextStyle{Italic: true}},
    })

    // Layout
    content := container.NewVBox(
		widget.NewLabel("FortiGate Details"),
        addressEntry,
        portEntry,
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
        container.NewHBox(saveButton, forceReLoginButton, logoutButton), // Buttons in horizontal layout
        widget.NewLabel("Magic ID: "+dynamicText.Text),
        richText,
        progressBar,
    )

    myWindow.SetContent(content)
    myWindow.Resize(fyne1.NewSize(300, 200))
    myWindow.ShowAndRun()
}