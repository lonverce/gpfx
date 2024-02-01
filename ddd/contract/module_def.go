package contract

import (
	"github.com/lonverce/gpfx"
	"github.com/lonverce/gpfx/ddd/shared"
)

var ModuleDef = &gpfx.Module{
	Name:     "gpfx.ddd.contract",
	DependOn: []*gpfx.Module{shared.ModuleDef},
}
