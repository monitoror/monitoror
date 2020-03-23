package system

import (
	"net"

	"golang.org/x/net/icmp"
)

func IsRawSocketAvailable() bool {
	_, err := icmp.ListenPacket("ip4:icmp", "")
	return err == nil
}

// ListLocalhostIpv4 list IP of every local network interfaces
func ListLocalhostIpv4() ([]string, error) {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			if ip, ok := addr.(*net.IPNet); ok && ip.IP.To4() != nil {
				ips = append(ips, ip.IP.String())
			}
		}
	}

	return ips, nil
}
