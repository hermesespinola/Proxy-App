package middlewares

import (
	"bufio"
	"encoding/json"
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

// RepoNode is a node in the heap
type RepoNode struct {
	Domain   string `json:"domain,omitempty"`
	Weight   int    `json:"weight,omitempty"`
	Priority int    `json:"priority,omitempty"`
}

// Value returns the priority value of the node
func (node RepoNode) Value() float64 {
	return float64(node.Weight) * float64(node.Priority) / 2.0
}

// Greater Implement lib.Node for RepoNode
func (node RepoNode) Greater(that interface{}) bool {
	if v, ok := that.(RepoNode); ok {
		return node.Value() > v.Value()
	}
	return false
}

func readData() {
	path, _ := filepath.Abs("api/middlewares/domain.txt")
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	currentNode := RepoNode{}
	for scanner.Scan() {
		if scanner.Text() == "" {
			repoHeap = repoHeap.Insert(currentNode)
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
	if currentNode.Domain != "" {
		repoHeap = repoHeap.Insert(currentNode)
	}
}

var domains = []string{"alpha", "omega", "beta"}

func contains(arr []string, element string) bool {
	for _, val := range arr {
		if val == element {
			return true
		}
	}
	return false
}

// PushNode is the middleware for our Proxy
func PushNode(c iris.Context) {
	domain := c.GetHeader("domain")
	if !contains(domains, domain) {
		c.JSON(iris.Map{"status": "domain error"})
		return
	}
	readData()
	newNode := RepoNode{domain, rand.Intn(100), rand.Intn(100)}
	repoHeap = repoHeap.Insert(newNode)
	nodeStr, _ := json.Marshal(newNode)
	c.JSON(iris.Map{
		"status": "ok",
		"domain": domain,
		"new":    string(nodeStr),
	})
}

// PopNode handles pop from queue
func PopNode(c iris.Context) {
	var node lib.Node
	repoHeap, node = repoHeap.Pop()
	var nodeStr []byte
	if v, ok := node.(RepoNode); ok {
		nodeStr, _ = json.Marshal(v)
	}
	c.JSON(iris.Map{
		"status": "ok",
		"popped": string(nodeStr),
	})
}

// Read is a handler that reads from file
func Read(c iris.Context) {
	readData()
	queueStr, _ := json.Marshal(repoHeap)
	c.JSON(iris.Map{
		"status": "ok",
		"queue":  string(queueStr),
	})
}
