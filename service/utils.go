package service

import (
	"reflect"
)

func Typeof[T any]() reflect.Type {
	var p *T
	return reflect.TypeOf(p).Elem()
}

func defaultConstructor[TImplement any](InterimProvider) any {
	return new(TImplement)
}

func Add[TImplement any](services Registry, lifetime Lifetime, types ...reflect.Type) {
	if lifetime == ExternalInstance {
		panic("Call AddInstanceOnly[TImplement]")
	}

	t := Typeof[*TImplement]()
	cType := Typeof[TImplement]()

	if cType.Kind() != reflect.Struct {
		panic("TImplement 必须为结构体类型")
	}

	if len(types) == 0 {
		types = []reflect.Type{t}
	}

	defaultCtor := defaultConstructor[TImplement]

	services.AddService(Registration{
		Constructor: defaultCtor,
		Injector:    CreateInjectorByReflection(cType),
		Lifetime:    lifetime,
	}, types...)
}

func AddTransient[TImplement any](services Registry, types ...reflect.Type) {
	Add[TImplement](services, Transient, types...)
}

func AddSingleton[TImplement any](services Registry, types ...reflect.Type) {
	Add[TImplement](services, Singleton, types...)
}
func AddScoped[TImplement any](services Registry, types ...reflect.Type) {
	Add[TImplement](services, Scoped, types...)
}

func AddInstanceOnly[TImplement any](services Registry, instance TImplement) {
	services.AddService(Registration{
		Lifetime: ExternalInstance,
		Instance: instance,
	}, Typeof[TImplement]())
}
