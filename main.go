package main

import "orderin-server/cmd"

//https://zhuanlan.zhihu.com/p/451756835

//go:generate swag init --instanceName=admin --generalInfo=internal/routers/admin.go --exclude=internal/api/app     -o docs/admin-api --parseInternal false --parseDepth 1
//go:generate swag init --instanceName=app    --generalInfo=internal/routers/app.go    --exclude=internal/api/admin  -o docs/app-api --parseInternal false --parseDepth 1

func main() {
	cmd.Execute()
}
