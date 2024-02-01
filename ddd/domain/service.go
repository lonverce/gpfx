package domain

import (
	"github.com/lonverce/gpfx"
	"github.com/lonverce/gpfx/clock"
	"github.com/lonverce/gpfx/service"
)

// Service 领域服务
type Service struct {
	LazyProvider gpfx.LazyServiceProvider `gpfx.inject:""`
}

func (srv *Service) ClockService() clock.Service {
	return service.MustGet[clock.Service](srv.LazyProvider)
}
