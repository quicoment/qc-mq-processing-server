package config

import (
	"github.com/bwmarrin/snowflake"
)

var (
	node, _ = snowflake.NewNode(1)
)

func getId() int64 {
	id := node.Generate().Int64()
	return id
}
