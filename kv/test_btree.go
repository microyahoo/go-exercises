package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/btree"
	"go.uber.org/zap"
)

var (
	ErrRevisionNotFound = errors.New("mvcc: revision not found")
	btreeDegree         = flag.Int("degree", 32, "B-Tree degree")
	ch                  = make(chan struct{})
)

func createBytesSlice(bytesN, sliceN int) [][]byte {
	rs := [][]byte{}
	for len(rs) != sliceN {
		v := make([]byte, bytesN)
		if _, err := rand.Read(v); err != nil {
			panic(err)
		}
		rs = append(rs, v)
	}
	return rs
}

func main() {
	tr := btree.New(*btreeDegree)
	for i := btree.Int(0); i < 10; i++ {
		tr.ReplaceOrInsert(i)
	}
	fmt.Println("len:       ", tr.Len())
	fmt.Println("get3:      ", tr.Get(btree.Int(3)))
	fmt.Println("get100:    ", tr.Get(btree.Int(100)))
	fmt.Println("del4:      ", tr.Delete(btree.Int(4)))
	fmt.Println("del100:    ", tr.Delete(btree.Int(100)))
	fmt.Println("replace5:  ", tr.ReplaceOrInsert(btree.Int(5)))
	fmt.Println("replace100:", tr.ReplaceOrInsert(btree.Int(100)))
	fmt.Println("min:       ", tr.Min())
	fmt.Println("delmin:    ", tr.DeleteMin())
	fmt.Println("max:       ", tr.Max())
	fmt.Println("delmax:    ", tr.DeleteMax())
	fmt.Println("len:       ", tr.Len())

	ti := newTreeIndex(zap.NewExample())
	size := 20
	// size := 1000000
	bytesN := 64
	keys := createBytesSlice(bytesN, size)
	for i := 1; i < size-10; i++ {
		ti.Put(keys[i], revision{main: int64(i), sub: int64(i)})
	}
	for i := 1; i < size-10; i++ {
		ti.Put(keys[i], revision{main: int64(i + 3), sub: int64(i)})
	}
	for i := 1; i < size-10; i++ {
		ti.Tombstone(keys[i], revision{main: int64(i + 5), sub: int64(i)})
	}
	for i := 1; i < size-10; i++ {
		ti.Put(keys[i], revision{main: int64(i + 7), sub: int64(i)})
	}
	for i := 1; i < size-10; i++ {
		ti.Tombstone(keys[i], revision{main: int64(i + 9), sub: int64(i)})
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	fmt.Println("--------- before compaction ----------------")
	ti.tree.Ascend(func(item btree.Item) bool {
		keyi := item.(*keyIndex)
		fmt.Println(keyi)
		return true
	})
	go func() {
		defer wg.Done()
		t1 := time.Now()
		ti.Compact(int64(15))
		// ti.Compact(int64(500000))
		fmt.Printf("The compaction time is: %v\n", time.Since(t1))
		fmt.Println("--------- after compaction ----------------")
		ti.tree.Ascend(func(item btree.Item) bool {
			keyi := item.(*keyIndex)
			fmt.Printf("***** %s\n", keyi)
			return true
		})
	}()

	t1 := time.Now().UnixNano()
	// for i := size - 10; i < size; i++ {
	// 	ti.Put(keys[i], revision{main: int64(i), sub: int64(i)})
	// }
	<-ch
	for i := 1; i < size-10; i++ {
		ti.Put(keys[i], revision{main: int64(i + 11), sub: int64(i)})
	}
	fmt.Println("--------- after put ----------------")
	ti.tree.Ascend(func(item btree.Item) bool {
		keyi := item.(*keyIndex)
		fmt.Printf("###### %s\n", keyi)
		return true
	})
	// ti.Put(keys[size-1], revision{main: int64(size - 1), sub: int64(size - 1)})
	t2 := time.Now().UnixNano() - t1
	if t2 > 150000000 {
		fmt.Printf("Run put time took too long! %vns\n", t2)
	} else {
		fmt.Printf("The put time after compaction: %vns\n", t2)
	}
	wg.Wait()
}

type treeIndex struct {
	sync.RWMutex
	tree *btree.BTree
	lg   *zap.Logger
}

func newTreeIndex(lg *zap.Logger) *treeIndex {
	return &treeIndex{
		tree: btree.New(32),
		lg:   lg,
	}
}

func (ti *treeIndex) Put(key []byte, rev revision) {
	keyi := &keyIndex{key: key}

	ti.Lock()
	defer ti.Unlock()
	item := ti.tree.Get(keyi)
	if item == nil {
		keyi.put(ti.lg, rev.main, rev.sub)
		ti.tree.ReplaceOrInsert(keyi)
		return
	}
	okeyi := item.(*keyIndex)
	okeyi.put(ti.lg, rev.main, rev.sub)
}

func (ti *treeIndex) Tombstone(key []byte, rev revision) error {
	keyi := &keyIndex{key: key}

	ti.Lock()
	defer ti.Unlock()
	item := ti.tree.Get(keyi)
	if item == nil {
		return ErrRevisionNotFound
	}

	ki := item.(*keyIndex)
	return ki.tombstone(ti.lg, rev.main, rev.sub)
}

func (ti *treeIndex) Get(key []byte, atRev int64) (modified, created revision, ver int64, err error) {
	keyi := &keyIndex{key: key}
	ti.RLock()
	defer ti.RUnlock()
	if keyi = ti.keyIndex(keyi); keyi == nil {
		return revision{}, revision{}, 0, ErrRevisionNotFound
	}
	return keyi.get(ti.lg, atRev)
}

func (ti *treeIndex) KeyIndex(keyi *keyIndex) *keyIndex {
	ti.RLock()
	defer ti.RUnlock()
	return ti.keyIndex(keyi)
}

func (ti *treeIndex) keyIndex(keyi *keyIndex) *keyIndex {
	if item := ti.tree.Get(keyi); item != nil {
		return item.(*keyIndex)
	}
	return nil
}

func (ti *treeIndex) Compact(rev int64) map[revision]struct{} {
	available := make(map[revision]struct{})
	ti.lg.Info("compact tree index", zap.Int64("revision", rev))
	ti.Lock()
	// 使用 copy on write clone 避免在遍历 treeIndex 的过程中一直持有锁
	// 具体参考: https://github.com/etcd-io/etcd/pull/9511
	clone := ti.tree.Clone()
	ti.Unlock()

	ch <- struct{}{} // 测试使用

	clone.Ascend(func(item btree.Item) bool {
		keyi := item.(*keyIndex)
		fmt.Printf("@@@@@@@@@ %s\n", keyi)
		// Lock is needed here to prevent modification to the keyIndex while
		// compaction is going on or revision added to empty before deletion
		ti.Lock()
		keyi.compact(ti.lg, rev, available) // 每个 key 对应一个 keyIndex，遍历并压缩 keyIndex
		if keyi.isEmpty() {
			item := ti.tree.Delete(keyi) // 如果 keyIndex 为空了，则从 treeIndex 中删除
			if item == nil {
				ti.lg.Panic("failed to delete during compaction")
			}
		}
		// time.Sleep(time.Millisecond * 50) // simulate high load
		ti.Unlock()
		return true
	})
	return available // 通过一个 map 记录 treeIndex 中有效的版本号返回给 boltdb 模块使用
}

type keyIndex struct {
	key         []byte
	modified    revision // the main rev of the last modification 最后一次修改的版本号
	generations []generation
}

// put puts a revision to the keyIndex.
func (ki *keyIndex) put(lg *zap.Logger, main int64, sub int64) {
	rev := revision{main: main, sub: sub}

	if !rev.GreaterThan(ki.modified) {
		lg.Panic(
			"'put' with an unexpected smaller revision",
			zap.Int64("given-revision-main", rev.main),
			zap.Int64("given-revision-sub", rev.sub),
			zap.Int64("modified-revision-main", ki.modified.main),
			zap.Int64("modified-revision-sub", ki.modified.sub),
		)
	}
	if len(ki.generations) == 0 {
		ki.generations = append(ki.generations, generation{})
	}
	g := &ki.generations[len(ki.generations)-1]
	if len(g.revs) == 0 { // create a new key
		g.created = rev
	}
	g.revs = append(g.revs, rev)
	g.ver++
	ki.modified = rev
}

func (ki *keyIndex) findGeneration(rev int64) *generation {
	lastg := len(ki.generations) - 1
	cg := lastg

	for cg >= 0 {
		if len(ki.generations[cg].revs) == 0 {
			cg--
			continue
		}
		g := ki.generations[cg]
		if cg != lastg {
			// 比下一代的最小版本要小，但是大于等于上一代的tomb版本
			if tomb := g.revs[len(g.revs)-1].main; tomb <= rev {
				return nil
			}
		}
		if g.revs[0].main <= rev { // 如果rev大于等于generation中的最小rev
			return &ki.generations[cg]
		}
		cg--
	}
	return nil
}

func (ki *keyIndex) compact(lg *zap.Logger, atRev int64, available map[revision]struct{}) {
	if ki.isEmpty() {
		lg.Panic(
			"'compact' got an unexpected empty keyIndex",
			zap.String("key", string(ki.key)),
		)
	}

	genIdx, revIndex := ki.doCompact(atRev, available)

	g := &ki.generations[genIdx]
	if !g.isEmpty() {
		// remove the previous contents.
		if revIndex != -1 {
			g.revs = g.revs[revIndex:]
		}
		// remove any tombstone
		// NOTE: 如果 g.revs 只剩最后一个 rev，但是 g 不是最新的 generation，则其为 tomb，可以删除
		if len(g.revs) == 1 && genIdx != len(ki.generations)-1 {
			delete(available, g.revs[0])
			genIdx++
		}
	}
	// remove the previous generations.
	ki.generations = ki.generations[genIdx:] // 删除 old generations
}

func (ki *keyIndex) doCompact(atRev int64, available map[revision]struct{}) (genIdx int, revIndex int) {
	// walk until reaching the first revision smaller or equal to "atRev",
	// and add the revision to the available map
	f := func(rev revision) bool {
		if rev.main <= atRev {
			available[rev] = struct{}{}
			return false
		}
		return true
	}

	genIdx, g := 0, &ki.generations[0]
	// find first generation includes atRev or created after atRev
	// 一代代遍历，直到找到大于指定 atRev 的 generation
	for genIdx < len(ki.generations)-1 {
		if tomb := g.revs[len(g.revs)-1].main; tomb > atRev {
			break
		}
		genIdx++
		g = &ki.generations[genIdx]
	}

	revIndex = g.walk(f) // 倒序遍历 revs

	return genIdx, revIndex
}

func (ki *keyIndex) get(lg *zap.Logger, atRev int64) (modified, created revision, ver int64, err error) {
	if ki.isEmpty() {
		lg.Panic(
			"'get' got an unexpected empty keyIndex",
			zap.String("key", string(ki.key)),
		)
	}
	g := ki.findGeneration(atRev)
	if g.isEmpty() {
		return revision{}, revision{}, 0, ErrRevisionNotFound
	}

	n := g.walk(func(rev revision) bool { return rev.main > atRev }) // 返回generation中版本号小于等于atRev的index
	if n != -1 {
		return g.revs[n], g.created, g.ver - int64(len(g.revs)-n-1), nil // 注意这个地方的ver的计算
	}

	return revision{}, revision{}, 0, ErrRevisionNotFound
}

// tombstone puts a revision, pointing to a tombstone, to the keyIndex.
// It also creates a new empty generation in the keyIndex.
// It returns ErrRevisionNotFound when tombstone on an empty generation.
func (ki *keyIndex) tombstone(lg *zap.Logger, main int64, sub int64) error {
	if ki.isEmpty() {
		lg.Panic(
			"'tombstone' got an unexpected empty keyIndex",
			zap.String("key", string(ki.key)),
		)
	}
	if ki.generations[len(ki.generations)-1].isEmpty() {
		return ErrRevisionNotFound
	}
	ki.put(lg, main, sub)
	ki.generations = append(ki.generations, generation{})
	return nil
}

func (ki *keyIndex) isEmpty() bool {
	return len(ki.generations) == 1 && ki.generations[0].isEmpty()
}
func (ki *keyIndex) Less(b btree.Item) bool {
	return bytes.Compare(ki.key, b.(*keyIndex).key) == -1
}

func (ki *keyIndex) String() string {
	var s string
	for _, g := range ki.generations {
		s += g.String()
	}
	return s
}

// generation contains multiple revisions of a key.
type generation struct {
	ver     int64
	created revision // when the generation is created (put in first revision).
	revs    []revision
}

func (g *generation) String() string {
	return fmt.Sprintf("g: created[%d] ver[%d], revs %#v\n", g.created, g.ver, g.revs)
}

type revision struct {
	// main is the main revision of a set of changes that happen atomically.
	main int64 // 其实对应事务ID，全局递增不重复，逻辑时钟

	// sub is the sub revision of a change in a set of changes that happen
	// atomically. Each change has different increasing sub revision in that
	// set.
	sub int64 // 一次事务中的不同的修改操作
}

func (a revision) GreaterThan(b revision) bool {
	if a.main > b.main {
		return true
	}
	if a.main < b.main {
		return false
	}
	return a.sub > b.sub
}

func (g *generation) isEmpty() bool { return g == nil || len(g.revs) == 0 }

// walk walks through the revisions in the generation in descending order.
// It passes the revision to the given function.
// walk returns until: 1. it finishes walking all pairs 2. the function returns false.
// walk returns the position at where it stopped. If it stopped after
// finishing walking, -1 will be returned.
func (g *generation) walk(f func(rev revision) bool) int { // 返回的是位置，索引
	l := len(g.revs)
	for i := range g.revs {
		ok := f(g.revs[l-i-1]) // 倒序遍历
		if !ok {
			return l - i - 1
		}
	}
	return -1
}
