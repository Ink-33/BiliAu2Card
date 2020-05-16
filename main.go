package main

import (
	_ "crypto/hmac"
	_ "crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

var cqsecret string = gjson.Get(readConfig(), "HttpAPIPosSecret").String()

//GetAu : Get audio number by regexp.
func GetAu(msg string) (au string) {
	if strings.Contains(msg, "CQ:rich") {
		return ""
	}

	reg, err := regexp.Compile("(?i)au[0-9]+")
	if err != nil {
		log.Fatalln(err)
	}

	str := strings.Join(reg.FindAllString(msg, 1), "")
	return str
}

// Au2Card : Handle meassage and send music card.
func Au2Card(MsgInfo *MsgInfo) {
	au := GetAu(MsgInfo.Message)

	if au != "" {
		log.Println("Created request for", au, "from:", MsgInfo.SenderID)
		Auinfo := getAuInfo(au)

		if !Auinfo.AuStatus {
			msgMake := "BiliAu2Card: AU" + Auinfo.AuNumber + Auinfo.AuMsg
			log.Println(msgMake)
			switch MsgInfo.MsgType {
			case "private":
				go cqSendPrivateMsg(MsgInfo.SenderID, msgMake)
				break
			case "group":
				go cqSendGroupMsg(MsgInfo.GroupID, msgMake)
				break
			}
		} else {
			cqCodeMake := "[CQ:music,type=custom,url=" + Auinfo.AuJumpURL + ",audio=" + Auinfo.AuURL + ",title=" + Auinfo.AuTitle + ",content=" + Auinfo.AuDesp + ",image=" + Auinfo.AuCoverURL + "@180w_180h]"
			switch MsgInfo.MsgType {
			case "private":
				go cqSendPrivateMsg(MsgInfo.SenderID, cqCodeMake)
				break
			case "group":
				go cqSendGroupMsg(MsgInfo.GroupID, cqCodeMake)
				break
			}
		}
	} else {
		log.Println("Ingore message:", MsgInfo.Message, "from:", MsgInfo.SenderID)
	}
}

//MsgHandler converts HTTP Post Body to MsgInfo Struct.
func MsgHandler(raw []byte) (Msg *MsgInfo) {
	var mi = new(MsgInfo)

	mi.MsgType = gjson.GetBytes(raw, "message_type").String()
	mi.GroupID = gjson.GetBytes(raw, "group_id").String()
	mi.Message = gjson.GetBytes(raw, "message").String()
	mi.SenderID = gjson.GetBytes(raw, "user_id").String()

	return mi
}

//HTTPhandler : Handle request type before handling message.
func HTTPhandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method != "POST" {
		w.WriteHeader(400)
		fmt.Fprint(w, "Bad request.")
	} else {
		//signature := r.Header.Get("X-Signature")

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}
		//TODO: Add HMAC-SHA1 signature verification.
		/*
			mac1 := hmac.New(sha1.New, []byte(cqsecret))
			fmt.Println(string(body[:]))
			mac1.Write(body)
			fmt.Printf("%x\n",mac1.Sum(nil))
			fmt.Println(signature)
			fmt.Println(hmac.Equal(mac1.Sum(nil), []byte(signature)))
		*/
		Au2Card(MsgHandler(body))
	}
}

func main() {
	config := readConfig()
	log.SetPrefix("BiliAu2Card: ")
	path := gjson.Get(config, "BiliAu2Card.0.ListeningPath").String()
	port := gjson.Get(config, "BiliAu2Card.0.ListeningPort").String()
	log.Println("Powered by Ink33")
	log.Println("Start listening", path, port)

	http.HandleFunc(path, HTTPhandler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
