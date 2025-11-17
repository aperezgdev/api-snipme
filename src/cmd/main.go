package main

import (
	"github.com/aperezgdev/api-snipme/src/cmd/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		panic(err)
	}
}
