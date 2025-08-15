package packet

import (
	"encoding/hex"
	"fmt"
)

type UUID []byte

func (u UUID) String() string {
	return hex.EncodeToString(u)
}

func (u UUID) FormattedString() string {
	uuid := hex.EncodeToString(u)

	// Construct the UUID string with hyphens
	formattedUUID := fmt.Sprintf("%s-%s-%s-%s-%s",
		uuid[0:8],
		uuid[8:12],
		uuid[12:16],
		uuid[16:20],
		uuid[20:32])

	return formattedUUID
}
