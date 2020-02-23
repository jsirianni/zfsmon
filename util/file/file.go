package file

import (
	"os"
	"io/ioutil"
)

// Write the json content to a given path, readable only by
// the user
func WriteFile(content []byte, filePath string) error {
	return ioutil.WriteFile(filePath, content, 0600)
}

// Read the file and return a byte array
func ReadFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
	    return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
