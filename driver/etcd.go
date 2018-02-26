package driver

import (
	"golang.org/x/net/context"
	"github.com/yyzybb537/ketty"
	"github.com/yyzybb537/ketty/driver/etcd"
)

type EtcdDriver struct {
}

func (this *EtcdDriver) DefaultPort() int {
	return 2379
}

func (this *EtcdDriver) Watch(url ketty.Url) (<-chan []Url, func(), error) {
	sess, err := etcd.GetEtcdMgr().GetSession(url.SAddr)
	if err != nil {
		return nil, nil, err
	}

	notify := make(chan []Url, 32)
	cb := func(*etcd.Session, string) error {
		_, values, err := sess.GetChildren(url.Path)
		if err != nil {
			return err
		}

		addrs := []Url{}
		for _, v := range values {
			addr, err := ketty.UrlFromDriverString(v)
			if err != nil {
				return err
			}
			
			addrs = append(addrs, addr)
		}
		notify <- addrs
		return nil
    }

	id := sess.WatchRecursive(url.Path, func(se *etcd.Session, s string){cb(se, s);})
	err = cb(sess, url.Path)
	if err != nil {
		return nil, err
	}
	stop := func(){
		sess.UnWatch(url.Path, id)
		close(notify)
	}
	return notify, stop, nil
}

func (this *EtcdDriver) Register(url, value ketty.Url) (err error) {
	sess, err := etcd.GetEtcdMgr().GetSession(url.SAddr)
	if err != nil {
		return
	}

	saddr := value.ToDriverString()
	err = sess.SetEphemeral(url.Path, saddr, context.Background())
	return 
}

