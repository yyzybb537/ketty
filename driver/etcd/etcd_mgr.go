package etcd

import (
	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"github.com/yyzybb537/ketty/log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const defaultTTL = 12

type WatchCallback func(*Session, string, WatchHandler)
type WatchHandlerT struct {
	id uint64
	cb WatchCallback
	MetaData interface{}
}
type WatchHandler *WatchHandlerT
type WatchMap map[uint64]WatchHandler

type Session struct {
	client  etcd.Client
	kapi    etcd.KeysAPI
	address string
	localIp string

	watcherHandlerIndex uint64
	watchers            map[string]WatchMap
	mutex               sync.RWMutex
}

func EtcdErrorCode(err error) int {
	etcdErr, ok := err.(etcd.Error)
	if !ok {
		return 0
	}
	return etcdErr.Code
}

func (this *Session) GetLocalIp() string {
	if this.localIp != "" {
		return this.localIp
	}

	remoteAddrs := strings.Split(this.address, ",")
	if len(remoteAddrs) == 0 {
		return "127.0.0.1"
	}

	conn, err := net.Dial("tcp", remoteAddrs[0])
	if err != nil {
		log.GetLog().Errorf("net.Dial returns error:%s", err.Error())
		return "127.0.0.1"
	}

	local := conn.LocalAddr()
	return strings.Split(local.String(), ":")[0]
}

func (this *Session) CreateInOrder(dir, value string, outPath *string, ctx context.Context) (err error) {
	//try to create dir
	_, _ = this.kapi.Set(ctx, dir, "", &etcd.SetOptions{Dir: true})
	//delete node in dir when value conflict
	resp, err := this.kapi.Get(ctx, dir, &etcd.GetOptions{Recursive: true})
	if err != nil {
		return
	}
	for _, respNode := range resp.Node.Nodes {
		if value == respNode.Value {
			log.GetLog().Warningf("exist node %s in %s", value, respNode.Key)
			_, err = this.kapi.Delete(ctx, respNode.Key, &etcd.DeleteOptions{PrevValue: value})
			if err != nil {
				return
			}
			log.GetLog().Infof("delete node %s success", respNode.Key)
		}
	}
	//create in order first time
	createInOrder_opt := etcd.CreateInOrderOptions{TTL: time.Second * defaultTTL}
	resp, err = this.kapi.CreateInOrder(ctx, dir, value, &createInOrder_opt)
	if err != nil {
		log.GetLog().Errorln("create etcd in order node error:", err)
		return
	}
	*outPath = resp.Node.Key

	// refresh ttl
	go func() {
		var interval time.Duration = 5
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * interval):
				refresh_opt := etcd.SetOptions{TTL: time.Second * defaultTTL, PrevValue: value, Refresh: true}
				_, err := this.kapi.Set(context.Background(), *outPath, "", &refresh_opt)
				if err != nil {
					log.GetLog().Errorf("Refresh etcd error. address=%s path=%s error=%s", this.address, *outPath, err)
					// set
					resp, err = this.kapi.CreateInOrder(context.Background(), dir, value, &createInOrder_opt)
					if err != nil {
						log.GetLog().Errorln("reset etcd in order node error:", err)
						interval = 2
					}
					*outPath = resp.Node.Key
				} else {
					interval = 5
				}
			}
		}
	}()

	return
}

func (this *Session) SetEphemeral(path, value string, ctx context.Context) (err error) {
	set_opt := etcd.SetOptions{PrevExist: etcd.PrevIgnore, TTL: time.Second * defaultTTL}
	_, err = this.kapi.Set(context.Background(), path, value, &set_opt)
	if err != nil {
		log.GetLog().Errorln("Set etcd ephemeral node error:", err)
		return
	}

	// refresh ttl
	go func() {
		var interval time.Duration = 5
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * interval):
				refresh_opt := etcd.SetOptions{PrevExist: etcd.PrevIgnore, TTL: time.Second * defaultTTL, PrevValue: "1", Refresh: true}
				_, err := this.kapi.Set(context.Background(), path, "", &refresh_opt)
				if err != nil {
					log.GetLog().Errorf("Refresh etcd error. address=%s path=%s error=%s", this.address, path, err.Error())
					// set
					_, err := this.kapi.Set(context.Background(), path, value, &set_opt)
					if err != nil {
						log.GetLog().Errorln("Reset etcd ephemeral node error:", err)
						interval = 2
					}
				} else {
					interval = 5
				}
			}
		}
	}()

	return
}

func (this *Session) Mkdir(path string) (err error) {
	_, err = this.kapi.Set(context.Background(), path, "", &etcd.SetOptions{Dir: true})
	return
}

func (this *Session) GetChildren(path string) (keys []string, values []string, err error) {
	opt := etcd.GetOptions{Recursive: true}
	rsp, err := this.kapi.Get(context.Background(), path, &opt)
	if err != nil {
		return
	}

	keys = []string{}
	values = []string{}
	for _, v := range rsp.Node.Nodes {
		child_node := v.Key[len(path):]
		if strings.HasPrefix(child_node, "/") {
			child_node = child_node[1:]
		}
		keys = append(keys, child_node)
		values = append(values, v.Value)
	}
	return
}

func (this *Session) Get(path string) (value string, err error) {
	opt := etcd.GetOptions{Recursive: false}
	rsp, err := this.kapi.Get(context.Background(), path, &opt)
	if err != nil {
		return
	}

	value = rsp.Node.Value
	return
}

func (this *Session) Del(path string, isDir bool, isRecursive bool) (err error) {
	opt := etcd.DeleteOptions{Recursive: isRecursive, Dir: isDir}
	_, err = this.kapi.Delete(context.Background(), path, &opt)
	return
}

func (this *Session) WatchRecursive(path string, cb WatchCallback) (watcherHandler WatchHandler) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	m, exists := this.watchers[path]
	if !exists {
		// If the path not exists, create it.
		_, _, err := this.GetChildren(path)
		if err != nil && EtcdErrorCode(err) == etcd.ErrorCodeKeyNotFound {
			err = this.Mkdir(path)
			if err != nil {
				log.GetLog().Errorf("Mkdir(%s) error:%v", path, err)
			} else {
				log.GetLog().Infof("Mkdir(%s) in WatchRecursive", path)
			}
		}

		this.watchers[path] = make(WatchMap)
		m, _ = this.watchers[path]
		go func() {
			var index uint64
			watcher := this.kapi.Watcher(path, &etcd.WatcherOptions{Recursive: true, AfterIndex: index})
			for {
				rsp, err := watcher.Next(context.Background())
				if err != nil {
					watcher = this.kapi.Watcher(path, &etcd.WatcherOptions{Recursive: true, AfterIndex: index})
					time.Sleep(time.Second * 6)
				} else {
					index = rsp.Index
				}

				log.GetLog().Infof("Watcher begin trigger path=%s. callback len(m)=%d. Index=%d. Err:%+v", path, len(m), index, err)

				for _, handler := range m {
					handler.cb(this, path, handler)
				}

				log.GetLog().Infof("Watcher end trigger path=%s. callback len(m)=%d", path, len(m))
			}
		}()
	}

	watcherHandlerIndex := atomic.AddUint64((*uint64)(&this.watcherHandlerIndex), 1)
	if watcherHandlerIndex == 0 {
		watcherHandlerIndex = atomic.AddUint64((*uint64)(&this.watcherHandlerIndex), 1)
	}
	watcherHandler = &WatchHandlerT{
		id : watcherHandlerIndex,
		cb : cb,
	}
	m[watcherHandlerIndex] = watcherHandler
	return
}

func (this *Session) UnWatch(path string, watcherHandler WatchHandler) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	m, exists := this.watchers[path]
	if !exists {
		return
	}

	delete(m, watcherHandler.id)
}

type etcdMgr struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
}

var etcd_mgr *etcdMgr = &etcdMgr{
	sessions: make(map[string]*Session),
}

func GetEtcdMgr() *etcdMgr {
	return etcd_mgr
}

func (this *etcdMgr) GetSession(address string) (sess *Session, err error) {
	sess, exists := this.getSession(address)
	if exists {
		return
	}

	s := &Session{
		watcherHandlerIndex: 0,
		watchers:            make(map[string]WatchMap),
	}

	var addrs []string = strings.Split(address, ",")
	for index, _ := range addrs {
		addrs[index] = "http://" + addrs[index]
	}

	cfg := etcd.Config{
		Endpoints:               addrs,
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second * 12,
	}

	s.client, err = etcd.New(cfg)
	if err != nil {
		return
	}
	s.kapi = etcd.NewKeysAPI(s.client)
	s.address = address

	sess = this.setSession(address, s)
	return
}

func (this *etcdMgr) getSession(address string) (*Session, bool) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	sess, exists := this.sessions[address]
	return sess, exists
}

func (this *etcdMgr) setSession(address string, sess *Session) *Session {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if s, exists := this.sessions[address]; exists {
		return s
	}

	this.sessions[address] = sess
	return sess
}

