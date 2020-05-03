package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

const biliAuAPIAddr string = "https://www.bilibili.com/audio/music-service-c/web/song/info/h5?sid="

func getAuInfo(au string) (Auinfo Auinfo) {
	reg, err := regexp.Compile("([0-9]+)")
	if err != nil {
		log.Fatalln(err)
	}
	Auinfo.AuNumber = strings.Join(reg.FindAllString(au, 1), "")
	requestAddr := biliAuAPIAddr + Auinfo.AuNumber
	body := getWbeContent(requestAddr)
	fmt.Println(string(body[:]))
	return Auinfo
}

//Auinfo contains some basic info of a Au number.
type Auinfo struct {
	AuNumber   string
	AuJumpURL  string
	AuCoberURL string
	AuURL      string
	AuTitle    string
	AuDesp     string
}
