package snowflake

import "testing"

func TestNextID(t *testing.T) {
	uid1 := NextID()
	uid2 := NextID()
	if uid1 == uid2 {
		t.Errorf("uid equals")
	}
}

func TestNextIDWorker(t *testing.T) {
	uid1 := NextIDWorker(1)
	uid2 := NextIDWorker(2)
	if uid1 == uid2 {
		t.Errorf("uid equals")
	}
}

func TestNextIDWorkerError(t *testing.T) {
	if v := NextIDWorker(-1); v != -1 {
		t.Errorf("id would be -1 but was %d when worker id is -1", v)
	}
	if v := NextIDWorker(0x3ff + 1); v != -1 {
		t.Errorf("id would be -1 but was %d when worker id is 1024", v)
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
}

func BenchmarkNextID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NextID()
	}
}

func BenchmarkNextIDParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var workerID int64
		for pb.Next() {
			NextIDWorker(workerID % 0x3ff)
			workerID++
		}
	})
}
