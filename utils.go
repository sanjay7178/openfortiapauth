package main

import (
	"fmt"
	"net"
)

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
