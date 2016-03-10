package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	// "github.com/kr/pretty"
)

func unproxyURL(req *http.Request) {
	fmt.Printf("req.RequestURI=%v\n", req.RequestURI)
	strURL := req.RequestURI

	if inHostList(strings.Join(arg.HttpsList, ","), req.Host) {
		strURL = changeHostToHttps(req.RequestURI)
		fmt.Printf("changeHostToHttps=%v\n", strURL)
	}

	target, err := url.Parse(strURL)
	fatal(err)
	req.URL = target
}

func inHostList(hostList, hostname string) bool {
	index := strings.Index(hostList, hostname)
	return (index >= 0)
}

func changeHostToHttps(endpoint string) string {
	return strings.Replace(endpoint, "http://", "https://", -1)
}

func fatal(err error) {
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatal(err)
		os.Exit(1)
	}
}

func captureExitProgram() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		data.WriteStub()

		println()
		println("end proxy...")
		os.Exit(1)
	}()
}

func getValueByKey(key string, data string) string {
	list := strings.FieldsFunc(data, func(r rune) bool {
		return r == '<' || r == '>'
	})

	// fmt.Printf("list[%v]=% v\n", len(list), pretty.Formatter(list))
	for i, s := range list {
		// if s == key {
		// 	if len(list) > i {
		// 		if list[i+1] != "/"+key {
		// 			return list[i+1]
		// 		}
		// 	}
		// }

		if s == key && len(list) > i && list[i+1] != "/"+key {
			return list[i+1]
		}

	}
	return ""
}

func generateKey(req Inbound) string {
	conditionField := getConditionField(req.Host+req.Path, arg.IncludeList)
	conditionValue := getConditionValue(conditionField, byteToStr(req.Body))

	// fmt.Printf("condition field=%v\n", conditionField)
	// fmt.Printf("condition value=%v\n", conditionValue)
	// return req.Method + "|" + req.RequestURI + "|" + conditionValue
	return req.Method + "|" + req.Host + req.Path + "|" + conditionValue
}

func getConditionField(endpoint string, fieldList Condition) string {
	if list, found := fieldList[endpoint]; found {
		return list
	}
	return ""
}

func getConditionValue(key, data string) string {
	var result []string

	list := strings.Split(key, ",")
	for _, value := range list {

		if v := getValueByKey(value, data); v != "" {
			result = append(result, v)
		}
	}

	// fmt.Printf("result[%v]=% v\n", len(result), pretty.Formatter(result))
	return strings.Join(result, "|")
}

func byteToStr(data []byte) string {
	return fmt.Sprintf("%s", data)
}
