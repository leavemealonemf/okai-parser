package okaiparsetools

import (
	"strings"
)

func CutPacket(packet string, sep string) string {
	idx := strings.Index(packet, sep)
	if idx == -1 {
		return packet
	}
	return packet[:idx+1]
}

func SplitParams(packet string, sep string) []string {
	res := strings.Split(packet, sep)
	return res
}
