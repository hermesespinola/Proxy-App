package middlewares

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hermesespinola/proxy-app/api/lib"
	"github.com/kataras/iris"
)

var repoHeap lib.PriorityQueue = lib.BHeap{}

type repoNode struct {
	Domain   string
	Weight   int
	Priority int
}

// Implement lib.Node for repoNode
func (node repoNode) Greater(that interface{}) bool {
	if v, ok := that.(repoNode); ok {
		return node.Weight > v.Weight || node.Priority > v.Priority
	}
	return false
}

func readData() {
	path, _ := filepath.Abs("")
	file, err := os.Open(path + "/api/middlewares/domain.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	currentNode := repoNode{}
	for scanner.Scan() {
		if scanner.Text() == "" {
			repoHeap = repoHeap.Insert(currentNode)
			currentNode = repoNode{}
			continue
		}
		result := strings.SplitN(scanner.Text(), ":", 2)
		switch result[0] {
		case "weight":
			weight, _ := strconv.ParseInt(strings.Trim(result[1], " "), 10, 64)
			currentNode.Weight = int(weight)
			break
		case "priority":
			priority, _ := strconv.ParseInt(strings.Trim(result[1], " "), 10, 64)
			currentNode.Priority = int(priority)
			break
		default:
			currentNode.Domain = result[0]
		}
	}
}

// PushNode is the middleware for our Proxy
func PushNode(c iris.Context) {
	domain := c.GetHeader("domain")
	readData()
	newNode := repoNode{domain, rand.Intn(100), rand.Intn(100)}
	repoHeap = repoHeap.Insert(newNode)
	nodeStr := fmt.Sprintf("%+v", repoNode(newNode))
	queueStr := fmt.Sprintf("%+v", repoHeap)
	c.JSON(iris.Map{
		"result": "ok",
		"domain": domain,
		"new":    nodeStr,
		"queue":  queueStr,
	})
}

// PopNode handles pop from queue
func PopNode(c iris.Context) {
	var node lib.Node
	repoHeap, node = repoHeap.Pop()
	queueStr := fmt.Sprintf("%+v", repoHeap)
	nodeStr := ""
	if v, ok := node.(repoNode); ok {
		nodeStr = fmt.Sprintf("%+v", v)
	}
	c.JSON(iris.Map{
		"result": "ok",
		"poped":  nodeStr,
		"queue":  queueStr,
	})
}

// Read is a handler that reads from file
func Read(c iris.Context) {
	readData()
	queueStr := fmt.Sprintf("%+v", repoHeap)
	c.JSON(iris.Map{
		"result": "ok",
		"queue":  queueStr,
	})
}
