package main

import (
	"context"
	"fmt"
	"github.com/filecoin-project/venus/venus-shared/types"

	v1 "github.com/filecoin-project/venus/venus-shared/api/chain/v1"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	app := &cli.App{
		Name:                 "notifytest",
		Usage:                "notifytest",
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
				Name:     "url",
				Usage:    "url for access venus",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "token",
				Usage:    "token for access venus",
				Required: true,
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
	node, closer, err := v1.NewFullNodeRPC(ctx, url, headers)
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
	notifs, err := node.ChainNotify(ctx)
	if err != nil {
		return err
	}

	hccurrent := <-notifs
	log.Println(fmt.Sprintf("recieve hc current event height %d", hccurrent[0].Val.Height()))
	for notif := range notifs {
		for _, event := range notif {
			if event.Type == types.HCApply {
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
