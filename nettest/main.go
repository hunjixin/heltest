package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/venus/app/client"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	app := &cli.App{
		Name:                 "force-batch-committer",
		Usage:                "Filecoin batch precommit/commit client",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "token",
				Usage: "token for access venus",
			},
		},
	}
	//app.Setup()
	app.Action = run
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(cctx *cli.Context) error {
	token := cctx.String("token")
	//try to connect venus
	ctx := context.Background()
	node := client.FullNodeStruct{}
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+token)
	closer, err := jsonrpc.NewClient(ctx, "wss://node.filincubator.com:81/rpc/v1", "Filecoin", &node, headers)
	if err != nil {
		return err
	}
	defer closer()

	ctx, _ = context.WithTimeout(ctx, time.Minute*5)
	for {
		t := time.Now()
		_, err := node.ChainHead(ctx)
		if err != nil {
			fmt.Println("cannt connect to venus")
		}
		fmt.Println("connect success", time.Now().Sub(t))
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return errors.New("exit command")
		}
	}
}
