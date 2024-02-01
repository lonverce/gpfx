package clock

import (
	"github.com/lonverce/gpfx/config"
	"time"
)

type Service interface {
	Now() time.Time
}

type Option struct {
	UseUTC bool
}

type DefaultService struct {
	Option config.Option[Option] `gpfx.inject:""`
}

func (srv *DefaultService) Now() time.Time {
	now := time.Now()
	if srv.Option.OnceValue().UseUTC {
		now = now.UTC()
	}
	return now
}
