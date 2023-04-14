package main

import (
	_ "ConfigManager/db"
	"ConfigManager/globals"
	_ "ConfigManager/kube"
	_ "ConfigManager/routers"
)

func main() {

	globals.GvaGinEngine.Run(":8009")

}
