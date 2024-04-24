package utils

import (
	"net"
	"strings"
)

// GetIPByInterface 通过网卡获取 IP 地址
func GetIPByInterface(name string) (ipv4 string, ipv6 string, err error) {
	var targetName string

	if name != "" {
		targetName = name
	} else {
		interfaces, err := net.Interfaces()
		if err != nil {
			return "", "", err
		}
		targetName = interfaces[0].Name
	}

	ipv4s, ipv6s, err := getIPByInterface(targetName)
	if err != nil {
		return "", "", nil
	}

	if len(ipv4s) > 0 {
		ipv4 = ipv4s[0]
	}

	if len(ipv6s) > 0 {
		ipv6 = ipv6s[0]
	}

	return ipv4, ipv6, nil
}

// getIPByInterface 通过网卡获取 IP 地址
func getIPByInterface(name string) ([]string, []string, error) {
	ipv4s := make([]string, 0)
	ipv6s := make([]string, 0)

	inter, err := net.InterfaceByName(name)
	if err != nil {
		return nil, nil, err
	}

	addrs, err := inter.Addrs()
	if err != nil {
		return nil, nil, err
	}

	for _, addr := range addrs {
		targetAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		if targetAddr.IP.IsLoopback() {
			continue
		}

		if targetAddr.IP.To4() != nil {
			ipv4s = append(ipv4s, strings.TrimSuffix(targetAddr.String(), "/24"))
		} else if targetAddr.IP.To16() != nil {
			if strings.HasPrefix("fe80::", targetAddr.String()) {
				continue
			}

			ipv6s = append(ipv6s, strings.TrimSuffix(targetAddr.String(), "/64"))
		}
	}

	return ipv4s, ipv6s, nil
}
