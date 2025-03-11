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

func GetAuthDataFile() (string, error) {
	configDir, err := GetConfigDirPath()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/auth.json", configDir), nil
}

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	// Do we need to check for other errors?
	return err == nil
}

func GetConfigDirPath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	config := fmt.Sprintf("%s/.zmail", userHomeDir)
	if !IsFileExists(config) {
		err = os.Mkdir(config, 0o755)
		if err != nil {
			return "", err
		}
	}
	return config, nil
}
