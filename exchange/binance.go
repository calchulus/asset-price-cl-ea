package exchange

import (
	"github.com/adshao/go-binance"
	"context"
	"strconv"
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"log"
)

type Binance struct {
	Exchange
}

type BinanceProduct struct {
	LastPrice string `json:"lastprice"`
	Volume string `json:"volume"`
}

var binancePairs []*Pair

func (exchange Binance) GetResponse(base, quote string) (*Response, *Error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET", fmt.Sprintf("%s/ticker/24hr?symbol=%s%s", exchange.GetConfig().BaseUrl, base, quote), nil)
	if err != nil {
		return nil, &Error{exchange.GetConfig().Name, "500 ERROR", "error on forming request to Binance"}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &Error{exchange.GetConfig().Name, resp.Status, err.Error()}
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &Error{exchange.GetConfig().Name, resp.Status, err.Error()}
	}
	binanceProduct := BinanceProduct{}
	err = json.Unmarshal(bodyBytes, &binanceProduct)
	if err != nil {
		return nil, &Error{exchange.GetConfig().Name, resp.Status, err.Error()}
	}
	currentPrice, err := strconv.ParseFloat(binanceProduct.LastPrice, 64)
	if err != nil {
		return nil, &Error{exchange.GetConfig().Name, resp.Status, err.Error()}
	}
	currentVolume, err := strconv.ParseFloat(binanceProduct.Volume, 64)
	if err != nil {
		return nil, &Error{exchange.GetConfig().Name, resp.Status, err.Error()}
	}
	return &Response{exchange.GetConfig().Name, currentPrice, currentVolume}, nil
}

func (exchange Binance) SetPairs() {
	clientInterface := exchange.GetConfig().Client
	client := clientInterface.(*binance.Client)

	exchangeInfo, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, product := range exchangeInfo.Symbols {
		binancePairs = append(binancePairs, &Pair{product.BaseAsset, product.QuoteAsset})
	}
}

func (exchange Binance) GetConfig() *Config {
	return &Config{
		Name: "Binance",
		BaseUrl: "https://www.binance.com/api/v1",
		Client: binance.NewClient("", ""),
		Pairs: binancePairs}
}
