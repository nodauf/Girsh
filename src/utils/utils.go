package utils

import (
	"fmt"
	"net"
)

// Backspace character
var Backspace = []byte{8}

// Newline character
var Newline = []byte{10}

// Nullbyte character
var Nullbyte = []byte{0}

// SliceStringContains check if a string is present in a slice of string
func SliceStringContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// SliceByteContains check if a byte is present in a slice of byte
func SliceByteContains(s []byte, e byte) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Min returns the minimal integer between two integer
func Min(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// Iface represents an interface on the system with its name and its IP
type Iface struct {
	Name string
	IP   string
}

// ListInterfaces returns a slice of struct iface (Name of interface with the IP) which are the interfaces on the system
func ListInterfaces() []Iface {
	interfaces := []Iface{}
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(fmt.Errorf("localAddresses: %+v", err.Error()))
		return nil
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(fmt.Errorf("localAddresses: %+v", err.Error()))
			continue
		}
		for _, a := range addrs {

			switch v := a.(type) {
			case *net.IPNet:
				// Test if it's ipv4
				if v.IP.To4() != nil {
					interfaces = append(interfaces, Iface{Name: i.Name, IP: v.IP.String()})
				}

			case *net.IPAddr:
				// Test if it's ipv4
				if v.IP.To4() != nil {
					interfaces = append(interfaces, Iface{Name: i.Name, IP: v.IP.String()})
				}
			}
		}
	}
	return interfaces
}
