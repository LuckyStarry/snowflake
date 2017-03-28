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

func TestNextSnowflakeID(t *testing.T) {
	sfid := NextSnowflakeID(0)
	expect := base62encode(sfid.ToInt64())
	actual := sfid.ToBase62()
	if expect != actual {
		t.Logf("sfid:%d", sfid.ToInt64())
		t.Logf("expect:%s", expect)
		t.Logf("actual:%s", actual)
		t.Errorf("base 62 generate failed")
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

// algorithm from network
var base = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

func base62encode(num int64) string {
	baseStr := ""
	for {
		if num <= 0 {
			break
		}

		i := num % 62
		baseStr += base[i]
		num = (num - i) / 62
	}
	return baseStr
}
