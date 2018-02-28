package driver

import (
	"golang.org/x/net/context"
	"github.com/yyzybb537/ketty"
	"github.com/yyzybb537/ketty/driver/etcd"
)

type EtcdDriver struct {
}

func init() {
	ketty.RegDriver("etcd", new(EtcdDriver))
}

func (this *EtcdDriver) DefaultPort() int {
	return 2379
}

func (this *EtcdDriver) Watch(url ketty.Url) (up, down <-chan []ketty.Url, stop func(), err error) {
	sess, err := etcd.GetEtcdMgr().GetSession(url.SAddr)
	if err != nil {
		return
	}

	upC := make(chan []ketty.Url, 32)
	downC := make(chan []ketty.Url, 32)
	up = upC
	down = downC
	cb := func(s *etcd.Session, str string, h etcd.WatchHandler) error {
		values, _, err := sess.GetChildren(url.Path)
		if err != nil {
			return err
		}

		upAddrs := []ketty.Url{}
		downAddrs := []ketty.Url{}
		lastNodes, ok := h.MetaData.(map[string]bool)
		if !ok {
			lastNodes = map[string]bool{}
			h.MetaData = lastNodes
		}
		newNodes := map[string]bool{}
		for _, v := range values {
			newNodes[v] = true

			// up
			if _, exists := lastNodes[v]; !exists {
				addr, err := ketty.UrlFromDriverString(v)
				if err == nil {
					upAddrs = append(upAddrs, addr)
					ketty.GetLog().Debugf("up url %s", addr.ToString())
				} else {
					ketty.GetLog().Warningf("unkown up url %s", addr.ToString())
                }
			}
		}

		// down
		for v, _ := range lastNodes {
			if _, exists := newNodes[v]; !exists {
				addr, err := ketty.UrlFromDriverString(v)
				if err == nil {
					downAddrs = append(downAddrs, addr)
					ketty.GetLog().Debugf("down url %s", addr.ToString())
				} else {
					ketty.GetLog().Warningf("unkown down url %s", addr.ToString())
                }
            }
        }

		h.MetaData = newNodes

		upC <- upAddrs
		downC <- downAddrs
		return nil
    }

	handler := sess.WatchRecursive(url.Path, func(se *etcd.Session, s string, h etcd.WatchHandler){cb(se, s, h);})
	stop = func(){
		sess.UnWatch(url.Path, handler)
		close(upC)
		close(downC)
	}

	err = cb(sess, url.Path, handler)
	if err != nil {
		stop()
		return
	}
	return
}

func (this *EtcdDriver) Register(url, value ketty.Url) (err error) {
	sess, err := etcd.GetEtcdMgr().GetSession(url.SAddr)
	if err != nil {
		return
	}

	saddr := value.ToDriverString()
	err = sess.SetEphemeral(url.Path + "/" + saddr, "1", context.Background())
	return 
}

