package main

import (
	"fmt"
	"log"
	"reflect"
	"userDemo/api"

	"github.com/ianschenck/envflag"
)

func main() {
	fmt.Println("Jai baba ri")
	//secretKey
	var secretKey = envflag.String("SECRET_KEY", "0123457890123457890123457890123456789", "secret key for JWT signing")

	const minSecretKeySize = 32
	if len(*secretKey) < minSecretKeySize {
		log.Fatalf("SECRET_KEY must be at least %d characters", minSecretKeySize)
	}


	//storage init
	conn, err := api.NewPostgreStore()
	log.Println(reflect.TypeOf(conn))
	if err != nil {
		log.Fatalf("Error is : %+v", err)
	}
	conn.Init()

	//route init
	listenAddr := ":3000"
	router := api.NewServerApi(listenAddr, conn, *secretKey)
	router.Run()
}