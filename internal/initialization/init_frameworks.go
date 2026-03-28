package initialization

import (
	"fmt"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/readertools"
)

// frameworkInitParams описывает параметры установки конкретного фреймворка.
type frameworkInitParams struct {
	name      string // человекочитаемое название фреймворка для сообщений пользователю
	install   func() error // функция, которая устанавливает фреймворк внутри контейнера
	withNode  bool // нужно ли запускать npm install после установки фреймворка
	downAfter bool // нужно ли останавливать контейнеры после установки (некоторые фреймворки требуют запущенного контейнера для установки, но в целом лучше оставлять их остановленными)
}

func InitLaravel() error {
	return initFramework(frameworkInitParams{
		name:      "Laravel",
		install:   installLaravelProject,
		withNode:  true,
		downAfter: true,
	})
}

func InitSymfony() error {
	return initFramework(frameworkInitParams{
		name:      "Symfony",
		install:   installSymfonyProject,
		withNode:  false,
		downAfter: true,
	})
}

func InitYii2(cfg *config.YamlConfig) error {
	return initFramework(frameworkInitParams{
		name:      "Yii2",
		install:   func() error { return installYii2Project(cfg) },
		withNode:  true,
		downAfter: true,
	})
}

func InitYii3() error {
	return initFramework(frameworkInitParams{
		name:      "Yii3",
		install:   installYii3Project,
		withNode:  true,
		downAfter: true,
	})
}

// initFramework — общий шаблон установки фреймворка:
// подтверждение → очистка директории → сборка контейнера → установка → (опционально) node → down.
func initFramework(p frameworkInitParams) error {
	if !readertools.AskYesNo(
		fmt.Sprintf("Устанавливать %s? (Если нет, будет создан пустой проект для ручной настройки)", p.name),
	) {
		return nil
	}

	siteDir := config.GetSiteDirPath()

	ok, err := ensureEmptyDir(siteDir)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	if err := globaltools.ExecDockerCompose([]string{"build", composefiletools.App}); err != nil {
		return err
	}

	if err := p.install(); err != nil {
		return err
	}

	if p.withNode {
		if err := setupNodePackages(siteDir); err != nil {
			return err
		}
	}

	if p.downAfter {
		globaltools.DownContainers()
	}

	return nil
}