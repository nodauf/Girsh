package utils

import (
	"fmt"
	"net"
)

var Backspace = []byte{8}
var Newline = []byte{10}

func SliceStringContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SliceByteContains(s []byte, e byte) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func Min(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

type iface struct {
	Name string
	IP   string
}

func ListInterfaces() []iface {
	interfaces := []iface{}
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
		return nil
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
			continue
		}
		for _, a := range addrs {

			switch v := a.(type) {
			case *net.IPNet:
				// Test if it's ipv4
				if v.IP.To4() != nil {
					interfaces = append(interfaces, iface{Name: i.Name, IP: v.IP.String()})
				}

			case *net.IPAddr:
				// Test if it's ipv4
				if v.IP.To4() != nil {
					interfaces = append(interfaces, iface{Name: i.Name, IP: v.IP.String()})
				}
			}
		}
	}
	return interfaces
}
