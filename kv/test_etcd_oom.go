// https://github.com/etcd-io/etcd/issues/12342
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	_ "net/http/pprof"

	go_bytesize "github.com/inhies/go-bytesize"
	go_metrics "github.com/rcrowley/go-metrics"
)

// txBuffer handles functionality shared between txWriteBuffer and txReadBuffer.
type txBuffer struct {
	buckets map[string]*bucketBuffer
}

func (txb *txBuffer) reset() {
	for k, v := range txb.buckets {
		if v.used == 0 {
			// demote
			delete(txb.buckets, k)
		}
		v.used = 0
	}
}

// txWriteBuffer buffers writes of pending updates that have not yet committed.
type txWriteBuffer struct {
	txBuffer
	seq bool
}

func (txw *txWriteBuffer) put(bucket, k, v []byte) {
	txw.seq = false
	txw.putSeq(bucket, k, v)
}

func (txw *txWriteBuffer) putSeq(bucket, k, v []byte) {
	b, ok := txw.buckets[string(bucket)]
	if !ok {
		b = newBucketBuffer()
		txw.buckets[string(bucket)] = b
	}
	b.add(k, v)
}

func (txw *txWriteBuffer) writeback(txr *txReadBuffer) {
	for k, wb := range txw.buckets {
		rb, ok := txr.buckets[k]
		if !ok {
			delete(txw.buckets, k)
			txr.buckets[k] = wb
			continue
		}
		if !txw.seq && wb.used > 1 {
			// assume no duplicate keys
			sort.Sort(wb)
		}
		rb.merge(wb)
	}
	txw.reset()
}

// txReadBuffer accesses buffered updates.
type txReadBuffer struct{ txBuffer }

func (txr *txReadBuffer) Range(bucketName, key, endKey []byte, limit int64) ([][]byte, [][]byte) {
	if b := txr.buckets[string(bucketName)]; b != nil {
		return b.Range(key, endKey, limit)
	}
	return nil, nil
}

func (txr *txReadBuffer) ForEach(bucketName []byte, visitor func(k, v []byte) error) error {
	if b := txr.buckets[string(bucketName)]; b != nil {
		return b.ForEach(visitor)
	}
	return nil
}

// unsafeCopy returns a copy of txReadBuffer, caller should acquire backend.readTx.RLock()
func (txr *txReadBuffer) unsafeCopy() txReadBuffer {
	txrCopy := txReadBuffer{
		txBuffer: txBuffer{
			buckets: make(map[string]*bucketBuffer, len(txr.txBuffer.buckets)),
		},
	}
	for bucketName, bucket := range txr.txBuffer.buckets {
		txrCopy.txBuffer.buckets[bucketName] = bucket.Copy()
	}
	return txrCopy
}

type kv struct {
	key []byte
	val []byte
}

// bucketBuffer buffers key-value pairs that are pending commit.
type bucketBuffer struct {
	buf []kv
	// used tracks number of elements in use so buf can be reused without reallocation.
	used int
}

func newBucketBuffer() *bucketBuffer {
	return &bucketBuffer{buf: make([]kv, 512), used: 0}
}

func (bb *bucketBuffer) Range(key, endKey []byte, limit int64) (keys [][]byte, vals [][]byte) {
	f := func(i int) bool { return bytes.Compare(bb.buf[i].key, key) >= 0 }
	idx := sort.Search(bb.used, f)
	if idx < 0 {
		return nil, nil
	}
	if len(endKey) == 0 {
		if bytes.Equal(key, bb.buf[idx].key) {
			keys = append(keys, bb.buf[idx].key)
			vals = append(vals, bb.buf[idx].val)
		}
		return keys, vals
	}
	if bytes.Compare(endKey, bb.buf[idx].key) <= 0 {
		return nil, nil
	}
	for i := idx; i < bb.used && int64(len(keys)) < limit; i++ {
		if bytes.Compare(endKey, bb.buf[i].key) <= 0 {
			break
		}
		keys = append(keys, bb.buf[i].key)
		vals = append(vals, bb.buf[i].val)
	}
	return keys, vals
}

func (bb *bucketBuffer) ForEach(visitor func(k, v []byte) error) error {
	for i := 0; i < bb.used; i++ {
		if err := visitor(bb.buf[i].key, bb.buf[i].val); err != nil {
			return err
		}
	}
	return nil
}

func (bb *bucketBuffer) add(k, v []byte) {
	bb.buf[bb.used].key, bb.buf[bb.used].val = k, v
	bb.used++
	if bb.used == len(bb.buf) {
		buf := make([]kv, (3*len(bb.buf))/2)
		copy(buf, bb.buf)
		bb.buf = buf
	}
}

// merge merges data from bbsrc into bb.
func (bb *bucketBuffer) merge(bbsrc *bucketBuffer) {
	for i := 0; i < bbsrc.used; i++ {
		bb.add(bbsrc.buf[i].key, bbsrc.buf[i].val)
	}
	if bb.used == bbsrc.used {
		return
	}
	if bytes.Compare(bb.buf[(bb.used-bbsrc.used)-1].key, bbsrc.buf[0].key) < 0 {
		return
	}

	sort.Stable(bb)

	// remove duplicates, using only newest update
	widx := 0
	for ridx := 1; ridx < bb.used; ridx++ {
		if !bytes.Equal(bb.buf[ridx].key, bb.buf[widx].key) {
			widx++
		}
		bb.buf[widx] = bb.buf[ridx]
	}
	bb.used = widx + 1
}

func (bb *bucketBuffer) Len() int { return bb.used }
func (bb *bucketBuffer) Less(i, j int) bool {
	return bytes.Compare(bb.buf[i].key, bb.buf[j].key) < 0
}
func (bb *bucketBuffer) Swap(i, j int) { bb.buf[i], bb.buf[j] = bb.buf[j], bb.buf[i] }

func (bb *bucketBuffer) Copy() *bucketBuffer {
	bbCopy := bucketBuffer{
		buf:  make([]kv, len(bb.buf)),
		used: bb.used,
	}
	copy(bbCopy.buf, bb.buf)
	return &bbCopy
}

func main() {

	cc := make(chan []int, 1)

	go func() {
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()

	go func() {
		registry := go_metrics.NewRegistry()
		go_metrics.RegisterRuntimeMemStats(registry)
		go_metrics.RegisterDebugGCStats(registry)

		go go_metrics.CaptureRuntimeMemStats(registry, 10*time.Second)
		go go_metrics.CaptureDebugGCStats(registry, 10*time.Second)

		t := time.NewTimer(10 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				for kk, vv := range registry.GetAll() {
					for k, v := range vv {
						var ff float64
						switch vvv := v.(type) {
						case int64:
							ff = float64(vvv)
							break
						case float64:
							ff = float64(vvv)
							break
						default:
							fmt.Printf("not handle this value type, kk:%s, k: %s, v: %s", kk, k, vvv)
						}

						bb := go_bytesize.New(ff)
						fmt.Printf("%s, %s, %s \n", kk, k, bb.String())
					}
				}
				fmt.Println("---------------------------------------")
				t.Reset(10 * time.Second)
			}
		}
	}()

	txReadBuf := &txReadBuffer{
		txBuffer{
			buckets: make(map[string]*bucketBuffer),
		},
	}
	for i := 0; i < 3000; i++ {
		txReadBuf.buckets[string(i)] = newBucketBuffer()
		for j := 0; j < 3000; j++ {
			txReadBuf.buckets[string(i)].add([]byte("aaaaa"+string(i)+"_"+string(j)), []byte("bbbbb"+string(i)+"_"+string(j)))
		}
	}

	wg := sync.WaitGroup{}

	allocCopy := func() {
		for k := 0; k < 100; k++ {
			wg.Add(1)
			go func(kk int) {
				defer wg.Done()
				a1 := txReadBuf.unsafeCopy()
				fmt.Printf("golang thread %d, len: %d\n", kk, len(a1.buckets))
			}(k)

			time.Sleep(200 * time.Millisecond)
		}
		return
	}

	allocCopy()

	wg.Wait()
	fmt.Printf("golang thread run finish!!! \n")
	fmt.Printf("single run start!!! \n")

	for k := 0; k < 10000; k++ {
		a1 := txReadBuf.unsafeCopy()
		fmt.Printf("single run time number %d, len: %d\n", k, len(a1.buckets))
		time.Sleep(200 * time.Millisecond)
	}

	select {
	case <-cc:
	}
}
