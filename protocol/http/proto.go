package http_proto

import (
	"strings"
	"fmt"
	P "github.com/yyzybb537/ketty/protocol"
)

type Proto struct {
	DefaultMethod    string
	DefaultMarshaler string
	DefaultTransport string
}

func ParseProto(protocol string) (*Proto, error) {
	opt := &Proto{}
	ss := strings.Split(protocol, ".")
	if len(ss) > 1 {
		ss = ss[1:]
		for _, s := range ss {
			s = strings.ToLower(s)
			if s == "get" || s == "post" {
				if opt.DefaultMethod != "" {
					return nil, fmt.Errorf("Http protocol has too many methods. protocol=%s", protocol)
				}
				opt.DefaultMethod = s
			}

			if P.MgrMarshaler.Get(s) != nil {
				if opt.DefaultMarshaler != "" {
					return nil, fmt.Errorf("Http protocol has too many marshaler. protocol=%s", protocol)
				}
				opt.DefaultMarshaler = s
            }

			if MgrTransport.Get(s) != nil {
				if opt.DefaultTransport != "" {
					return nil, fmt.Errorf("Http protocol has too many transport. protocol=%s", protocol)
				}
				opt.DefaultTransport = s
            }
		}
	}

	// default value
	if opt.DefaultMethod == "" {
		opt.DefaultMethod = "post"
	}

	if opt.DefaultMarshaler == "" {
		opt.DefaultMarshaler = "pb"
    }

	if opt.DefaultTransport == "" {
		opt.DefaultTransport = "body"
    }

	return opt, nil
}
