package stm

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/beevik/etree"
)

type BufferPool struct {
	sync.Pool
}

func NewBufferPool() *BufferPool {
	return &BufferPool{
		Pool: sync.Pool{
			New: func() any {
				b := bytes.NewBuffer(make([]byte, 256))
				b.Reset()
				return b
			},
		},
	}
}

func (bp *BufferPool) Get() *bytes.Buffer {
	return bp.Pool.Get().(*bytes.Buffer)
}

func (bp *BufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	bp.Pool.Put(b)
}

func SetBuilderElementValue(elm *etree.Element, data [][]any, baseKey string) (*etree.Element, bool) {
	var child *etree.Element

	key := baseKey
	ts, tk := spaceDecompose(elm.Tag)
	_, sk := spaceDecompose(elm.Space)

	if elm.Tag != "" && ts != "" && tk != "" {
		key = fmt.Sprintf("%s:%s", elm.Space, baseKey)
	} else if sk != "" {
		key = fmt.Sprintf("%s:%s", sk, baseKey)
	}

	var values any
	var found bool
	for _, v := range data {
		if v[0] == baseKey {
			values = v[1]
			found = true
			break
		}
	}
	if !found {
		return child, false
	}

	switch value := values.(type) {
	case nil:
	default:
		child = elm.CreateElement(key)
		child.SetText(fmt.Sprint(value))
	case int:
		child = elm.CreateElement(key)
		child.SetText(fmt.Sprint(value))
	case string:
		child = elm.CreateElement(key)
		child.SetText(value)
	case float64, float32:
		child = elm.CreateElement(key)
		child.SetText(fmt.Sprint(value))
	case time.Time:
		child = elm.CreateElement(key)
		child.SetText(value.Format(time.RFC3339))
	case bool:
		_ = elm.CreateElement(fmt.Sprintf("%s:%s", key, key))
	case []int:
		for _, v := range value {
			child = elm.CreateElement(key)
			child.SetText(fmt.Sprint(v))
		}
	case []string:
		for _, v := range value {
			child = elm.CreateElement(key)
			child.SetText(v)
		}
	case []Attr:
		for _, attr := range value {
			child = elm.CreateElement(key)
			for k, v := range attr {
				child.CreateAttr(k, v)
			}
		}
	case Attrs:
		val, attrs := value[0], value[1]

		child, _ = SetBuilderElementValue(elm, URL{[]any{baseKey, val}}, baseKey)
		switch attr := attrs.(type) {
		case map[string]string:
			for k, v := range attr {
				child.CreateAttr(k, v)
			}
		case Attr:
			for k, v := range attr {
				child.CreateAttr(k, v)
			}
		}

	case any:
		var childKey string
		if sk == "" {
			childKey = fmt.Sprintf("%s:%s", key, key)
		} else {
			childKey = fmt.Sprint(key)
		}

		switch value := values.(type) {
		case []URL:
			for _, val := range value {
				child := elm.CreateElement(childKey)
				for _, v := range val {
					SetBuilderElementValue(child, val, v[0].(string))
				}
			}
		case URL:
			child := elm.CreateElement(childKey)
			for _, v := range value {
				SetBuilderElementValue(child, value, v[0].(string))
			}
		}
	}
	return child, true
}

func MergeMap(src, dst [][]any) [][]any {
	for _, v := range dst {
		found := false
		for _, vSrc := range src {
			if v[0] == vSrc[0] {
				found = true
				break
			}
		}
		if !found {
			src = append(src, v)
		}
	}
	return src
}

// ToLowerString converts lower strings from including capital or upper strings.
func ToLowerString(in []string) (out []string) {
	for _, name := range in {
		out = append(out, strings.ToLower(name))
	}
	return out
}

func URLJoin(src string, joins ...string) string {
	var u *url.URL
	lastNum := len(joins)
	base, _ := url.Parse(src)

	for i, j := range joins {
		if !strings.HasSuffix(j, "/") && lastNum > (i+1) {
			j = j + "/"
		}

		u, _ = url.Parse(j)
		base = base.ResolveReference(u)
	}

	return base.String()
}

// spaceDecompose is separating strings for the SetBuilderElementValue
func spaceDecompose(str string) (space, key string) {
	colon := strings.IndexByte(str, ':')
	if colon == -1 {
		return "", str
	}
	return str[:colon], str[colon+1:]
}
