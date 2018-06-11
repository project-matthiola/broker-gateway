package matcher

import (
	"container/heap"
	"fmt"

	"github.com/rudeigerc/broker-gateway/model"
	"github.com/shopspring/decimal"
)

type Level struct {
	Price decimal.Decimal `json:"price"`
	Order []*model.Order  `json:"orders"`
}

type LevelHeap []Level

func (h LevelHeap) Len() int           { return len(h) }
func (h LevelHeap) Less(i, j int) bool { return h[i].Price.LessThan(h[j].Price) }
func (h LevelHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *LevelHeap) Push(x interface{}) {
	for index, level := range *h {
		if x.(Level).Price.Equal(level.Price) {
			(*h)[index].Order = append((*h)[index].Order, x.(Level).Order...)
			return
		}
	}
	*h = append(*h, x.(Level))
}

func (h *LevelHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Peek returns the peek of the heap.
func (h *LevelHeap) Peek() *Level {
	if h.Len() == 0 {
		return nil
	}
	return &(*h)[0]
}

func (h *LevelHeap) Remove(o model.Order) error {
	for index, level := range *h {
		if o.Price.Equal(level.Price) {
			for i, order := range level.Order {
				if o.OrderID == order.OrderID {
					(*h)[index].Order = append((*h)[index].Order[:i], (*h)[index].Order[i+1:]...)
					if len((*h)[index].Order) == 0 {
						heap.Remove(h, index)
					}
					return nil
				}
			}
		}
	}
	return fmt.Errorf("[matcher.level_heap.Remove] [ERROR] Removing order failed")
}

// MinHeap defines a min heap.
type MinHeap struct {
	LevelHeap
}

func NewMinHeap() *MinHeap {
	return &MinHeap{LevelHeap{}}
}

func (h MinHeap) Less(i, j int) bool { return h.LevelHeap[i].Price.LessThan(h.LevelHeap[j].Price) }

// MaxHeap defines a max heap.
type MaxHeap struct {
	LevelHeap
}

func NewMaxHeap() *MaxHeap {
	return &MaxHeap{LevelHeap{}}
}

func (h MaxHeap) Less(i, j int) bool { return h.LevelHeap[i].Price.GreaterThan(h.LevelHeap[j].Price) }
