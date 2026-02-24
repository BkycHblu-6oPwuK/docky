package commands

import (
	"fmt"
	"os"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"

	"github.com/spf13/cobra"
)

var workdir string
var composerCmd = &cobra.Command{
	Use:   "composer",
	Short: "Запускает composer команду в контейнере " + composefiletools.App,
	Args:  cobra.ArbitraryArgs,
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		globaltools.ValidateWorkDir()
		if err := execComposerInContainer(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(composerCmd)
	composerCmd.Flags().StringVarP(&workdir, "workdir", "w", "", "Рабочая директория внутри контейнера")
}

func execComposerInContainer(args []string) error {
	baseArgs := []string{
		"exec", "-it",
		"--user", "docky",
		"-e", "XDEBUG_MODE=off",
	}

	resolvedWorkdir := workdir
	if resolvedWorkdir == "" {
		resolvedWorkdir = os.Getenv("COMPOSER_WORKDIR")
	}
	if resolvedWorkdir != "" {
		baseArgs = append(baseArgs, "-w", resolvedWorkdir)
	}

	baseArgs = append(baseArgs,
		composefiletools.App,
		"composer",
	)
	execArgs := append(baseArgs, args...)

	return globaltools.ExecDockerCompose(execArgs)
}
