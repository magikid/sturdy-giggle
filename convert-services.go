package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

const inputDirectory = "Raw Services"
const inputExtension = "wav"
const outputDirectory = "Converted Services"
const outputExtension = "mp3"

func readFilesInDir(dirName string, extension string) []string {
	files, err := os.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}

	lowercaseExtension := strings.ToLower(extension)
	var matchingFiles []string
	for _, f := range files {
		if strings.Contains(strings.ToLower(f.Name()), lowercaseExtension) {
			matchingFiles = append(matchingFiles, f.Name())
		}
	}

	return matchingFiles
}

func getFilenamesWithoutExtensions(filenames []string) []string {
	var filenamesWithoutExtensions []string

	for _, f := range filenames {
		fileName := f[:len(f)-len(filepath.Ext(f))]
		filenamesWithoutExtensions = append(filenamesWithoutExtensions, fileName)
	}

	return filenamesWithoutExtensions
}

func convertService(filename string) {
	inputFileWithExtension := fmt.Sprintf("%s.%s", filename, inputExtension)
	inputFileWithPath := filepath.Join(inputDirectory, inputFileWithExtension)

	outputFileWithExtension := fmt.Sprintf("%s.%s", filename, outputExtension)
	outputFileWithPath := filepath.Join(outputDirectory, outputFileWithExtension)

	err := ffmpeg_go.Input(inputFileWithPath).Output(outputFileWithPath, ffmpeg_go.KwArgs{"vn": "", "ar": 44100, "ac": 2, "channel_layout": "stereo", "q:a": 2}).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		log.Printf("Error processing %s\n", inputFileWithExtension)
		log.Println(err)
	}
}

// found https://stackoverflow.com/a/45428032
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func convertServices(services []string) {
	for _, service := range services {
		convertService(service)
	}
}

func convertAllServices() {
	rawServices := readFilesInDir(inputDirectory, inputExtension)
	log.Printf("raw services: %s\n", strings.Join(rawServices, ", "))
	log.Printf("Converting the following services: %s\n\n", strings.Join(rawServices, ", "))
	convertServices(rawServices)
}

func convertNewServices() {
	convertedServices := readFilesInDir(outputDirectory, outputExtension)
	rawServices := readFilesInDir(inputDirectory, inputExtension)

	convertedServicesWithoutExtensions := getFilenamesWithoutExtensions(convertedServices)
	rawServicesWithoutExtensions := getFilenamesWithoutExtensions(rawServices)

	servicesToConvert := difference(rawServicesWithoutExtensions, convertedServicesWithoutExtensions)

	if len(servicesToConvert) == 0 {
		log.Println("No new services to convert")
		return
	}

	log.Printf("Converting the following services: %s\n\n", strings.Join(servicesToConvert, ", "))

	convertServices(servicesToConvert)
}

func main() {
	allRawServices := flag.Bool("all", false, "Convert all raw services found")
	flag.Parse()

	if *allRawServices {
		convertAllServices()
	} else {
		convertNewServices()
	}
}
