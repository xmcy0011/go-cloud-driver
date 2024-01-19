package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xmcy0011/go-cloud-driver/internal/logics/interfaces"
)

func TestClosures2TreeNode(t *testing.T) {
	closures := []interfaces.MetadataNode{
		{ObjectId: "test", Ancestor: "test", Descendant: "test", Name: "n-test", Depth: 0, ParentId: "gns"},
		{ObjectId: "A2", Ancestor: "test", Descendant: "A2", Name: "n-a2", Depth: 1, ParentId: "test"},
		{ObjectId: "A3", Ancestor: "test", Descendant: "A3", Name: "n-a3", Depth: 1, ParentId: "test"},
		{ObjectId: "A1", Ancestor: "test", Descendant: "A1", Name: "n-a1", Depth: 2, ParentId: "A2"},
		{ObjectId: "B2", Ancestor: "test", Descendant: "B2", Name: "n-b2", Depth: 2, ParentId: "A2"},
		{ObjectId: "B1", Ancestor: "test", Descendant: "B1", Name: "n-b1", Depth: 3, ParentId: "A1"},
	}

	r, err := closures2TreeNode("test", closures)
	assert.NoError(t, err)
	printTree(t, r, 0)
}

// 递归打印树结构
func printTree(t *testing.T, node *interfaces.MetadataTreeNode, depth int) {
	t.Logf("%s%s(%s)\n", getIndentation(depth), node.Name, node.ObjectId)

	for _, child := range node.Children {
		printTree(t, child, depth+1)
	}
}

// 获取缩进字符串
func getIndentation(level int) string {
	indentation := ""

	for i := 0; i < level; i++ {
		indentation += "\t"
	}

	return indentation
}
