package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/venus/venus-shared/api"
	v1 "github.com/filecoin-project/venus/venus-shared/api/chain/v1"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

func main() {
	go func() {
		//start pprof
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	logging.SetLogLevel("*", "debug")
	app := &cli.App{
		Name:                 "net test",
		Usage:                "nettest",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Usage:    "url for access venus",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "token",
				Usage:    "token for access venus",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "duration",
				Usage: "time for test",
				Value: "5m",
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
	url := cctx.String("url")
	token := cctx.String("token")
	//try to connect venus
	ctx := context.Background()

	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+token)

	var node v1.FullNodeStruct
	closer, err := jsonrpc.NewMergeClient(ctx, url, "Filecoin", api.GetInternalStructs(&node), headers, jsonrpc.WithRetry(true), jsonrpc.WithReconnectBackoff(time.Millisecond*100, time.Second*10))
	if err != nil {
		return err
	}
	defer closer()

	if cctx.IsSet("duration") {
		du, err := time.ParseDuration(cctx.String("duration"))
		if err != nil {
			return err
		}
		ctx, _ = context.WithTimeout(ctx, du)
	}

	for {
		t := time.Now()
		_, err := node.ChainHead(ctx)
		if err != nil {
			fmt.Println("cannt connect to venus")
		}
		fmt.Println("connect success", time.Now().Sub(t))
		select {
		case <-time.After(5 * time.Second):
		case <-ctx.Done():
			return errors.New("exit command")
		}
	}
}
