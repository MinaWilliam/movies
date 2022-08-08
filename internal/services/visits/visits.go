package visits

import "sync"

var (
	mu    sync.Mutex
	count uint64
)

func RecordVisit() {
	mu.Lock()
	count++
	mu.Unlock()
}

func GetVisitsCount() uint64 {
	return count
}
