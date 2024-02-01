package service

import (
	"fmt"
	"reflect"
)

type baseInjector struct {
	FieldNum int
	Handler  func(field reflect.Value, provider InterimProvider)
}

func injectorWithLoad(field reflect.Value, provider InterimProvider) {
	field.Set(reflect.ValueOf(provider.MustGet(field.Type())))
}

func injectorWithOwner(field reflect.Value, provider InterimProvider) {
	field.Set(reflect.ValueOf(provider.GetOwner()))
}

func injectorWithMakeMap(field reflect.Value, provider InterimProvider) {
	field.Set(reflect.MakeMap(field.Type()))
}

func injectorWithMakeSlice(field reflect.Value, provider InterimProvider) {
	field.Set(reflect.MakeSlice(field.Type(), 0, 0))
}

func injectorWithMakeChan(field reflect.Value, provider InterimProvider) {
	field.Set(reflect.MakeChan(field.Type(), 0))
}

func CreateInjectorByReflection(instanceType reflect.Type) Injector {
	fieldCnt := instanceType.NumField()

	var fieldInjectors []*baseInjector

	for i := 0; i < fieldCnt; i++ {
		field := instanceType.Field(i)
		injectTag, ok := field.Tag.Lookup("gpfx.inject")

		if !ok {
			continue
		}

		if !field.IsExported() {
			panic(fmt.Sprintf("The field '%s.%s' has 'gpfx.inject' tag but is not exported", instanceType.String(), field.Name))
		}

		injector := &baseInjector{
			FieldNum: i,
		}
		switch injectTag {
		case "":
			if field.Type == Typeof[Provider]() {
				injector.Handler = injectorWithOwner
			} else {
				injector.Handler = injectorWithLoad
			}
		case "make":
			switch field.Type.Kind() {
			case reflect.Map:
				injector.Handler = injectorWithMakeMap
			case reflect.Slice:
				injector.Handler = injectorWithMakeSlice
			case reflect.Chan:
				injector.Handler = injectorWithMakeChan
			default:
				panic(fmt.Sprintf("The field '%s.%s' is not supported by gpfx.inject to make", instanceType.String(), field.Name))
			}
		}
		fieldInjectors = append(fieldInjectors, injector)
	}

	if len(fieldInjectors) == 0 {
		return nil
	}

	return func(instance any, provider InterimProvider) {
		val := reflect.ValueOf(instance).Elem()

		for _, injector := range fieldInjectors {
			field := val.Field(injector.FieldNum)
			injector.Handler(field, provider)
		}
	}
}
