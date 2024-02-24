package mulung

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
)

func pretest(t *testing.T) []byte {
	file, err := os.Open("benin.html")
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	file_stat, err := file.Stat()
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	file_len := file_stat.Size()
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	buff := make([]byte, file_len)	
	read, err := file.Read(buff)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	if int64(read) != file_len {
		assert.FailNow(t, "Not all read")
	}
	repl := regexp.MustCompile("[(\n\r)|(\n)|(\r)]")
	return repl.ReplaceAll(buff, []byte(""))
}

func TestInfoName(t *testing.T) {
	buff := pretest(t)
	info_name := NewInfoName()
	country_name := info_name.GetInfo(string(buff))
	assert.Equal(t, "Republic of Benin", country_name)
}

func TestInfoCapital(t *testing.T) {
	buff := pretest(t)
	info_capital := NewInfoCapital()
	country_capital := info_capital.GetInfo(string(buff))
	assert.Equal(t, "Porto-Novo", country_capital)
}

func TestInfoLang(t *testing.T) {
	buff := pretest(t)
	info_lang := NewInfoLang()
	country_lang := info_lang.GetInfo(string(buff))
	assert.Equal(t, "French Language", country_lang)
}

func TestInfoArea(t *testing.T) {
	buff := pretest(t)
	info_area := NewInfoArea()
	country_area := info_area.GetInfo(string(buff))
	assert.Equal(t, 114763, country_area)
}

func TestInfoPop(t *testing.T) {
	buff := pretest(t)
	info_pop := NewInfoPop()
	country_pop := info_pop.GetInfo(string(buff))
	assert.Equal(t, 11733059, country_pop)
}

func TestInfoCur(t *testing.T) {
	buff := pretest(t)
	info_cur := NewInfoCur()
	country_cur := info_cur.GetInfo(string(buff))
	assert.Equal(t, "West African CFA franc not yet started", country_cur)
}

func TestInfoDriv(t *testing.T) {
	buff := pretest(t)
	info_driv := NewInfoDriv()
	country_driv := info_driv.GetInfo(string(buff))
	assert.Equal(t, "right", country_driv)
}

func TestInfoCall(t *testing.T) {
	buff := pretest(t)
	info_call := NewInfoCall()
	country_call := info_call.GetInfo(string(buff))
	assert.Equal(t, 229, country_call)
}

func TestInfoTLD(t *testing.T) {
	buff := pretest(t)
	info_tld := NewInfoTLD()
	country_tld := info_tld.GetInfo(string(buff))
	assert.Equal(t, ".bj", country_tld)
}

func TestInfoBuilder(t *testing.T) {
	left := &CountryInfo {
		Name: "Republic of Benin",
		Capital: "Porto-Novo",
		OfficialLanguages: []string{"France Language"},
		Area: 114763,
		Population: 11733059,
		Currency: "West African CFA franc not yet started",
		DrivingSide: "right",
		CallingCode: 229,
		TLD: ".bj"}
	country_builder := CountryInfoBuilderNew()
	right, _ := country_builder.Src(string(pretest(t))).Build()
	assert.Equal(t, *left, *right)
}
