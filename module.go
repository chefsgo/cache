package cache

import (
	"sync"
	"time"

	. "github.com/chefsgo/base"
	"github.com/chefsgo/chef"
	"github.com/chefsgo/util"
)

func init() {
	chef.Register(module)
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

	Configs map[string]Config
	Config  struct {
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

// Driver 注册驱动
func (this *Module) Driver(name string, driver Driver, override bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if driver == nil {
		panic("Invalid cache driver: " + name)
	}

	if override {
		this.drivers[name] = driver
	} else {
		if this.drivers[name] == nil {
			this.drivers[name] = driver
		}
	}
}

func (this *Module) Config(name string, config Config, override bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if name == "" {
		name = chef.DEFAULT
	}

	if override {
		this.configs[name] = config
	} else {
		if _, ok := this.configs[name]; ok == false {
			this.configs[name] = config
		}
	}
}
func (this *Module) Configs(name string, config Configs, override bool) {
	for key, val := range config {
		this.Config(key, val, override)
	}
}
