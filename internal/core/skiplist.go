package core

import (
	"math/rand"
	"strings"
)

const SkipListMaxLevel = 32

/*
/level 2: span=2 | forward\ --------------------------------------> /span=0 | forward\ ----> NULL
|level 1: span=1 | forward| --------> /span=1 | forward\ ---------> |span=0 | forward| ----> NULL
|ele                      |           |ele             |            |ele             |
|score                    |           |score           |            |score           |
|backward                 | <-------- |backward        | <--------- |backward        |
\node1                    /           \node2           /            \node3           /
*/
type SkiplistLevel struct {
	forward *SkipListNode
	span    uint32 // span is number of nodes between current node and the node forward at current level
}

type SkipListNode struct {
	ele      string
	score    float64
	backward *SkipListNode
	levels   []SkiplistLevel
}

type Skiplist struct {
	head   *SkipListNode
	tail   *SkipListNode
	length uint32
	level  int
}

func (sl *Skiplist) randomLevel() int {
	level := 1
	for rand.Intn(2) == 1 {
		level++
	}

	if level > SkipListMaxLevel {
		return SkipListMaxLevel
	}
	return level
}

func (sl *Skiplist) CreateNode(level int, score float64, ele string) *SkipListNode {
	res := &SkipListNode{
		ele:      ele,
		score:    score,
		backward: nil,
	}

	res.levels = make([]SkiplistLevel, level)

	return res
}

func CreateSkipList() *Skiplist {
	sl := Skiplist{
		length: 0,
		level:  1,
	}

	sl.head = sl.CreateNode(SkipListMaxLevel, 0, "")
	sl.head.backward = nil
	sl.tail = nil
	return &sl
}

func (sl *Skiplist) Insert(score float64, ele string) *SkipListNode {
	// update store nodes we have to cross to reach the insert position
	// rank scores the corresponding "rank" of each node in update. Skiplist head's rank = 0
	update := [SkipListMaxLevel]*SkipListNode{}
	rank := [SkipListMaxLevel]uint32{}
	x := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		for x.levels[i].forward != nil && (x.levels[i].forward.score < score ||
			x.levels[i].forward.score == score && strings.Compare(x.levels[i].forward.ele, ele) == -1) {
			rank[i] += x.levels[i].span
			x = x.levels[i].forward
		}

		update[i] = x
	}

	level := sl.randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.head
			update[i].levels[i].span = sl.length
		}

		sl.level = level
	}

	x = sl.CreateNode(level, score, ele)
	for i := 0; i < level; i++ {
		x.levels[i].forward = update[i].levels[i].forward
		update[i].levels[i].forward = x
		x.levels[i].span = update[i].levels[i].span - (rank[0] - rank[i]) // rank[0] - rank[i] = distance to insertion position
		update[i].levels[i].span = rank[0] - rank[i] + 1
	}

	// increase span for untouched level because we have a new node
	for i := level; i < sl.level; i++ {
		update[i].levels[i].span++
	}

	if update[0] == sl.head {
		x.backward = nil
	} else {
		x.backward = update[0]
	}

	if x.levels[0].forward != nil { // x not the end of list
		x.levels[0].forward.backward = x
	} else { // x is the end item of list
		sl.tail = x
	}

	sl.length++
	return x
}

func (sl *Skiplist) UpdateScore(curScore float64, ele string, newScore float64) *SkipListNode {
	update := [SkipListMaxLevel]*SkipListNode{}
	x := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil && (x.levels[i].forward.score < curScore ||
			(x.levels[i].forward.score == curScore && strings.Compare(x.levels[i].forward.ele, ele) < 0)) {
			x = x.levels[i].forward
		}
		update[i] = x
	}
	x = x.levels[0].forward
	// if x is head or tail, or newScore still greater than backward node
	if (x.backward == nil || x.backward.score < newScore) &&
		(x.levels[0].forward == nil || x.levels[0].forward.score > newScore) {
		x.score = newScore
		return x
	}

	sl.DeleteNode(x, update)
	newNode := sl.Insert(newScore, ele)
	return newNode
}

func (sl *Skiplist) Delete(score float64, ele string) int {
	update := [SkipListMaxLevel]*SkipListNode{}
	x := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil && (x.levels[i].forward.score < score ||
			(x.levels[i].forward.score == score && strings.Compare(x.levels[i].forward.ele, ele) < 0)) {
			x = x.levels[i].forward
		}
		update[i] = x
	}

	x = x.levels[0].forward
	if x != nil && x.score == score && strings.Compare(x.ele, ele) == 0 {
		sl.DeleteNode(x, update)
		return 1
	}
	return 0
}

func (sl *Skiplist) DeleteNode(x *SkipListNode, update [SkipListMaxLevel]*SkipListNode) {
	for i := 0; i < sl.level; i++ {
		if update[i].levels[i].forward == x {
			update[i].levels[i].span += x.levels[i].span - 1
			update[i].levels[i].forward = x.levels[i].forward
		} else {
			update[i].levels[i].span--
		}
	}

	// not the end of list
	if x.levels[0].forward != nil {
		x.levels[0].forward.backward = x.backward
	} else {
		sl.tail = x.backward
	}

	for sl.level > 1 && sl.head.levels[sl.level-1].forward == nil {
		sl.level--
	}
	sl.length--
}

func (sl *Skiplist) GetRank(score float64, ele string) uint32 {
	x := sl.head
	var rank uint32 = 0

	for i := sl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil && (x.levels[i].forward.score < score ||
			(x.levels[i].forward.score == score &&
				strings.Compare(x.levels[i].forward.ele, ele) <= 0)) {
			rank += x.levels[i].span
			x = x.levels[i].forward
		}
		if x.score == score && strings.Compare(x.ele, ele) == 0 {
			return rank
		}
	}

	return 0
}
