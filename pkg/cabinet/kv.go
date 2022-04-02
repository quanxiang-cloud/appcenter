package cabinet

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// KV kv
type KV struct {
	data map[string]string
}

// New new
func New() *KV {
	kv := &KV{
		data: make(map[string]string),
	}
	flag.Var(kv, "kv", "kv used to obtain customize params. usage: --kv k=v")

	return kv
}

// Set set
func (kv *KV) Set(val string) error {
	index := strings.Index(val, "=")
	if index > 0 {
		kv.data[val[0:index]] = val[index+1:]
		return nil
	}
	return fmt.Errorf("value of KV not contains '='")
}

func (kv *KV) String() string {
	return fmt.Sprint(kv.data)
}

// Get get
func (kv *KV) Get(k string) string {
	return kv.data[k]
}

// GetString get string
func (kv *KV) GetString(k string) (string, error) {
	return kv.Get(k), nil
}

// GetInt GetInt
func (kv *KV) GetInt(k string) (int, error) {
	v := kv.Get(k)
	return strconv.Atoi(v)
}

// GetBool GetBool
func (kv *KV) GetBool(k string) (bool, error) {
	v := kv.Get(k)
	return strconv.ParseBool(v)
}
