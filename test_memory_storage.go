package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

func main() {
	m := NewMemoryStorage()
	fmt.Printf("first index: %d\n", m.FirstIndex())
	fmt.Printf("last index: %d\n", m.LastIndex())
	fmt.Println("Print entries:")
	for _, e := range m.ents {
		fmt.Printf("\t%#v\n", e)
	}

	peers := 5
	ents := make([]Entry, peers)
	for i := 0; i < peers; i++ {
		ents[i] = Entry{Type: EntryConfChange, Term: 1, Index: uint64(i + 1)}
	}
	err := m.Append(ents)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("after append entries, first index: %d\n", m.FirstIndex())
	fmt.Printf("after append entries, last index: %d\n", m.LastIndex())
	fmt.Println("Print entries:")
	for _, e := range m.ents {
		fmt.Printf("\t%#v\n", e)
	}

	m.Compact(3)
	fmt.Println("After compact, print entries:")
	for _, e := range m.ents {
		fmt.Printf("\t%#v\n", e)
	}

}

var ErrCompacted = errors.New("requested index is unavailable due to compaction")

// EntryType describes entry type
type EntryType int32

// entry types
const (
	EntryNormal       EntryType = 0
	EntryConfChange   EntryType = 1
	EntryConfChangeV2 EntryType = 2
)

// Entry describes the raft log entry
type Entry struct {
	Term  uint64
	Index uint64
	Type  EntryType
	Data  []byte
}

// MemoryStorage implements the Storage interface backed by an
// in-memory array.
type MemoryStorage struct {
	// Protects access to all fields. Most methods of MemoryStorage are
	// run on the raft goroutine, but Append() is run on an application
	// goroutine.
	sync.Mutex

	// hardState pb.HardState
	// snapshot  pb.Snapshot

	// ents[i] has raft log position i+snapshot.Metadata.Index
	ents []Entry
}

// NewMemoryStorage creates an empty MemoryStorage.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		// When starting from scratch populate the list with a dummy entry at term zero.
		ents: make([]Entry, 1), //初始化时第一个是 dumpy entry, 如果有快照的话，第一个entry对应快照的index和term信息
	}
}

// LastIndex implements the Storage interface.
func (ms *MemoryStorage) LastIndex() uint64 {
	ms.Lock()
	defer ms.Unlock()
	return ms.lastIndex()
}

func (ms *MemoryStorage) lastIndex() uint64 {
	return ms.ents[0].Index + uint64(len(ms.ents)) - 1
}

// FirstIndex implements the Storage interface.
func (ms *MemoryStorage) FirstIndex() uint64 {
	ms.Lock()
	defer ms.Unlock()
	return ms.firstIndex()
}

func (ms *MemoryStorage) firstIndex() uint64 {
	return ms.ents[0].Index + 1
}

// Append the new entries to strorage
func (ms *MemoryStorage) Append(entries []Entry) error {
	if len(entries) == 0 {
		return nil
	}

	ms.Lock()
	defer ms.Unlock()

	first := ms.firstIndex()
	last := entries[0].Index + uint64(len(entries)) - 1

	// shortcut if there is no new entry.
	if last < first {
		return nil
	}
	// truncate compacted entries
	if first > entries[0].Index {
		entries = entries[first-entries[0].Index:]
	}

	offset := entries[0].Index - ms.ents[0].Index
	switch {
	case uint64(len(ms.ents)) > offset:
		ms.ents = append([]Entry{}, ms.ents[:offset]...)
		ms.ents = append(ms.ents, entries...)
	case uint64(len(ms.ents)) == offset:
		ms.ents = append(ms.ents, entries...)
	default:
		log.Panicf("missing log entry [last: %d, append at: %d]",
			ms.lastIndex(), entries[0].Index)
	}
	return nil
}

func (ms *MemoryStorage) Compact(compactIndex uint64) error {
	ms.Lock()
	defer ms.Unlock()
	offset := ms.ents[0].Index
	if compactIndex <= offset {
		return ErrCompacted
	}
	if compactIndex > ms.lastIndex() {
		log.Panic(fmt.Sprintf("compact %d is out of bound lastindex(%d)", compactIndex, ms.lastIndex()))
	}

	i := compactIndex - offset
	ents := make([]Entry, 1, 1+uint64(len(ms.ents))-i)
	ents[0].Index = ms.ents[i].Index
	ents[0].Term = ms.ents[i].Term
	ents = append(ents, ms.ents[i+1:]...)
	ms.ents = ents
	return nil
}
