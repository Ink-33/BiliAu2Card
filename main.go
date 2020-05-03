package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func main() {
	var input string
	log.Println("Start")
	for {
		fmt.Print("Input:")
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Println(err)
		}
		au := getAu(input)
		if au != ""{
			fmt.Println(au)
			_ = getAuInfo(au)
		}
	}
}

func getAu(msg string) (au string) {
	reg, err := regexp.Compile("([auAU]+)([0-9]+)")
	if err != nil {
		log.Fatalln(err)
	}
	str := strings.Join(reg.FindAllString(msg, 1),"")
	return str
}
