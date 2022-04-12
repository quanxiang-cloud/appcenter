package exec

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

// post http post
func post(ctx context.Context, client *http.Client, uri string, params interface{}, entity interface{}) error {
	if reflect.ValueOf(entity).Kind() != reflect.Ptr {
		return errors.New("the entity type must be a pointer")
	}

	paramByte, err := json.Marshal(params)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(paramByte)
	req, err := http.NewRequest("POST", uri, reader)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(header.GetRequestIDKV(ctx).Wreck())
	req.Header.Add(header.GetTimezone(ctx).Wreck())
	req.Header.Add(header.GetTenantID(ctx).Wreck())

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("expected state value is 200, actually %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return decomposeBody(body, entity)
}

func decomposeBody(body []byte, entity interface{}) error {
	r := new(resp.Resp)
	r.Data = entity

	err := json.Unmarshal(body, r)
	if err != nil {
		return err
	}

	return nil
}
