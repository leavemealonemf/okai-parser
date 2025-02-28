package okaiparsetools

import (
	"fmt"
	"strings"
)

func CutPacket(packet string, sep string) (string, error) {
	_, res, found := strings.Cut(packet, sep)
	if !found {
		return "", fmt.Errorf("failed to cut packet by sep")
	}

	return res, nil
}

func SplitParams(packet string, sep string) []string {
	res := strings.Split(packet, sep)
	return res
}
