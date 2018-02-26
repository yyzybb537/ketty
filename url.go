package ketty

import (
	"fmt"
	"strings"
	"strconv"
	"regexp"
)

var re *regexp.Regexp

func init() {
	re = regexp.MustCompile("(\\w+)://([^/]+)(.*)")
}

type Addr struct {
	Host string
	Port int
	MetaData interface{}
}

func AddrFromString(hostport string, protocol string) (a Addr, err error) {
	dPort := 0
	proto, err := GetProtocol(protocol)
	if err != nil {
		var driver Driver
		driver, err = GetDriver(protocol)
		if err != nil {
			return 
		}
		dPort = driver.DefaultPort()
	} else {
		dPort = proto.DefaultPort()
    }

	sIndex := strings.Index(hostport, ":")
	if sIndex >= 0 {
		a.Host = hostport[:sIndex]
		a.Port, err = strconv.Atoi(hostport[sIndex+1:])
	} else {
		a.Host = hostport
		a.Port = dPort
    }
	return
}

func FormatAddr(hostport string, protocol string) string {
	addr, err := AddrFromString(hostport, protocol)
	if err != nil {
		return hostport
	}
	return addr.ToString()
}

func (this *Addr) ToString() string {
	return fmt.Sprintf("%s:%d", this.Host, this.Port)
}

type Url struct {
	Protocol string
	SAddr	 string
	//Addrs    []Addr
	Path     string
	MetaData interface{}
}

// @url:   protocol://[balance@]ip[:port][,ip[:port]]/balancePath
func UrlFromString(url string) (Url, error) {
	u := Url{}
	if url == "" {
		return u, nil
	}

	matchs := re.FindStringSubmatch(url)
	if len(matchs) != 4 {
		return u, fmt.Errorf("Url regex parse error. url=%s. matchs=%d", url, len(matchs))
	}

	u.Protocol = matchs[1]
	u.SAddr = matchs[2]
	u.Path = matchs[3]
	//saddrs := strings.Split(u.SAddr, ",")
	//for _, saddr := range saddrs {
		//addr, err := AddrFromString(saddr, u.Protocol)
		//if err != nil {
			//return nil, err
		//}
		//u.Addrs = append(u.Addrs, addr)
    //}
	return u, nil
}

func UrlFromDriverString(url string) (Url, error) {
	return UrlFromString(strings.Replace(url, "|", "/", -1))
}

func (this Url) GetAddrs() []string {
	return strings.Split(this.SAddr, ",")
}

func (this Url) ToString() string {
	if this.IsEmpty() {
		return ""
    }

	if this.Path == "" {
		return fmt.Sprintf("%s://%s", this.Protocol, this.SAddr)
	} else {
		return fmt.Sprintf("%s://%s/%s", this.Protocol, this.SAddr, this.Path)
    }
}

func (this Url) ToDriverString() string {
	return strings.Replace(this.ToString(), "/", "|", -1)
}

func (this Url) IsEmpty() bool {
	return this.Protocol == ""
}

