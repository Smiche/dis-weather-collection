package main

import (
	"fmt"
	db "weather_client/db"
)

func main() {
	//get_available_conf("")
	fmt.Println("Hello")
	config, _ := db.Get_available_conf()
	fmt.Println(config)
}
