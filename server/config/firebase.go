package config

import (
	"cloud.google.com/go/firestore"
	"context"
	_ "embed"
	firebase "firebase.google.com/go"
	"log"
)

//go:embed firebase.yml
var firebaseFile []byte

type firebaseConfig struct {
	AuthDomain        string `yaml:"auth_domain"`
	ProjectID         string `yaml:"project_id"`
	StorageBucket     string `yaml:"storage_bucket"`
	MessagingSenderID string `yaml:"messaging_sender_id"`
	AppID             string `yaml:"app_id"`
	MeasurementID     string `yaml:"measurement_id"`
}

var FirestoreClient *firestore.Client

func init() {
	cfg := new(firebaseConfig)
	if err := loadEnv(EnvLoader{DefaultENV: firebaseFile}, cfg); err != nil {
		log.Fatalf("error loading firebase configuration: %v\n", err)
	}

	ctx := context.Background()

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: cfg.ProjectID})
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	firestoreApp, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing firestore app: %v\n", err)
	}

	FirestoreClient = firestoreApp
}
