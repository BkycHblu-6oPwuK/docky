package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"

	"github.com/spf13/cobra"
)

var shareCmd = &cobra.Command{
	Use:                "share",
	Short:              "Туннелирование через Cloudpub",
	Long:               "Создаёт публичный URL через Cloudpub",
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
    token := "Uy2ZPTdV_Gu_2yrHh_vmNZ7meW3TV5YY0x_bxGW7rJo"
    networkName := getDockyNetwork()

    cmdArgs := []string{
        "run", "--rm",
        "--network", networkName,
        "-e", "TOKEN=" + token,
        "cloudpub/cloudpub:latest",
        "publish", "http", "nginx:80",
    }

    cmdArgs = append(cmdArgs, args...)

    cmd := exec.Command("docker", cmdArgs...)
    cmd.Dir = config.GetWorkDirPath()
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    return cmd.Run()
}

func getDockyNetwork() string {
    workDir := config.GetWorkDirPath()
    projectName := strings.ToLower(filepath.Base(workDir))
    return projectName + "_docky"
}