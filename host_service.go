package gpfx

// HostedService 模块初始化后立即启动的服务
type HostedService interface {
	// Start 在后台启动服务, 该方法在主线程中调用, 实现方不能阻塞此方法
	Start()

	// Stop 关闭服务, 该方法在主线程中调用
	Stop()
}

type hostedServiceManager struct {
	services []HostedService
}

func (h *hostedServiceManager) Inject(provider InterimServiceContext) {
	h.services = LoadAllServices[HostedService](provider)
}

func (h *hostedServiceManager) StartAllServices() {
	for _, service := range h.services {
		service.Start()
	}
}

func (h *hostedServiceManager) StopAllServices() {
	for i := len(h.services); i > 0; i-- {
		h.services[i-1].Stop()
	}
}
