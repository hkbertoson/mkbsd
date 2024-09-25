package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	url := "https://storage.googleapis.com/panels-api/data/20240916/media-1a-i-p~s"

	delay := func(ms int) {
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}

	response, err := http.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		fmt.Printf("⛔ Failed to fetch JSON file: %s\n", err.Error())
		return
	}
	defer response.Body.Close()

	var jsonData struct {
		Data map[string]struct {
			Dhd string `json:"dhd"`
		} `json:"data"`
	}

	err = json.NewDecoder(response.Body).Decode(&jsonData)
	if err != nil {
		fmt.Printf("⛔ Error parsing JSON: %s\n", err.Error())
		return
	}

	if jsonData.Data == nil {
		fmt.Println("⛔ JSON does not have a \"data\" property at its root.")
		return
	}


	downloadDir := filepath.Join(".", "downloads")
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		err = os.Mkdir(downloadDir, os.ModePerm)
		if err != nil {
			fmt.Printf("⛔ Error creating directory: %s\n", err.Error())
			return
		}
		fmt.Printf("📁 Created directory: %s\n", downloadDir)
	}

	fileIndex := 1 // Initialize the file index

	for key, subproperty := range jsonData.Data {
		if subproperty.Dhd != "" {
			imageUrl := subproperty.Dhd

			fmt.Println("🔍 Found image URL!")
			delay(100)


			imageUrlWithoutParams := strings.Split(imageUrl, "?")[0]

			ext := filepath.Ext(imageUrlWithoutParams)
			if ext == "" {
				ext = ".jpg"
			}

			filename := fmt.Sprintf("%d%s", fileIndex, ext)
			filePath := filepath.Join(downloadDir, filename)

			err := downloadImage(imageUrl, filePath)
			if err != nil {
				fmt.Printf("⛔ Error downloading image %s: %s\n", key, err.Error())
				continue
			}

			fmt.Printf("🖼️ Saved image to %s\n", filePath)
			fileIndex++
			delay(250)
		}
	}
}

func downloadImage(url, filePath string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image: %s", response.Status)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("error saving image: %w", err)
	}

	return nil
}

func asciiArt() {
	fmt.Println(`
 /$$      /$$ /$$   /$$ /$$$$$$$   /$$$$$$  /$$$$$$$
| $$$    /$$$| $$  /$$/| $$__  $$ /$$__  $$| $$__  $$
| $$$$  /$$$$| $$ /$$/ | $$  \ $$| $$  \__/| $$  \ $$
| $$ $$/$$ $$| $$$$$/  | $$$$$$$ |  $$$$$$ | $$  | $$
| $$  $$$| $$| $$  $$  | $$__  $$ \____  $$| $$  | $$
| $$\  $ | $$| $$\  $$ | $$  \ $$ /$$  \ $$| $$  | $$
| $$ \/  | $$| $$ \  $$| $$$$$$$/|  $$$$$$/| $$$$$$$/
|__/     |__/|__/  \__/|_______/  \______/ |_______/`)
	fmt.Println()
	fmt.Println("🤑 Starting downloads from your favorite sellout grifter's wallpaper app...")
}

func init() {
	asciiArt()
	time.Sleep(5 * time.Second)
	go main()
}
