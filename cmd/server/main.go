package main

import "github.com/felipedias-dev/fullcycle-go-expert-basic-api/configs"

func main() {
	config, _ := configs.LoadConfig(".")
	println(config.DBDriver)
}
