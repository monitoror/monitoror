package system

import (
	"net"

	"golang.org/x/net/icmp"
)

func IsRawSocketAvailable() bool {
	_, err := icmp.ListenPacket("ip4:icmp", "")
	return err == nil
}

func GetNetworkIP() string {
	ip := "0.0.0.0"

	conn, err := net.Dial("udp", "255.255.255.255:80")
	if err == nil {
		defer conn.Close()
		networkIP := conn.LocalAddr().(*net.UDPAddr).IP
		ip = networkIP.To4().String()
	}

	return ip
}
