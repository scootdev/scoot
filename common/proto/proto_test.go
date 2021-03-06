package proto

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/empty"
	remoteexecution "github.com/twitter/scoot/bazel/remoteexecution"

	"github.com/twitter/scoot/bazel"
)

func TestGetEmptySha256(t *testing.T) {
	e := empty.Empty{}
	s, l, err := GetSha256(&e)
	if err != nil {
		t.Fatalf("GetSha256 failure: %v", err)
	}
	if s != bazel.EmptySha {
		t.Fatalf("Expected known sha for nil/empty data %s, got: %s", bazel.EmptySha, s)
	}
	if l != 0 {
		t.Fatalf("Expected zero length data, got: %d", l)
	}
}

// This is a canary test of sorts for generating digests of Action data from bazel
// ExecuteRequests. If this starts to fail, it indicates an instability in hashing Action messages.
func TestGetActionSha256(t *testing.T) {
	a := &remoteexecution.Action{
		CommandDigest:   &remoteexecution.Digest{Hash: "abc123", SizeBytes: 10},
		InputRootDigest: &remoteexecution.Digest{Hash: "def456", SizeBytes: 20},
		Timeout:         GetDurationFromMs(60000),
		DoNotCache:      true,
	}
	s, _, err := GetSha256(a)
	if err != nil {
		t.Fatalf("GetSha256 failure: %v", err)
	}
	expectedSha := "ee30b3288c00f2b0e89d49138ca1ff2c0e8223b0de59b491b8de555595f22586"
	if s != expectedSha {
		t.Fatalf("Expected known sha for message data: %s, got: %s", expectedSha, s)
	}
}

func TestMsDuration(t *testing.T) {
	d := duration.Duration{Seconds: 3, Nanos: 5000004}
	ms := GetMsFromDuration(&d)
	if ms != 3005 {
		t.Fatalf("Expected 3005, got: %dms", ms)
	}

	dp := GetDurationFromMs(ms)
	if dp == nil {
		t.Fatalf("Unexpected nil result from GetDurationFromMs(%d)", ms)
	}
	if dp.GetSeconds() != 3 || dp.GetNanos() != 5000000 {
		t.Fatalf("Expected 3s 5000000ns, got: %ds %dns", dp.GetSeconds(), dp.GetNanos())
	}
}

func TestTimeTimestamps(t *testing.T) {
	now := time.Now()
	ts := GetTimestampFromTime(now)
	t2 := GetTimeFromTimestamp(ts)

	if !now.Equal(t2) {
		t.Fatalf("Time converted from timestamp did not match, got: %v, expected: %v", t2, now)
	}
}
