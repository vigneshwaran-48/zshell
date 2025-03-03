package utils

import (
	"fmt"
	"os"
)

func EnsureDirExists(dir string) error {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0o755)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("'%s' exists but is not a directory", dir)
	}
	return nil
}
