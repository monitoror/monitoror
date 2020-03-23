package system

import "golang.org/x/net/icmp"

func IsRawSocketAvailable() bool {
	_, err := icmp.ListenPacket("ip4:icmp", "")
	return err == nil
}
