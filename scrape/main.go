package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
 * Url example: https://www.googleapis.com/customsearch/v1?key=AIzaSyDD-LtlrvwUSlo-PDyIyJ-h4UNCUepe934&cx=e3f537cd794f84eaf&q=sunflower&searchType=image
 * Image example: https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSgUnFs-IARzAn9IGn9jUMsVjZWELuhvArVjA&usqp=CAU
 * APIKEY : AIzaSyDD-LtlrvwUSlo-PDyIyJ-h4UNCUepe934
**/

const APIKEY string = "AIzaSyCo9RpDw43ZcEbBMUQgCriM1WgwCz7NoY0"

func main() {
	imageCount := 200 // MAX: 200
	now := 1

	for {
		result := getImageUrls("dandelion flower", 10, now)
		log.Println("Success getting image urls")
		for _, res := range result {
			saveImageFromUrl(res, "../datasets/dandelion")
		}

		now += 10

		if now >= imageCount {
			break
		}
	}

}

type ImageResult struct {
	url string
	ext string
}

func getImageUrls(q string, num int, start int) []ImageResult {
	q = strings.ReplaceAll(q, " ", "%20")

	if num > 10 {
		log.Panicln("Error: num must be less or equal than 10")
	}

	var result []ImageResult
	// TODO: "&imgDominantColor=pink" implement imgDominantColor Later
	url := "https://www.googleapis.com/customsearch/v1?key=" + APIKEY + "&cx=e3f537cd794f84eaf&q=" + q + "&searchType=image&start=" + strconv.Itoa(start) + "&num=" + strconv.Itoa(num)

	response, err := http.Get(url)
	if err != nil {
		log.Panicln("Error while getting data", url, "-", err)
	}

	if response.StatusCode != 200 {
		log.Panicln("Error while getting data", url, "-", response.StatusCode)
	}

	defer response.Body.Close()

	read, err := io.ReadAll(response.Body)
	if err != nil {
		log.Panicln("Error while read data", url, "-", response.StatusCode)
	}

	var data map[string]interface{}
	json.Unmarshal(read, &data)

	items := data["items"].([]interface{})

	for _, item := range items {
		dummy := item.(map[string]interface{})
		image := dummy["image"].(map[string]interface{})

		ext := strings.Split(dummy["fileFormat"].(string), "/")[1]

		result = append(result, ImageResult{
			url: image["thumbnailLink"].(string),
			ext: ext,
		})
	}

	return result
}

// saveImageFromUrl downloads an image from a url and saves it to a path
func saveImageFromUrl(image ImageResult, path string) {
	response, err := http.Get(image.url)
	if err != nil {
		log.Println("Error while downloading", image.url, "-", err)
	}

	var filename string
	if response.Header.Get("Content-Disposition") != "" {
		filename = response.Header.Get("Content-Disposition")
	} else {
		if image.ext == "" {
			image.ext = "jpeg"
		}

		if image.ext == "webp" { // TODO: remove this
			image.ext = "jpeg"
		}

		filename = strconv.FormatInt(time.Now().UnixNano(), 10) + "." + image.ext
	}

	defer response.Body.Close()

	file, err := os.Create(path + "/" + filename)
	if err != nil {
		log.Println("Error while creating", filename, "-", err)
	}

	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		os.Remove(path + "/" + filename)
		log.Println("Error while write file", image.url, "-", err)
	}

	log.Println("Downloaded", image.url, "to", filename, "in", path, "folder")
}
