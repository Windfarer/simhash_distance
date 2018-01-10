package simhash_distance

import (
	"fmt"
	"math/bits"
	"sync"

	list "github.com/emirpasic/gods/lists/arraylist"
	log "github.com/sirupsen/logrus"
)

const (
	segment         = 4
	hammingDistance = segment - 1
)

type SimHashStore struct {
	Hashs []*sync.Map
}

func NewSimHashStore() *SimHashStore {
	store := new(SimHashStore)
	store.Hashs = make([]*sync.Map, segment)
	for i := 0; i < segment; i++ {
		store.Hashs[i] = &sync.Map{}
	}
	return store
}

func ConvertSimHashToHexStrSegs(simHash uint64) (r []string) {
	h := fmt.Sprintf("%016X", simHash)
	segSize := 16 / segment
	for i := 0; i < segment; i++ {
		r = append(r, h[i*segSize:i+segSize])
	}
	log.Debug(r)
	return r
}

func (e *SimHashStore) CheckSimHash(simHash uint64) (hit bool, sh uint64) {
	hashSegs := ConvertSimHashToHexStrSegs(simHash)
	log.Debug(simHash)

	ch := make(chan *uint64, segment)

	for i := 0; i < segment; i++ {
		go func(i int) {
			if value, ok := e.Hashs[i].Load(hashSegs[i]); ok {
				l := value.(*list.List)
				log.Debugf("hit: %v", i)
				iter := l.Iterator()
				for iter.Next() {
					val := iter.Value().(*uint64)
					log.Debug(val)
					if bits.OnesCount64(simHash^*val) <= hammingDistance {
						ch <- val
						return
					}
				}
			}
			ch <- nil
			return
		}(i)
	}

	for i := 0; i < segment; i++ {
		r := <-ch
		if r != nil {
			return true, *r
		}
	}
	return false, 0
}

func (e *SimHashStore) AddSimHash(simHash uint64) {
	hashSegs := ConvertSimHashToHexStrSegs(simHash)
	for i := 0; i < segment; i++ {
		actual, _ := e.Hashs[i].LoadOrStore(hashSegs[i], list.New())
		actual.(*list.List).Add(&simHash)
	}
}
