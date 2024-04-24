package main

import (
	"partielgo/pkg/cli"
	"partielgo/pkg/db"
	"partielgo/pkg/web"
)

func main() {
	database := db.New()
	go web.Serve(database)
	cli.Run(database)
}
