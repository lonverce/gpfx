package gpfx

import (
	"fmt"
	"log"
)

type moduleStateEnum int

const (
	blocking moduleStateEnum = iota
	completed
)

// moduleLoader 模块加载器
type moduleLoader struct {
	moduleState   map[*Module]moduleStateEnum
	solvedModules []*Module
}

func (ctx *moduleLoader) Load(current *Module, prefix string) {
	if current == nil {
		return
	}

	log.Printf("[gpfx] - %s%s\n", prefix, current.GetString())
	if current.DependOn == nil || len(current.DependOn) == 0 {
		ctx.moduleState[current] = completed

		ctx.solvedModules = append(ctx.solvedModules, current)
		return
	}

	prefix = "  " + prefix
	ctx.moduleState[current] = blocking
	for _, child := range current.DependOn {
		if child == nil {
			continue
		}
		childState, found := ctx.moduleState[child]
		if found {
			switch childState {
			case blocking:
				panic(fmt.Sprintf("DependModule(%s)发现环依赖", child.GetString()))
			case completed:
				continue
			default:
				panic("unhandled module state")
			}
		}

		ctx.Load(child, prefix)
	}

	ctx.moduleState[current] = completed
	ctx.solvedModules = append(ctx.solvedModules, current)
}
