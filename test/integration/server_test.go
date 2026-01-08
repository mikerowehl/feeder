/*
Copyright (c) Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
// Small http server used as the target to fetch from during integration tests
package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func startTestFeedServer(t *testing.T) *httptest.Server {
	t.Helper()

	fs := http.FileServer(http.Dir("../feeds"))
	server := httptest.NewServer(fs)

	t.Cleanup(func() {
		server.Close()
	})

	return server
}

func getTestFeedURL(server *httptest.Server, filename string) string {
	return server.URL + "/" + filename
}
