package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

const biliAuAPIAddr string = "https://www.bilibili.com/audio/music-service-c/web/song/info/h5?sid="
const biliAudioPlayURL string = "https://api.bilibili.com/audio/music-service-c/shareUrl/redirectHttp?songid="
const biliAudioJumpURL string = "https://www.bilibili.com/audio/au"

func getAuInfo(au string) (Auinfo Auinfo) {
	reg, err := regexp.Compile("[0-9]+")
	if err != nil {
		log.Fatalln(err)
	}

	Auinfo.AuNumber = strings.Join(reg.FindAllString(au, 1), "")
	Auinfo.AuURL = biliAudioPlayURL + Auinfo.AuNumber
	Auinfo.AuJumpURL = biliAudioJumpURL + Auinfo.AuNumber

	requestAddr := biliAuAPIAddr + Auinfo.AuNumber
	body := string(getWbeContent(requestAddr)[:])

	Auinfo.AuMsg = gjson.Get(body, "msg").String()
	if Auinfo.AuMsg != "success" {
		Auinfo.AuStatus = false
		return
	}
	Auinfo.AuStatus = true
	Auinfo.AuCoverURL = gjson.Get(body, "data.h5Songs.cover_url").String()
	Auinfo.AuTitle = gjson.Get(body, "data.h5Songs.title").String()
	Auinfo.AuDesp = gjson.Get(body, "data.h5Songs.author").String()
	return
}
