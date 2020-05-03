package main

import (
	"github.com/tidwall/gjson"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

//Read config file.
func readConfig() string {
	file, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	result := string(file)
	return result
}

//Get audio number by regexp.
func getAu(msg string) (au string) {
	reg, err := regexp.Compile("(?i)au[0-9]+")
	if err != nil {
		log.Fatalln(err)
	}
	str := strings.Join(reg.FindAllString(msg, 1), "")
	return str
}

//Input main function
func inputmain() {
	var input string
	fmt.Print("Input:")

	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Println(err)
	}

	au := getAu(input)

	if au != "" {
		log.Println("BiliAu2Card: Created request for", au)
		Auinfo := getAuInfo(au)

		if !Auinfo.AuStatus {
			log.Println("BiliAu2Card: AU", Auinfo.AuNumber, Auinfo.AuMsg)
		} else {
			cqCodeMake := "[CQ:music,type=custom,url=" + Auinfo.AuJumpURL + ",audio=" + Auinfo.AuURL + ",title=" + Auinfo.AuTitle + ",content=" + Auinfo.AuDesp + ",image=" + Auinfo.AuCoverURL + "@180w_180h]"
			cqSendPrivateMsg("", cqCodeMake)
		}

	} else {
		log.Println("BiliAu2Card: Ingore message:", input)
	}
}

//Handle request type before handling message.
func handleMsg(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method !="POST"{
		w.WriteHeader(400)
		fmt.Fprint(w,"Bad request.")
	}
}

func main() {
	config := readConfig()
	path :=gjson.Get(config,"Bili2Card.0.ListeningPath").String()
	port :=gjson.Get(config,"Bili2Card.0.ListeningPort").String()
	log.Println("BiliAu2Card: Start listening",path,port)

	http.HandleFunc(path, handleMsg)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
