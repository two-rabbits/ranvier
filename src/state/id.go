package state

import (
	"config-manager/src/model"
	"strconv"
	"strings"
)

const IdServiceKey = "IdService"

const IdSeparator = "~"

type IdService interface {
	Name(id model.Id) string
	GlobalId() model.Id
	ClusterId(name string) model.Id
	NamespaceId(name, cluster string) model.Id
	ApplicationId(name, namespace, cluster string) model.Id
	VersionedId(id model.Id, version int) model.Id
}

type idServiceImpl struct {
}

func (i *idServiceImpl) Name(id model.Id) string {
	meta := strings.Split(id.String(), IdSeparator)
	return meta[len(meta)-1]
}

func (i *idServiceImpl) GlobalId() model.Id {
	return GlobalId
}

func (i *idServiceImpl) ClusterId(name string) model.Id {
	return i.toId(name)
}

func (i *idServiceImpl) NamespaceId(name, cluster string) model.Id {
	return i.toId(cluster, name)
}

func (i *idServiceImpl) ApplicationId(name, namespace, cluster string) model.Id {
	return i.toId(cluster, namespace, name)
}

func (i *idServiceImpl) VersionedId(id model.Id, version int) model.Id {
	ver := strconv.Itoa(version)
	idSlice := append([]string{id.String()}, ver)
	return model.Id(strings.Join(idSlice, IdSeparator))
}

func (i *idServiceImpl) toId(names ...string) model.Id {
	return model.Id(strings.Join(names, IdSeparator))
}

func (i *idServiceImpl) ToName(id model.Id) string {
	names := strings.Split(string(id), IdSeparator)
	if _, err := strconv.Atoi(names[len(names)-1]); err != nil {
		return names[len(names)-1]
	} else {
		return names[len(names)-2]
	}
}
