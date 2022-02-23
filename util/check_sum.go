package util

// udp checksum
func CheckSum(buf []byte) uint16 {
	sum := 0
	for n := 1; n < len(buf)-1; n += 2 {
		sum += (int(buf[n])<<8) + int(buf[n+1])
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += sum >> 16
	return (uint16(^sum))
}
