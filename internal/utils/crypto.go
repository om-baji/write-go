package utils

import "github.com/snksoft/crc"

func GenerateCRC(data []byte) uint64 {
	hash := crc.NewHash(crc.XMODEM)
	return hash.CalculateCRC([]byte(data))
}
