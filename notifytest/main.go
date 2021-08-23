package main

import (
	"context"
	"fmt"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/venus/app/client"
	"github.com/filecoin-project/venus/pkg/chain"
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
				Name:  "duration",
				Usage: "time for test",
				Value: "5m",
			},
			&cli.StringFlag{
				Name:  "wait",
				Usage: "wait time to test websockets",
				Value: "1m",
			},
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

	du, err := time.ParseDuration(cctx.String("duration"))
	if err != nil {
		return err
	}
	ctx, _ = context.WithTimeout(ctx, du)

	head, err := node.ChainHead(ctx)
	if err != nil {
		return err
	}
	log.Println("receive head head", head.String())

	wait, err := time.ParseDuration(cctx.String("wait"))
	if err != nil {
		return err
	}
	time.Sleep(wait)
	notifs := node.ChainNotify(ctx)

	hccurrent := <-notifs
	log.Println(fmt.Sprintf("recieve hc current event height %d", hccurrent[0].Val.Height()))
	for notif := range notifs {
		for _, event := range notif {
			if event.Type == chain.HCApply {
				log.Println("receive event ", event.Type, event.Val.Height())
				for _, blk := range event.Val.Blocks() {
					_, err := node.ChainGetBlockMessages(ctx, blk.Cid())
					if err != nil {
						log.Fatal(err)
					} else {
						log.Println("get block message successfully")
					}
				}
			}
		}
	}
	return nil
}
