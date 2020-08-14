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

	"gopkg.in/check.v1"
)

type GCSStreamSuite struct {
	handler *Handler
	test.HandlerSuite
}

var _ = check.Suite(&GCSStreamSuite{})

func (s *GCSStreamSuite) SetUpSuite(c *check.C) {
	utils.InitLoggerForTests(testing.Verbose())

	config := Config{}
	uri := os.Getenv(teleport.GCSTestURI)
	if uri == "" {
		c.Skip(
			fmt.Sprintf("Skipping GCS tests, set env var %q, details here: https://gravitational.com/teleport/docs/gcp_guide/",
				teleport.GCSTestURI))
	}
	u, err := url.Parse(uri)
	c.Assert(err, check.IsNil)

	err = config.SetFromURL(u)
	c.Assert(err, check.IsNil)

	s.handler, err = DefaultNewHandler(config)
	c.Assert(err, check.IsNil)

	s.HandlerSuite.Handler = s.handler
}

func (s *GCSStreamSuite) TestStream(c *check.C) {
	s.StreamManyParts(c)
}

func (s *GCSStreamSuite) TearDownSuite(c *check.C) {
	if s.handler != nil {
		s.handler.Close()
	}
}
