package initialization

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/filetools"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/readertools"
)

// ensureEmptyDir проверяет, пуста ли директория.
// Если нет — спрашивает пользователя, очищать ли её.
// Возвращает false, если установку нужно прервать.
func ensureEmptyDir(dir string) (bool, error) {
	if filetools.IsDirEmpty(dir) {
		return true, nil
	}
	if !readertools.AskYesNo("Директория с сайтом не пуста. Удалить всё и продолжить?") {
		return false, nil
	}
	return true, recreateDir(dir)
}

// удаляет директорию и создаёт её заново.
func recreateDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("не удалось очистить директорию: %w", err)
	}
	return filetools.InitDirs(dir)
}

// проверяет наличие package.json и, если он есть, устанавливает node-модули внутри контейнера.
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

// проверяет, существует ли docker-compose.yml, и при необходимости переименовывает его.
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

// перемещает содержимое src в dst, если src существует.
func moveDirIfExists(src, dst string) error {
	if exists, _ := filetools.FileIsExists(src); exists {
		return filetools.MoveDirContents(src, dst)
	}
	return nil
}

// применяет конфигуратор для выбранного фреймворка.
// Если фреймворк не зарегистрирован — используется defaultConfigurator.
func configureFramework(cfg *config.YamlConfig) error {
	c, ok := frameworkConfigurators[cfg.FrameworkName]
	if !ok {
		c = &defaultConfigurator{}
	}
	return c.configure(cfg)
}