package cloner

import (
	"fmt"
	"os"
)

// CheckThatPathDoesntExist checks provided path doesn't exist
func CheckThatPathDoesntExist(value string) error {
	if _, err := os.Stat(value); err == nil {
		return fmt.Errorf("path %s already exists", value)
	}
	return nil
}
