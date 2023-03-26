/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
	"jacobCloudAdapter/model"
)

// setParamsCmd represents the setParams command
var setParamsCmd = &cobra.Command{
	Use:   "setParams",
	Short: "Set configuration parameters for the program",
	Long: `This command will set configuration parameters for the program in the current directory.
`,
	Run: setConfigs,
}

func init() {
	rootCmd.AddCommand(setParamsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	setParamsCmd.PersistentFlags().String("basePath", "", "The base path of the program")

	setParamsCmd.PersistentFlags().String("prodUsername", "", "The username of the prod environment db user")
	setParamsCmd.PersistentFlags().String("prodPwd", "", "The password of the prod environment db user")
	setParamsCmd.PersistentFlags().String("prodUrl", "", "The url of the prod environment db")

	setParamsCmd.PersistentFlags().String("collUsername", "", "The username of the coll environment db user")
	setParamsCmd.PersistentFlags().String("collPwd", "", "The password of the coll environment db user")
	setParamsCmd.PersistentFlags().String("collUrl", "", "The url of the coll environment db")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setParamsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func setConfigs(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		println("No arguments passed!")
		return
	}

	configFile, err := configFS.ReadFile("config.yml")
	if err != nil {
		println("Error: " + err.Error())
		return
	}

	var config model.ApplicationConfig

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		println("Error: " + err.Error())
		return
	}

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		setParameter(flag.Name, flag.Value.String(), &config)
	})

	fmt.Println(config)

}

func setParameter(paramName string, paramValue string, config *model.ApplicationConfig) {
	switch paramName {
	case "basePath":
		config.Base.Path = paramValue
		break
	case "prodUsername":
		config.Prod.Username = paramValue
		break
	case "prodPwd":
		config.Prod.Password = paramValue
		break
	case "prodUrl":
		config.Prod.Url = paramValue
		break
	case "collUsername":
		config.Coll.Username = paramValue
		break
	case "collPwd":
		config.Coll.Password = paramValue
		break
	case "collUrl":
		config.Coll.Url = paramValue
		break
	default:
		println("Parameter not found!")
	}
}
