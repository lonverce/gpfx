package service

import "testing"

type SomeExternalInstance struct {
	Name string
}

func TestInstanceServiceMaintainer_InjectForInstance(t *testing.T) {
	instance := &SomeExternalInstance{
		Name: "123",
	}

	v := &instanceServiceMaintainer{
		instance: instance,
	}

	resolvedInstance, needInject := v.CreateServiceInstance(nil)

	if resolvedInstance != instance {
		t.Error("Should return a same instance")
	}

	if needInject {
		t.Error("external instance never need to be injected")
	}

	v.Clear()
	v.ReUse()

	resolvedInstance2, _ := v.CreateServiceInstance(nil)
	if resolvedInstance2 != instance {
		t.Error("Should return a same instance")
	}

	v2 := v.Fork(&defaultProvider{})

	resolvedInstance3, _ := v2.CreateServiceInstance(nil)

	if resolvedInstance3 != instance {
		t.Error("Should return a same instance")
	}

	{
		defer func() {
			err := recover()
			if err == nil {
				t.Error("instanceServiceMaintainer should panic when invoked for inject")
			}
		}()

		v2.InjectForInstance(instance, nil)
	}
}
