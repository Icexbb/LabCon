package utils

import (
	"os"
)

func FileExist(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}
