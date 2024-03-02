package pantry

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
)

type Pantry struct {
	GetDetails    func() PantryInfo
	UpdateDetails func(info UpdatedInfo) PantryInfo

	CreateOrReplaceBasket func(name string, data any) bool
	UpdateBasketContent   func(name string, data any) any
	GetBasketContent      func(name string, format any) any
	DeleteBasket          func(name string) bool
}

type BasketInfo struct {
	Name string `json:"name"`
	Ttl  string `json:"ttl"`
}

type PantryInfo struct {
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	Errors        []string     `json:"errors"`
	Notifications bool         `json:"notifications"`
	PercentFull   int8         `json:"percentFull"`
	Baskets       []BasketInfo `json:"baskets"`
}

type UpdatedInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreatePantry(apiKey string) Pantry {
	var url = "https://getpantry.cloud/apiv1/pantry/" + apiKey

	var client = &http.Client{}

	var pantry = Pantry{
		GetDetails: func() PantryInfo {
			resp, err := http.Get(url)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			var info = PantryInfo{}

			json.Unmarshal(body, &info)

			return info
		},

		UpdateDetails: func(info UpdatedInfo) PantryInfo {
			reqBody, err := json.Marshal(info)
			if err != nil {
				panic(err)
			}

			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
			if err != nil {
				panic(err)
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			var newInfo = PantryInfo{}

			json.Unmarshal(body, &newInfo)

			return newInfo
		},

		CreateOrReplaceBasket: func(name string, data any) bool {
			if reflect.TypeOf(data).Kind() != reflect.Struct {
				panic("data must be a struct but got " + reflect.TypeOf(data).Kind().String())
			}

			reqBody, err := json.Marshal(data)
			if err != nil {
				panic(err)
			}

			resp, err := http.Post(url+"/basket/"+name, "application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				panic(err)
			}

			return resp.StatusCode == 200
		},

		UpdateBasketContent: func(name string, data any) any {
			if reflect.TypeOf(data).Kind() != reflect.Struct {
				panic("data must be a struct but got " + reflect.TypeOf(data).Kind().String())
			}

			reqBody, err := json.Marshal(data)
			if err != nil {
				panic(err)
			}

			req, err := http.NewRequest(http.MethodPut, url+"/basket/"+name, bytes.NewBuffer(reqBody))
			if err != nil {
				panic(err)
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			json.Unmarshal(body, &data)

			return data
		},

		GetBasketContent: func(name string, format any) any {
			resp, err := http.Get(url + "/basket/" + name)
			if err != nil {
				panic(err)
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			json.Unmarshal(body, &format)

			return format
		},

		DeleteBasket: func(name string) bool {
			req, err := http.NewRequest(http.MethodDelete, url+"/basket/"+name, nil)
			if err != nil {
				panic(err)
			}

			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			return resp.StatusCode == 200
		},
	}

	return pantry
}
