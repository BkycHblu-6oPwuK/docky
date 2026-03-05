package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"

	"github.com/spf13/cobra"
)

var shareCmd = &cobra.Command{
	Use:                "share",
	Short:              "Туннелирование локального сайта",
	Long:               "Туннелирование происходит с помощью Expose и вы можете прокидывать все флаги что принимает Expose",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		globaltools.ValidateWorkDir()
		if err := share(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(shareCmd)
}

func share(args []string) error {
	hasFlag := func(name string) bool {
		prefix := "--" + name
		for _, a := range args {
			if a == prefix || strings.HasPrefix(a, prefix+"=") {
				return true
			}
		}
		return false
	}

	cmdArgs := []string{
		"run", "--init", "--rm",
		"-p", "4040:4040",
		"-t",
		"beyondcodegmbh/expose-server:latest",
		"share",
	}

	if !hasFlag("server-host") {
		cmdArgs = append(cmdArgs, "--server-host=laravel-sail.site")
	}

	if !hasFlag("server-port") {
		cmdArgs = append(cmdArgs, "--server-port=8080")
	}

	if !hasFlag("domain") {
		cmdArgs = append(cmdArgs, "--domain=laravel-sail.site")
	}

	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("docker", cmdArgs...)
	cmd.Dir = config.GetWorkDirPath()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
