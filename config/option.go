package config

// Option 是针对配置项的抽象
type Option[TOption any] interface {
	// Value 获取配置项值
	Value() *TOption

	OnceValue() *TOption
}
