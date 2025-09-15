package helpers

import (
	"math/rand"
	"strconv"
	"time"
)

// Create a local generator seeded once at package init
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomTravelImage() (path string, alt string) {
	imageNum := rng.Intn(44) + 1
	path = "/static/media/images/35mm/" + strconv.Itoa(imageNum) + ".JPG"
	alt = "Travel photo " + strconv.Itoa(imageNum)
	return
}
