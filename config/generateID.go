package config

import (
	"github.com/bwmarrin/snowflake"
)

var (
	node, _ = snowflake.NewNode(1)
)

func GetId() string {
	id := node.Generate().Base2()
	return id
}
