package zhsig

import (
	"github.com/xorvercom/util/pkg/json"
)

const (
	groupTitle       = "title"
	groupDescription = "description"
	groupDocs        = "docs"
)

// GroupInfo はグループの表示情報です。
type GroupInfo struct {
	// Title はドキュメントグループの表示名です
	Title string
	// Description はドキュメントグループの内容の説明です
	Description string
	// Docs はドキュメントのマップです
	Docs map[string]*DocInfo
}

// DocNames はドキュメントの識別子を返します。
func (g *GroupInfo) DocNames() []string {
	ret := []string{}
	for key := range g.Docs {
		ret = append(ret, key)
	}
	return ret
}

// JSON はjsonエレメントを返します。
func (g *GroupInfo) JSON() json.ElemObject {
	ret := json.NewElemObject()
	ret.Put(groupTitle, json.NewElemString(g.Title))
	ret.Put(groupDescription, json.NewElemString(g.Description))
	docs := json.NewElemObject()
	for key, val := range g.Docs {
		docs.Put(key, val.JSON())
	}
	ret.Put(groupDocs, docs)
	return ret
}

// DecodeGroupInfo はグループの情報をデコードします。
func DecodeGroupInfo(group json.ElemObject, groupname string) *GroupInfo {
	ret := &GroupInfo{Docs: map[string]*DocInfo{}}
	// 表題
	if t, ok := group.Child(groupTitle).AsString(); ok {
		ret.Title = t.Text()
	}
	if ret.Title == "" {
		ret.Title = groupname
	}
	// 注釈
	if t, ok := group.Child(groupDescription).AsString(); ok {
		ret.Description = t.Text()
	}
	if ret.Description == "" {
		ret.Description = ret.Title
	}
	// ドキュメント情報
	if edoc, ok := group.Child(groupDocs).AsObject(); ok {
		ret.Docs = DecodeDocInfoMap(edoc)
	}
	return ret
}

// DecodeGroupInfoMap はグループの情報一覧をデコードします。
func DecodeGroupInfoMap(groups json.ElemObject) map[string]*GroupInfo {
	ret := map[string]*GroupInfo{}
	for _, groupname := range groups.Keys() {
		if eg, ok := groups.Child(groupname).AsObject(); ok {
			ret[groupname] = DecodeGroupInfo(eg, groupname)
		}
	}
	return ret
}
