package main

import (
	"context"
	"fmt"
	"github.com/filecoin-project/go-jsonrpc"
	"net/http"
)

func main() {
	type Func struct {
		Add func(uint64, uint64) (uint64, error)
	}

	node := Func{}
	closer, err := jsonrpc.NewClient(context.TODO(), "http://127.0.0.1:8888", "PROOF", &node, http.Header{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer closer()
	fmt.Println(node.Add(1, 1))
}
