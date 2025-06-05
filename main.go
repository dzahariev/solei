package main

import (
	"context"
	"log"

	"github.com/dzahariev/respite/api"
	"github.com/dzahariev/respite/auth"
	"github.com/dzahariev/respite/basemodel"
	"github.com/dzahariev/respite/cfg"
	"github.com/dzahariev/solei/model"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	ctx := context.Background()

	var loggerCfg cfg.Logger
	if err := envconfig.Process(ctx, &loggerCfg); err != nil {
		log.Fatal(err)
	}

	var databaseCfg cfg.DataBase
	if err := envconfig.Process(ctx, &databaseCfg); err != nil {
		log.Fatal(err)
	}

	var serverCfg cfg.Server
	if err := envconfig.Process(ctx, &serverCfg); err != nil {
		log.Fatal(err)
	}

	var authCfg cfg.Keycloak
	if err := envconfig.Process(ctx, &authCfg); err != nil {
		log.Fatal(err)
	}

	objects := []basemodel.Object{
		&model.Category{},
		&model.Meal{},
		&model.Order{},
		&model.OrderItem{},
	}

	rolesToPermissions := map[string][]string{
		"Customer": {
			"user.read",
			"address.read",
			"address.write",
			"meal.read",
			"category.read",
			"order.read",
			"order.write",
			"orderitem.read",
			"orderitem.write",
		},
		"Chef": {
			"user.read",
			"orderitem.global",
			"orderitem.read",
			"orderitem.write",
		},
		"Courier": {
			"user.read",
			"order.global",
			"order.read",
			"order.write",
		},
		"Owner": {
			"user.read",
			"meal.read",
			"meal.write",
			"category.read",
			"category.write",
			"order.global",
			"order.read",
			"orderitem.global",
			"orderitem.read",
		},
	}
	authClient := auth.NewClient(authCfg)
	server, err := api.NewServer(serverCfg, loggerCfg, databaseCfg, objects, authClient, rolesToPermissions)
	if err != nil {
		log.Fatal(err)
	}
	server.Run()
}
