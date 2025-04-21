package main

import (
	"flkcli/cmd"
	_ "flkcli/cmd/list"
	_ "flkcli/cmd/login"
	_ "flkcli/cmd/setup"
	_ "flkcli/cmd/upload"
)

func main() {
	cmd.Execute()
}
