package cmd

/*
Copyright © 2023 Nicola Iacovelli <nicolaiacovelli98@gmail.com>

*/

import (
	"embed"
	"gopkg.in/yaml.v3"
	"jacobCloudAdapter/model"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

//go:embed kustomizationBaseTemplate.tmpl
var kbtFS embed.FS

//go:embed kustomizationOverlaysTemplate.tmpl
var kotFS embed.FS

//go:embed config.yml
var configFS embed.FS

// adaptCmd represents the adapt command
var adaptCmd = &cobra.Command{
	Use:   "adapt",
	Short: "Adapt the program in the current directory to run in cloud environment",
	Long: `This command will adapt the program in the current directory to run in cloud environment creating k8s folder
with coll and prod environment folder and kustomization.yaml file with the program file property 
copied from src/main/resources/${programName}.yml file.
`,
	Run: func(cmd *cobra.Command, args []string) {
		programName, _ := cmd.Flags().GetString("programName")
		isNewProgram, _ := cmd.Flags().GetBool("isNewProgram")

		doTheMagic(programName, isNewProgram)
	},
}

func init() {
	rootCmd.AddCommand(adaptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	adaptCmd.PersistentFlags().Int("jacobVersion", 3, "The version of jacob program to adapt")

	adaptCmd.PersistentFlags().BoolP("isNewProgram", "n", false, "If the program is new or not")

	adaptCmd.PersistentFlags().String("programName", "", "The name of the program to adapt")
	err := adaptCmd.MarkPersistentFlagRequired("programName")
	if err != nil {
		println("programName flag is REQUIRED!")
		return
	}

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adaptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func doTheMagic(programName string, isNewProgram bool) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	baseTemplateModel := model.KustomizationBaseTemplateModel{ProgramName: programName}

	overlaysTemplateModel := model.KustomizationOverlaysTemplateModel{
		ProgramName: programName,
		Filename:    strings.ToUpper(programName) + ".yml",
	}

	k8sBaseDir := workingDir + "/k8s"

	baseDir := k8sBaseDir + "/base"
	overlaysDir := k8sBaseDir + "/overlays"
	productionDir := overlaysDir + "/prod"
	collDir := overlaysDir + "/coll"

	baseTemplateFilePath := "kustomizationBaseTemplate.tmpl"
	overlaysTemplateFilePath := "kustomizationOverlaysTemplate.tmpl"

	kustomizationBaseFilePath := baseDir + "/kustomization.yaml"
	kustomizationProdFilePath := productionDir + "/kustomization.yaml"
	kustomizationCollFilePath := collDir + "/kustomization.yaml"

	propertyFilePath := workingDir + "/src/main/resources/" + overlaysTemplateModel.Filename
	var propertyMap map[string]interface{}

	propertyFile, err := os.ReadFile(propertyFilePath)
	if err != nil {
		println("Error: " + err.Error())
		return err
	}

	err = yaml.Unmarshal(propertyFile, &propertyMap)
	if err != nil {
		println("Error: " + err.Error())
		return err
	}

	prodProperty := propertyMap
	collProperty := propertyMap

	dataSourceProperty, ok := propertyMap["dataSourceProperties"]
	if !ok {
		println("Error: dataSourceProperties not found in " + overlaysTemplateModel.Filename)
		return err
	}

	collDataSourceProperty, ok := dataSourceProperty.(map[string]interface{})["MAIN"]
	prodDataSourceProperty, ok := dataSourceProperty.(map[string]interface{})["MAIN"]
	if !ok {
		println("Error: MAIN not found in dataSourceProperties")
		return err
	}

	basePath := os.Getenv("JACOB_ADAPTER_BASE_PATH") + programName

	prodDataSourceProperty.(map[string]interface{})["url"] = os.Getenv("JACOB_ADAPTER_PROD_URL")
	prodDataSourceProperty.(map[string]interface{})["user"] = os.Getenv("JACOB_ADAPTER_PROD_USERNAME")
	prodDataSourceProperty.(map[string]interface{})["password"] = os.Getenv("JACOB_ADAPTER_PROD_PASSWORD")

	prodProperty["dataSourceProperties"] = prodDataSourceProperty
	prodProperty["basePath"] = basePath

	productionPropertyFile, err := yaml.Marshal(prodProperty)
	if err != nil {
		println("Error: " + err.Error())
		return err
	}

	collDataSourceProperty.(map[string]interface{})["url"] = os.Getenv("JACOB_ADAPTER_COLL_URL")
	collDataSourceProperty.(map[string]interface{})["user"] = os.Getenv("JACOB_ADAPTER_COLL_USERNAME")
	collDataSourceProperty.(map[string]interface{})["password"] = os.Getenv("JACOB_ADAPTER_COLL_PASSWORD")

	collProperty["dataSourceProperties"] = collDataSourceProperty
	collProperty["basePath"] = basePath

	collPropertyFile, err := yaml.Marshal(collProperty)
	if err != nil {
		println("Error: " + err.Error())
		return err
	}

	if isNewProgram {
		err = os.Mkdir(k8sBaseDir, 0755)
		err = os.Mkdir(baseDir, 0755)
		err = os.Mkdir(overlaysDir, 0755)
		err = os.Mkdir(productionDir, 0755)
		err = os.Mkdir(collDir, 0755)

		if err != nil {
			println("Error: " + err.Error())
			return err
		}
	}

	err = os.WriteFile(productionDir+"/"+overlaysTemplateModel.Filename, productionPropertyFile, 0755)
	if err != nil {
		println("Error: " + err.Error())
		return err
	}

	err = os.WriteFile(collDir+"/"+overlaysTemplateModel.Filename, collPropertyFile, 0755)
	if err != nil {
		println("Error: " + err.Error())
		return err
	}

	if isNewProgram {
		kustomizationBaseFile, err := os.Create(kustomizationBaseFilePath)
		if err != nil {
			println("Error: " + err.Error())
			return err
		}
		kustomizationProdFile, err := os.Create(kustomizationProdFilePath)
		if err != nil {
			println("Error: " + err.Error())
			return err
		}
		kustomizationCollFile, err := os.Create(kustomizationCollFilePath)
		if err != nil {
			println("Error: " + err.Error())
			return err
		}

		tmplBase, err := template.ParseFS(kbtFS, baseTemplateFilePath)
		if err != nil {
			println("Error: " + err.Error())
			return err
		}
		tmplOverlays, err := template.ParseFS(kotFS, overlaysTemplateFilePath)
		if err != nil {
			println("Error: " + err.Error())
			return err
		}

		err = tmplBase.Execute(kustomizationBaseFile, baseTemplateModel)
		if err != nil {
			return err
		}
		err = tmplOverlays.Execute(kustomizationProdFile, overlaysTemplateModel)
		if err != nil {
			return err
		}

		err = tmplOverlays.Execute(kustomizationCollFile, overlaysTemplateModel)
		if err != nil {
			return err
		}
	}

	return nil
}
