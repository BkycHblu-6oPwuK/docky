package composefiletools

import (
	"fmt"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config/framework"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile/service"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile/volume"
)

// функция для публикации изменений в compose файле. Загружает текущий compose, передает его билдеру, который изменяет его, затем билдит и сохраняет обратно
func publishWithBuilder(modifier func(builder *composefile.ComposeFileBuilder) error) error {
	path := config.GetDockerComposeFilePath()
	compose, err := composefile.Load(path)
	if err != nil {
		return err
	}

	builder := composefile.NewComposeFileBuilderFrom(*compose)
	if err := modifier(builder); err != nil {
		return err
	}

	final := builder.Build()
	composefile.SetCurrentYaml(&final)
	return final.Save(path)
}

// функции публикации для разных сервисов. Внутри используют publishWithBuilder, передавая ему нужные изменения. Например, публикация mysql сервиса удаляет mariadb и postgres, если они есть, и добавляет mysql, если его нет
func publishDatabaseService(target string, builderFunc func() service.Service) error {
	alternatives := []string{Mysql, Mariadb, Postgres}
	removeVolumes := map[string]string{
		Mysql:    Mysql_data,
		Mariadb:  Mariadb_data,
		Postgres: Postgres_data,
	}

	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if !b.HasService(target) {
			if appService, exists := b.GetService(App); exists {
				serviceBuilder := service.NewServiceBuilderFrom(appService)
				deps := serviceBuilder.GetDependenciesBuilder()
				for _, alt := range alternatives {
					if alt != target && b.HasService(alt) {
						deps.RewriteDependency(alt, target)
					}
				}
				b.AddService(App, serviceBuilder.Build())
			}
			b.AddService(target, builderFunc()).
				AddVolume(removeVolumes[target], volume.Volume{})
		}

		for _, alt := range alternatives {
			if alt != target && b.HasService(alt) {
				b.RemoveService(alt).
					RemoveVolume(removeVolumes[alt])
			}
		}
		return nil
	})
}

// функции публикации для каждого сервиса. Вызывают publishDatabaseService с нужным билдером, или publishWithBuilder с нужными изменениями
func PublishMysqlService() error {
	return publishDatabaseService(Mysql, buildMysqlService)
}

// Публикует mariadb сервис. Удаляет mysql и postgres, если они есть, и добавляет mariadb, если его нет
func PublishMariadbService() error {
	return publishDatabaseService(Mariadb, buildMariadbService)
}

// Публикует postgres сервис. Удаляет mysql и mariadb, если они есть, и добавляет postgres, если его нет
func PublishPostgresService() error {
	return publishDatabaseService(Postgres, buildPostgresService)
}

// Публикует node сервис. Добавляет node, если его нет
func PublishNodeService() error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if !b.HasService(Node) {
			b.AddService(Node, buildNodeService())
		}
		return nil
	})
}

// Публикует sphinx сервис. Добавляет sphinx, если его нет
func PublishSphinxService() error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if !b.HasService(Sphinx) {
			b.AddService(Sphinx, buildSphinxService()).
				AddVolume(Sphinx_data, volume.Volume{})
		}
		return nil
	})
}

// Публикует redis сервис. Добавляет redis, если его нет
func PublishRedisService() error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if !b.HasService(Redis) {
			b.AddService(Redis, buildRedisService()).
				AddVolume(Redis_data, volume.Volume{})
		}
		return nil
	})
}

// Публикует memcached сервис. Добавляет memcached, если его нет
func PublishMemcachedService() error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if !b.HasService(Memcached) {
			b.AddService(Memcached, buildMemcachedService())
		}
		return nil
	})
}

// Публикует mailhog сервис. Добавляет mailhog, если его нет
func PublishMailhogService() error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if !b.HasService(Mailhog) {
			b.AddService(Mailhog, buildMailHogService())
		}
		return nil
	})
}

// Публикует phpmyadmin сервис. Добавляет phpmyadmin, если его нет, и связывает его с mysql или mariadb, если они есть. Если mysql и mariadb нет, возвращает ошибку
func PublishPhpMyAdminService() error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		publish := func(host string) {
			if !b.HasService(PhpMyAdmin) {
				b.AddService(PhpMyAdmin, buildPhpMyAdminService(host))
			}
		}
		if b.HasService(Mysql) {
			publish(Mysql)
			return nil
		} else if b.HasService(Mariadb) {
			publish(Mariadb)
			return nil
		}
		return fmt.Errorf("phpmyadmin работает только с mysql. В docker-compose не найден сервис %s или %s", Mysql, Mariadb)
	})
}

func PublishNginxYii2BackendService(yamlConfig *config.YamlConfig) error {
	if yamlConfig.FrameworkName != framework.Yii2 {
		return fmt.Errorf("nginx yii2 backend сервис может быть опубликован только для проектов на yii2. Текущий фреймворк: %s", yamlConfig.FrameworkName)
	}
	if !yamlConfig.Yii2Advanced {
		return fmt.Errorf("nginx yii2 backend сервис может быть опубликован только для yii2 advanced. Текущий проект не является yii2 advanced")
	}
	yamlConfig.Yii2Backend = true
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if !b.HasService(NginxBackend) {
			b.AddService(NginxBackend, buildNginxService(yamlConfig, map[string]string{"8000": "80", "4430": "443"}))
		}
		return nil
	})
}

// volumes map serviceName>>[]string volumes
// modifier функция, которая принимает билдер сервиса и может его изменить. Если возвращает isContinue=false, то к этому сервису не будут применены тома из volumes
// Публикует тома для сервисов. Для каждого сервиса из volumes, если он есть в compose, вызывает modifier для его билдера, и если modifier возвращает isContinue=true, то добавляет тома из volumes к этому сервису
func PublishVolumes(volumes map[string][]string, modifier func(b *service.ServiceBuilder) (isContinue bool, err error)) error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		for serviceName, volumes := range volumes {
			if curService, exists := b.GetService(serviceName); exists {
				serviceBuilder := service.NewServiceBuilderFrom(curService)
				if modifier != nil {
					isContinue, err := modifier(serviceBuilder)
					if err != nil {
						return fmt.Errorf("ошибка при модификации сервиса %s: %w", serviceName, err)
					}
					if !isContinue {
						continue
					}
				}
				for _, vol := range volumes {
					serviceBuilder.SetVolume(vol)
				}
				b.AddService(serviceName, serviceBuilder.Build())
			}
		}
		return nil
	})
}

// Публикует dockerfile для сервиса. Если сервиса нет, возвращает ошибку. Если сервис есть, заменяет его dockerfile на переданный
func PublishDockerfile(serviceName, dockerfile string) error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		config := config.GetYamlConfig()
		services := []string{serviceName}
		if serviceName == Nginx && config.FrameworkName == framework.Yii2 && config.Yii2Advanced {
			services = append(services, NginxBackend)
		}
		addServiceHandler := func(serviceName string) error {
			if curService, exists := b.GetService(serviceName); exists {
				if curService.Build.Dockerfile != "" {
					curService.Build.Dockerfile = dockerfile
					b.AddService(serviceName, curService)
				}
			} else {
				return fmt.Errorf("сервис %s не найден", serviceName)
			}
			return nil
		}
		for _, serviceName := range services {
			if err := addServiceHandler(serviceName); err != nil {
				return err
			}
		}
		return nil
	})
}
