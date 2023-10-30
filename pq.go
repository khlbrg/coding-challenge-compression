package main

import (
	"container/heap"
	"strings"
)

type FrequencyTable map[string]int

type Item struct {
	Value     string `json:"v,omitempty"`
	Prio      int    `json:"-"`
	Index     int    `json:"-"`
	LeftNode  *Item  `json:"l,omitempty"`
	RightNode *Item  `json:"r,omitempty"`
}

// Serialize is a depth first recursive function that creates a serialisation of a tree
// key is user as a null separator to know when a lead node is hit
func Serialize(item *Item) string {
	res := []string{}
	nullKey := byte(0x1E)
	separator := string(byte(0x1F))

	var dfs func(i *Item)
	dfs = func(i *Item) {
		if i == nil {
			res = append(res, string(nullKey))
			return
		}
		res = append(res, i.Value)

		dfs(i.LeftNode)
		dfs(i.RightNode)
	}
	dfs(item)
	return strings.Join(res, separator)
}

func Deserialize(data string) *Item {
	separator := string(byte(0x1F))
	nodes := strings.Split(data, separator)
	i := 0

	nullKey := byte(0x1E)

	var dfs func() *Item
	dfs = func() *Item {
		if i >= len(nodes) {
			return nil
		}

		if nodes[i] == string(nullKey) {
			// We reach a 'null' node, which means we are done
			// with this leaf
			i++
			return nil
		}

		item := Item{
			Value: nodes[i],
		}
		i++
		item.LeftNode = dfs()
		item.RightNode = dfs()

		return &item
	}

	return dfs()
}

type PrioQueue []*Item

func (pq PrioQueue) Len() int {
	return len(pq)
}
func (pq PrioQueue) Less(i, j int) bool {
	return pq[i].Prio < pq[j].Prio
}

func (pq PrioQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}
func (pq *PrioQueue) Push(i any) {
	n := len(*pq)
	item := i.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PrioQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
func (pq *PrioQueue) update(item *Item, value string, priority int) {
	item.Value = value
	item.Prio = priority
	heap.Fix(pq, item.Index)
}

func NewPriorityQueue(ft FrequencyTable) PrioQueue {
	pq := make(PrioQueue, len(ft))

	i := 0
	for k, v := range ft {
		pq[i] = &Item{
			Value: string(k),
			Prio:  v,
			Index: i,
		}
		i++
	}

	heap.Init(&pq)

	return pq
}
func (pq PrioQueue) generateNodeTree() Item {
	for pq.Len() > 1 {
		item_min_freq := heap.Pop(&pq).(*Item)
		item_second_min_freq := heap.Pop(&pq).(*Item)

		combinedItem := Item{
			Prio:      item_min_freq.Prio + item_second_min_freq.Prio,
			LeftNode:  item_min_freq,
			RightNode: item_second_min_freq,
		}

		pq.Push(&combinedItem)
	}
	return *pq[0]
}
