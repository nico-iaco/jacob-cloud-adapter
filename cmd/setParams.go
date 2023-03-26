/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
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

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		setParameter(flag.Name, flag.Value.String())
	})

}

func setParameter(paramName string, paramValue string) {
	switch paramName {
	case "basePath":
		err := os.Setenv("JACOB_ADAPTER_BASE_PATH", paramValue)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		break
	case "prodUsername":
		err := os.Setenv("JACOB_ADAPTER_PROD_USERNAME", paramValue)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		break
	case "prodPwd":
		err := os.Setenv("JACOB_ADAPTER_PROD_PWD", paramValue)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		break
	case "prodUrl":
		err := os.Setenv("JACOB_ADAPTER_PROD_URL", paramValue)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		break
	case "collUsername":
		err := os.Setenv("JACOB_ADAPTER_COLL_USERNAME", paramValue)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		break
	case "collPwd":
		err := os.Setenv("JACOB_ADAPTER_COLL_PWD", paramValue)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		break
	case "collUrl":
		err := os.Setenv("JACOB_ADAPTER_COLL_URL", paramValue)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		break
	default:
		println("Parameter not found!")
	}
}
