package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/reujab/wallpaper"
)

func main() {
	updateIntervalPtr := flag.Int("int", 60, "Background update interval in minutes.")
	imageDirectoryPtr := flag.String("imgdir", "./images", "Location of folder with background wallpapers.")
	flag.Parse()

	fmt.Println("Starting the background changer process ...")
	scheduleBackgroundUpdate(*updateIntervalPtr, *imageDirectoryPtr)
}

// scheduleBackgroundUpdate schedules a job triggering background image update.
func scheduleBackgroundUpdate(updateInterval int, imageDirectory string) {
	for {
		fmt.Println("Changing background wallpaper ...")
		changeBackground(imageDirectory)
		time.Sleep(time.Duration(updateInterval) * time.Minute)
	}
}

// changeBackground fetches a new wallpaper image which is different from the current
// background image and sets it as the new background image.
func changeBackground(imageDirectory string) {
	background, err := wallpaper.Get()
	if err != nil {
		fmt.Println("Failed to fetch current wallpaper")
		fmt.Println(err.Error())

		return
	}

	for {
		newImagePath, err := getRandomImageFilepath(imageDirectory)
		if err != nil {
			fmt.Println("Failed to fetch a new image")
			fmt.Println(err.Error())

			return
		}

		if newImagePath == background {
			imageFiles, _ := ioutil.ReadDir(imageDirectory)
			if len(imageFiles) <= 1 {
				fmt.Println("The set of provided images is too small, please add some more")

				return
			}

			continue
		}

		wallpaper.SetFromFile(newImagePath)
		break
	}

	fmt.Println("Wallpaper changed successfully")
}

// getRandomImageFilepath retrieves the filepath of a random image from a given image directory.
func getRandomImageFilepath(imageDirectory string) (string, error) {
	files, err := ioutil.ReadDir(imageDirectory)
	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })

	for _, file := range files {
		fullFilePath := imageDirectory + "/" + file.Name()

		isImage, err := isImageFile(fullFilePath)
		if err != nil {
			return "", err
		}

		if isImage {
			return fullFilePath, nil
		}
	}

	return "", errors.New("Images folder is empty or it doesn't contain an image file")
}

// checkFileIfImage checks whether a given filepath corresponds to an image file.
func isImageFile(imageFilepath string) (bool, error) {
	file, err := os.Open(imageFilepath)
	if err != nil {
		fmt.Println("Failed to open file")
		fmt.Println(err)

		return false, err
	}

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		fmt.Println("Failed to read from file")
		fmt.Println(err)

		return false, err
	}

	switch filetype := http.DetectContentType(buff); filetype {
	case "image/jpeg", "image/jpg", "image/png":
		return true, nil
	default:
		return false, nil
	}
}
