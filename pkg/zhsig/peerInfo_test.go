package zhsig

import (
	"testing"
)

func TestGetPeerInfo(t *testing.T) {
	host := "notexists.co.jp"
	p, e := GetPeerInfo(host)
	if e == nil {
		t.Error("must error")
	}
	if IsOrganization(host) {
		t.Error(p)
	}
	t.Log(p)

	host = "www.toyota.co.jp"
	p, e = GetPeerInfo(host)
	if e != nil {
		t.Error(e)
	}
	if false == IsOrganization(host) {
		t.Error(p)
	}
	t.Log(p)
	t.Log(DecodePeerInfo(p.JSON()))

	host = "ziphttpd.com"
	p, e = GetPeerInfo(host)
	if e != nil {
		t.Error(e)
	}
	if IsOrganization(host) {
		t.Error(p)
	}
	t.Log(DecodePeerInfo(p.JSON()))
	t.Log(p)
}
