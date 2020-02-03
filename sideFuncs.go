package pconfig

import (
	"crypto/rand"
	"encoding/hex"
	"os"
)

// FileExists returns true if the file exists, false otherwise.
func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateDirIfNotExist creates the directory path if it doesn't exist.
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

// FilterUnique takes an array of strings and returns an array with unique entries.
func FilterUnique(vals []string) []string {
	var tmp []string
	dupe := make(map[string]bool)
	for _, v := range vals {
		if !dupe[v] {
			dupe[v] = true
			tmp = append(tmp, v)
		}
	}
	return tmp
}

// MakeHex returns a random Hex string based on n length.
func MakeHex(n int) string {
	b := randomBytes(n)
	hexstring := hex.EncodeToString(b)
	return hexstring
}

func randomBytes(n int) []byte {
	return makeByte(n)
}

func makeByte(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}
