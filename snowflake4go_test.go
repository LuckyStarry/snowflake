package snowflake4go

import "testing"

func TestNextID(t *testing.T) {
	uid1 := NextID()
	uid2 := NextID()
	if uid1 == uid2 {
		t.Errorf("uid equals")
	} else {
		t.Logf("UID Duplicate detection pass")
	}
}

func TestNextIDWorker(t *testing.T) {
	uid1 := NextIDWorker(1)
	uid2 := NextIDWorker(2)
	if uid1 == uid2 {
		t.Errorf("uid equals")
	} else {
		t.Logf("UID Duplicate detection pass")
	}
}

func TestDuplicate(t *testing.T) {
	times := 100
	for i := 0; i < times; i++ {
		o := NextID()
		for j := 0; j <= 0xfff; j++ {
			v := NextID()
			if o == v {
				t.Errorf("uid equals")
			}
		}
	}
	t.Logf("UID Duplicate detection pass, times %d", times)
}

func BenchmarkNextID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NextID()
	}
}

func BenchmarkNextIDParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			NextID()
		}
	})
}
