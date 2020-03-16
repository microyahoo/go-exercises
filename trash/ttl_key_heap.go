package trash

import (
	"container/heap"
)

// An ttlKeyHeap is a min-heap according to expiration time
type ttlKeyHeap struct {
	array  []Job
	keyMap map[Job]int
}

func newTTLKeyHeap() *ttlKeyHeap {
	h := &ttlKeyHeap{
		array:  make([]Job, 0),
		keyMap: make(map[Job]int),
	}
	heap.Init(h)
	return h
}

func (h ttlKeyHeap) Len() int {
	return len(h.array)
}

func (h ttlKeyHeap) Less(i, j int) bool {
	return h.array[i].GetExpireTime().Before(h.array[j].GetExpireTime())
}

func (h ttlKeyHeap) Swap(i, j int) {
	// swap job
	h.array[i], h.array[j] = h.array[j], h.array[i]

	// update map
	h.keyMap[h.array[i]] = i
	h.keyMap[h.array[j]] = j
}

func (h *ttlKeyHeap) Push(x interface{}) {
	j, _ := x.(Job)
	h.keyMap[j] = len(h.array)
	h.array = append(h.array, j)
}

func (h *ttlKeyHeap) Pop() interface{} {
	old := h.array
	n := len(old)
	x := old[n-1]
	old[n-1] = nil
	h.array = old[0 : n-1]
	delete(h.keyMap, x)
	return x
}

func (h *ttlKeyHeap) top() Job {
	if h.Len() != 0 {
		return h.array[0]
	}
	return nil
}

func (h *ttlKeyHeap) pop() Job {
	x := heap.Pop(h)
	j, _ := x.(Job)
	return j
}

func (h *ttlKeyHeap) push(x interface{}) {
	heap.Push(h, x)
}

func (h *ttlKeyHeap) update(j Job) {
	index, ok := h.keyMap[j]
	if ok {
		heap.Remove(h, index)
		heap.Push(h, j)
	}
}

func (h *ttlKeyHeap) remove(j Job) {
	index, ok := h.keyMap[j]
	if ok {
		heap.Remove(h, index)
	}
}
