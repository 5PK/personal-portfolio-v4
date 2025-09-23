package helpers

import (
	"math/rand"
	"time"
)

// Create a local generator seeded once at package init
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomTravelImage() (pathsArr [6]string, altArr [6]string) {
	imageNames := []string{
		"000061740001.jpg", "000061740016.jpg", "000061740018.jpg", "000061740022.jpg",
		"000061740023.jpg", "000061740029.jpg", "000061750024.jpg", "000061750031.jpg",
		"000061760011.jpg", "000061760012.jpg", "000061760022.jpg", "000061770004.jpg",
		"000061770018.jpg", "000061770021.jpg", "000061780016.jpg", "000061780018.jpg",
		"000061780020.jpg", "000061790010.jpg", "000061790026.jpg", "000061820009.jpg",
		"000089950005.jpg",
	}

	for i := 0; i < 6; i++ {
		randomIndex := rng.Intn(len(imageNames))
		selectedImage := imageNames[randomIndex]
		imageNames = append(imageNames[:randomIndex], imageNames[randomIndex+1:]...)
		pathsArr[i] = "/static/media/images/35mm/" + selectedImage
		altArr[i] = "Travel photo " + selectedImage
	}

	return
}
