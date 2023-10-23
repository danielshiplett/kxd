package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/radiusmethod/kxd/src/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var fileCmd = &cobra.Command{
	Use:     "file",
	Short:   "Kubeconfig file command",
	Aliases: []string{"f"},
	Long:    "This is the default file command.",
}

var currentFileCmd = &cobra.Command{
	Use:     "current",
	Short:   "Shows currently set kubeconfig",
	Aliases: []string{"c"},
	Long:    "This shows the current set kubeconfig file.",
	Run: func(cmd *cobra.Command, args []string) {
		err := runGetCurrentConfig()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var switchFileCmd = &cobra.Command{
	Use:     "switch",
	Short:   "Switch kubeconfig",
	Aliases: []string{"s"},
	Long:    "This allows for switching of your kubeconfig.",
	Run: func(cmd *cobra.Command, args []string) {
		err := runConfigSwitcher()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var listFileCmd = &cobra.Command{
	Use:     "list",
	Short:   "List kubeconfigs",
	Aliases: []string{"l"},
	Long:    "This displays a simple list of your kubeconfigs.",
	Run: func(cmd *cobra.Command, args []string) {
		err := runConfigLister()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	fileCmd.AddCommand(switchFileCmd, currentFileCmd, listFileCmd)
	rootCmd.AddCommand(fileCmd)
}

func runConfigSwitcher() error {
	configs := utils.GetConfigs()
	err := utils.TouchFile(filepath.Join(utils.GetHomeDir(), ".kxd"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(utils.NoticeColor, "Kubeconfig Switcher\n")
	prompt := promptui.Select{
		Label:        fmt.Sprintf(utils.PromptColor, "Choose a config"),
		Items:        configs,
		HideHelp:     true,
		HideSelected: true,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   fmt.Sprintf("%s {{ . | cyan }}", promptui.IconSelect),
			Inactive: "  {{.}}",
			Selected: "  {{ . | cyan }}",
		},
		Searcher:          utils.NewPromptUISearcher(configs),
		StartInSearchMode: true,
		Stdout:            &utils.BellSkipper{},
	}

	_, result, err := prompt.Run()
	if err != nil {
		utils.CheckError(err)
	}

	fmt.Printf(utils.PromptColor, "Choose a config")
	fmt.Printf(utils.NoticeColor, "? ")
	fmt.Printf(utils.CyanColor, result)
	fmt.Println("")

	if result == "default" {
		result = "config"
	}
	utils.WriteFile(result, utils.GetHomeDir())

	return nil
}

func runGetCurrentConfig() error {
	kubeconfigPath := utils.GetEnv("KUBECONFIG", filepath.Join(utils.GetHomeDir(), ".kube/config"))
	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		log.Fatal("No current kubeconfig found.")
	} else if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(kubeconfigPath)
	}
	return nil
}

func runConfigLister() error {
	configs := utils.GetConfigs()
	for _, c := range configs {
		fmt.Println(c)
	}
	return nil
}
