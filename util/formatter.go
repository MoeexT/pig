package util

import (
    "fmt"
    "strings"
)

// convert a byte array with specified length to *string
func BytesString(bytes []byte, limit int) string {
    length := len(bytes)
    if length == 0 || limit == 0 {
        return "[null]"
    }

    min := Min(length, limit) - 1

    var sb strings.Builder
    sb.WriteRune('[')

    for i := 0; i < min; i++ {
        sb.WriteString(fmt.Sprintf("%02X", bytes[i]))
        sb.WriteRune(' ')
    }
    sb.WriteString(fmt.Sprintf("%02X", bytes[min]))
    if limit < length {
        sb.WriteString("...")
    }
    sb.WriteRune(']')
    return sb.String()
}
