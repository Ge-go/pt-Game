package email

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type HttpResponse struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
	Seq string `json:"seq"`
}

func (e *EmailSender) sendHttpRequest(ctx context.Context, method, url string, headers map[string]string, body []byte) (*HttpResponse, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var ret HttpResponse
	err = json.NewDecoder(res.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
