package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	defaultDownloadTime = 5
	serverPort          = 8088
)

var (
	imageNames   []string
	baseImageUrl string
	textFileURL  string
	log          = logrus.New()
)

func init() {
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

func downloadTextFile(url string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download text file: %v", err)
	}
	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)
	imageNames = nil
	for scanner.Scan() {
		line := scanner.Text()
		imageNames = append(imageNames, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read text file: %v", err)
	}

	return nil
}

func getRandomImageURL() string {
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
		if err := downloadTextFile(textFileURL); err != nil {
			log.Errorf("Error downloading the text file: %v", err)
		} else {
			log.Info("Image number is ", len(imageNames))
			log.Info("Downloaded the text file at", time.Now())
		}
	}
}

func setupRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", getRandomImageURLHandler)
	return router
}

func main() {
	textFileURL = os.Getenv("FILE_URL")
	USER := os.Getenv("USER")
	REPO := os.Getenv("REPO")
	FILE_PATH := os.Getenv("FILE_PATH")

	downloadTimeStr := os.Getenv("downloadTime")
	downloadTime, err := strconv.Atoi(downloadTimeStr)
	if err != nil || downloadTime <= 0 {
		downloadTime = defaultDownloadTime
	}

	baseImageUrl = fmt.Sprintf("https://cdn.jsdelivr.net/gh/%s/%s/%s", USER, REPO, FILE_PATH)

	go startDownloadTimer(time.Duration(downloadTime) * time.Hour)

	err = downloadTextFile(textFileURL)
	if err != nil {
		log.Errorf("Error downloading the text file: %v", err)
		return
	}

	router := setupRoutes()
	log.Infof("FILE_URL = %s, USER = %s, REPO = %s, FILE_PATH = %s, downloadTime = %d", textFileURL, USER, REPO, FILE_PATH, downloadTime)
	log.Infof("Server started at :%d", serverPort)
	log.Info("Image number is ", len(imageNames))
	http.ListenAndServe(fmt.Sprintf(":%d", serverPort), router)
}
