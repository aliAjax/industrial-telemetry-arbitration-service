package arbitration

import "testing"

func TestLeaseSnapshotCopiesMetadata(t *testing.T) {
	lease := NewLease("lease-1", map[string]string{"site": "north"})
	snapshot := lease.Snapshot()
	snapshot.Metadata["site"] = "changed"
	if lease.Snapshot().Metadata["site"] != "north" {
		t.Fatal("metadata escaped")
	}
}
