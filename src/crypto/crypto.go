package crypto

import (
	"encoding/hex"

	"github.com/kalafut/imohash"
)

func CompareFileHashes(path string, hash string) (bool, *string, error) {
	fileHash, err := HashFile(path)

	if err != nil {
		return false, nil, err
	}

	if *fileHash == hash {
		return true, fileHash, nil
	}

	return false, fileHash, nil
}

func CompareStringHashes(s string, hash string) bool {
	stringHash := HashString(s)
	return stringHash == hash
}

func HashFile(path string) (*string, error) {
	hasher := imohash.New()

	checkbytes, err := hasher.SumFile(path)

	if err != nil {
		return nil, err
	}

	hexString := EncodeToHex(ToRegularBytes(checkbytes))

	return &hexString, nil
}

func HashString(s string) string {
	hasher := imohash.New()

	checkbyte := hasher.Sum([]byte(s))

	hexString := EncodeToHex(ToRegularBytes(checkbyte))

	return hexString
}

func ToRegularBytes(bytes [16]byte) []byte {
	regularBytes := make([]byte, 16)
	for i := 0; i < 16; i++ {
		regularBytes[i] = bytes[i]
	}
	return regularBytes
}

func EncodeToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}
