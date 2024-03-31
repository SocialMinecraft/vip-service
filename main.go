package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"runtime"
	"vip-service/database"
	"vip-service/eventHandlers"
)

var (
	nc   *nats.Conn
	subs []*nats.Subscription
	db   *database.Db
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
		return
	}

	config, err := getConfig()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	db, err = database.Connect(config.PostgresUrl)
	if err != nil {
		log.Fatalln(err)
		return
	}

	nc, err = nats.Connect(config.NatsUrl)
	if err != nil {
		log.Fatalln(err)
		return
	}

	var sub *nats.Subscription

	sub, err = nc.Subscribe("vip.claim", func(msg *nats.Msg) {
		if err := eventHandlers.Claim(db, msg); err != nil {
			log.Println(err)
		}
	})
	checkError(err)
	subs = append(subs, sub)

	sub, err = nc.Subscribe("vip.get", func(msg *nats.Msg) {
		if err := eventHandlers.Get(db, msg); err != nil {
			log.Println(err)
		}
	})
	checkError(err)
	subs = append(subs, sub)

	sub, err = nc.Subscribe("kofi.payment", func(msg *nats.Msg) {
		if err := eventHandlers.KofiPayment(db, msg); err != nil {
			log.Println(err)
		}
	})
	checkError(err)
	subs = append(subs, sub)

	sub, err = nc.Subscribe("vip.sync", func(msg *nats.Msg) {
		if err := eventHandlers.SyncAccounts(nc, db, msg); err != nil {
			log.Println(err)
		}
	})
	checkError(err)
	subs = append(subs, sub)

	sub, err = nc.Subscribe("accounts.minecraft.changed", func(msg *nats.Msg) {
		if err := eventHandlers.ChangeAcconut(db, msg); err != nil {
			log.Println(err)
		}
	})
	checkError(err)
	subs = append(subs, sub)

	log.Println("Running")
	runtime.Goexit()
}

func checkError(err error) {
	if err == nil {
		return
	}

	log.Fatalln(err)
	os.Exit(2)
}
