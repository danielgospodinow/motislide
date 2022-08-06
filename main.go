package main

import (
	"danielgospodinow/motislide/utils"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/reujab/wallpaper"
)

func main() {
	updateIntervalPtr := flag.Int("interval", 60, "Background update interval in minutes.")
	imageDirectoryPtr := flag.String("imgdir", "./images", "Location of folder with background wallpapers.")
	flag.Parse()

	log.Println("Starting the background changer process...")
	start(*updateIntervalPtr, *imageDirectoryPtr)
	log.Println("Application finished successfully.")
}

// start schedules a job triggering background image update.
func start(updateInterval int, imageDirectory string) {
	execute := make(chan struct{}, 1)
	execute <- struct{}{}
	ticker := time.NewTicker(time.Duration(updateInterval) * time.Minute)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	originalBackground, err := getBackground()
	if err != nil {
		log.Fatalln("Failed to get original background.", err)
	}

	imageFiles, _ := ioutil.ReadDir(imageDirectory)
	if len(imageFiles) <= 1 {
		log.Fatalln("The set of provided images is too small, please add some more.")
	}

	for {
		select {
		case <-execute:
			log.Println("Changing background wallpaper...")
			changeBackgroundFromDirectory(imageDirectory)
		case <-ticker.C:
			execute <- struct{}{}
		case <-exit:
			log.Println("Stopping the background change job...")

			log.Println("Restoring original background image...")
			err := setBackground(originalBackground)
			if err != nil {
				log.Println("Failed to restore original background image.", err)
			}

			return
		}
	}
}

// changeBackgroundFromDirectory fetches a new random wallpaper image which is different from the current background image and sets it as the new background image.
func changeBackgroundFromDirectory(imageDirectory string) {
	background, err := getBackground()
	if err != nil {
		return
	}

	for {
		newImagePath, err := utils.GetRandomImageAbsPath(imageDirectory)
		if err != nil {
			log.Println("Failed to fetch a new image.", err)
			return
		}

		if newImagePath == background {
			continue
		}

		err = setBackground(newImagePath)
		if err != nil {
			log.Printf("Failed to set background to '%s'. %s", newImagePath, err)
			return
		}

		log.Println("Wallpaper changed successfully.")
		break
	}
}

// setBackground changes the current background with a new one.
func setBackground(file string) error {
	err := wallpaper.SetFromFile(file)
	if err != nil {
		log.Printf("Failed to set wallpaper '%s'.", file)
		return err
	}

	return nil
}

// getBackground retrieves the current background image's file path.
func getBackground() (string, error) {
	background, err := wallpaper.Get()
	if err != nil {
		log.Println("Failed to get current wallpaper.", err)
		return "", err
	}

	return background, nil
}
