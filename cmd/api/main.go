package main

import (
	"fmt"

	"github.com/pt010104/api-golang/config"
	"github.com/pt010104/api-golang/internal/appconfig/mongo"
	httpServer "github.com/pt010104/api-golang/internal/httpserver"
	pkgLog "github.com/pt010104/api-golang/pkg/log"

	"log"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("the uri: %s", cfg.Mongo.URI)
	client, db, err := mongo.ConnectDB(cfg.Mongo.URI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer client.Disconnect(nil)

	log.Printf("Connected to database: %s", db.Name())

	fmt.Printf("Mongo URI: %s\n", cfg.Mongo.URI)
	l := pkgLog.InitializeZapLogger(pkgLog.ZapConfig{
		Level:    cfg.Logger.Level,
		Mode:     cfg.Logger.Mode,
		Encoding: cfg.Logger.Encoding,
	})

	srv := httpServer.New(l, httpServer.Config{
		Port:         cfg.HTTPServer.Port,
		JWTSecretKey: cfg.JWT.SecretKey,
		Mode:         cfg.HTTPServer.Mode,
		Database:     *db,
	})
	if err := srv.Run(); err != nil {
		panic(err)
	}

}
