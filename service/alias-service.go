package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/vigneshwaran-48/zshell/models"
	"github.com/vigneshwaran-48/zshell/utils"
)

func AddAlias(alias *models.Alias) error {
	if alias == nil {
		return errors.New("alias is required")
	}
	if alias.Name == "" {
		return errors.New("alias name is required")
	}
	if alias.Command == "" {
		return errors.New("alias command is required")
	}
	existingAlias, err := FindAliasByName(alias.Name)
	if err != nil {
		return err
	}
	if existingAlias != nil {
		return fmt.Errorf("Alias already exists with name %s", alias.Name)
	}
	return storeAlias(alias)
}

func FindAliasByName(aliasName string) (*models.Alias, error) {
	aliases, err := FindAllAlias()
	if err != nil {
		return nil, err
	}
	for _, alias := range aliases {
		if alias.Name == aliasName {
			return &alias, nil
		}
	}
	return nil, nil
}

func FindByCommand(command string) (*models.Alias, error) {
	aliases, err := FindAllAlias()
	if err != nil {
		return nil, err
	}
	for _, alias := range aliases {
		if alias.Command == command {
			return &alias, nil
		}
	}
	return nil, nil
}

func FindAllAlias() ([]models.Alias, error) {
	aliasFilePath, err := getAliasStoreFile()
	if err != nil {
		return nil, err
	}
	if !utils.IsFileExists(aliasFilePath) {
		err = CreateDefaultAliasData()
		if err != nil {
			return nil, err
		}
	}
	data, err := os.ReadFile(aliasFilePath)
	if err != nil {
		return nil, err
	}
	var aliases []models.Alias
	err = json.Unmarshal(data, &aliases)
	if err != nil {
		return nil, err
	}
	return aliases, nil
}

func RemoveAlias(name string) error {
	aliases, err := FindAllAlias()
	if err != nil {
		return err
	}
	aliases = utils.Filter(aliases, func(alias models.Alias) bool {
		return alias.Name != name
	})
	err = storeAliases(aliases)
	if err != nil {
		return err
	}
	return nil
}

func storeAlias(alias *models.Alias) error {
	if alias == nil {
		return errors.New("alias is required")
	}
	aliases, err := FindAllAlias()
	if err != nil {
		return err
	}
	aliases = append(aliases, *alias)
	err = storeAliases(aliases)
	if err != nil {
		return err
	}
	return nil
}

func storeAliases(aliases []models.Alias) error {
	data, err := json.Marshal(aliases)
	if err != nil {
		return err
	}
	aliasFilePath, err := getAliasStoreFile()
	if err != nil {
		return err
	}
	err = os.WriteFile(aliasFilePath, data, 0o755)
	if err != nil {
		return err
	}
	return nil
}

func getAliasStoreFile() (string, error) {
	configDir, err := utils.GetConfigDirPath()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/alias.json", configDir), nil
}

func CreateDefaultAliasData() error {
	aliasFilePath, err := getAliasStoreFile()
	if err != nil {
		return err
	}
	content := "[]"
	err = os.WriteFile(aliasFilePath, []byte(content), 0o755)
	if err != nil {
		return err
	}
	return nil
}
