package simhash_distance

import (
	"testing"
)

func TestSimHashStore_TestAddSimHash(t *testing.T) {
	s := NewSimHashStore()
	s.AddSimHash(4310046566681409340)
	s.AddSimHash(4310046566681409341)
	s.AddSimHash(1310046566681409341)

}

func TestSimHashStore_CheckSimHash(t *testing.T) {
	s := NewSimHashStore()
	var sampleHash uint64 = 4310046566681409340
	s.AddSimHash(sampleHash)
	hit, dupHash := s.CheckSimHash(sampleHash+3)
	if hit != true {
		t.Error("distance calculating failed")
	}
	if dupHash != sampleHash {
		t.Error("get wrong simhash")
	}

	hit, dupHash = s.CheckSimHash(sampleHash+4)
	if hit != false {
		t.Error("should be false")
	}
}