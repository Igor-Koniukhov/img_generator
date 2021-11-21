package main

import (
	"flag"
	"github.com/igor-koniukhov/img_generator/configs"
	"github.com/igor-koniukhov/img_generator/inernal/server"
	"log"
)

var confPath = flag.String("conf-path", "./configs/.env", "Path to config env.")

func main()  {
	conf, err := configs.New(*confPath)
	if err !=nil {
		log.Fatalln(err)
	}
	server.Run(conf)

}