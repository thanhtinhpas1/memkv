package main

import (
	"fmt"
	"log"
	"net"

	"memkv/internal/core"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	conn.Write(core.Encode([]string{"PING"}, true))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(buf[:n]))

	// ZADD - Add members with scores to a sorted set
	conn.Write(core.Encode([]string{"ZADD", "leaderboard", "100", "alice", "80", "bob", "95", "carol"}, true))
	printResponse(conn)

	// ZSCORE - Get the score of a member
	conn.Write(core.Encode([]string{"ZSCORE", "leaderboard", "alice"}, true))
	printResponse(conn)

	// ZRANK - Get the rank of a member (ordered by score from low to high)
	conn.Write(core.Encode([]string{"ZRANK", "leaderboard", "bob"}, true))
	printResponse(conn)

	// ZCARD - Get the number of members in the sorted set
	conn.Write(core.Encode([]string{"ZCARD", "leaderboard"}, true))
	printResponse(conn)

	// ZREM - Remove a member from the sorted set
	conn.Write(core.Encode([]string{"ZREM", "leaderboard", "bob"}, true))
	printResponse(conn)

	// Check ZCARD again after removal
	conn.Write(core.Encode([]string{"ZCARD", "leaderboard"}, true))
	printResponse(conn)
}

func printResponse(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf[:n]))
}
