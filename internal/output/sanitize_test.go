/*
Copyright (c) Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package output_test

import (
	"testing"

	"github.com/mikerowehl/feeder/internal/output"

	"github.com/stretchr/testify/assert"
)

func TestSanitize_SafeURL(t *testing.T) {
	assert.Equal(t, "http://rowehl.com", output.SafeURL("http://rowehl.com"))
	assert.Equal(t, "#", output.SafeURL(`"><script>alert('XSS')</script>`))
}
