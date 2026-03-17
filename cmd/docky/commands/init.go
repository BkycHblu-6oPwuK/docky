package commands

import (
	"fmt"
	"os"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config/framework"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/initialization"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Инициализация проекта",
	Run: func(cmd *cobra.Command, args []string) {
		err := initialization.InitDockerComposeFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
		err = globaltools.InitSiteDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
		yamlConfig := config.GetYamlConfig()
		switch yamlConfig.FrameworkName {
		case framework.Laravel:
			fmt.Println("Инициализация ларавел")
			initialization.InitLaravel()
		case framework.Symfony:
			fmt.Println("Инициализация симфони")
			initialization.InitSymfony()
		case framework.Yii2:
			fmt.Println("Инициализация yii2")
			initialization.InitYii2(yamlConfig)
		case framework.Yii3:
			fmt.Println("Инициализация yii3")
			initialization.InitYii3()
		}

		fmt.Println("✅ Инициализация проекта завершена!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
