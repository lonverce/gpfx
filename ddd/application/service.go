package application

import (
	"github.com/lonverce/gpfx"
	"github.com/lonverce/gpfx/authority"
	"github.com/lonverce/gpfx/clock"
	"github.com/lonverce/gpfx/event"
	"github.com/lonverce/gpfx/security"
	"github.com/lonverce/gpfx/service"
	"github.com/lonverce/gpfx/uow"
)

// Service 应用层服务基类
type Service struct {
	LazyProvider gpfx.LazyServiceProvider `gpfx.inject:""`
}

func (srv *Service) ClockService() clock.Service {
	return service.MustGet[clock.Service](srv.LazyProvider)
}

func (srv *Service) UnitOfWorkManager() uow.Manager {
	return service.MustGet[uow.Manager](srv.LazyProvider)
}

func (srv *Service) LocalEventBus() event.LocalEventBus {
	return service.MustGet[event.LocalEventBus](srv.LazyProvider)
}

func (srv *Service) CurrentUser() security.CurrentUser {
	return service.MustGet[security.CurrentUser](srv.LazyProvider)
}

func (srv *Service) AuthorityService() authority.Service {
	return service.MustGet[authority.Service](srv.LazyProvider)
}
