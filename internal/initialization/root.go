package initialization

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config/framework"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/filetools"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/readertools"
)

func InitDockerComposeFile() error {
	if err := handleExistingComposeFile(); err != nil {
		return err
	}

	yamlConfig := config.GetYamlConfig()
	yamlConfig.FrameworkName = framework.ParseFramework(readertools.GetOrChoose("Ваш фреймворк: ", "", framework.GetAllStrings()))
	if yamlConfig.FrameworkName == framework.Yii2 {
		yamlConfig.Yii2Advanced = readertools.AskYesNo("Использовать advanced шаблон Yii2? (Если нет, будет установлен basic шаблон)")
	}
	yamlConfig.PhpVersion = readertools.GetOrChoose("Выберите версию php: ", "", composefiletools.GetAvailableVersions(composefiletools.App, yamlConfig))

	switch yamlConfig.FrameworkName {
	case framework.Laravel:
		if err := initLaravelConfig(yamlConfig); err != nil {
			return err
		}
	case framework.Vanilla:
		initVanillaConfig(yamlConfig)
	case framework.Symfony:
		initSymfonyConfig(yamlConfig)
	case framework.BitrixNuxt:
		initBitrixNuxtConfig(yamlConfig)
	case framework.Yii2:
		initYii2Config(yamlConfig)
	default:
		initDefaultConfig(yamlConfig)
	}

	if err := globaltools.InitEnvFile(yamlConfig); err != nil {
		return err
	}

	return composefiletools.BuildYaml(yamlConfig).Save(config.GetDockerComposeFilePath())
}

func InitLaravel() error {
	isInstall := readertools.AskYesNo("Устанавливать Laravel? (Если нет, будет создан пустой проект для ручной настройки)")
	if !isInstall {
		return nil
	}

	siteDir := config.GetSiteDirPath()

	if !filetools.IsDirEmpty(siteDir) {
		if !readertools.AskYesNo("Директория с сайтом не пуста. Удалить всё и установить Laravel?") {
			return nil
		}
		if err := recreateDir(siteDir); err != nil {
			return err
		}
	}

	if err := globaltools.ExecDockerCompose([]string{"build", composefiletools.App}); err != nil {
		return err
	}

	if err := installLaravelProject(); err != nil {
		return err
	}

	if err := setupNodePackages(siteDir); err != nil {
		return err
	}

	globaltools.DownContainers()
	return nil
}

func InitSymfony() error {
	isInstall := readertools.AskYesNo("Устанавливать Symfony? (Если нет, будет создан пустой проект для ручной настройки)")
	if !isInstall {
		return nil
	}

	siteDir := config.GetSiteDirPath()

	if !filetools.IsDirEmpty(siteDir) {
		if !readertools.AskYesNo("Директория с сайтом не пуста. Удалить всё и установить Symfony?") {
			return nil
		}
		if err := recreateDir(siteDir); err != nil {
			return err
		}
	}

	if err := globaltools.ExecDockerCompose([]string{"build", composefiletools.App}); err != nil {
		return err
	}

	if err := installSymfonyProject(); err != nil {
		return err
	}

	globaltools.DownContainers()
	return nil
}

func initBitrixNuxtConfig(yamlConfig *config.YamlConfig) {
	yamlConfig.DbType = composefiletools.Mysql
	if yamlConfig.MysqlVersion == "" {
		yamlConfig.MysqlVersion = readertools.GetOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, composefiletools.GetAvailableVersions(composefiletools.Mysql, yamlConfig))
	}

	yamlConfig.CreateNode = true
	yamlConfig.NodePath = "/var/www/nuxt"
	globaltools.InitNode(yamlConfig)
	yamlConfig.CreateSphinx = readertools.AskYesNo("Добавлять sphinx?")
}

func initYii2Config(yamlConfig *config.YamlConfig) {
	chooseDbAndCache(yamlConfig)
	chooseNode(yamlConfig)
}

func handleExistingComposeFile() error {
	composeFilePath := config.GetDockerComposeFilePath()
	if exists, _ := filetools.FileIsExists(composeFilePath); !exists {
		return nil
	}

	if !readertools.AskYesNo("Файл docker-compose.yml уже существует, создать новый?") {
		return nil
	}
	return os.Rename(composeFilePath, composeFilePath+config.GetTimeStamp())
}

func chooseDbAndCache(yamlConfig *config.YamlConfig) {
	yamlConfig.DbType = readertools.GetOrChoose("Выберите базу данных: ", "", composefiletools.AvailableDb[:])
	switch yamlConfig.DbType {
	case composefiletools.Mysql:
		yamlConfig.MysqlVersion = readertools.GetOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, composefiletools.GetAvailableVersions(composefiletools.Mysql, yamlConfig))
	case composefiletools.Mariadb:
		yamlConfig.MariadbVersion = readertools.GetOrChoose("Выберите версию mariadb: ", yamlConfig.MariadbVersion, composefiletools.GetAvailableVersions(composefiletools.Mariadb, yamlConfig))
	case composefiletools.Postgres:
		yamlConfig.PostgresVersion = readertools.GetOrChoose("Выберите версию postgres: ", yamlConfig.PostgresVersion, composefiletools.GetAvailableVersions(composefiletools.Postgres, yamlConfig))
	}

	cache := readertools.GetOrChoose("Выберите сервер кеширования: ", "", append(composefiletools.AvailableServerCache[:], "Пропуск"))
	if cache != "Пропуск" {
		yamlConfig.ServerCache = cache
	}
}

func chooseNode(yamlConfig *config.YamlConfig) {
	if readertools.AskYesNo("Добавлять node js?") {
		yamlConfig.CreateNode = true
		globaltools.InitNode(yamlConfig)
	}
}

func initLaravelConfig(yamlConfig *config.YamlConfig) error {
	if _, err := globaltools.IsDockerComposeAvailable(); err != nil {
		return err
	}

	chooseDbAndCache(yamlConfig)
	yamlConfig.CreateNode = true
	globaltools.InitNode(yamlConfig)
	return nil
}

func initVanillaConfig(yamlConfig *config.YamlConfig) {
	chooseDbAndCache(yamlConfig)
	chooseNode(yamlConfig)
}

func initSymfonyConfig(yamlConfig *config.YamlConfig) {
	chooseDbAndCache(yamlConfig)
}

func initDefaultConfig(yamlConfig *config.YamlConfig) {
	yamlConfig.DbType = composefiletools.Mysql
	if yamlConfig.MysqlVersion == "" {
		yamlConfig.MysqlVersion = readertools.GetOrChoose("Выберите версию mysql: ", yamlConfig.MysqlVersion, composefiletools.GetAvailableVersions(composefiletools.Mysql, yamlConfig))
	}

	chooseNode(yamlConfig)
	yamlConfig.CreateSphinx = readertools.AskYesNo("Добавлять sphinx?")
}

func recreateDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("не удалось очистить директорию: %w", err)
	}
	return filetools.InitDirs(dir)
}

func InitYii2(yamlConfig *config.YamlConfig) error {
	isInstall := readertools.AskYesNo("Устанавливать Yii2? (Если нет, будет создан пустой проект для ручной настройки)")
	if !isInstall {
		return nil
	}
	siteDir := config.GetSiteDirPath()

	if !filetools.IsDirEmpty(siteDir) {
		if !readertools.AskYesNo("Директория с сайтом не пуста. Удалить всё и установить Yii2?") {
			return nil
		}
		if err := recreateDir(siteDir); err != nil {
			return err
		}
	}

	if err := globaltools.ExecDockerCompose([]string{"build", composefiletools.App}); err != nil {
		return err
	}
	if err := installYii2Project(yamlConfig); err != nil {
		return err
	}

	if err := setupNodePackages(siteDir); err != nil {
		return err
	}
	return nil
}

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

	newPath := filepath.Join(siteDir, dir)
	if exists, _ := filetools.FileIsExists(newPath); exists {
		return filetools.MoveDirContents(newPath, siteDir)
	}
	return nil
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

	if !isCli {
		if err := globaltools.ExecDockerCompose([]string{
			"run", "--rm",
			"--user", "docky", "-e", "XDEBUG_MODE=off", "--entrypoint", "composer",
			composefiletools.App, "require", "webapp",
		}); err != nil {
			return err
		}
	}

	return nil
}

func installYii2Project(yamlConfig *config.YamlConfig) error {
	template := "yiisoft/yii2-app-basic"
	if yamlConfig.Yii2Advanced {
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

	siteDir := config.GetSiteDirPath()
	newPath := filepath.Join(siteDir, dir)
	if exists, _ := filetools.FileIsExists(newPath); exists {
		return filetools.MoveDirContents(newPath, siteDir)
	}
	globaltools.DownContainers()
	return nil
}

func setupNodePackages(siteDir string) error {
	if exists, _ := filetools.FileIsExists(filepath.Join(siteDir, "package.json")); !exists {
		return nil
	}

	if err := globaltools.ExecDockerCompose([]string{"build", composefiletools.Node}); err != nil {
		return err
	}

	return globaltools.ExecDockerCompose([]string{
		"run", "--rm",
		"--user", "docky", "--entrypoint", "npm",
		composefiletools.Node, "install",
	})
}
