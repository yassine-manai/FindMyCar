package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"fmc/config"
	_ "fmc/docs"
	"fmc/pkg"
	"fmc/routes"
)

// @title           Find Your Car
// @version         1.0
// @securityDefinitions.basic  BasicAuth

func main() {
	config.InitLogger()
	fmt.Println("----------------------------- # START PROGRAM # ------------------------------")

	var configvar config.ConfigFile

	if err := configvar.Load(); err != nil {
		log.Err(err).Msgf("Error loading config: %v", err)
	} else {
		fmt.Println("Success fetching config data")
	}

	fmt.Printf("Server running on %s:%d ", configvar.Server.Host, configvar.Server.Port)
	fmt.Printf("Database connecting to %s:%d", configvar.Database.Host, configvar.Database.Port)

	var dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", configvar.Database.User, configvar.Database.Password, configvar.Database.Host, configvar.Database.Port, configvar.Database.Name)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	pkg.Dbg = bun.NewDB(sqldb, pgdialect.New())

	if err := pkg.Dbg.Ping(); err != nil {
		fmt.Println(err)
	}

	/* // List of table models to be created
	models := []interface{}{
		(*pkg.Carpark)(nil),
		(*pkg.PresentCar)(nil),
		(*pkg.Zone)(nil),
		(*pkg.ImageZone)(nil),
		(*pkg.Camera)(nil),
		(*pkg.CarDetail)(nil),
		(*pkg.ApiManage)(nil),
		(*pkg.PresentCarHistory)(nil),
	}
	ctx := context.Background()
	if err := functions.CreateTables(ctx, pkg.Dbg, models); err != nil {
		fmt.Printf("Failed to create tables: %v", err)
	} */
	ctx := context.Background()

	_, PresentCarError := pkg.Dbg.NewCreateTable().Model((*pkg.Settings)(nil)).IfNotExists().Exec(ctx)
	if PresentCarError != nil {
		panic(PresentCarError)
	}

	fmt.Println("Tables created successfully")

	// Data in list in startup
	fmt.Println("-------------------------------- # DATA LIST START # ------------------------------")
	pkg.Loadzonelist()
	fmt.Println("Zonelist Length", len(pkg.Zonelist))

	pkg.LoadCameralist()
	fmt.Println("CameraList Length", len(pkg.CameraList))

	fmt.Println("-------------------------------- #  DATA LIST END # ------------------------------")

	r := routes.SetupRouter()

	var host = fmt.Sprintf("%s:%d", configvar.Server.Host, configvar.Server.Port)
	if err := r.Run(host); err != nil {
		log.Err(err).Msgf("Failed to run server: %v", err)
	}

	log.Debug().Msgf("-------------------------------- # END PROGRAM # ------------------------------")
}
