/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package mock

import (
	"io"
	"net/http"
	"strings"
)

type MockRoundTripper func(req *http.Request) *http.Response

func (m MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m(req), nil
}

func NewMockClient(body string) *http.Client {
	mockTransport := MockRoundTripper(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}
	})

	return &http.Client{Transport: mockTransport}
}
