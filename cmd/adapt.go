package cmd

/*
Copyright © 2023 Nicola Iacovelli <nicolaiacovelli98@gmail.com>

*/

import (
	"embed"
	"errors"
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

	adaptCmd.PersistentFlags().Bool("isNewProgram", false, "If the program is new or not (For feature use)")

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

	configFile, err := configFS.ReadFile("config.yml")
	if err != nil {
		println("Error: " + err.Error())
		return err
	}

	var config model.ApplicationConfig

	err = yaml.Unmarshal(configFile, &config)
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

	_, ok = dataSourceProperty.(map[string]interface{})["H2"]
	if ok {
		dataSourceProperty.(map[string]interface{})["H2"].(map[string]interface{})["url"] = "jdbc:h2:file:" + config.Base.Path + programName + "/" + programName + "_db"
	}

	collDataSourceProperty := dataSourceProperty
	prodDataSourceProperty := dataSourceProperty

	prodDataSourceProperty.(map[string]interface{})["MAIN"].(map[string]interface{})["url"] = config.Prod.Url
	prodDataSourceProperty.(map[string]interface{})["MAIN"].(map[string]interface{})["user"] = config.Prod.Username
	prodDataSourceProperty.(map[string]interface{})["MAIN"].(map[string]interface{})["password"] = config.Prod.Password

	prodProperty["basePath"] = config.Base.Path + programName

	productionPropertyFile, err := yaml.Marshal(prodProperty)
	if err != nil {
		println("Error: " + err.Error())
		return err
	}

	collDataSourceProperty.(map[string]interface{})["MAIN"].(map[string]interface{})["url"] = config.Coll.Url
	collDataSourceProperty.(map[string]interface{})["MAIN"].(map[string]interface{})["user"] = config.Coll.Username
	collDataSourceProperty.(map[string]interface{})["MAIN"].(map[string]interface{})["password"] = config.Coll.Password

	collProperty["dataSourceProperties"] = collDataSourceProperty
	collProperty["basePath"] = config.Base.Path + programName

	collPropertyFile, err := yaml.Marshal(collProperty)
	if err != nil {
		println("Error: " + err.Error())
		return err
	}

	err = createDirIfNotExists(k8sBaseDir)
	err = createDirIfNotExists(baseDir)
	err = createDirIfNotExists(overlaysDir)
	err = createDirIfNotExists(productionDir)
	err = createDirIfNotExists(collDir)

	if err != nil {
		println("Error: " + err.Error())
		return err
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

	return nil
}

func createDirIfNotExists(dir string) error {
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
