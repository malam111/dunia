package mulung

import (
	"net/http"
	"io/ioutil"
	"os"
	"regexp"
	"fmt"
)

const BASE string = "https://simple.wikipedia.org"
const COUNTRY_LIST string = "/wiki/List_of_countries"
var client *http.Client

func scrap() {
	client := &http.Client{}	
	resp, err := client.Get(BASE + COUNTRY_LIST)
	if err != nil {
		panic("Error")
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.StatusCode)
	}
	buff, _ := ioutil.ReadAll(resp.Body)
	var re *regexp.Regexp
	re = regexp.MustCompile(`[\n\r]`)
	cleaned := re.ReplaceAll(buff, []byte(" "))
	os.WriteFile("wiki.html", cleaned, 0755)
	re = regexp.MustCompile(`mw-content-ltr mw-parser-output.*Disputed countries`)
	scoped := re.FindString(string(cleaned))
	re = regexp.MustCompile(`</span><a href="\/wiki\/[A-Z][a-z]+`)
	countries := re.FindAllString(scoped, -1)
	fmt.Print(countries)
	//re_link := regexp.MustCompile("/wiki/[A-Z][a-z]+")
	//for _, txt := range countries {
		//countryInfo(re_link.FindString(txt))
	//}

}
