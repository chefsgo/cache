package cache

import (
	"time"

	. "github.com/chefsgo/base"
)

func (this *Module) Read(key string) (Any, error) {
	locate := this.hashring.Locate(key)

	if inst, ok := this.instances[locate]; ok {
		key := inst.config.Prefix + key //加前缀
		return inst.connect.Read(key)
	}

	return nil, errInvalidCacheConnection
}

func (this *Module) Exists(key string) (bool, error) {
	locate := this.hashring.Locate(key)

	if inst, ok := this.instances[locate]; ok {
		key := inst.config.Prefix + key //加前缀
		return inst.connect.Exists(key)
	}

	return false, errInvalidCacheConnection
}

// Write 写缓存
func (this *Module) Write(key string, val Any, expiries ...time.Duration) error {
	locate := this.hashring.Locate(key)

	if inst, ok := this.instances[locate]; ok {
		expiry := inst.config.Expiry
		if len(expiries) > 0 {
			expiry = expiries[0]
		}

		//KEY加上前缀
		key := inst.config.Prefix + key

		return inst.connect.Write(key, val, expiry)
	}

	return errInvalidCacheConnection
}

// Delete 删除缓存
func (this *Module) Delete(key string) error {
	locate := this.hashring.Locate(key)

	if inst, ok := this.instances[locate]; ok {
		key := inst.config.Prefix + key
		return inst.connect.Delete(key)
	}

	return errInvalidCacheConnection
}

// Serial 生成序列编号
func (this *Module) Serial(key string, start, step int64) (int64, error) {
	locate := this.hashring.Locate(key)

	if inst, ok := this.instances[locate]; ok {
		key := inst.config.Prefix + key
		return inst.connect.Serial(key, start, step)
	}

	return -1, errInvalidCacheConnection
}

// Keys 获取所有前缀的KEYS
func (this *Module) Keys(prefix string) ([]string, error) {
	keys := make([]string, 0)

	for _, inst := range this.instances {
		prefix := inst.config.Prefix + prefix
		temps, err := inst.connect.Keys(prefix)
		if err == nil {
			keys = append(keys, temps...)
		}
	}

	return keys, nil
}

// Clear 按前缀清理缓存
func (this *Module) Clear(prefix string) error {
	for _, inst := range this.instances {
		prefix := inst.config.Prefix + prefix
		inst.connect.Clear(prefix)
	}

	return errInvalidCacheConnection
}
