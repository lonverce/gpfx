package config

import "time"

// Provider 配置
type Provider interface {
	GetString(key string) (val string, exist bool)
	GetStringSlice(key string) (val []string, exist bool)
	GetStringMap(key string) (map[string]any, bool)
	GetStringMapString(key string) (map[string]string, bool)
	GetStringMapStringSlice(key string) (map[string][]string, bool)
	GetBool(key string) (val bool, exist bool)
	GetInt(key string) (val int, exist bool)
	GetIntSlice(key string) (val []int, exist bool)
	GetFloat64(key string) (val float64, exist bool)
	GetTime(key string) (val time.Time, exist bool)
	GetDuration(key string) (val time.Duration, exist bool)
	IsSet(key string) bool
	Sub(key string) Provider
	Unmarshal(key string, structPointer any) bool
}

func GetStringOrDefault(c Provider, key string, defaultVal string) string {
	val, ok := c.GetString(key)
	if !ok {
		return defaultVal
	}
	return val
}

func GetBoolOrDefault(c Provider, key string, defaultVal bool) bool {
	val, ok := c.GetBool(key)
	if !ok {
		return defaultVal
	}
	return val
}

func GetIntOrDefault(c Provider, key string, defaultVal int) int {
	val, ok := c.GetInt(key)
	if !ok {
		return defaultVal
	}
	return val
}

func GetFloat64OrDefault(c Provider, key string, defaultVal float64) float64 {
	val, ok := c.GetFloat64(key)
	if !ok {
		return defaultVal
	}
	return val
}

func GetStruct[T any](c Provider, key string) *T {
	s := new(T)
	if !c.Unmarshal(key, s) {
		return nil
	}
	return s
}
