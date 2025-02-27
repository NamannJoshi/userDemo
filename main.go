package main

import (
	"fmt"
	"log"
	"reflect"
	"userDemo/api"
)

func main() {
	fmt.Println("Jai baba ri")

	//storage init
	conn, err := api.NewPostgreStore()
	log.Println(reflect.TypeOf(conn))
	if err != nil {
		log.Fatalf("Error is : %+v", err)
	}
	conn.Init()

	//route init
	listenAddr := ":3000"
	router := api.NewServerApi(listenAddr, conn)
	router.Run()
}