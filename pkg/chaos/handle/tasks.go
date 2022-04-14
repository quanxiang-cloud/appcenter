package handle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"syscall"
)

type taskQueue struct {
	cachePath string
}

func newTaskQueue(cachePath string) (*taskQueue, bool, error) {
	var init bool = false
	if _, err := os.Stat(cachePath); err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(cachePath)
			if err != nil {
				return nil, false, err
			}
			init = true
		} else {
			return nil, false, err
		}
	}

	return &taskQueue{
		cachePath: cachePath,
	}, init, nil
}

func (tq *taskQueue) put(data interface{}) error {
	f, err := os.OpenFile(tq.cachePath, os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		return err
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_SH); err != nil {
		return err
	}

	cache, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := f.WriteString(string(cache) + "\n"); err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return err
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_UN); err != nil {
		return err
	}

	return nil
}

func (tq *taskQueue) pop(n int) ([][]byte, error) {
	f, err := os.OpenFile(tq.cachePath, os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return nil, err
	}

	defer func() {
		syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	}()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if fi.Size() > 0 {
		buf := bytes.NewBuffer(make([]byte, 0, fi.Size()))

		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}

		if _, err := io.Copy(buf, f); err != nil {
			return nil, err
		}

		ret := make([][]byte, 0, n)
		for i := 0; i < n; i++ {
			line, err := buf.ReadBytes('\n')
			if err != nil && err != io.EOF {
				return nil, err
			}
			if err == io.EOF {
				break
			}

			ret = append(ret, line)
		}

		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}

		nw, err := io.Copy(f, buf)
		if err != nil {
			return nil, err
		}

		if err := f.Truncate(nw); err != nil {
			return nil, err
		}

		if err := f.Sync(); err != nil {
			return nil, err
		}

		return ret, nil
	}

	return nil, fmt.Errorf("no data to pop")
}

type serializeCTX struct {
	RequestID interface{} `json:"requestID"`
	Timezone  interface{} `json:"timezone"`
	TenantID  interface{} `json:"tenantID"`
}

const (
	_requestID = "Request-Id"
	_timezone  = "Timezone"
	_tenantID  = "Tenant-Id"
)

func marshalCTXHeader(c context.Context) serializeCTX {
	return serializeCTX{
		RequestID: c.Value(_requestID),
		Timezone:  c.Value(_timezone),
		TenantID:  c.Value(_tenantID),
	}
}

func unmarshalCTXHeader(c serializeCTX) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, _requestID, c.RequestID)
	ctx = context.WithValue(ctx, _timezone, c.Timezone)
	ctx = context.WithValue(ctx, _tenantID, c.TenantID)
	return ctx
}
