package mulung

import (
	"regexp"
	"fmt"
	"errors"
	"strings"
	"strconv"
)


type CountryInfo struct {
	Name string `json:"name"` 
	Capital string `json:"capital"` 
	OfficialLanguages []string `json:"officialLanguage"` 
	Area int `json:"area"` 
	Population int `json:"population"`
	Currency string `json:"currency"`
	DrivingSide string `json:"drivingSide"`
	CallingCode int `json:"callingCode`
	TLD string `json:"tld"`
}

func (self *CountryInfo) Serialize() string {
	return fmt.Sprintf("%s,%s,%s,%d,%d,%s,%s,%d,%s", 
		self.Name, 
		self.Capital, 
		self.OfficialLanguages, 
		self.Area, 
		self.Population, 
		self.Currency, 
		self.DrivingSide, 
		self.CallingCode, 
		self.TLD)
}

func (self *CountryInfo) IntoArray() []string {
	return []string {
		self.Name, 
		self.Capital, 
		strings.Join(self.OfficialLanguages, ";"), 
		strconv.Itoa(self.Area), 
		strconv.Itoa(self.Population),
		self.Currency, 
		self.DrivingSide, 
		strconv.Itoa(self.CallingCode), 
		self.TLD}
}

type CountryInfoBuilder struct {
	src string
	NameExtractor *InfoName
	CapitalExtractor *InfoCapital
	LangExtractor *InfoLang
	AreaExtractor *InfoArea
	PopulationExtractor *InfoPop
	CurrencyExtractor *InfoCur
	DrivingSideExtractor *InfoDriv
	CallingCodeExtractor *InfoCall
	TLDExtractor *InfoTLD
	
}

func CountryInfoBuilderNew() *CountryInfoBuilder {
	return &CountryInfoBuilder {
		src: "",
		NameExtractor: NewInfoName(),
		CapitalExtractor: NewInfoCapital(),
		LangExtractor: NewInfoLang(),
		AreaExtractor: NewInfoArea(),
		PopulationExtractor: NewInfoPop(),
		CurrencyExtractor: NewInfoCur(),
		DrivingSideExtractor: NewInfoDriv(),
		CallingCodeExtractor: NewInfoCall(),
		TLDExtractor: NewInfoTLD(),
	}
}

func (self *CountryInfoBuilder) Src(src string) *CountryInfoBuilder {
	self.src = src
	return self
}

func (self *CountryInfoBuilder) Build() (*CountryInfo, error) {
	if self.src == "" {
		return nil, errors.New("src is empty")
	}
	if self.NameExtractor == nil {
		return nil, errors.New("NameExtractor not Found")
	}
	if self.CapitalExtractor == nil {
		return nil, errors.New("CapitalExtractor not Found")
	}
	if self.LangExtractor == nil {
		return nil, errors.New("LangExtractor not Found")
	}
	if self.AreaExtractor == nil {
		return nil, errors.New("AreaExtractor not Found")
	}
	if self.PopulationExtractor == nil {
		return nil, errors.New("PopulationExtractor not Found")
	}
	if self.CurrencyExtractor == nil {
		return nil, errors.New("CurrencyExtractor not Found")
	}
	if self.DrivingSideExtractor == nil {
		return nil, errors.New("DrivingSideExtractor not Found")
	}
	if self.CallingCodeExtractor == nil {
		return nil, errors.New("CallingCodeExtractor not Found")
	}
	if self.TLDExtractor == nil {
		return nil, errors.New("TLDExtractor not Found")
	}
	
	return &CountryInfo {
		Name: self.NameExtractor.GetInfo(self.src),
		Capital: self.CapitalExtractor.GetInfo(self.src),
		OfficialLanguages: strings.Split(self.LangExtractor.GetInfo(self.src), string(',')),
		Area: self.AreaExtractor.GetInfo(self.src),
		Population: self.PopulationExtractor.GetInfo(self.src),
		Currency: self.CurrencyExtractor.GetInfo(self.src),
		DrivingSide: self.DrivingSideExtractor.GetInfo(self.src),
		CallingCode: self.CallingCodeExtractor.GetInfo(self.src),
		TLD: self.TLDExtractor.GetInfo(self.src)}, nil 
}

type Extract interface {
	GetInfo() string	
}

type InfoName struct {
	regexes *regexp.Regexp		
}

func NewInfoName() *InfoName {
	return &InfoName {
		// get this
		regexes: regexp.MustCompile(`country-name">[\w ]+</div>`)}
		// remove extras
}

func (self *InfoName) GetInfo(src string) string {
	temp := self.regexes.FindString(src)
	if temp == "" {
		return ""
	}
	temp = temp[14:len(temp)-6]
	return temp
}

type InfoCapital struct {
	regexes []*regexp.Regexp
}

func NewInfoCapital() *InfoCapital {
	return &InfoCapital {
		regexes: []*regexp.Regexp {
			regexp.MustCompile(`Capital.*Official.{0,10}languages`),
			regexp.MustCompile(`title="[\w-]+`)}}
}

func (self *InfoCapital) GetInfo(src string) string {
	temp := self.regexes[0].FindIndex([]byte(src))
	if temp == nil {
		return ""
	}
	capital := self.regexes[1].FindString(src[temp[0]:temp[1]])
	return capital[7:]
}

type InfoLang struct {
	regexes []*regexp.Regexp
}

func NewInfoLang() *InfoLang {
	return &InfoLang {
		regexes: []*regexp.Regexp {
			// get this area
			regexp.MustCompile(`Official.{0,10}languages.*Ethnic.{0,10}groups`),
			// get every one of these
			regexp.MustCompile(`title="[\w ]+`)}}
			// remove extras
}

func (self *InfoLang) GetInfo(src string) string {
	temp := self.regexes[0].FindIndex([]byte(src))
	if temp == nil {
		return ""
	}
	titles := self.regexes[1].FindAllString(src[temp[0]:temp[1]], -1)
	if titles == nil {
		return ""
	}
	langs := ""
	last := len(titles)
	for idx, lang := range(titles) {
		langs = langs + lang[7:]
		if idx != last-1 {
			langs = langs + ","
		}
	}
	return langs
}

type InfoArea struct {
	regexes []*regexp.Regexp
}

func NewInfoArea() *InfoArea {
	return &InfoArea {
		regexes: []*regexp.Regexp{
			regexp.MustCompile(`Area.*Population`),
			regexp.MustCompile(`[\d,]+.{0,10}km`),
			regexp.MustCompile(`,`),
			regexp.MustCompile(`[\d]+`)}}
}

func (self *InfoArea) GetInfo(src string) int {
	temp := self.regexes[0].FindIndex([]byte(src))
	if temp == nil {
		return -1
	}
	str_area := self.regexes[1].FindString(src[temp[0]:temp[1]])
	nocom := self.regexes[2].ReplaceAll([]byte(str_area), []byte(""))
	area := self.regexes[3].FindString(string(nocom))
	ret, err := strconv.Atoi(area)
	if err != nil {
		return 0
	}
	return ret
}

type InfoPop struct {
	regexes []*regexp.Regexp
}

func NewInfoPop() *InfoPop {
	return &InfoPop {
		regexes: []*regexp.Regexp {	
			regexp.MustCompile(`Population.*Density`),
			regexp.MustCompile(`\d+,[\d,]+`),
			regexp.MustCompile(`,`)}}
}

func (self *InfoPop) GetInfo(src string) int {
	temp := self.regexes[0].FindIndex([]byte(src))
	if temp == nil {
		return -1
	}
	pop_raw := self.regexes[1].FindString(src[temp[0]:temp[1]])
	pop := self.regexes[2].ReplaceAll([]byte(pop_raw), []byte(""))
	ret, err := strconv.Atoi(string(pop))
	if err != nil {
		return 0
	}
	return ret
}

type InfoCur struct {
	regexes []*regexp.Regexp
}

func NewInfoCur() *InfoCur {
	return &InfoCur {
		regexes: []*regexp.Regexp {
			regexp.MustCompile(`Currency.*Time.{0,10}zone`),
			regexp.MustCompile(`title="[^(TLD)][\w \(\)]+`)}}
}

func (self *InfoCur) GetInfo(src string) string {
	temp := self.regexes[0].FindIndex([]byte(src))
	if temp == nil {
		return ""
	}
	curr := self.regexes[1].FindString(src[temp[0]:temp[1]])
	return curr[7:]
}

type InfoDriv struct {
	regexes []*regexp.Regexp
}

func NewInfoDriv() *InfoDriv {
	return &InfoDriv {
		regexes: []*regexp.Regexp {
			regexp.MustCompile(`Driving.{0,10}side.*Calling.{0,10}code`),
			regexp.MustCompile(`(left)|(right)`)}}
}

func (self *InfoDriv) GetInfo(src string) string {
	temp := self.regexes[0].FindIndex([]byte(src))
	if temp == nil {
		return ""
	}
	side := self.regexes[1].FindString(src[temp[0]:temp[1]])
	return side
}

type InfoCall struct {
	regexes []*regexp.Regexp
}

func NewInfoCall() *InfoCall {
	return &InfoCall {
		regexes: []*regexp.Regexp {
			regexp.MustCompile(`Calling.{0,10}code.*3166`),
			regexp.MustCompile(`[\+]\d+`)}}
}

func (self *InfoCall) GetInfo(src string) int {
	temp := self.regexes[0].FindIndex([]byte(src))
	if temp == nil {
		return -1
	}
	call := self.regexes[1].FindString(src[temp[0]:temp[1]])
	if call == "" {
		return -1
	}
	ret, err := strconv.Atoi(call[1:])
	if err != nil {
		return -1
	}
	return ret
}

type InfoTLD struct {
	regexes []*regexp.Regexp
}

func NewInfoTLD() *InfoTLD {
	return &InfoTLD {
		regexes: []*regexp.Regexp {
			regexp.MustCompile(`Internet.{0,10}TLD.*title="\.[a-z]{2,3}`),
			regexp.MustCompile(`\.\w{2,3}`)}}
}

func (self *InfoTLD) GetInfo(src string) string {
	temp := self.regexes[0].FindIndex([]byte(src))
	if temp == nil {
		return ""
	}
	tld := self.regexes[1].FindString(src[temp[0]:temp[1]])
	return tld
}
