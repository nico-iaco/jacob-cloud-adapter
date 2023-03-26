package cmd

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"jacobCloudAdapter/model"
	"os"

	"github.com/spf13/cobra"
)

/*
Copyright © 2023 Nicola Iacovelli <nicolaiacovelli98@gmail.com>

*/

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jacobCloudAdapter",
	Short: "Adapt jacob program to run in cloud environment",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jacobCloudAdapter.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	configFile, err := configFS.ReadFile("config.yml")
	if err != nil {
		println("Error: " + err.Error())
		return
	}

	var config model.ApplicationConfig

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	err = os.Setenv("JACOB_ADAPTER_BASE_PATH_ENV", config.Base.Path)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	err = os.Setenv("JACOB_ADAPTER_PROD_USERNAME_ENV", config.Prod.Username)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	err = os.Setenv("JACOB_ADAPTER_PROD_PWD_ENV", config.Prod.Password)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	err = os.Setenv("JACOB_ADAPTER_PROD_URL_ENV", config.Prod.Url)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	err = os.Setenv("JACOB_ADAPTER_COLL_USERNAME_ENV", config.Coll.Username)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	err = os.Setenv("JACOB_ADAPTER_COLL_PWD_ENV", config.Coll.Password)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	err = os.Setenv("JACOB_ADAPTER_COLL_URL_ENV", config.Coll.Url)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

}
