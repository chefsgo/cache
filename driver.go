package cache

import (
	"time"

	. "github.com/chefsgo/base"
)

type (
	// Driver 数据驱动
	Driver interface {
		Connect(name string, config Config) (Connect, error)
	}

	// Connect 会话连接
	Connect interface {
		Open() error
		Close() error

		Read(string) (Any, error)
		Write(key string, val Any, expiry time.Duration) error
		Exists(key string) (bool, error)
		Delete(key string) error
		Serial(key string, start, step int64) (int64, error)
		Keys(prefix string) ([]string, error)
		Clear(prefix string) error
	}
)
