package service

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

func (h *httpService) Get(ctx context.Context, api string, params map[string]interface{}) (*Response, error) {
	if h.cache != nil {
		return h.cache.Get(ctx, api, params)
	}

	return h.call(ctx, http.MethodGet, api, params, nil, nil)
}

func (h *httpService) Post(ctx context.Context, api string, params map[string]interface{}, body []byte) (*Response, error) {
	return h.call(ctx, http.MethodPost, api, params, body, nil)
}

func (h *httpService) Put(ctx context.Context, api string, params map[string]interface{}, body []byte) (*Response, error) {
	return h.call(ctx, http.MethodPut, api, params, body, nil)
}

func (h *httpService) Patch(ctx context.Context, api string, params map[string]interface{}, body []byte) (*Response, error) {
	return h.call(ctx, http.MethodPatch, api, params, body, nil)
}

func (h *httpService) Delete(ctx context.Context, api string, body []byte) (*Response, error) {
	return h.call(ctx, http.MethodDelete, api, nil, body, nil)
}

func (h *httpService) GetWithHeaders(ctx context.Context, api string, params map[string]interface{},
	headers map[string]string) (*Response, error) {
	if h.cache != nil {
		return h.cache.GetWithHeaders(ctx, api, params, headers)
	}

	return h.call(ctx, "GET", api, params, nil, headers)
}

func (h *httpService) PostWithHeaders(ctx context.Context, api string, params map[string]interface{}, body []byte,
	headers map[string]string) (*Response, error) {
	return h.call(ctx, http.MethodPost, api, params, body, headers)
}

func (h *httpService) PutWithHeaders(ctx context.Context, api string, params map[string]interface{}, body []byte,
	headers map[string]string) (*Response, error) {
	return h.call(ctx, http.MethodPut, api, params, body, headers)
}

func (h *httpService) PatchWithHeaders(ctx context.Context, api string, params map[string]interface{}, body []byte,
	headers map[string]string) (*Response, error) {
	return h.call(ctx, http.MethodPatch, api, params, body, headers)
}

func (h *httpService) DeleteWithHeaders(ctx context.Context, api string, body []byte, headers map[string]string) (*Response, error) {
	return h.call(ctx, http.MethodDelete, api, nil, body, headers)
}

// Bind takes Response and binds it to i based on content-type.
func (h *httpService) Bind(resp []byte, i interface{}) error {
	var err error

	h.mu.Lock()
	contentType := h.contentType
	h.mu.Unlock()

	switch contentType {
	case XML:
		err = xml.NewDecoder(bytes.NewBuffer(resp)).Decode(&i)
	case TEXT:
		v, ok := i.(*string)
		if ok {
			*v = fmt.Sprintf("%s", resp)
		}
	case HTML, JSON:
		err = json.NewDecoder(bytes.NewBuffer(resp)).Decode(&i)
	}

	return err
}

func (h *httpService) BindStrict(resp []byte, i interface{}) error {
	var err error

	h.mu.Lock()
	contentType := h.contentType
	h.mu.Unlock()

	switch contentType {
	case XML:
		err = xml.NewDecoder(bytes.NewBuffer(resp)).Decode(&i)
	case TEXT:
		v, ok := i.(*string)
		if ok {
			*v = fmt.Sprintf("%s", resp)
		}
	case HTML, JSON:
		dec := json.NewDecoder(bytes.NewBuffer(resp))
		dec.DisallowUnknownFields()
		err = dec.Decode(&i)
	}

	return err
}
