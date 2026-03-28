package initialization

import (
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config/framework"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/globaltools"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/readertools"
)

// frameworkConfigurator заполняет YamlConfig специфичными для фреймворка параметрами.
type frameworkConfigurator interface {
	configure(cfg *config.YamlConfig) error
}

// defaultConfigurator используется для фреймворков без явной записи в реестре.
type defaultConfigurator struct{}
type laravelConfigurator struct{}
type vanillaConfigurator struct{}
type symfonyConfigurator struct{}
type bitrixNuxtConfigurator struct{}
type yii2Configurator struct{}
type yii3Configurator struct{}

// frameworkConfigurators — реестр конфигураторов по фреймворку.
// Добавление нового фреймворка не требует изменения существующего кода.
var frameworkConfigurators = map[framework.Framework]frameworkConfigurator{
	framework.Laravel:    &laravelConfigurator{},
	framework.Vanilla:    &vanillaConfigurator{},
	framework.Symfony:    &symfonyConfigurator{},
	framework.BitrixNuxt: &bitrixNuxtConfigurator{},
	framework.Yii2:       &yii2Configurator{},
	framework.Yii3:       &yii3Configurator{},
}

func (d *defaultConfigurator) configure(cfg *config.YamlConfig) error {
	cfg.DbType = composefiletools.Mysql
	if cfg.MysqlVersion == "" {
		cfg.MysqlVersion = readertools.GetOrChoose(
			"Выберите версию mysql: ", cfg.MysqlVersion,
			composefiletools.GetAvailableVersions(composefiletools.Mysql, cfg),
		)
	}
	chooseNode(cfg)
	cfg.CreateSphinx = readertools.AskYesNo("Добавлять sphinx?")
	return nil
}

func (l *laravelConfigurator) configure(cfg *config.YamlConfig) error {
	if _, err := globaltools.IsDockerComposeAvailable(); err != nil {
		return err
	}
	chooseDbAndCache(cfg)
	cfg.CreateNode = true
	globaltools.InitNode(cfg)
	return nil
}

func (v *vanillaConfigurator) configure(cfg *config.YamlConfig) error {
	chooseDbAndCache(cfg)
	chooseNode(cfg)
	return nil
}

func (s *symfonyConfigurator) configure(cfg *config.YamlConfig) error {
	chooseDbAndCache(cfg)
	return nil
}

func (b *bitrixNuxtConfigurator) configure(cfg *config.YamlConfig) error {
	cfg.DbType = composefiletools.Mysql
	if cfg.MysqlVersion == "" {
		cfg.MysqlVersion = readertools.GetOrChoose(
			"Выберите версию mysql: ", cfg.MysqlVersion,
			composefiletools.GetAvailableVersions(composefiletools.Mysql, cfg),
		)
	}
	cfg.CreateNode = true
	cfg.NodePath = "/var/www/nuxt"
	globaltools.InitNode(cfg)
	cfg.CreateSphinx = readertools.AskYesNo("Добавлять sphinx?")
	return nil
}

func (y *yii2Configurator) configure(cfg *config.YamlConfig) error {
	chooseDbAndCache(cfg)
	chooseNode(cfg)
	return nil
}

func (y *yii3Configurator) configure(cfg *config.YamlConfig) error {
	chooseDbAndCache(cfg)
	chooseNode(cfg)
	return nil
}

func chooseDbAndCache(cfg *config.YamlConfig) {
	cfg.DbType = readertools.GetOrChoose("Выберите базу данных: ", "", composefiletools.AvailableDb[:])
	switch cfg.DbType {
	case composefiletools.Mysql:
		cfg.MysqlVersion = readertools.GetOrChoose(
			"Выберите версию mysql: ", cfg.MysqlVersion,
			composefiletools.GetAvailableVersions(composefiletools.Mysql, cfg),
		)
	case composefiletools.Mariadb:
		cfg.MariadbVersion = readertools.GetOrChoose(
			"Выберите версию mariadb: ", cfg.MariadbVersion,
			composefiletools.GetAvailableVersions(composefiletools.Mariadb, cfg),
		)
	case composefiletools.Postgres:
		cfg.PostgresVersion = readertools.GetOrChoose(
			"Выберите версию postgres: ", cfg.PostgresVersion,
			composefiletools.GetAvailableVersions(composefiletools.Postgres, cfg),
		)
	}

	cache := readertools.GetOrChoose(
		"Выберите сервер кеширования: ", "",
		append(composefiletools.AvailableServerCache[:], "Пропуск"),
	)
	if cache != "Пропуск" {
		cfg.ServerCache = cache
	}
}

func chooseNode(cfg *config.YamlConfig) {
	if readertools.AskYesNo("Добавлять node js?") {
		cfg.CreateNode = true
		globaltools.InitNode(cfg)
	}
}