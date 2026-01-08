/*
Copyright (c) Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package mock

import (
	"io"
	"net/http"
	"strings"
)

type MockRoundTripper func(req *http.Request) (*http.Response, error)

func (m MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m(req)
}

func NewMockClient(body string, statusCode int) *http.Client {
	mockTransport := MockRoundTripper(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: statusCode,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	return &http.Client{Transport: mockTransport}
}

func NewMockClientWithError(err error) *http.Client {
	mockTransport := MockRoundTripper(func(req *http.Request) (*http.Response, error) {
		return nil, err
	})

	return &http.Client{Transport: mockTransport}
}
