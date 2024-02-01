package service

import (
	"testing"
)

type TestServiceA interface {
	ADoJob() string
}

type TestServiceB interface {
	BDoJob() string
}

type TestServiceImpl struct {
	id int
}

func (t *TestServiceImpl) ADoJob() string {
	return "TestServiceImpl.ADoJob"
}

func (t *TestServiceImpl) BDoJob() string {
	return "TestServiceImpl.BDoJob"
}

type TestServiceImpl3 struct {
	name string
}

func (t *TestServiceImpl3) BDoJob() string {
	return "TestServiceImpl3.BDoJob"
}

func TestTransientService(t *testing.T) {
	registry := NewRegistry()
	AddTransient[TestServiceImpl](registry, Typeof[*TestServiceImpl]())

	root := registry.Build().GetProvider()

	var srv, srv2 *TestServiceImpl

	MustLoad(root, &srv)

	if srv.ADoJob() != "TestServiceImpl.ADoJob" {
		t.Error("MustLoad Transient Service Failed.")
	}
	MustLoad(root, &srv2)

	if srv == srv2 {
		t.Error("Transient Service Should MustGet Different Instance in different load.")
	}
}

func TestSingletonService(t *testing.T) {
	registry := NewRegistry()
	AddSingleton[TestServiceImpl](registry, Typeof[*TestServiceImpl]())

	root := registry.Build().GetProvider()

	var srv, srv2 *TestServiceImpl

	MustLoad(root, &srv)
	if srv.ADoJob() != "TestServiceImpl.ADoJob" {
		t.Error("MustLoad Transient Service Failed.")
	}

	MustLoad(root, &srv2)
	if srv != srv2 {
		t.Error("Singleton Service Should MustGet Same Instance in different load.")
	}
}

func TestSingletonServiceForkBeforeCreate(t *testing.T) {
	registry := NewRegistry()
	AddSingleton[TestServiceImpl](registry, Typeof[*TestServiceImpl]())

	root := registry.Build().GetProvider()
	child := root.CreateScope()

	instanceFromChild := MustGet[*TestServiceImpl](child.GetProvider())
	instanceFromRoot := MustGet[*TestServiceImpl](root)

	if instanceFromChild != instanceFromRoot {
		t.Error("Singleton instance should be same in different scope")
	}
}

func TestSingletonServiceForkAfterCreate(t *testing.T) {
	registry := NewRegistry()
	AddSingleton[TestServiceImpl](registry, Typeof[*TestServiceImpl]())

	root := registry.Build().GetProvider()
	instanceFromRoot := MustGet[*TestServiceImpl](root)

	child := root.CreateScope()

	instanceFromChild := MustGet[*TestServiceImpl](child.GetProvider())

	if instanceFromChild != instanceFromRoot {
		t.Error("Singleton instance should be same in different scope")
	}
}

func TestScopedService(t *testing.T) {
	registry := NewRegistry()
	AddScoped[TestServiceImpl](registry, Typeof[*TestServiceImpl]())

	root := registry.Build().GetProvider()
	var srv, srv2, srv3, srv4 *TestServiceImpl
	MustLoad(root, &srv)
	if srv.ADoJob() != "TestServiceImpl.ADoJob" {
		t.Error("MustLoad Transient Service Failed.")
	}
	MustLoad(root, &srv2)
	if srv != srv2 {
		t.Error("Scoped Service Should MustGet Same Instance in same scope.")
	}

	childScope := root.CreateScope()
	childProvider := childScope.GetProvider()
	defer childScope.Close()

	MustLoad(childProvider, &srv3)

	if srv == srv3 {
		t.Error("Scoped Service Should MustGet Different Instance in Different scope.")
	}

	MustLoad(childProvider, &srv4)
	if srv4 != srv3 {
		t.Error("Scoped Service Should MustGet Same Instance in Same child scope.")
	}
}

func TestProviderReuse(t *testing.T) {
	registry := NewRegistry()
	AddScoped[TestServiceImpl](registry, Typeof[*TestServiceImpl]())

	root := registry.Build().GetProvider()

	child := root.CreateScope()

	instanceFromChild := MustGet[*TestServiceImpl](child.GetProvider())

	child.Close()

	child2 := root.CreateScope()

	if child2 != child {
		t.Error("Provider instance should be reused")
	}

	instanceFromChild2 := MustGet[*TestServiceImpl](child2.GetProvider())

	if instanceFromChild == instanceFromChild2 {
		t.Error("Reused provider should not maintains same scope-service instance")
	}
}

func TestInstanceService(t *testing.T) {
	registry := NewRegistry()
	srv := TestServiceImpl{
		id: 1,
	}

	AddInstanceOnly[*TestServiceImpl](registry, &srv)

	root := registry.Build().GetProvider()
	var srv2 *TestServiceImpl

	MustLoad(root, &srv2)
	if srv2 != &srv {
		t.Error("Instance Service Should MustGet the same instance as register in any situation")
	}
}

func TestSliceMaintainer_CreateServiceInstance(t *testing.T) {
	registry := NewRegistry()
	AddTransient[TestServiceImpl](registry, Typeof[TestServiceA](), Typeof[TestServiceB]())
	AddTransient[TestServiceImpl3](registry, Typeof[TestServiceB]())

	root := registry.Build().GetProvider()
	srvList := MustGet[[]TestServiceB](root)

	if len(srvList) != 2 {
		t.Error("Should has 2 services")
	}

	if srvList[0].BDoJob() != "TestServiceImpl.BDoJob" {
		t.Error("The first service is supposed to be TestServiceImpl")
	}
	if srvList[1].BDoJob() != "TestServiceImpl3.BDoJob" {
		t.Error("The second service is supposed to be TestServiceImpl3")
	}
}

func TestLazyMaintainer_CreateServiceInstance(t *testing.T) {
	registry := NewRegistry()
	AddTransient[TestServiceImpl](registry, Typeof[TestServiceA](), Typeof[TestServiceB]())
	AddTransient[TestServiceImpl3](registry, Typeof[TestServiceB]())

	root := registry.Build().GetProvider()
	srvList := MustGet[[]func(provider Provider) TestServiceB](root)

	if len(srvList) != 2 {
		t.Error("Should has 2 services")
	}

	if srvList[0](root).BDoJob() != "TestServiceImpl.BDoJob" {
		t.Error("The first service is supposed to be TestServiceImpl")
	}
	if srvList[1](root).BDoJob() != "TestServiceImpl3.BDoJob" {
		t.Error("The second service is supposed to be TestServiceImpl3")
	}
}

func BenchmarkLoadTransient(b *testing.B) {
	registry := NewRegistry()
	AddTransient[TestServiceImpl](registry, Typeof[*TestServiceImpl]())

	root := registry.Build().GetProvider()

	for i := 0; i < b.N; i++ {
		var v *TestServiceImpl
		MustLoad(root, &v)
		v.id = i
	}
}

func BenchmarkLoadTransient_Compare2(b *testing.B) {
	registry := NewRegistry()
	AddTransient[TestServiceImpl](registry, Typeof[*TestServiceImpl]())

	root := registry.Build().GetProvider()

	var creator func(provider Provider) *TestServiceImpl
	MustLoad(root, &creator)

	for i := 0; i < b.N; i++ {
		v := creator(root)
		v.id = i
	}
}

func BenchmarkLoadSingleton(b *testing.B) {
	registry := NewRegistry()
	AddSingleton[TestServiceImpl](registry, Typeof[*TestServiceImpl]())

	root := registry.Build().GetProvider()

	for i := 0; i < b.N; i++ {
		var v *TestServiceImpl
		MustLoad(root, &v)
	}
}

type TestServiceImpl2 struct {
	SrvImpl1 *TestServiceImpl `gpfx.inject:""`
	Provider Provider         `gpfx.inject:""`
	Slice    []int            `gpfx.inject:"make"`
	Map      map[string]int   `gpfx.inject:"make"`
	Chan     chan int         `gpfx.inject:"make"`
}

type TestServiceWithUnexportedField struct {
	srvImpl1 *TestServiceImpl `gpfx.inject:""`
}

type TestServiceWithUnsupportedMakeField struct {
	Id int `gpfx.inject:"make"`
}

func (t *TestServiceImpl2) BDoJob() string {
	return "TestServiceImpl2.BDoJob"
}

func TestOverwriteService(t *testing.T) {
	registry := NewRegistry()

	Add[TestServiceImpl](registry, Transient, Typeof[*TestServiceImpl](), Typeof[TestServiceB]())
	Add[TestServiceImpl2](registry, Transient, Typeof[*TestServiceImpl2](), Typeof[TestServiceB]())

	root := registry.Build().GetProvider()
	var srv TestServiceB
	MustLoad(root, &srv)

	if srv.BDoJob() != "TestServiceImpl2.BDoJob" {
		t.Error("TestServiceImpl2 should overwrite TestServiceImpl in TestServiceB")
	}
}

func TestAutoMake(t *testing.T) {
	registry := NewRegistry()
	AddTransient[TestServiceImpl](registry, Typeof[*TestServiceImpl]())
	AddTransient[TestServiceImpl2](registry, Typeof[*TestServiceImpl2]())

	root := registry.Build().GetProvider()
	srv := MustGet[*TestServiceImpl2](root)

	if srv.Slice == nil {
		t.Error("TestServiceImpl2 should be injected a slice")
	}

	if srv.Map == nil {
		t.Error("TestServiceImpl2 should be injected a map")
	}

	if srv.Chan == nil {
		t.Error("TestServiceImpl2 should be injected a chan")
	}
}

func TestInjectUnexportedField(t *testing.T) {
	registry := NewRegistry()

	defer func() {
		if err := recover(); err == nil {
			t.Error("Add TestServiceWithUnexportedField should panic")
		}
	}()
	AddTransient[TestServiceWithUnexportedField](registry)
}

func TestInjectUnsupportedField(t *testing.T) {
	registry := NewRegistry()

	defer func() {
		if err := recover(); err == nil {
			t.Error("Add TestServiceWithUnsupportedMakeField should panic")
		}
	}()
	AddTransient[TestServiceWithUnsupportedMakeField](registry)
}

func TestAutoInjectTransient(t *testing.T) {
	registry := NewRegistry()

	Add[TestServiceImpl](registry, Transient, Typeof[*TestServiceImpl]())
	Add[TestServiceImpl2](registry, Transient, Typeof[*TestServiceImpl2]())

	root := registry.Build().GetProvider()
	var srv *TestServiceImpl2
	MustLoad(root, &srv)
	if srv.SrvImpl1 == nil {
		t.Error("TestServiceImpl2 should be injected")
	}
	if srv.Provider != root {
		t.Error("TestServiceImpl2 should be injected a provider")
	}
}

func TestAutoInjectSingleton(t *testing.T) {
	registry := NewRegistry()

	Add[TestServiceImpl](registry, Transient, Typeof[*TestServiceImpl]())
	Add[TestServiceImpl2](registry, Singleton, Typeof[*TestServiceImpl2]())

	root := registry.Build().GetProvider()
	var srv *TestServiceImpl2
	MustLoad(root, &srv)

	if srv.SrvImpl1 == nil {
		t.Error("TestServiceImpl2 should be injected")
	}
}

func TestAutoInjectScoped(t *testing.T) {
	registry := NewRegistry()

	Add[TestServiceImpl](registry, Transient, Typeof[*TestServiceImpl]())
	Add[TestServiceImpl2](registry, Scoped, Typeof[*TestServiceImpl2]())

	root := registry.Build().GetProvider()
	var srv *TestServiceImpl2
	MustLoad(root, &srv)

	if srv.SrvImpl1 == nil {
		t.Error("TestServiceImpl2 should be injected")
	}
}
