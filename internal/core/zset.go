package core

const (
	ZAddInNX = 1 << 1 // Only add new elements. Don't update already existing elements
	ZAddInXX = 1 << 2 // Only update elements that already exist. Don't add new element
)
const (
	ZAddOutNop     = 1 << 0 // Operation not performed because of conditional
	ZAddOutAdded   = 1 << 1 // The element was new and was added
	ZAddOutUpdated = 1 << 2 // The element already existed, score updated
)

type ZSet struct {
	zskiplist *Skiplist
	// map from ele to score
	dict map[string]float64
}

func (zs *ZSet) Add(score float64, ele string, flag int) (int, int) {
	nx := flag & ZAddInNX
	xx := flag & ZAddInXX

	if len(ele) == 0 {
		return 0, ZAddOutNop
	}

	if curScore, exist := zs.dict[ele]; exist {
		if nx != 0 {
			return 1, ZAddOutNop
		}
		if curScore != score {
			znode := zs.zskiplist.UpdateScore(curScore, ele, score)
			zs.dict[ele] = znode.score
			return 1, ZAddOutUpdated
		}
		return 1, ZAddOutNop
	}

	if xx != 0 {
		return 1, ZAddOutNop
	}

	znode := zs.zskiplist.Insert(score, ele)
	zs.dict[ele] = znode.score
	return 1, ZAddOutAdded
}

func (zs *ZSet) Del(ele string) int {
	score, exists := zs.dict[ele]
	if !exists {
		return 0
	}
	delete(zs.dict, ele)
	zs.zskiplist.Delete(score, ele)
	return 1
}

func (zs *ZSet) GetRank(ele string, reverse bool) (rank int64, score float64) {
	setSize := zs.zskiplist.length
	score, exists := zs.dict[ele]
	if !exists {
		return -1, 0
	}
	rank = int64(zs.zskiplist.GetRank(score, ele))
	if reverse {
		rank = int64(setSize) - rank
	} else {
		rank--
	}
	return rank, score
}

func (zs *ZSet) GetScore(ele string) (int, float64) {
	score, exist := zs.dict[ele]
	if !exist {
		return -1, 0
	}
	return 0, score
}

func (zs *ZSet) Len() int {
	return len(zs.dict)
}

func CreateZSet() *ZSet {
	zs := ZSet{
		zskiplist: CreateSkipList(),
		dict:      map[string]float64{},
	}
	return &zs
}
