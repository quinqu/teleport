package test

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gravitational/teleport/lib/events"
	"github.com/gravitational/teleport/lib/fixtures"
	"github.com/gravitational/teleport/lib/session"

	"gopkg.in/check.v1"
)

// StreamParams configures parameters of a stream test suite
type StreamParams struct {
	// PrintEvents is amount of print events to generate
	PrintEvents int64
	// ConcurrentUploads is amount of concurrent uploads
	ConcurrentUploads int
	// MinUploadBytes is minimum required upload bytes
	MinUploadBytes int64
}

// StreamSinglePart tests stream upload and subsequent download and reads the results
func (s *HandlerSuite) StreamSinglePart(c *check.C) {
	s.StreamWithParameters(c, StreamParams{
		PrintEvents:    1024,
		MinUploadBytes: 1024 * 1024,
	})
}

// Stream tests stream upload and subsequent download and reads the results
func (s *HandlerSuite) Stream(c *check.C) {
	s.StreamWithParameters(c, StreamParams{
		PrintEvents:       1024,
		MinUploadBytes:    1024,
		ConcurrentUploads: 2,
	})
}

// StreamManyParts tests stream upload and subsequent download and reads the results
func (s *HandlerSuite) StreamManyParts(c *check.C) {
	s.StreamWithParameters(c, StreamParams{
		PrintEvents:       8192,
		MinUploadBytes:    1024,
		ConcurrentUploads: 64,
	})
}

// StreamWithParameters tests stream upload and subsequent download and reads the results
func (s *HandlerSuite) StreamWithParameters(c *check.C, params StreamParams) {
	ctx := context.TODO()

	inEvents := events.GenerateSession(params.PrintEvents)
	sid := session.ID(inEvents[0].(events.SessionMetadataGetter).GetSessionID())

	streamer, err := events.NewProtoStreamer(events.ProtoStreamerConfig{
		Uploader:          s.Handler,
		MinUploadBytes:    params.MinUploadBytes,
		ConcurrentUploads: params.ConcurrentUploads,
	})
	c.Assert(err, check.IsNil)

	stream, err := streamer.CreateAuditStream(ctx, sid)
	c.Assert(err, check.IsNil)

	for _, event := range inEvents {
		err := stream.EmitAuditEvent(ctx, event)
		c.Assert(err, check.IsNil)
	}

	err = stream.Complete(ctx)
	c.Assert(err, check.IsNil)

	dir := c.MkDir()
	f, err := os.Create(filepath.Join(dir, string(sid)))
	c.Assert(err, check.IsNil)
	defer f.Close()

	err = s.Handler.Download(ctx, sid, f)
	c.Assert(err, check.IsNil)

	_, err = f.Seek(0, 0)
	c.Assert(err, check.IsNil)

	reader := events.NewProtoReader(f)
	out, err := reader.ReadAll(ctx)
	c.Assert(err, check.IsNil)

	stats := reader.GetStats()
	c.Assert(stats.SkippedEvents, check.Equals, int64(0))
	c.Assert(stats.OutOfOrderEvents, check.Equals, int64(0))
	c.Assert(stats.TotalEvents, check.Equals, int64(len(inEvents)))

	fixtures.DeepCompareSlices(c, inEvents, out)
}
