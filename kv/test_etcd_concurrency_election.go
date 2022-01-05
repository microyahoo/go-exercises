package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

const prefix = "/election-demo"
const prop = "local"

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:12379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	campaign(cli, prefix, prop)

}

func campaign(c *clientv3.Client, election string, prop string) {
	for {
		// 租约到期时间：5s
		s, err := concurrency.NewSession(c, concurrency.WithTTL(5))
		if err != nil {
			fmt.Println(err)
			continue
		}
		e := concurrency.NewElection(s, election)
		ctx := context.TODO()

		log.Println("开始竞选")

		err = e.Campaign(ctx, prop)
		if err != nil {
			log.Println("竞选 leader失败，继续")
			switch {
			case err == context.Canceled:
				return
			default:
				continue
			}
		}

		log.Println("获得leader")
		if err := doCrontab(); err != nil {
			log.Println("调用主方法失败，辞去leader，重新竞选")
			_ = e.Resign(ctx)
			continue
		}
		return
	}
}

func doCrontab() error {
	for {
		fmt.Println("doCrontab")
		time.Sleep(time.Second * 4)
		//return fmt.Errorf("sss")
	}
}
