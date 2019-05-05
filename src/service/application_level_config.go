package service

import (
	"config-manager/src/model"
	"config-manager/src/state"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

const ApplicationLevelConfigServiceKey = "ApplicationLevelConfigService"

type ApplicationLevelConfigService interface {
	Query(cluster string, namespace string, application string, query string) (*model.LevelConfig, error)
	Create(cluster string, namespace string, application string, data []byte) (model.LevelConfig, error)
	Rollback(cluster string, namespace string, application string, version int) (config model.LevelConfig, err error)
	Update(cluster, namespace, application string, data []byte) (config model.LevelConfig, err error)
}

type applicationLevelConfigServiceImpl struct {
	LevelConfigService LevelConfigService `inject:"LevelConfigService"`
	IdService          state.IdService    `inject:"IdService"`
}

func (a *applicationLevelConfigServiceImpl) Update(cluster, namespace, application string, data []byte) (config model.LevelConfig, err error) {
	if err := a.exists(cluster, namespace); err != nil {
		return config, err
	}

	id := a.IdService.ApplicationId(application, namespace, cluster)

	return a.LevelConfigService.Update(model.Application, id, data)
}

func (a *applicationLevelConfigServiceImpl) Query(cluster string, namespace string, application string, query string) (*model.LevelConfig, error) {
	if err := a.exists(cluster, namespace); err != nil {
		return nil, err
	}

	id := a.IdService.ApplicationId(application, namespace, cluster)

	return a.LevelConfigService.Query(model.Application, id, query)
}

func (a *applicationLevelConfigServiceImpl) Create(cluster string, namespace string, application string, data []byte) (config model.LevelConfig, err error) {
	if err = a.exists(cluster, namespace); err != nil {
		return
	}

	id := a.IdService.ApplicationId(application, namespace, cluster)

	return a.LevelConfigService.Create(model.Application, id, data)
}

func (a *applicationLevelConfigServiceImpl) Rollback(cluster string, namespace string, application string, version int) (config model.LevelConfig, err error) {
	if err = a.exists(cluster, namespace); err != nil {
		return
	}

	id := a.IdService.ApplicationId(application, namespace, cluster)

	return a.LevelConfigService.Rollback(model.Application, id, version)
}

func (a *applicationLevelConfigServiceImpl) exists(cluster string, namespace string) error {
	if !a.LevelConfigService.Exists(model.Cluster, a.IdService.ClusterId(cluster)) {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not a valid cluster", cluster))
	}
	if !a.LevelConfigService.Exists(model.Namespace, a.IdService.NamespaceId(namespace, cluster)) {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("%s is not a valid namespace", namespace))
	}
	return nil
}