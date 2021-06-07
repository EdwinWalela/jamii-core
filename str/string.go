package str

import "encoding/hex"

func bytetoHex(src []byte) string {
	return hex.Dump(src)
}
