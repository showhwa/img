package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var imageNames []string
var baseImageUrl string
var textFileURL string

func downloadTextFile(url string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)
	imageNames = nil
	for scanner.Scan() {
		line := scanner.Text()
		imageNames = append(imageNames, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func getRandomImageURL() string {
	rand.Seed(time.Now().UnixNano())
	randomImageName := imageNames[rand.Intn(len(imageNames))]
	return fmt.Sprintf("%s/%s", baseImageUrl, randomImageName)
}

func getRandomImageURLHandler(w http.ResponseWriter, r *http.Request) {
	imageURL := getRandomImageURL()
	http.Redirect(w, r, imageURL, http.StatusFound)
}

func startDownloadTimer(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		err := downloadTextFile(textFileURL)
		if err != nil {
			fmt.Println("Error downloading the text file:", err)
		} else {
			fmt.Println("Downloaded the text file at", time.Now())
		}
	}
}

func main() {

	textFileURL = os.Getenv("FILE_URL")
	USER := os.Getenv("USER")
	REPO := os.Getenv("REPO")
	FILE_PATH := os.Getenv("FILE_PATH")

	baseImageUrl = fmt.Sprintf("https://cdn.jsdelivr.net/gh/%s/%s/%s", USER, REPO, FILE_PATH)

	go startDownloadTimer(5 * time.Hour)

	err := downloadTextFile(textFileURL)
	if err != nil {
		fmt.Println("Error downloading the text file:", err)
		return
	}

	router := http.NewServeMux()
	router.HandleFunc("/", getRandomImageURLHandler)

	fmt.Println("Server started at :8088")
	http.ListenAndServe(":8088", router)
}
