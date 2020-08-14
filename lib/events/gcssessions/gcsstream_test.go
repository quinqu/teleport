/*
Copyright 2020 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package gcssessions

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/gravitational/teleport"
	"github.com/gravitational/teleport/lib/events/test"
	"github.com/gravitational/teleport/lib/utils"

	"github.com/stretchr/testify/assert"
)

// TestStreams tests various streaming upload scenarios
func TestStreams(t *testing.T) {
	utils.InitLoggerForTests(testing.Verbose())

	uri := os.Getenv(teleport.GCSTestURI)
	if uri == "" {
		t.Skip(
			fmt.Sprintf("Skipping GCS tests, set env var %q, details here: https://gravitational.com/teleport/docs/gcp_guide/",
				teleport.GCSTestURI))
	}
	u, err := url.Parse(uri)
	assert.Nil(t, err)

	config := Config{}
	err = config.SetFromURL(u)
	assert.Nil(t, err)

	handler, err := DefaultNewHandler(config)
	assert.Nil(t, err)
	defer handler.Close()

	// Stream with handler and many parts
	t.Run("StreamManyParts", func(t *testing.T) {
		test.StreamManyParts(t, handler)
	})
	t.Run("UploadDownload", func(t *testing.T) {
		test.UploadDownload(t, handler)
	})
	t.Run("DownloadNotFound", func(t *testing.T) {
		test.DownloadNotFound(t, handler)
	})
}
