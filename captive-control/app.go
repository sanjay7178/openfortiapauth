// package main

// import (
// 	"context"
// 	"fmt"
// )

// // App struct
// type App struct {
// 	ctx context.Context
// }

// // NewApp creates a new App application struct
// func NewApp() *App {
// 	return &App{}
// }

// // startup is called when the app starts. The context is saved
// // so we can call the runtime methods
// func (a *App) startup(ctx context.Context) {
// 	a.ctx = ctx
// }

// // Greet returns a greeting for the given name
// func (a *App) Greet(name string) string {
// 	return fmt.Sprintf("Hello %s, It's show time!", name)
// }
package main

import (
	// "crypto/md5"
	// "encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// App struct
type App struct{}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Login function to handle user authentication
func (a *App) Login(username, password string) bool {
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
	payload.Set("username", username)
	payload.Set("password", password)

	url2 := "https://172.18.10.10:1000"

	req, _ = http.NewRequest("POST", url2, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.99 Safari/537.36")
	req.PostForm = payload

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) == "Failed" {
		return false
	} else {
		fmt.Println("Successfully authenticated....!")
		return true
	}
}
