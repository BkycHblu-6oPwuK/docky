package config

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config/framework"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/dotenv"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/filetools"
)

type YamlConfig struct {
	FrameworkName   framework.Framework // bitrix, laravel и т.д.
	DbType          string              // mysql, postgres, sqlite
	PhpVersion      string
	MysqlVersion    string
	MariadbVersion  string
	PostgresVersion string
	NodeVersion     string
	NodePath        string
	CreateNode      bool
	CreateSphinx    bool
	ServerCache     string // memcached, redis
	UserGroup       string
	Yii2Advanced    bool
	Yii2Backend     bool
}

const (
	UserGroupVarName       string = "USERGROUP"
	DockyFrameworkVarName  string = "DOCKY_FRAMEWORK"
	DockerPathVarName      string = "DOCKER_PATH"
	ConfPathVarName        string = "CONF_PATH"
	PhpVersionVarName      string = "PHP_VERSION"
	MysqlVersionVarName    string = "MYSQL_VERSION"
	MariadbVersionVarName  string = "MARIADB_VERSION"
	PostgresVersionVarName string = "POSTGRES_VERSION"
	NodeVersionVarName     string = "NODE_VERSION"
	SitePathVarName        string = "SITE_PATH"
	NodePathVarName        string = "NODE_PATH"
	Yii2AdvancedVarName    string = "YII2_ADVANCED"
)

var (
	cfg  *YamlConfig
	once sync.Once
)

func GetYamlConfig() *YamlConfig {
	once.Do(func() {
		if err := loadEnvFile(); err != nil {
			panic("Ошибка при загрузке файла .env")
		}
		var frameworkValue framework.Framework
		yii2AdvancedRawValue := os.Getenv(Yii2AdvancedVarName)
		yii2Advanced := yii2AdvancedRawValue == "true" || yii2AdvancedRawValue == "1"
		if frameworkEnv := os.Getenv(DockyFrameworkVarName); frameworkEnv == "" {
			frameworkValue = ""
		} else {
			frameworkValue = framework.ParseFramework(frameworkEnv)
		}
		if !yii2Advanced && frameworkValue == framework.Yii2 {
			backandPath := filepath.Join(GetSiteDirPath(), "backend")
			backandExists, _ := filetools.FileIsExists(backandPath)
			frontendPath := filepath.Join(GetSiteDirPath(), "frontend")
			frontendExists, _ := filetools.FileIsExists(frontendPath)
			if backandExists && frontendExists {
				yii2Advanced = true
			}
		}
		nodeVersion := os.Getenv(NodeVersionVarName)
		if nodeVersion == "" {
			nodeVersion = "24"
		}
		cfg = &YamlConfig{
			FrameworkName:   frameworkValue,
			PhpVersion:      os.Getenv(PhpVersionVarName),
			MysqlVersion:    os.Getenv(MysqlVersionVarName),
			MariadbVersion:  os.Getenv(MariadbVersionVarName),
			PostgresVersion: os.Getenv(PostgresVersionVarName),
			NodeVersion:     nodeVersion,
			NodePath:        os.Getenv(NodePathVarName),
			UserGroup:       os.Getenv(UserGroupVarName),
			Yii2Advanced:    yii2Advanced,
		}
	})
	return cfg
}

func loadEnvFile() error {
	envPath := GetEnvFilePath()
	if fileExists, _ := filetools.FileIsExists(envPath); fileExists {
		return dotenv.Load(envPath)
	}
	return nil
}

func GetUserGroup() string {
	userGroup := GetYamlConfig().UserGroup
	if userGroup == "" {
		userGroup = strconv.Itoa(os.Getegid())
		if userGroup == "0" {
			userGroup = "1000"
		}
	}
	return userGroup
}

func GetCurFramework() framework.Framework {
	curFramework := GetYamlConfig().FrameworkName
	if curFramework == "" {
		curFramework = framework.Bitrix
	}
	return curFramework
}
