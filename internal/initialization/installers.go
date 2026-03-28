package initialization

import (
	"path/filepath"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/readertools"
)

func installLaravelProject() error {
	siteDir := config.GetSiteDirPath()
	dir := "laravel"

	args := []string{
		"run", "--rm",
		"--user", "docky", "--entrypoint", "php",
		composefiletools.App, "/home/docky/.config/composer/vendor/bin/laravel", "new", dir,
	}
	if err := globaltools.ExecDockerCompose(args); err != nil {
		return err
	}

	return moveDirIfExists(filepath.Join(siteDir, dir), siteDir)
}

func installSymfonyProject() error {
	isCli := readertools.AskYesNo("Вы создаете консольное приложение Symfony?")

	args := []string{
		"run", "--rm",
		"--user", "docky",
		"-e", "XDEBUG_MODE=off",
		"--entrypoint", "composer",
		composefiletools.App,
		"create-project", "symfony/skeleton", ".",
	}
	if err := globaltools.ExecDockerCompose(args); err != nil {
		return err
	}

	if isCli {
		return nil
	}

	return globaltools.ExecDockerCompose([]string{
		"run", "--rm",
		"--user", "docky", "-e", "XDEBUG_MODE=off", "--entrypoint", "composer",
		composefiletools.App, "require", "webapp",
	})
}

func installYii2Project(cfg *config.YamlConfig) error {
	template := "yiisoft/yii2-app-basic"
	if cfg.Yii2Advanced {
		template = "yiisoft/yii2-app-advanced"
	}

	dir := "yii2"
	args := []string{
		"run", "--rm",
		"--user", "docky",
		"-e", "XDEBUG_MODE=off",
		"--entrypoint", "composer",
		composefiletools.App,
		"create-project",
		"--prefer-dist",
		"--no-install",
		template,
		dir,
	}
	if err := globaltools.ExecDockerCompose(args); err != nil {
		return err
	}

	return moveDirIfExists(filepath.Join(config.GetSiteDirPath(), dir), config.GetSiteDirPath())
}

func installYii3Project() error {
	dir := "yii3"
	args := []string{
		"run", "--rm",
		"--user", "docky",
		"-e", "XDEBUG_MODE=off",
		"--entrypoint", "composer",
		composefiletools.App,
		"create-project",
		"yiisoft/app",
		dir,
	}
	if err := globaltools.ExecDockerCompose(args); err != nil {
		return err
	}

	return moveDirIfExists(filepath.Join(config.GetSiteDirPath(), dir), config.GetSiteDirPath())
}