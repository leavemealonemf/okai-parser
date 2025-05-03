package mg

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

func Seed(client *mongo.Client, ctx context.Context) {
	scootColl := client.Database("iot").Collection("okai_scooters")

	if scootColl == nil {
		err := client.Database("iot").CreateCollection(ctx, "okai_scooters")
		if err != nil {
			log.Fatalln("Failed to create mongo scooter coolection", err.Error())
		}
	} else {
		fmt.Println("Mongo scooter coll already exist")
	}

	configsColl := client.Database("iot").Collection("okai_configs")

	if configsColl == nil {
		err := client.Database("iot").CreateCollection(ctx, "okai_configs")
		if err != nil {
			log.Fatalln("Failed to create mongo scooter_configs coolection", err.Error())
		}
	} else {
		fmt.Println("Mongo scooter_configs coll already exist")
	}

	cmdsColl := client.Database("iot").Collection("okai_commands")

	if cmdsColl == nil {
		err := client.Database("iot").CreateCollection(ctx, "okai_commands")
		if err != nil {
			log.Fatalln("Failed to create mongo okai_commands coolection", err.Error())
		}
	} else {
		fmt.Println("Mongo okai_commands coll already exist")
	}

	fmt.Println("[MONGO] successfull seed.")
}
