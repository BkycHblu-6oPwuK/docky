package initialization

import (
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config/framework"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/readertools"
)

// InitDockerComposeFile интерактивно создаёт docker-compose.yml.
func InitDockerComposeFile() error {
	if err := handleExistingComposeFile(); err != nil {
		return err
	}

	cfg := config.GetYamlConfig()

	cfg.FrameworkName = framework.ParseFramework(
		readertools.GetOrChoose("Ваш фреймворк: ", "", framework.GetAllStrings()),
	)
	if cfg.FrameworkName == framework.Yii2 {
		cfg.Yii2Advanced = readertools.AskYesNo(
			"Использовать advanced шаблон Yii2? (Если нет, будет установлен basic шаблон)",
		)
	}
	cfg.PhpVersion = readertools.GetOrChoose(
		"Выберите версию php: ", "",
		composefiletools.GetAvailableVersions(composefiletools.App, cfg),
	)

	if err := configureFramework(cfg); err != nil {
		return err
	}

	if err := globaltools.InitEnvFile(cfg); err != nil {
		return err
	}

	return composefiletools.BuildYaml(cfg).Save(config.GetDockerComposeFilePath())
}