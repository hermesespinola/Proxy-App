package middlewares

import (
	"fmt"
	"math/rand"

	"github.com/hermesespinola/proxy-app/api/lib"
	"github.com/kataras/iris"
)

var repoHeap lib.PriorityQueue = lib.BHeap{}

type repoNode struct {
	Domain   string
	Weight   int
	Priority int
}

func (node repoNode) value() int {
	return node.Weight
}

// Implement lib.Node for repoNode
func (node repoNode) Greater(that interface{}) bool {
	if v, ok := that.(repoNode); ok {
		return node.value() > v.value()
	}
	return false
}

// PushNode is the middleware for our Proxy
func PushNode(c iris.Context) {
	domain := c.GetHeader("domain")
	newNode := repoNode{domain, rand.Intn(100), rand.Intn(100)}
	repoHeap = repoHeap.Insert(newNode)
	c.JSON(iris.Map{
		"result": "ok",
		"domain": domain,
		"weight": newNode.Weight,
	})
	fmt.Printf("my new node: %+v\n", newNode)
	fmt.Printf("my queue [PushNode]: %+v\n", repoHeap)
}

// PopNode handles pop from queue
func PopNode(c iris.Context) {
	var node lib.Node
	repoHeap, node = repoHeap.Pop()
	c.JSON(iris.Map{
		"result": "ok",
		"weight": node.(repoNode).Weight,
	})
	fmt.Printf("my queue [PopNode]: %+v\n", repoHeap)
}
