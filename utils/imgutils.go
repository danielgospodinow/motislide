package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// GetRandomImageFilepath retrieves the filepath of a random image from a given image directory.
func GetRandomImageFilepath(imageDirectory string) (string, error) {
	files, err := ioutil.ReadDir(imageDirectory)
	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })

	for _, file := range files {
		fullFilePath := imageDirectory + "/" + file.Name()

		isImage, err := IsImageFile(fullFilePath)
		if err != nil {
			return "", err
		}

		if isImage {
			return fullFilePath, nil
		}
	}

	return "", errors.New("Images folder is empty or it doesn't contain an image file")
}

// IsImageFile checks whether a given filepath corresponds to an image file.
func IsImageFile(imageFilepath string) (bool, error) {
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
