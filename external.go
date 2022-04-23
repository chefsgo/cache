package cache

import (
	"time"

	. "github.com/chefsgo/base"
)

func Read(key string) (Any, error) {
	return module.Read(key)

}
func Exists(key string) (bool, error) {
	return module.Exists(key)
}

func Write(key string, value Any, expiries ...time.Duration) error {
	return module.Write(key, value, expiries...)
}

func Delete(key string) error {
	return module.Delete(key)
}

func Serial(key string, start, step int64) (int64, error) {
	return module.Serial(key, start, step)
}

func Keys(prefix string) ([]string, error) {
	return module.Keys(prefix)
}
func Clear(prefix string) error {
	return module.Clear(prefix)
}
