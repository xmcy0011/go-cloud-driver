package common

import (
	"errors"
	"strings"

	"github.com/oklog/ulid/v2"
)

const (
	Prefix string = "url://"
)

var (
	ErrInvalidUrlFormat = errors.New("invalid url format")
	ErrEmptyLeafNode    = errors.New("empty leaf node")
)

// IsValidUrl 判断元数据 url 格式是否合法
// example:
// - url://
// - url://01HMQDV23J3HMV4CBH2MHSRN45
// - url://01HMQDV23J3HMV4CBH2MHSRN45/01HMQE6KCPQ0PNKFDNP5MEMZSP
func IsValidUrl(url string) bool {
	_, err := GetUrlNodes(url)
	return err == nil
}

// GetUrlNodes 获取节点集
// url:// -> []
// url://01HMQDV23J3HMV4CBH2MHSRN45 -> ["01HMQDV23J3HMV4CBH2MHSRN45"]
func GetUrlNodes(url string) ([]string, error) {
	if url == "" || len(url) < len(Prefix) {
		return nil, ErrInvalidUrlFormat
	}

	prefix := url[:len(Prefix)]
	if prefix != Prefix {
		return nil, ErrInvalidUrlFormat
	}

	items := url[len(Prefix):]

	// empty nodes
	if items == "" {
		return []string{}, nil
	}

	strs := strings.Split(items, "/")
	nodes := make([]string, len(strs))
	for i, str := range strs {
		if _, err := ulid.Parse(str); err != nil {
			return nil, ErrInvalidUrlFormat
		}
		nodes[i] = str
	}
	return nodes, nil
}

// GetUrlNode 获取末尾节点
// url:// -> 空
// url://01HMQE6KCPQ0PNKFDNP5MEMZSP -> 01HMQE6KCPQ0PNKFDNP5MEMZSP
// url://01HMQE6KCPQ0PNKFDNP5MEMZSP/01HMQFJFKTM8513YSQXZYVF8S7 -> 01HMQFJFKTM8513YSQXZYVF8S7
func GetUrlLeafNode(url string) (string, error) {
	nodes, err := GetUrlNodes(url)
	if err != nil {
		return "", nil
	}
	if len(nodes) == 0 {
		return "", ErrEmptyLeafNode
	}

	return nodes[len(nodes)-1], nil
}
