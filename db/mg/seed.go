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
	fmt.Println("[MONGO] successfull seed.")
}
