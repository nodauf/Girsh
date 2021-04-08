package utils

import (
	"fmt"
	"io"
	"net"
)

func ReadUntil(con io.Reader, buf []byte, untilByte []byte) (int, error) {

	bufsize := len(buf)
	n := 0
	for { // copy loop
		// non buffered byte-by-byte
		if _, err := con.Read(buf[n : n+1]); err != nil {
			return n, err
		}

		if SliceByteContains(untilByte, buf[n]) || n == bufsize { // reached CR or size limit
			return n, nil

		} else {
			n++
		}

	}
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
