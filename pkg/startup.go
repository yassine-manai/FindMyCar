package pkg

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

var ctx = context.Background()
var Zonelist []int
var CarParkList []int
var CameraList []int

func Loadzonelist(db *bun.DB) {
	log.Debug().Msgf("Prepare Zone List \n")
	ZoneService, _ := GetAllZone(ctx, db)
	for _, v := range ZoneService {
		Zonelist = append(Zonelist, v.ID)
	}
}

func LoadCarparklist(db *bun.DB) {
	log.Debug().Msgf("Prepare Car Park List \n")
	CarService, _ := GetAllCarparks(ctx, db)
	for _, v := range CarService {
		CarParkList = append(CarParkList, v.ID)

	}
}

func LoadCameralist(db *bun.DB) {
	log.Debug().Msgf("Prepare Camera List \n")
	CameraService, _ := GetAllCamera(ctx, db)
	for _, v := range CameraService {
		CameraList = append(CameraList, v.ID)

	}
}
