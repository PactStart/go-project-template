package utils

import (
	"github.com/bwmarrin/snowflake"
)

func init() {
	var err error
	idGenerator, err = snowflake.NewNode(getNodeNum())
	if err != nil {
		panic(err)
	}
}

func getNodeNum() int64 {
	return 1
}

var idGenerator *snowflake.Node

func GenID() int64 {
	return idGenerator.Generate().Int64()
}
