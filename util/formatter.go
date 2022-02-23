package util

import (
	"fmt"
	"strings"
)

// convert a byte array with specified length to *string
func BytesString(bytes []byte, length int) string {
	var sb strings.Builder
	sb.WriteRune('[')
	
	for i := 0; i < length - 1; i++ {
		sb.WriteString(fmt.Sprintf("%02X", bytes[i]))
		sb.WriteRune(' ')
	}
	sb.WriteString(fmt.Sprintf("%02X", bytes[length - 1]))
	sb.WriteRune(']')
	return sb.String()
}
