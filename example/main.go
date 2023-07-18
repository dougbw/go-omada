package main

import (
	"fmt"
	"log"
	"os"

	omada "github.com/dougbw/go-omada"
)

func main() {

	// variables
	controllerUrl := "https://10.0.0.10"
	user, present := os.LookupEnv("OMADA_USERNAME")
	if !present {
		log.Fatal("⛔ required environment variable not set: OMADA_USERNAME")
		os.Exit(1)
	}
	pass, present := os.LookupEnv("OMADA_PASSWORD")
	if !present {
		log.Fatal("⛔ required environment variable not set: OMADA_PASSWORD")
		os.Exit(1)
	}

	siteName := "Home"

	// setup
	omada := omada.New(controllerUrl)
	err := omada.GetControllerInfo()
	if err != nil {
		log.Fatal(err)
	}

	// login
	err = omada.Login(user, pass)
	if err != nil {
		log.Fatal(err)
	}
	// set site
	err = omada.SetSite(siteName)
	if err != nil {
		log.Fatal(err)
	}

	// get clients
	clients, err := omada.GetClients()
	if err != nil {
		log.Printf("error getting clients: %v", err)
	}
	for _, client := range clients {
		fmt.Printf("client ip: %s, dnsName: %s, name: %s\n", client.Ip, client.DnsName, client.Name)
	}

	// get devices
	devices, err := omada.GetDevices()
	if err != nil {
		log.Printf("error getting devices: %v", err)
	}
	for _, device := range devices {
		fmt.Printf("device name: %s, dnsName: %s,  ip: %s\n", device.Name, device.DnsName, device.IP)
	}

	// get networks
	networks, err := omada.GetNetworks()
	if err != nil {
		log.Printf("error getting devices: %v", err)
	}
	for _, network := range networks {
		fmt.Printf("network name: %s, subnet: %s, domain: %s\n", network.Name, network.Subnet, network.Domain)
	}

	dhcpr, err := omada.GetDhcpReservations()
	if err != nil {
		log.Printf("error getting dhcp reservations: %v", err)
	}

	for _, r := range dhcpr {
		fmt.Printf("dhcp ip: %s, description: %s\n", r.IP, r.Description)
	}

}
