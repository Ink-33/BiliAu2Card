package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

//Send Group message by using CoolQ HttpAPI
func cqSendGroupMsg(id, msg string) {
	conf := readConfig()
	cqAddr := gjson.Get(conf, "CoolQ.0.Api.HttpAPIAddr").String()
	cqToken := gjson.Get(conf, "CoolQ.0.Api.HttpAPIToken").String()
	getWbeContent(cqAddr + "/send_group_msg?access_token=" + cqToken + "&group_id=" + id + "&message=" + url.QueryEscape(msg))
}

//Send private message by using CoolQ HttpAPI
func cqSendPrivateMsg(id, msg string) {
	conf := readConfig()
	cqAddr := gjson.Get(conf, "CoolQ.0.Api.HttpAPIAddr").String()
	cqToken := gjson.Get(conf, "CoolQ.0.Api.HttpAPIToken").String()
	getWbeContent(cqAddr + "/send_private_msg?access_token=" + cqToken + "&user_id=" + id + "&message=" + url.QueryEscape(msg))
}

//Get web Content by using GET request.
func getWbeContent(url string) (body []byte) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4117.2 Safari/537.36")
	if err != nil {
		log.Fatalln(err)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return content
}

//Read config file.
func readConfig() string {
	file, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	result := string(file)
	return result
}
