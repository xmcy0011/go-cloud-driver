package logics

import (
	"fmt"
	"log"
	"testing"
)

type closure struct {
	a        string
	d        string
	name     string
	depth    int
	parentId string
}

type TreeNode struct {
	ID       string
	Name     string
	Children []*TreeNode
}

func TestXxx(t *testing.T) {
	closures := []closure{
		{a: "test", d: "test", name: "n-test", depth: 0, parentId: "gns"},
		{a: "test", d: "A2", name: "n-a2", depth: 1, parentId: "test"},
		{a: "test", d: "A3", name: "n-a3", depth: 1, parentId: "test"},
		{a: "test", d: "A1", name: "n-a1", depth: 2, parentId: "A2"},
		{a: "test", d: "B2", name: "n-b2", depth: 2, parentId: "A2"},
		{a: "test", d: "B1", name: "n-b1", depth: 3, parentId: "A1"},
	}

	t.Log("dddd")
	fmt.Println("ddd")
	log.Println("dddsss")

	// 把 closures 转换为 TreeNode

	// 构建闭包表映射
	closureMap := make(map[string][]closure)
	for _, c := range closures {
		if _, ok := closureMap[c.parentId]; !ok {
			closureMap[c.parentId] = make([]closure, 0)
		}
		closureMap[c.parentId] = append(closureMap[c.parentId], c)
	}

	// 构建树
	rootID := "test"
	rootNode := buildTree(rootID, closureMap)

	// 输出树结构
	printTree(t, rootNode, 0)
}

// 递归构建树
func buildTree(nodeID string, closureMap map[string][]closure) *TreeNode {
	node := &TreeNode{ID: nodeID}

	if closures, ok := closureMap[nodeID]; ok {
		for _, c := range closures {
			childNode := buildTree(c.d, closureMap)
			node.Children = append(node.Children, childNode)
		}
	}

	return node
}

// 递归打印树结构
func printTree(t *testing.T, node *TreeNode, level int) {
	//fmt.Printf("%s%s\n", getIndentation(level), node.ID)
	t.Logf("%s%s(%s)\n", getIndentation(level), node.ID, node.Name)

	for _, child := range node.Children {
		printTree(t, child, level+1)
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
