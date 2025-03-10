package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vigneshwaran-48/zshell/utils"
)

var rootCmd = &cobra.Command{
	Use:     "zshell",
	Version: "1.0",
	Long:    "Zoho shell",
}

func createDefaultConfig(path string) error {
	clientId := os.Getenv("ZSHELL_CLIENT_ID")
	clientSecret := os.Getenv("ZSHELL_CLIENT_SECRET")

	if clientId == "" {
		cobra.CheckErr("Environment variable ZSHELL_CLIENT_ID is not configured")
	}
	if clientSecret == "" {
		cobra.CheckErr("Environment variable ZSHELL_CLIENT_SECRET is not configured")
	}

	viper.Set(utils.APP_NAME, "ZShell")
	viper.Set(utils.DEFAULT_DC, "zoho.com")
	viper.Set(utils.PORT, 31200)
	viper.Set(utils.SCOPE, "ZohoMail.accounts.ALL,ZohoMail.organization.accounts.ALL,ZohoMail.messages.ALL,ZohoMail.attachments.ALL,ZohoMail.tags.ALL,ZohoMail.folders.ALL,ZohoMail.tasks.ALL,ZohoMail.notes.ALL,ZohoMail.links.ALL,ZohoMail.settings.ALL,ZohoMail.search.ALL,ZohoMail.partner.organization.ALL")
	viper.Set(utils.CLIENT_ID, clientId)
	viper.Set(utils.CLIENT_SECRET, clientSecret)
	viper.Set(utils.REDIRECT_URI, "http://localhost:31200/oauth/callback")

	if err := viper.WriteConfigAs(path); err != nil {
		return fmt.Errorf("Error creating config file: %v", err)
	}
	return nil
}

func initViperConfig() error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		cobra.CheckErr(err)
	}
	err = utils.EnsureDirExists(fmt.Sprintf("%s/.zmail", userHomeDir))
	if err != nil {
		cobra.CheckErr(err)
	}

	configPath := fmt.Sprintf("%s/.zmail", userHomeDir)

	viper.SetConfigName("config")
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		err := createDefaultConfig(fmt.Sprintf("%s/config.yaml", configPath))
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	err := initViperConfig()
	if err != nil {
		cobra.CheckErr(err)
	}

	rootCmd.PersistentFlags().String("dc", viper.GetString(utils.DEFAULT_DC), "Which dc to use like zoho.com, zoho.in, zoho.eu, etc")
	rootCmd.PersistentFlags().Int64("account", 0, "Account Id")

	rootCmd.Flags().String("password", "", "Password to encrypt/decrypt access tokens")
}
