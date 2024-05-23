package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

var basePath = os.Getenv("TESTPATH")

func TestHashFile(t *testing.T) {
	const testHashString = "44860d92477931a27a960090edde8546"

	resolved, err := filepath.Abs(filepath.Join(basePath, "test.txt"))

	if err != nil {
		t.Fatal(err.Error())
	}

	hash, err := HashFile(resolved)

	if err != nil {
		t.Fatal(err.Error())
	}

	if *hash != testHashString {
		t.Errorf("Expected hash to be %s, got %s", testHashString, *hash)
		t.Fatal("Hashes don't match")
	}
}

func TestHashString(t *testing.T) {
	const testHashString = "047d28cc74bde19d9a128231f9bd4d82"
	hash := HashString("test")

	if hash != testHashString {
		t.Errorf("Expected hash to be %s, got %s", testHashString, hash)
		t.Fatal("Hashes don't match")
	}
}

func TestCompareStringHashes(t *testing.T) {
	const testHashString = "047d28cc74bde19d9a128231f9bd4d82"

	match := CompareStringHashes("test", testHashString)

	if !match {
		t.Errorf("Expected hash to be %s, got %s", testHashString, testHashString)
		t.Fatal("Hashes don't match")
	}
}

func TestCompareFileHashes(t *testing.T) {
	const testHashString = "44860d92477931a27a960090edde8546"

	resolved, err := filepath.Abs(filepath.Join(basePath, "test.txt"))

	if err != nil {
		t.Fatal(err.Error())
	}

	match, hash, err := CompareFileHashes(resolved, testHashString)

	if err != nil {
		t.Fatal(err.Error())
	}

	if !match {
		t.Errorf("Expected hash to be %s, got %s", testHashString, *hash)
		t.Fatal("Hashes don't match")
	}
}
