package epay

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// Checksum calculates the Values checksum using the epay specific format.
func Checksum(q url.Values, secret string) string {
	keys := make([]string, 0, len(q))
	for k := range q {
		if k == "CHECKSUM" {
			continue
		}
		keys = append(keys, k)
	}

	sort.Strings(keys)
	message := ""
	for _, k := range keys {
		message += fmt.Sprintf("%s%s\n", k, q[k][0])
	}

	key := []byte(secret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// IsContractCode validates a 7-digit code with Luhn checksum.
// Used to determine if an IDN is a Contract ID (for UCRM query search)
// or a User ID (for UCRM userIdent search).
func IsContractCode(code string) bool {
	const length = 7

	// Check if the length of the input is valid
	if len(code) != length {
		return false
	}

	// Separate the base number (neid) and the check digit (ncrc)
	neid := code[:length-1]
	ncrc, err := strconv.Atoi(code[length-1:])
	if err != nil {
		return false
	}

	// Reverse the neid string
	reversedNeid := reverseString(neid)

	// Calculate the checksum
	sum := 0
	for i := 0; i < len(reversedNeid); i++ {
		digit, _ := strconv.Atoi(string(reversedNeid[i]))
		if i%2 == 0 { // Odd positions (in the reversed string)
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	// Compute the expected check digit
	crc := (10 - (sum % 10)) % 10

	// Check if the computed check digit matches the provided one
	return crc == ncrc
}

// reverseString reverses the input string
func reverseString(s string) string {
	var sb strings.Builder
	for i := len(s) - 1; i >= 0; i-- {
		sb.WriteByte(s[i])
	}
	return sb.String()
}
