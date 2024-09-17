package main

import (
	"context"
	"database/sql"
	"fmc/config"
	_ "fmc/docs"
	"fmc/functions"
	"fmc/pkg"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// @title           Find Your Car
// @version         1.0

func main() {
	fmt.Println("----------------------------- # START # ------------------------------")

	var configvar config.ConfigFile
	if err := configvar.Load(); err != nil {
		log.Err(err).Msgf("Error loading config: %v", err)
	} else {
		fmt.Println("Success fetching config data")
	}

	fmt.Printf("Server running on %s:%d ", configvar.Server.Host, configvar.Server.Port)
	fmt.Printf("Database connecting to %s:%d", configvar.Database.Host, configvar.Database.Port)

	var dbv = configvar.Database
	var dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbv.User, dbv.Password, dbv.Host, dbv.Port, dbv.Name)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	ctx := context.Background()

	if err := db.Ping(); err != nil {
		fmt.Println(err)
	}

	// List of table models to be created
	models := []interface{}{
		(*pkg.Carpark)(nil),
		(*pkg.PresentCar)(nil),
		(*pkg.Zone)(nil),
		(*pkg.ImageZone)(nil),
		(*pkg.Camera)(nil),
		(*pkg.CarDetail)(nil),
		(*pkg.ApiManage)(nil),
	}

	// Call the CreateTables function to create all tables
	if err := functions.CreateTables(ctx, db, models); err != nil {
		fmt.Printf("Failed to create tables: %v", err)
	}

	fmt.Println("Tables created successfully")
	log.Debug().Msgf("-------------------------------- # VALUES START # ------------------------------")

	pkg.Loadzonelist(db)
	pkg.LoadCarparklist(db)
	pkg.LoadCameralist(db)

	fmt.Println(len(pkg.Zonelist))
	fmt.Println(len(pkg.CarParkList))
	fmt.Println(len(pkg.CameraList))
	log.Debug().Msgf("-------------------------------- # VALUES END # ------------------------------")

	r := pkg.SetupRouter(db)

	var host = fmt.Sprintf("%s:%d", configvar.Server.Host, configvar.Server.Port)
	if err := r.Run(host); err != nil {
		log.Err(err).Msgf("Failed to run server: %v", err)
	}

	log.Debug().Msgf("-------------------------------- # END # ------------------------------")
}
