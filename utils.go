package main

import (
	"fmt"
	"net"
	"strings"

	// "net"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	// "time"
)
func Magic(fgtIP string, port string) (string, error) {
	// The URL to send the GET request to
	// url := "http://172.18.10.10:1000/login?"
	url := "http://" + fgtIP + ":" + port + "/logout?"

	// Send the GET request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Received non-200 response: %d %s", resp.StatusCode, resp.Status)
	}

	// Parse the HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	// Find the hidden input element with name="magic" and extract its value
	magicValue, exists := doc.Find(`input[name="magic"]`).Attr("value")
	if !exists {
		log.Fatalf("Magic value not found")
		return "", fmt.Errorf("Magic value not found")
	}
	return magicValue, nil

}

func Login(fgtIP, uname, passw, sessionID, port string) bool {
	postData := url.Values{}
	// magic, err := Magic(fgtIP)
	// if err != nil {
	// 	log.Error("Error:", err)
	// }
	// sessionID = magic
	log.Debug(sessionID)
	postData.Set("magic", sessionID)
	postData.Set("username", uname)
	postData.Set("password", passw)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://"+fgtIP+":"+port+"/logout?", strings.NewReader(postData.Encode()))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Successfully authenticated....!")
		return true
	} else {
		fmt.Println("Failed to authenticate.")
		return false
	}
}

func Logout(fgtIP string, port string ,userAgent string) error {
	// Define the URL
	url := "http://" + fgtIP + ":" + port + "/logout?"

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Set the headers
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Priority", "u=0, i")

	// Create an HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("logout failed: %s", resp.Status)
	}

	fmt.Println("Logout successful")
	return nil
}


/*
Detect the addresses which can  be used for dispatching in non-tunnelling mode.
Alternate to ipconfig/ifconfig
*/
func detect_interfaces() (map[string]string, error) {
	fmt.Println("--- Listing the available adresses for dispatching")
	ifaces, _ := net.Interfaces()

	if len(ifaces) == 0 {
		fmt.Println("No interfaces found")
		return nil, nil
	}

	dict := map[string]string{}

	for _, iface := range ifaces {
		if (iface.Flags&net.FlagUp == net.FlagUp) && (iface.Flags&net.FlagLoopback != net.FlagLoopback) {
			addrs, _ := iface.Addrs()
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						fmt.Printf("[+] %s, IPv4:%s\n", iface.Name, ipnet.IP.String())
						dict[iface.Name] = ipnet.IP.String()
					}
				}
			}
		}
	}
	return dict, nil

}

