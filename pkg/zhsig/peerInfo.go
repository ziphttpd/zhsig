package zhsig

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"

	"github.com/xorvercom/util/pkg/json"
)

// PeerInfo は証明書の情報です。
type PeerInfo struct {
	Host         string
	CommonName   string
	Country      []string
	Organization []string
}

// IsOrganization は EV, OV 証明書であった場合に真を返します。
func IsOrganization(host string) bool {
	p, e := GetPeerInfo(host)
	if e != nil {
		return false
	}
	return p.IsOrganization()
}

// IsOrganization は EV, OV 証明書であった場合に真を返します。
func (c *PeerInfo) IsOrganization() bool {
	return c.Organization != nil && len(c.Organization) > 0
}

func (c *PeerInfo) String() string {
	//return fmt.Sprintf("%s \"%s\" (%+v, %+v)", c.Host, c.CommonName, c.Country, c.Organization)
	return c.JSON().Text()
}

// DecodePeerInfo はjsonから復号します。
func DecodePeerInfo(j json.ElemObject) *PeerInfo {
	c := &PeerInfo{}
	if s, ok := json.QueryElemString(j, "host"); ok {
		c.Host = s.Text()
	}
	if s, ok := json.QueryElemString(j, "commonname"); ok {
		c.CommonName = s.Text()
	}
	if a, ok := json.QueryElemArray(j, "country"); ok {
		arr := []string{}
		for i := 0; i < a.Size(); i++ {
			if s, ok := a.Child(i).AsString(); ok {
				arr = append(arr, s.Text())
			}
		}
		c.Country = arr
	}
	if a, ok := json.QueryElemArray(j, "organization"); ok {
		arr := []string{}
		for i := 0; i < a.Size(); i++ {
			if s, ok := a.Child(i).AsString(); ok {
				arr = append(arr, s.Text())
			}
		}
		c.Organization = arr
	}
	return c
}

// JSON はjsonに変換します。
func (c *PeerInfo) JSON() json.ElemObject {
	j := json.NewElemObject()
	j.Put("host", json.NewElemString(c.Host))
	j.Put("commonname", json.NewElemString(c.CommonName))
	arr := json.NewElemArray()
	if c.Country != nil {
		for _, str := range c.Country {
			arr.Append(json.NewElemString(str))
		}
	}
	j.Put("country", arr)
	arr = json.NewElemArray()
	if c.Organization != nil {
		for _, str := range c.Organization {
			arr.Append(json.NewElemString(str))
		}
	}
	j.Put("organization", arr)
	return j
}

// GetPeerInfo は証明書情報を返します。
func GetPeerInfo(host string) (*PeerInfo, error) {
	conf := &tls.Config{}
	conn, err := tls.Dial("tcp", host+":443", conf)
	if err != nil {
		return nil, err
	}

	certs := conn.ConnectionState().PeerCertificates
	defer conn.Close()

	c := filter(host, certs)
	if c == nil {
		return nil, fmt.Errorf("unmatch: %s", host)
	}

	return &PeerInfo{
		Host:         host,
		CommonName:   c.Issuer.CommonName,
		Country:      c.Subject.Country,
		Organization: c.Subject.Organization,
	}, nil
}

func filter(host string, certs []*x509.Certificate) *x509.Certificate {
	for _, cert := range certs {
		if chackName(host, cert.Subject.CommonName) {
			return cert
		}
		for _, name := range cert.DNSNames {
			if chackName(host, name) {
				return cert
			}
		}
	}
	return nil
}
func chackName(host, name string) bool {
	return host == name || (strings.HasPrefix(name, "*.") && strings.HasSuffix(host, name[1:]))
}
