package gpfx

import (
	"reflect"
)

func Typeof[T any]() reflect.Type {
	var p *T = nil
	return reflect.TypeOf(p).Elem()
}

type AutoInjectedObject interface {
	Inject(provider InterimServiceContext)
}

func defaultConstructor[TImplement any]() any {
	return new(TImplement)
}

func defaultInject(instance any, provider InterimServiceContext) {
	instance.(AutoInjectedObject).Inject(provider)
}

func AddService[TImplement any](services ServiceRegistry, lifetime ServiceLifetime, types ...reflect.Type) {
	t := Typeof[*TImplement]()
	defaultCtor := defaultConstructor[TImplement]

	var injector ServiceInjector = nil
	if t.Implements(Typeof[AutoInjectedObject]()) {
		injector = defaultInject
	}

	var destructor ServiceDestructor = nil
	if lifetime == Scoped && t.Implements(Typeof[ISupportScopeReuse]()) {
		destructor = func(v any) {
			v.(ISupportScopeReuse).HandleClearBeforeReuse()
		}
	}

	services.AddService(RegistrationItem{
		Constructor:        defaultCtor,
		Injector:           injector,
		ScopedReuseHandler: destructor,
		Lifetime:           lifetime,
	}, types...)
}

func AddTransient[TImplement any](services ServiceRegistry, types ...reflect.Type) {
	AddService[TImplement](services, Transient, types...)
}
func AddSingleton[TImplement any](services ServiceRegistry, types ...reflect.Type) {
	AddService[TImplement](services, Singleton, types...)
}
func AddScoped[TImplement any](services ServiceRegistry, types ...reflect.Type) {
	AddService[TImplement](services, Scoped, types...)
}

func AddInstanceOnly[TImplement any](services ServiceRegistry, instance TImplement) {
	services.AddService(RegistrationItem{
		UseInstance: true,
		Instance:    instance,
	}, Typeof[TImplement]())
}
