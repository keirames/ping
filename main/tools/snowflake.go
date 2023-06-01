package tools

import "github.com/bwmarrin/snowflake"

var Snowflake *snowflake.Node

func NewNode() error {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return err
	}

	Snowflake = node

	return nil
}
