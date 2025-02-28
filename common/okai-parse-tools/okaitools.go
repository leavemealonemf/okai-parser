package okaiparsetools

import (
	"strings"
)

func CutPacket(packet string, sep string) string {
	pkt := make([]byte, 4096)

	for i := 0; i < len(packet); i++ {
		pkt[i] = packet[i]
		if packet[i] == []byte(sep)[0] {
			break
		}
	}

	return string(pkt)
}

func SplitParams(packet string, sep string) []string {
	res := strings.Split(packet, sep)
	return res
}
