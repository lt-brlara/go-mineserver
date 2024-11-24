package packet

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

// readUUID parses a 128-bit integer and returns in hyphenated hexadecimal
// format.
func readUUID(data *bytes.Buffer) (UUID, error) {
	arr := make([]byte, 16)
	_, err := data.Read(arr)
	if err != nil {
		return UUID(""), err
	}

	uuid := hex.EncodeToString(arr)
	uuidStr := fmt.Sprintf("%s-%s-%s-%s-%s", uuid[:8], uuid[8:12], uuid[12:16], uuid[16:20], uuid[20:])

	return UUID(uuidStr), nil
}
