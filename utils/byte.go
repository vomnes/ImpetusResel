package utils

import "strconv"

// ByteArrayJoin joins each element of a byte array separating with a given separator
// and return result as a string
func ByteArrayJoin(array []byte, separator string) string {
	var result string
	len := len(array)
	for i, byte := range array {
		result += strconv.Itoa(int(byte))
		if i < len-1 {
			result += separator
		}
	}
	return result
}
