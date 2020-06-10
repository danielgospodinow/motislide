package main

import (
	"danielgospodinow/motislide/utils"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
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
		newImageRelPath, err := utils.GetRandomImageFilepath(imageDirectory)
		if err != nil {
			fmt.Println("Failed to fetch a new image")
			fmt.Println(err.Error())

			return
		}

		if newImageRelPath == background {
			imageFiles, _ := ioutil.ReadDir(imageDirectory)
			if len(imageFiles) <= 1 {
				fmt.Println("The set of provided images is too small, please add some more")

				return
			}

			continue
		}

		newImageAbsPath, _ := filepath.Abs(newImageRelPath)
		wallpaper.SetFromFile(newImageAbsPath)
		break
	}

	fmt.Println("Wallpaper changed successfully")
}
