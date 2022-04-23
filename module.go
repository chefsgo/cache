package cache

import (
	"sync"
	"time"

	. "github.com/chefsgo/base"
	"github.com/chefsgo/chef"
	"github.com/chefsgo/util"
)

func init() {
	chef.Register(NAME, module)
}

var (
	module = &Module{
		configs:   make(map[string]Config, 0),
		drivers:   make(map[string]Driver, 0),
		instances: make(map[string]Instance, 0),
	}
)

type (
	Module struct {
		mutex sync.Mutex

		connected, initialized, launched bool

		configs map[string]Config
		drivers map[string]Driver

		instances map[string]Instance

		weights  map[string]int
		hashring *util.HashRing
	}

	Config struct {
		Driver  string
		Weight  int
		Prefix  string
		Expiry  time.Duration
		Setting Map
	}
	Instance struct {
		name    string
		config  Config
		connect Connect
	}
)

func (this *Module) Read(key string) (Any, error) {
	locate := module.hashring.Locate(key)

	if inst, ok := module.instances[locate]; ok {
		key := inst.config.Prefix + key //加前缀
		return inst.connect.Read(key)
	}

	return nil, errInvalidCacheConnection
}

func (this *Module) Exists(key string) (bool, error) {
	locate := module.hashring.Locate(key)

	if inst, ok := module.instances[locate]; ok {
		key := inst.config.Prefix + key //加前缀
		return inst.connect.Exists(key)
	}

	return false, errInvalidCacheConnection
}

// Write 写缓存
func (this *Module) Write(key string, val Any, expiries ...time.Duration) error {
	locate := module.hashring.Locate(key)

	if inst, ok := module.instances[locate]; ok {
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
	locate := module.hashring.Locate(key)

	if inst, ok := module.instances[locate]; ok {
		key := inst.config.Prefix + key
		return inst.connect.Delete(key)
	}

	return errInvalidCacheConnection
}

// Serial 生成序列编号
func (this *Module) Serial(key string, start, step int64) (int64, error) {
	locate := module.hashring.Locate(key)

	if inst, ok := module.instances[locate]; ok {
		key := inst.config.Prefix + key
		return inst.connect.Serial(key, start, step)
	}

	return -1, errInvalidCacheConnection
}

// Keys 获取所有前缀的KEYS
func (this *Module) Keys(prefix string) ([]string, error) {
	keys := make([]string, 0)

	for _, inst := range module.instances {
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
	for _, inst := range module.instances {
		prefix := inst.config.Prefix + prefix
		inst.connect.Clear(prefix)
	}

	return errInvalidCacheConnection
}
