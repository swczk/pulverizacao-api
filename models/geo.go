package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type GeoPoint struct {
	Latitude  float64   `json:"latitude" bson:"latitude"`
	Longitude float64   `json:"longitude" bson:"longitude"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Altitude  *float64  `json:"altitude,omitempty" bson:"altitude,omitempty"`
	Speed     *float64  `json:"speed,omitempty" bson:"speed,omitempty"`
	Accuracy  *float64  `json:"accuracy,omitempty" bson:"accuracy,omitempty"`
}

type GeoTrajetoria struct {
	AplicacaoID         bson.ObjectID `json:"aplicacaoId" bson:"aplicacao_id"`
	PontoInicial        GeoPoint           `json:"pontoInicial" bson:"ponto_inicial"`
	PontoFinal          *GeoPoint          `json:"pontoFinal,omitempty" bson:"ponto_final,omitempty"`
	Trajetoria          []GeoPoint         `json:"trajetoria" bson:"trajetoria"`
	AreaCobertura       float64            `json:"areaCobertura" bson:"area_cobertura"`
	DistanciaPercorrida float64            `json:"distanciaPercorrida" bson:"distancia_percorrida"`
	CreatedAt           time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt           time.Time          `json:"updatedAt" bson:"updated_at"`
}
