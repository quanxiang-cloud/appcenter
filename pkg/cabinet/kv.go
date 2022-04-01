package cabinet

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type KV struct {
	data map[string]string
}

func New() *KV {
	kv := &KV{
		data: make(map[string]string),
	}
	flag.Var(kv, "kv", "kv used to obtain customize params. usage: --kv k=v")

	return kv
}

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

func (kv *KV) Get(k string) string {
	return kv.data[k]
}

func (kv *KV) GetString(k string) (string, error) {
	return kv.Get(k), nil
}

func (kv *KV) GetInt(k string) (int, error) {
	v := kv.Get(k)
	return strconv.Atoi(v)
}

func (kv *KV) GetBool(k string) (bool, error) {
	v := kv.Get(k)
	return strconv.ParseBool(v)
}
