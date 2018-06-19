package main

import (
	"log"
	"net/http"
	. "./app"
)

var ctx *Context
var metaweblogHandler MetaweblogHandler
var blogListHandler BlogListHandler

type RequestHandler interface {
	Process(ctx *Context, w http.ResponseWriter, r *http.Request) (error)
}

func main() {
	ctx, _ = InitContext()
	log.Println("GitPage path: " + ctx.Conf.GitPage.Root)

	http.HandleFunc("/blog", blogHandler)
	http.HandleFunc("/xmlrpc", metablogApiHandler)
	if err := http.ListenAndServe(":" + ctx.Conf.Server.Port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func metablogApiHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	metaweblogHandler.Process(ctx, w, r)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	blogListHandler.Process(ctx, w, r)
}