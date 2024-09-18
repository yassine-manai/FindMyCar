package pkg

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

var Zonelist []int
var CarParkList []int
var CameraList []int
var Dbg *bun.DB

func Loadzonelist() {
	log.Debug().Msgf("Prepare Zone List \n")
	ctx := context.Background()

	ZoneService, _ := GetAllZone(ctx)
	for _, v := range ZoneService {
		Zonelist = append(Zonelist, *v.ZoneID)
	}
}

func LoadCarparklist() {
	log.Debug().Msgf("Prepare Car Park List \n")
	ctx := context.Background()

	CarService, _ := GetAllCarparks(ctx)
	for _, v := range CarService {
		CarParkList = append(CarParkList, v.ID)
	}
}

func LoadCameralist() {
	log.Debug().Msgf("Prepare Camera List \n")
	ctx := context.Background()

	CameraService, _ := GetAllCamera(ctx)
	for _, v := range CameraService {
		CameraList = append(CameraList, v.ID)
	}
}
