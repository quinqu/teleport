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

package memsessions

import (
	"testing"

	"github.com/gravitational/teleport/lib/events"
	"github.com/gravitational/teleport/lib/events/test"
	"github.com/gravitational/teleport/lib/utils"

	"gopkg.in/check.v1"
)

func TestGCS(t *testing.T) { check.TestingT(t) }

type MemSuite struct {
	test.HandlerSuite
}

var _ = check.Suite(&MemSuite{})

func (s *MemSuite) SetUpSuite(c *check.C) {
	utils.InitLoggerForTests(testing.Verbose())
	s.HandlerSuite.Handler = events.NewMemoryUploader()
}

func (s *MemSuite) TestStream(c *check.C) {
	s.StreamManyParts(c)
}
