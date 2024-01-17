package gpfx

// Option 是针对配置项的抽象
type Option[TOption any] interface {
	// Value 获取配置项值
	Value() *TOption
}

type optionBuilder interface {
	PublishToRegistry(registry ServiceRegistry)
}

type typedOptionBuilder[TOption any] struct {
	actions []func(option *TOption)
}

type defaultOption[TOption any] struct {
	actions []func(option *TOption)
}

func (d *defaultOption[TOption]) Value() *TOption {
	v := new(TOption)
	for _, action := range d.actions {
		action(v)
	}
	return v
}

func (d *typedOptionBuilder[TOption]) PublishToRegistry(registry ServiceRegistry) {
	actions := d.actions[:]
	registry.AddService(RegistrationItem{
		Lifetime: Transient,
		Constructor: func() any {
			opt := &defaultOption[TOption]{
				actions: actions,
			}
			return opt
		},
	}, Typeof[Option[TOption]]())
}
