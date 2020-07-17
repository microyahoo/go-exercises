package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	peerURL1 := "http://127.0.0.1:8081"
	peerURL2 := "http://127.0.0.1:8082"

	id1 := getPid(peerURL1)
	id2 := getPid(peerURL2)

	//开启节点1
	tr1 := &Transport{}
	tr1.Start(int64(id1))
	go func() {
		err := http.ListenAndServe(":8081", tr1.Handler())
		log.Fatal(err)
	}()

	//开启节点2
	tr2 := &Transport{}
	tr2.Start(int64(id2))
	go func() {
		err := http.ListenAndServe(":8082", tr2.Handler())
		log.Fatal(err)
	}()

	// time.Sleep(time.Second * 3)

	//节点1添加节点2
	tr1.AddPeer(int64(id2), peerURL2)

	//节点2添加节点1
	tr2.AddPeer(int64(id1), peerURL1)

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				peers := tr1.GetPeers()
				for i := range peers {
					peers[i].send(&Message{MsgType: msgTypeProp, MsgBody: "Hello, I am tr1"})
				}
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				peers := tr2.GetPeers()
				for i := range peers {
					peers[i].send(&Message{MsgType: msgTypeProp, MsgBody: "Hello, I am tr2"})
				}
			}
		}
	}()

	time.Sleep(time.Minute * 10)

	tr1.Stop()
	tr2.Stop()
}

type Transport struct {
	ClusterID int64
	ID        int64 // local member ID  当前节点自己的ID
	streamRt  http.RoundTripper
	mu        sync.RWMutex
	peers     map[int64]*peer // peers map
}

func (tr *Transport) Start(id int64) error {
	tr.ID = id
	tr.streamRt = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
			// value taken from http.DefaultTransport
			KeepAlive: 30 * time.Second,
		}).Dial,
	}
	tr.peers = make(map[int64]*peer)

	return nil
}

func (tr *Transport) GetPeers() (result []*peer) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	for k := range tr.peers {
		result = append(result, tr.peers[k])
	}

	return
}

func (tr *Transport) AddPeer(id int64, peerURL string) {
	tr.mu.RLock()
	if _, ok := tr.peers[id]; ok {
		tr.mu.RUnlock()
		return
	}

	tr.mu.RUnlock()

	peer := startPeer(tr, peerURL, tr.ID, id)

	tr.mu.Lock()
	tr.peers[id] = peer
	tr.mu.Unlock()
}

func (tr *Transport) Handler() http.Handler {
	streamHandler := newStreamHandler(tr, tr.ID, tr.ClusterID)
	mux := http.NewServeMux()
	mux.Handle("/raft/stream"+"/", streamHandler)
	return mux
}

type streamHandler struct {
	tr  *Transport //关联的rafthttp.Transport实例
	id  int64      //当前节点ID
	cid int64      //当前集群ID
}

func newStreamHandler(tr *Transport, id, cid int64) http.Handler {
	return &streamHandler{
		tr:  tr,
		id:  id,
		cid: cid,
	}
}

func (h *streamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//请求参数校验，如Method是否是GET，检验集群ID
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.Header.Get("PeerID")
	if id == "" {
		w.Header().Set("PeerID", "Must")
		http.Error(w, "PeerID is not allow empty", http.StatusMethodNotAllowed)
		return
	}

	pid, _ := strconv.ParseUint(id, 10, 64)

	p, ok := h.tr.peers[int64(pid)]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		time.Sleep(time.Second * 5)
		return
	}

	w.WriteHeader(http.StatusOK) //返回状态码 200

	w.(http.Flusher).Flush() //调用Flush()方法将响应数据发送到对端节点

	c := newCloseNotifier()
	conn := &outgoingConn{ //创建outgoingConn实例
		Writer:  w,
		Flusher: w.(http.Flusher),
		localID: h.tr.ID,
		Closer:  c,
		peerID:  h.id,
	}
	p.attachOutgoingConn(conn) //建立连接,将outgoingConn实例与对应的streamWriter实例绑定
	<-c.closeNotify()
}

type closeNotifier struct {
	done chan struct{}
}

func newCloseNotifier() *closeNotifier {
	return &closeNotifier{
		done: make(chan struct{}),
	}
}

func (n *closeNotifier) Close() error {
	close(n.done)
	return nil
}

func (n *closeNotifier) closeNotify() <-chan struct{} { return n.done }

func (tr *Transport) Stop() {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	for _, v := range tr.peers {
		v.stop()
	}
	tr.peers = nil
}

const (
	msgTypeHeartbeat = "01" //心跳
	msgTypeProp      = "02" //prop消息
)

type Message struct {
	MsgType string
	MsgBody string
}

type peer struct {
	localID int64 //当前节点ID
	// id of the remote raft peer node
	id int64 //该peer实例对应的节点ID，对端ID

	writer *streamWriter //负责向Stream消息通道中写消息

	msgAppReader *streamReader //负责从Stream消息通道中读消息

	msgc  chan *Message
	stopc chan struct{}
}

func startPeer(t *Transport, peerURL string, localID, peerID int64) *peer {
	pr := &peer{
		localID:      localID,
		id:           peerID,
		writer:       newStreamWriter(localID, peerID),
		msgAppReader: newStreamReader(localID, peerID, peerURL, t),
		msgc:         make(chan *Message, 20),
		stopc:        make(chan struct{}),
	}

	go func() {
		for msg := range pr.msgc {
			select {
			case pr.writer.msgc <- msg:
			default:
				log.Printf("write to writer error msg is %v", msg)
			}
		}
	}()

	return pr
}

func (pr *peer) stop() {
	close(pr.stopc)
}

func (pr *peer) send(msg *Message) bool {
	select {
	case pr.msgc <- msg:
		return true
	default:
		return false
	}
}

func (pr *peer) attachOutgoingConn(conn *outgoingConn) {
	select {
	case pr.writer.connc <- conn:

	default:
		log.Printf("attachOutgoingConn error")
	}
}

type streamWriter struct {
	localID int64 //本端的ID
	peerID  int64 //对端节点的ID

	closer io.Closer //负责关闭底层的长连接

	mu sync.Mutex

	enc   *messageEncoder
	msgc  chan *Message      //Peer会将待发送的消息写入到该通道，streamWriter则从该通道中读取消息并发送出去
	connc chan *outgoingConn //通过该通道获取当前streamWriter实例关联的底层网络连接，  outgoingConn其实是对网络连接的一层封装，其中记录了当前连接使用的协议版本，以及用于关闭连接的Flusher和Closer等信息。
	stopc chan struct{}
}

func newStreamWriter(localID, peerID int64) *streamWriter {
	sw := &streamWriter{
		localID: localID,
		peerID:  peerID,
		msgc:    make(chan *Message, 20),
		connc:   make(chan *outgoingConn),
		stopc:   make(chan struct{}),
	}

	go sw.run()
	return sw
}

func (sw *streamWriter) writec() chan<- *Message {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.msgc
}

func (sw *streamWriter) run() {
	var (
		heartbeatC <-chan time.Time

		flusher http.Flusher //负责刷新底层连接，将数据真正发送出去

		msgc chan *Message
	)

	tickc := time.NewTicker(time.Second * 7) //发送心跳的定时器
	defer tickc.Stop()

	for {
		select {
		case msg := <-msgc:
			err := sw.enc.encode(msg)
			if err != nil {
				log.Printf("Send to peer peerID is %d fail, error is: %v", sw.peerID, err)
			} else {
				flusher.Flush()
				log.Printf("Send to peer peerID is %d success, MsgType is: %s, MsgBody is: %s",
					sw.peerID, msg.MsgType, msg.MsgBody)
			}

		case <-heartbeatC: //向对端发送心跳消息
			err := sw.enc.encode(&Message{
				MsgType: msgTypeHeartbeat,
				MsgBody: time.Now().Format("2006-01-02 15:04:05"),
			})
			if err != nil {
				log.Printf("Send to peer heartbeat data fail, peerID is %d, error is %v",
					sw.peerID, err)
			} else {
				flusher.Flush()
				log.Printf("Send to peer heartbeat data success peerID is %d ", sw.peerID)
			}
		case conn := <-sw.connc:
			sw.enc = &messageEncoder{w: conn.Writer}
			flusher = conn.Flusher
			sw.closer = conn.Closer

			heartbeatC, msgc = tickc.C, sw.msgc

		case <-sw.stopc:
			log.Println("msgWriter stop!")
			sw.closer.Close()
			return
		}
	}
}

func (sw *streamWriter) stop() {
	close(sw.stopc)
}

type streamReader struct {
	localID int64
	peerID  int64      //对端节点的ID
	tr      *Transport //关联的rafthttp.Transport实例
	peerURL string     //对端URL

	mu     sync.Mutex
	closer io.Closer //负责关闭底层的长连接

	done chan struct{}
}

func newStreamReader(localID, peerID int64, peerURL string, tr *Transport) *streamReader {
	sr := &streamReader{
		localID: localID,
		peerID:  peerID,
		peerURL: peerURL,
		tr:      tr,
		done:    make(chan struct{}),
	}
	go sr.run()

	return sr
}

func (sr *streamReader) run() {
	// time.Sleep(time.Second * 5)
	for {
		readColser, err := sr.dial()
		if err != nil {
			log.Printf("Dial peer error, peerID is %d, err is: %v", sr.peerID, err)
			time.Sleep(time.Second * 10)
			continue
		}
		sr.closer = readColser

		err = sr.decodeLoop(readColser)
		if err != nil {
			log.Printf("DecodeLoop error, peerID is %d, error is %v", sr.peerID, err)
		}
		sr.closer.Close()
	}
}

func (sr *streamReader) dial() (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", sr.peerURL+"/raft/stream/dial", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PeerID", fmt.Sprintf("%d", sr.localID))

	resp, err := sr.tr.streamRt.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (sr *streamReader) decodeLoop(rc io.ReadCloser) error {
	dec := &messageDecoder{rc}
	for {
		msg, err := dec.decode()
		if err != nil {
			log.Printf("\t**Read decodeLoop error, peerID is %d, err is %v", sr.peerID, err)
			continue
		}

		log.Printf("\t**Read from peer MsgType is %s, MsgBody is %s", msg.MsgType, msg.MsgBody)
	}
}

type outgoingConn struct {
	io.Writer
	http.Flusher
	io.Closer
	localID int64
	peerID  int64
}

type messageEncoder struct {
	w io.Writer
}

func (m *messageEncoder) encode(msg *Message) error {
	byts, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	dataLen := make([]byte, 8)
	binary.BigEndian.PutUint64(dataLen, uint64(len(byts)))

	sendData := append(dataLen, byts...)
	_, err = m.w.Write(sendData)
	if err != nil {
		return err
	}

	return nil
}

type messageDecoder struct {
	r io.Reader
}

func (dec *messageDecoder) decode() (*Message, error) {
	var m Message
	var l uint64
	if err := binary.Read(dec.r, binary.BigEndian, &l); err != nil {
		return nil, err
	}

	buf := make([]byte, int(l))
	if _, err := io.ReadFull(dec.r, buf); err != nil {
		return &m, err
	}
	err := json.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func getPid(purl string) int64 {
	index := strings.LastIndex(purl, ":")
	if index > 0 {
		id, err := strconv.ParseInt(purl[index+1:], 10, 64)
		if err != nil {
			println(err)
		}
		return id
	}

	return 0
}
