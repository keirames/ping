package keygen

import "github.com/bwmarrin/snowflake"

var node *snowflake.Node

func New() {
	n, err := snowflake.NewNode(1)
	if err != nil {
		panic("Cannot create new node")
	}

	node = n
}

func Snowflake() int64 {
	id := node.Generate()

	return id.Int64()
}
