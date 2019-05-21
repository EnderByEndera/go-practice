package main

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type res struct {
	X        int `bson:"Y"`
	Y        int `bson:"X"`
	Priority int `bson:"Priority"`
	ID       int `bson:"ID"`
} // position of the resource

type priores struct {
	resource []res
	dic      [num][num]float64 // distance from two resources
	N        int               // real resources
	seq      sequence          // the resource sequences of the best path
	answer   ans               // the length of the best path
}

// Three calculating resources with mutex
type message struct {
	mess [num][num]float64
	mux  sync.Mutex
}

type ans struct {
	answer float64
	mux    sync.Mutex
}

type sequence struct {
	seq [num]int
	mux sync.Mutex
}

const (
	width      int = 1000 // width of the office
	height     int = 1000 // height of the office
	num        int = 1000 // max resource number
	comma      int = 44   // value of comma
	retrn      int = 10   // value of '\n'
	initnumber     = 100  // initial number
)

var (
	resource1 priores
	resource2 priores
	resource3 priores
	resource0 priores
	resource  priores
	count     = 0
	choice    = 0
	ratio     float64
	router    = gin.Default()
)
