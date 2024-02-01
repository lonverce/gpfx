package service

import (
	"reflect"
)

type resolvingNode struct {
	Target   maintainer
	Created  bool
	Instance any
}

type interimProvider struct {
	owner        *defaultProvider
	instanceList []*resolvingNode
}

func (p *interimProvider) MustGet(srvType reflect.Type) any {
	target := p.owner.getMaintainerIdByServiceType(srvType)
	maintainer := p.owner.maintainers[target]
	return p.loadServiceByMaintainer(maintainer)
}

func (p *interimProvider) GetOwner() Provider {
	return p.owner
}

func (p *interimProvider) loadServiceByMaintainer(maintainer maintainer) any {

	var currentNode *resolvingNode = nil

	for _, node := range p.instanceList {
		if node.Target == maintainer {
			if !node.Created {
				panic("Found circle-dependency")
			}
			currentNode = node
		}
	}

	if currentNode == nil {
		currentNode = &resolvingNode{
			Target:  maintainer,
			Created: false,
		}
		p.instanceList = append(p.instanceList, currentNode)
	}

	newInstance, needInject := maintainer.CreateServiceInstance(p)
	currentNode.Instance = newInstance
	currentNode.Created = true

	if needInject {
		maintainer.InjectForInstance(newInstance, p)
	}

	return newInstance
}
