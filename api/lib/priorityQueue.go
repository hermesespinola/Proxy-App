package lib

import "fmt"

// Node must implement the Greater method to be able
// to order our priority queue
type Node interface {
	// return true if `this` is greater than `that`, false otherwise.
	Greater(that interface{}) bool
}

// PriorityQueue is a collection of nodes
type PriorityQueue interface {
	IsEmpty() bool
	Peek() Node
	Pop() (PriorityQueue, Node)
	Replace(Node) (PriorityQueue, Node)
	Insert(Node) PriorityQueue
}

// BHeap is a max Binary Heap
type BHeap []Node

// IsEmpty tells if the queue is empty
func (queue BHeap) IsEmpty() bool {
	return len(queue) == 0
}

// Peek returns the max element of the queue, without removing it
func (queue BHeap) Peek() Node {
	return queue[0]
}

// Pop removes the max element from the queue and removes it
func (queue BHeap) Pop() (PriorityQueue, Node) {
	n := len(queue)
	if n == 0 {
		panic("Pop of empty queue")
	}

	// store and remove the first element in slice
	tmp := queue[0]
	// place last element in the root and remove first
	queue[0] = queue[n-1]
	// resize slice
	queue = queue[:n-1]

	if len(queue) > 1 {
		// restore heap
		queue.siftDown(0)
	}
	return queue, tmp
}

// Replace pops the first element, replace it with the new
// element, and restore the priority property.
// This is more efficient than performint `Pop` followed by
// `Insert`
func (queue BHeap) Replace(node Node) (PriorityQueue, Node) {
	tmp := queue[0]
	queue[0] = node
	queue.siftDown(0)
	return queue, tmp
}

// Insert adds a node to the priority queue
func (queue BHeap) Insert(node Node) PriorityQueue {
	queue = append(queue, node)
	queue.siftUp(len(queue) - 1)
	return queue
}

// Move node `i` down the heap three
func (queue BHeap) siftDown(i int) {
	n := len(queue)
	for {
		left, right := 2*i+1, 2*i+2
		fmt.Println("[siftDown]", i, left, right)
		if left >= n {
			// We already are at the last level of the three
			break
		} else if right >= n {
			// There's only a left branch
			if !queue[i].Greater(queue[left]) {
				queue[i], queue[left] = queue[left], queue[i]
				i = left
			} else {
				break
			}
		} else {
			// Move up the max child
			var max int
			if queue[left].Greater(queue[right]) {
				max = left
			} else {
				max = right
			}

			if !queue[i].Greater(queue[max]) {
				queue[i], queue[max] = queue[max], queue[i]
				i = max
			} else {
				break
			}
		}
	}
}

// Move node `i` up the heap three
func (queue BHeap) siftUp(i int) {
	for {
		parent := (i - 1) / 2
		fmt.Println("i", i, "parent", parent)
		if parent == i || queue[parent].Greater(queue[i]) {
			break
		} else {
			fmt.Println("Swap i and parent")
			queue[i], queue[parent] = queue[parent], queue[i]
			i = parent
		}
	}
}

// MockQueue creates a queue with `node` repeated `n` times.
func MockQueue(node Node, n uint) PriorityQueue {
	queue := BHeap{}
	for i := uint(0); i < n; i++ {
		queue = append(queue, node)
	}
	return queue
}
