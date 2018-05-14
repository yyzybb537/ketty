package driver

import (
	"golang.org/x/net/context"
	U "github.com/yyzybb537/ketty/url"
	"github.com/yyzybb537/ketty/log"
	"github.com/yyzybb537/ketty/driver/etcd"
)

type EtcdDriver struct {
}

func init() {
	RegDriver("etcd", new(EtcdDriver))
	U.RegDefaultPort("etcd", 2379)
}

func (this *EtcdDriver) Watch(url U.Url) (up, down <-chan []U.Url, stop func(), err error) {
	sess, err := etcd.GetEtcdMgr().GetSession(url.SAddr)
	if err != nil {
		return
	}

	upC := make(chan []U.Url, 32)
	downC := make(chan []U.Url, 32)
	up = upC
	down = downC
	cb := func(s *etcd.Session, str string, h etcd.WatchHandler) error {
		values, _, err := sess.GetChildren(url.Path)
		if err != nil {
			return err
		}

		upAddrs := []U.Url{}
		downAddrs := []U.Url{}
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
				addr, err := U.UrlFromDriverString(v)
				if err == nil {
					upAddrs = append(upAddrs, addr)
					log.GetLog().Debugf("up url %s", addr.ToString())
				} else {
					log.GetLog().Warningf("unkown up url %s", addr.ToString())
                }
			}
		}

		// down
		for v, _ := range lastNodes {
			if _, exists := newNodes[v]; !exists {
				addr, err := U.UrlFromDriverString(v)
				if err == nil {
					downAddrs = append(downAddrs, addr)
					log.GetLog().Debugf("down url %s", addr.ToString())
				} else {
					log.GetLog().Warningf("unkown down url %s", addr.ToString())
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

func (this *EtcdDriver) Register(url, value U.Url) (err error) {
	sess, err := etcd.GetEtcdMgr().GetSession(url.SAddr)
	if err != nil {
		return
	}

	localIp := sess.GetLocalIp()
	saddr, err := this.toDriverString(localIp, value)
	if err != nil {
		return
	}

	err = sess.SetEphemeral(url.Path + "/" + saddr, "1", context.Background())
	return 
}

func (this *EtcdDriver) toDriverString(localIp string, value U.Url) (s string, err error) {
	saddrs := value.GetAddrs()
	newSaddrs := []string{}
	for _, saddr := range saddrs {
		var addr U.Addr
		addr, err = U.AddrFromString(saddr, value.GetMainProtocol())
		if err != nil {
			return
		}
		switch addr.Host {
		case "":
			fallthrough
		case "0.0.0.0":
			addr.Host = localIp
		}
		newSaddrs = append(newSaddrs, addr.ToString())
	}
	value.SetAddrs(newSaddrs)
	return value.ToDriverString(), nil
}
