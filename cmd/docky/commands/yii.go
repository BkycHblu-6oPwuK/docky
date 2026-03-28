package commands

import (
	"fmt"
	"os"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config/framework"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"

	"github.com/spf13/cobra"
)

var yiiCmd = &cobra.Command{
	Use:                "yii",
	Short:              "Запускает yii команду в контейнере " + composefiletools.App,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		globaltools.ValidateWorkDir()
		if err := execYiiInContainer(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
	},
}

func init() {
	switch config.GetCurFramework() {
	case framework.Yii2, framework.Yii3:
		rootCmd.AddCommand(yiiCmd)
	}
}

func execYiiInContainer(args []string) error {
	execArgs := append([]string{
		"yii",
	}, args...)

	return execPhpInContainer(execArgs)
}
