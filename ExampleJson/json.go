package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
)

type bcContext struct {
	Code            int          `json:"code"`
	Source         string        `json:"source"`
	Time           float32       `json:"time"`
	//...
}

type bcDataCoin struct {
	Bitcoin        bcDataCoinData `json:"bitcoin"`
	BitcoinSV      bcDataCoinData `json:"bitcoin-sv"`
	BitcoinCash    bcDataCoinData `json:"bitcoin-cash"`
	//...
}

type bcDataCoinData struct {
	Data           bcDataCoinInfo  `json:"data"`
	Context        bcContext       `json:"context"`
	//..
}

type bcDataCoinInfo struct {
	Blocks          int            `json:"blocks"`
	Hashrate_24h    string         `json:"hashrate_24h"`
	//..
}

type blockchairStruct struct {
	Data            bcDataCoin     `json:"data"`
	Context         bcContext      `json:"context"`
}

func main() {
	//http 请求
	response, err := http.Get("https://api.blockchair.com/stats")
	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	//http请求完成

	/*
	 * 在已知json 结构的情况下，构建相应的 struct 结构
	 */
	jsonStr := string(body)

	var blockchair blockchairStruct

	err = json.Unmarshal([]byte(jsonStr), &blockchair)

	if  err != nil {
		fmt.Println(err)
	}

	fmt.Println(blockchair.Data.BitcoinSV.Data)

	/**
	 * 第二种使用 simplejson 库
	 */
	bcJson, err := simplejson.NewJson([]byte(jsonStr)) //反序列化

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(bcJson)

	/**
	 * 在未知请求 json 结构的情况下
	 */
	var bc map[string]interface{}

	err = json.Unmarshal([]byte(jsonStr), &bc)

	if  err != nil {
		fmt.Println(err)
	}

	fmt.Println(11111)
	fmt.Println(bc["data"])
	bitcoin := bc["data"].(map[string]interface{})

	fmt.Println(bitcoin["bitcoin-sv"])
}
