package mulung

import (
	"net/http"
	"io/ioutil"
	"os"
	"regexp"
	_ "fmt"
	"encoding/csv"
	_ "errors"
	"log"
)

const BASE string = "https://simple.wikipedia.org"
const COUNTRY_LIST string = "/wiki/List_of_countries"
var client *http.Client
var thread = 2;

func scrap_base() {
	client := &http.Client{}	
	resp, err := client.Get(BASE + COUNTRY_LIST)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.StatusCode)
	}
	buff, _ := ioutil.ReadAll(resp.Body)
	var re *regexp.Regexp
	re = regexp.MustCompile(`[\n\r]`)
	cleaned := re.ReplaceAll(buff, []byte(""))
	os.WriteFile("wiki.html", cleaned, 0755)
	re = regexp.MustCompile(`mw-content-ltr mw-parser-output.*Disputed countries`)
	scoped := re.FindString(string(cleaned))
	re = regexp.MustCompile(`</span><a href="\/wiki\/[A-Za-z_]+`)
	countries := re.FindAllString(scoped, -1)
	builder := CountryInfoBuilderNew()
	csv_out, _ := os.OpenFile("countries.csv", os.O_RDWR | os.O_CREATE, 0666)
	csv_writer := csv.NewWriter(csv_out)	
	defer csv_writer.Flush()
	for _, link := range countries {
		source, err := scrap_countries(BASE + link[16:])
		if err != nil {
			log.Fatal(err)
		}
		country_info, err := builder.Src(source).Build()
		if err != nil {
			log.Fatal("nooooo" + source)
		}
		//if idx > 22 { continue }
		//if idx == 2 { break }
		csv_writer.Write(country_info.IntoArray())
	}
	
}

func scrap_countries(link string) (string, error) {
	re := regexp.MustCompile(`\n`)
	if client == nil {
		client = &http.Client{}	
	}
	resp, err := client.Get(link)
	if err != nil {
		return "", err	
	}
	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err	
	}
	cleaned := re.ReplaceAll(buff, []byte(""))
	return string(cleaned), nil
}
