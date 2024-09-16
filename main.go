package main

import (
	"context"
	"database/sql"
	"fmc/config"
	_ "fmc/docs"
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

	_, CarparkError := db.NewCreateTable().Model((*pkg.Carpark)(nil)).IfNotExists().Exec(ctx)
	if CarparkError != nil {
		panic(CarparkError)
	}

	_, PresentCarError := db.NewCreateTable().Model((*pkg.PresentCar)(nil)).IfNotExists().Exec(ctx)
	if PresentCarError != nil {
		panic(PresentCarError)
	}
	_, ZoneError := db.NewCreateTable().Model((*pkg.Zone)(nil)).IfNotExists().Exec(ctx)
	if ZoneError != nil {
		panic(ZoneError)
	}

	_, ZoneImageError := db.NewCreateTable().Model((*pkg.ImageZone)(nil)).IfNotExists().Exec(ctx)
	if ZoneImageError != nil {
		panic(ZoneImageError)
	}

	_, CameraError := db.NewCreateTable().Model((*pkg.Camera)(nil)).IfNotExists().Exec(ctx)
	if CameraError != nil {
		panic(CameraError)
	}

	_, CarDetailError := db.NewCreateTable().Model((*pkg.CarDetail)(nil)).IfNotExists().Exec(ctx)
	if CarDetailError != nil {
		panic(CarDetailError)
	}

	_, ClientCredError := db.NewCreateTable().Model((*pkg.ApiManage)(nil)).IfNotExists().Exec(ctx)
	if ClientCredError != nil {
		panic(ClientCredError)
	}

	r := pkg.SetupRouter(db)

	var host = fmt.Sprintf("%s:%d", configvar.Server.Host, configvar.Server.Port)
	if err := r.Run(host); err != nil {
		log.Err(err).Msgf("Failed to run server: %v", err)
	}

	log.Debug().Msgf("-------------------------------- # END # ------------------------------")
}
