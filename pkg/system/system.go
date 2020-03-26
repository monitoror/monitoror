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
	conn, _ := net.Dial("udp", "255.255.255.255:80")
	defer conn.Close()
	networkIP := conn.LocalAddr().(*net.UDPAddr).IP
	return networkIP.To4().String()
}
