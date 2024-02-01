package gpfx

import "github.com/lonverce/gpfx/service"

// HostedService 模块初始化后立即启动的服务
type HostedService interface {
	// Start 在后台启动服务, 该方法在主线程中调用, 实现方不能阻塞此方法
	Start()

	// Stop 关闭服务, 该方法在主线程中调用
	Stop()
}

type HostedServiceManager struct {
	services        []HostedService
	ServiceBuilders []func(provider service.Provider) HostedService `gpfx.inject:""`
	Provider        service.Provider                                `gpfx.inject:""`
}

func (h *HostedServiceManager) StartAllServices() {
	for _, builder := range h.ServiceBuilders {
		srv := builder(h.Provider)
		srv.Start()
		h.services = append(h.services, srv)
	}
}

func (h *HostedServiceManager) StopAllServices() {
	for i := len(h.services); i > 0; i-- {
		h.services[i-1].Stop()
	}
}
