package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type DoRequestParams struct {
	Method       string
	Path         string
	QueryParams  KV
	RequestBody  interface{}
	ResponseBody interface{}
	Headers      KV
}

func (c HTTPClient) DoRequest(ctx context.Context, params DoRequestParams) (*http.Request, *http.Response, error) {
	var request *http.Request
	var err error

	url := c.baseURL
	url.Path += params.Path
	query := url.Query()
	for _, qParam := range params.QueryParams {
		query.Add(qParam.Key, qParam.Value)
	}
	url.RawQuery = query.Encode()

	if params.RequestBody != nil {
		requestBodyBytes := &bytes.Buffer{}

		if b, ok := params.RequestBody.([]byte); ok {
			requestBodyBytes.Write(b)
		} else {
			if marshallable, ok := params.RequestBody.(marshallable); ok {
				b, err = marshallable.MarshalJSON()
			} else {
				b, err = json.Marshal(params.RequestBody)
			}
			if err != nil {
				return nil, nil, err
			}
			requestBodyBytes.Write(b)
			request.Header.Add("Content-Type", "application/json")
		}

		request, err = http.NewRequestWithContext(ctx, params.Method, url.String(), requestBodyBytes)
	} else {
		request, err = http.NewRequestWithContext(ctx, params.Method, url.String(), nil)
	}
	if err != nil {
		return request, nil, err
	}

	for _, kv := range params.Headers {
		request.Header.Set(kv.Key, kv.Value)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return request, response, err
	}

	defer response.Body.Close()

	if params.ResponseBody != nil {
		var b []byte
		var err error

		if unmarshallable, ok := params.ResponseBody.(unmarshallable); ok {
			err = unmarshallable.UnmarshalJSON(b)
		} else {
			err = json.Unmarshal(b, params.ResponseBody)
		}
		if err != nil {
			return request, response, err
		}
	}

	return request, response, nil
}
