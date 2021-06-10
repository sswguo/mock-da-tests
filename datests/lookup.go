package datests

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"io"
)

type config struct{
	PncRest string `yaml:"pnc_rest_url"`
	IndyUrl string `yaml:"indy_url"`
	DAGroup string `yaml:"da_group"`
	MaxConcurrentGoroutines int `yaml:"max_concurrent_goroutines"`
}

func LoadConfig() config {
	fmt.Println("load config from config.yaml")

	configFile, err := ioutil.ReadFile("config-local.yaml")
	if err != nil {
		fmt.Println("configFile.Get err #%v ", err)
	}

	c := config{}
	err = yaml.Unmarshal(configFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	fmt.Println(c)
	return c
}

func GetAlignLog(url string) string {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "text/plain")

	var c http.Client
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(responseData)
}

func LookupMetadata(gav string, url string) string {
	fmt.Println(url)
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

	// dump the metadata file locally for verifying
	tempArray := strings.Split(gav, "=")
	file := strings.ReplaceAll(tempArray[0], ":", "-")
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