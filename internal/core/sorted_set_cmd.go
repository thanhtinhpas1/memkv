package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"memkv/internal/constants"
)

func cmdZADD(args []string) []byte {
	if len(args) < 3 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'ZADD' command"), false)
	}
	key := args[0]
	scoreIndex := 1
	flags := 0
	for scoreIndex < len(args) {
		if strings.ToLower(args[scoreIndex]) == "nx" {
			flags |= ZAddInNX
		} else if strings.ToLower(args[scoreIndex]) == "xx" {
			flags |= ZAddInXX
		} else {
			break
		}
		scoreIndex++
	}
	nx := (flags & ZAddInNX) != 0
	xx := (flags & ZAddInXX) != 0
	if nx && xx {
		return Encode(errors.New("(error) Cannot have both NN and XX flag for 'ZADD' command"), false)
	}
	numScoreEleArgs := len(args) - scoreIndex
	if numScoreEleArgs%2 == 1 || numScoreEleArgs == 0 {
		return Encode(errors.New(fmt.Sprintf("(error) Wrong number of (score, member) arg: %d", numScoreEleArgs)), false)
	}

	zset, exist := zsetStore[key]
	if !exist {
		zset = CreateZSet()
		zsetStore[key] = zset
	}

	count := 0
	for i := scoreIndex; i < len(args); i += 2 {
		ele := args[i+1]
		score, err := strconv.ParseFloat(args[i], 64)
		if err != nil {
			return Encode(errors.New("(error) Score must be floating point number"), false)
		}
		ret, outFlag := zset.Add(score, ele, flags)
		if ret != 1 {
			return Encode(errors.New("error when adding element"), false)
		}
		if outFlag != ZAddOutNop {
			count++
		}
	}
	return Encode(count, false)
}

func cmdZRANK(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'ZRANK' command"), false)
	}
	key, member := args[0], args[1]
	zset, exist := zsetStore[key]
	if !exist {
		return constants.RespNil
	}
	rank, _ := zset.GetRank(member, false)
	return Encode(rank, false)
}

func cmdZREM(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'ZREM' command"), false)
	}
	key := args[0]
	zset, exist := zsetStore[key]
	if !exist {
		return constants.RespZero
	}
	deleted := 0
	for i := 1; i < len(args); i++ {
		ret := zset.Del(args[i])
		if ret == 1 {
			deleted++
		}
		if zset.Len() == 0 {
			delete(zsetStore, key)
			break
		}
	}
	return Encode(deleted, false)
}

func cmdZSCORE(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'ZSCORE' command"), false)
	}
	key, member := args[0], args[1]
	zset, exist := zsetStore[key]
	if !exist {
		return constants.RespNil
	}
	ret, score := zset.GetScore(member)
	if ret == -1 {
		return constants.RespNil
	}
	return Encode(fmt.Sprintf("%f", score), false)
}

func cmdZCARD(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'ZCARD' command"), false)
	}
	key := args[0]
	zset, exist := zsetStore[key]
	if !exist {
		return constants.RespZero
	}
	return Encode(zset.Len(), false)
}
