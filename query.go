package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type config struct{
    IndyUrl string `yaml:"indy_url"`
	DAGroup string `yaml:"da_group"`
	MaxConcurrentGoroutines int `yaml:"max_concurrent_goroutines"`
}

type entry struct {
	Name    string   `yaml:"name"`
	Version string   `yaml:"version"`
	Gavs    []string `yaml:"gavs"`
}

type dataset struct {
	Artifacts []entry `yaml:"artifacts"`
}

func loadData() dataset {

	fmt.Println("load data from dataset.yaml")

	yamlFile, err := ioutil.ReadFile("dataset.yaml")
	if err != nil {
		fmt.Println("yamlFile.Get err   #%v ", err)
	}

	ds := dataset{}
	err = yaml.Unmarshal(yamlFile, &ds)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return ds
}

func loadConfig() config {
	fmt.Println("load config from config.yaml")

	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("configFile.Get err   #%v ", err)
	}

	c := config{}
	err = yaml.Unmarshal(configFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	fmt.Println(c)
	return c
}

func deleteMetadata(url string) string{
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ5emxOa0tBUmZlMUVzcHZJbU9rdkVTeUttVGl6N05MTWp2Z3lSSGJEZHVBIn0.eyJqdGkiOiI5ZGExOTFjNC05YmRjLTRhNzUtYjY2ZC05NzZmNjI2ODJlZmYiLCJleHAiOjE2MjMyMDg2ODMsIm5iZiI6MCwiaWF0IjoxNjIzMDM1ODgzLCJpc3MiOiJodHRwczovL3NlY3VyZS1zc28tbmV3Y2FzdGxlLXN0YWdlLnBzaS5yZWRoYXQuY29tL2F1dGgvcmVhbG1zL3BuY3JlZGhhdCIsImF1ZCI6WyJwbmN3ZWIiLCJwbmNpbmR5IiwiYWNjb3VudCIsInBuY3Jlc3QiXSwic3ViIjoiZmYzNWFiODMtZTE5ZS00ZTVmLWFlZjctNjc0YjUyZDRmZWUzIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoicG5jaW5keXVpIiwibm9uY2UiOiIxZjgxYWY0Ni03NzNhLTRhZjUtOTIyNi1kOTMyNmRkZDExNzgiLCJhdXRoX3RpbWUiOjE2MjMwMzU4ODEsInNlc3Npb25fc3RhdGUiOiJmZmU1OWNlMi0wZTRjLTRiZWQtOWUxZS0xNjI4YzNlZDhkNzIiLCJhY3IiOiIxIiwiYWxsb3dlZC1vcmlnaW5zIjpbIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbInVzZXIiXX0sInJlc291cmNlX2FjY2VzcyI6eyJwbmN3ZWIiOnsicm9sZXMiOlsidXNlciJdfSwicG5jaW5keXVpIjp7InJvbGVzIjpbInBuY2luZHlhZG1pbiIsInBuY2luZHl1c2VyIl19LCJwbmNpbmR5Ijp7InJvbGVzIjpbInBuY2luZHlhZG1pbiIsInBvd2VyLXVzZXIiLCJwbmNpbmR5dXNlciJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19LCJwbmNyZXN0Ijp7InJvbGVzIjpbInVzZXIiXX19LCJzY29wZSI6Im9wZW5pZCIsInByZWZlcnJlZF91c2VybmFtZSI6IndndW8iLCJlbWFpbCI6IndndW9AaXBhLnJlZGhhdC5jb20ifQ.arH2j_x1JpFv9AgaKnsscPh8GDHld_9I8sBV3Kz23O529cCeh9LHeD3czZG1LB4H09zAe2QRpOKYbZkSEt7lrmBmHQJF2AD5k-IIu-Ds3roe046yWI4rd5xKGPiqlQc035sndu9CAKRZPOjyfmSveWtpt8PQKkbGi8Vq2-RCWrR_uE1VTINAtcsh_5k8_2yAqoktZZqV_SFfy3P2U-2_yVrG9sq6odzWWTcaEcybuU_OGPFXTns1WT3fpJ2ff2N2NBmudYWGND8qX67-thaRFPqK-sVAf4XfsS89S63nGFX-46xSsOk8gyYSWuYCXStTI79NAlY2ONv5GCaCIzShWQ")

	var c http.Client
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
	return "Done"
}

func fetchMetadata(gav string, url string) string {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/xml")

	var c http.Client
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("XML:")
	//io.Copy(os.Stdout, resp.Body)

	// Create new file
	file := strings.ReplaceAll(gav, ":", "-")
	tmp, err := os.Create("results/" + file + ".xml")
	if err != nil {
		log.Fatal(err)
	}
	defer tmp.Close()

	bytesWritten, err := io.Copy(tmp, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bytes Written: %d\n", bytesWritten)

	return "Done"
}

func doRun(c config, ds dataset){
	indyUrl := c.IndyUrl
	daGroup := c.DAGroup

	results := make(chan string)

	jobs := 0
	var urls [10000]string
	var gavs [10000]string

	for idx, element := range ds.Artifacts {
		fmt.Println(idx, element.Name)
		for idx_, gav := range element.Gavs {
		    gavs[jobs] = gav
			fmt.Println(idx_, gav)
			s := strings.Split(gav, ":")
			groupId := strings.ReplaceAll(s[0], ".", "/")
			artifactId := s[1]
			url := fmt.Sprintf("%s/api/content/maven/group/%s/%s/%s/maven-metadata.xml", indyUrl, daGroup, groupId, artifactId)
			urls[jobs] = url
			jobs = jobs+1
		}
	}

	fmt.Println("Total jobs:")
    fmt.Println(jobs)

	concurrentGoroutines := make(chan struct{}, c.MaxConcurrentGoroutines)
	var wg sync.WaitGroup

	for i := 0; i < jobs; i++ {
		concurrentGoroutines <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println("doing", i)
			start := time.Now()
			fetchMetadata(gavs[i], urls[i])
			//deleteMetadata(urls[i])
			elapsed := time.Since(start)
			fmt.Println("finished", i)
			fmt.Println("took", elapsed)
			<-concurrentGoroutines
		}(i)
	}

	for i := 0; i < jobs; i++ {
		fmt.Println(<-results)
	}

	wg.Wait()
}

func mainA() {
	c := loadConfig()
	ds := loadData()
	doRun(c, ds)
}
