package pantry

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"slices"
	"time"

	"golang.org/x/time/rate"
)

type Pantry struct {
	// GetDetails returns the details of the pantry
	GetDetails func() (PantryInfo, error)

	// UpdateDetails updates the details of the pantry
	UpdateDetails func(info UpdatedInfo) (PantryInfo, error)

	// CreateOrReplaceBasket creates or replaces a basket in the pantry
	// with the given name and data and returns true if successful
	CreateOrReplaceBasket func(name string, data any) (bool, error)

	// UpdateBasketContent updates the content of the basket
	// with the given name and returns the updated data
	UpdateBasketContent func(name string, data any) (any, error)

	// GetBasketContent returns the content of the basket
	// with the given name in the given format
	GetBasketContent func(name string, format any) (any, error)

	// DeleteBasket deletes the basket with the given name
	// and returns true if successful
	// THIS WILL DELETE ALL THE DATA IN THE BASKET
	DeleteBasket func(name string) (bool, error)

	// HasBasket checks if the pantry has a basket with the given name
	HasBasket func(name string) (bool, error)
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

// CreatePantry creates a new rateLimited Pantry API Wrapper with the given API key
func CreateRateLimitedPantry(apiKey string) Pantry {
	var PANTRY_URL = "https://getpantry.cloud/apiv1/pantry/" + apiKey
	var pantry = Pantry{
		GetDetails: func() (PantryInfo, error) {
			req, _ := http.NewRequest(http.MethodGet, PANTRY_URL, nil)
			resp, err := doRateLimitedRequest(req)
			if err != nil {
				return PantryInfo{}, err
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return PantryInfo{}, err
			}

			var info = PantryInfo{}

			json.Unmarshal(body, &info)

			return info, nil
		},

		UpdateDetails: func(info UpdatedInfo) (PantryInfo, error) {
			reqBody, err := json.Marshal(info)
			if err != nil {
				return PantryInfo{}, err
			}

			req, err := http.NewRequest(http.MethodPut, PANTRY_URL, bytes.NewBuffer(reqBody))
			if err != nil {
				return PantryInfo{}, err
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := doRateLimitedRequest(req)
			if err != nil {
				return PantryInfo{}, err
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return PantryInfo{}, err
			}

			var newInfo = PantryInfo{}

			json.Unmarshal(body, &newInfo)

			return newInfo, nil
		},

		CreateOrReplaceBasket: func(name string, data any) (bool, error) {
			if reflect.TypeOf(data).Kind() != reflect.Struct {
				return false, errors.New("data must be a struct but got " + reflect.TypeOf(data).Kind().String())
			}

			reqBody, err := json.Marshal(data)
			if err != nil {
				return false, err
			}

			req, _ := http.NewRequest(http.MethodGet, PANTRY_URL+"/basket/"+name, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := doRateLimitedRequest(req)
			if err != nil {
				return false, err
			}

			return resp.StatusCode == 200, nil
		},

		UpdateBasketContent: func(name string, data any) (any, error) {
			if reflect.TypeOf(data).Kind() != reflect.Struct {
				return nil, errors.New("data must be a struct but got " + reflect.TypeOf(data).Kind().String())
			}

			reqBody, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}

			req, err := http.NewRequest(http.MethodPut, PANTRY_URL+"/basket/"+name, bytes.NewBuffer(reqBody))
			if err != nil {
				return nil, err
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := doRateLimitedRequest(req)
			if err != nil {
				return nil, err
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(body, &data)

			return data, err
		},

		GetBasketContent: func(name string, format any) (any, error) {
			req, _ := http.NewRequest(http.MethodGet, PANTRY_URL+"/basket/", nil)
			req.Header.Set("Content-Type", "application/json")
			resp, err := doRateLimitedRequest(req)
			if err != nil {
				return nil, err
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(body, &format)

			return format, err
		},

		DeleteBasket: func(name string) (bool, error) {
			req, err := http.NewRequest(http.MethodDelete, PANTRY_URL+"/basket/"+name, nil)
			if err != nil {
				return false, err
			}

			resp, err := doRateLimitedRequest(req)
			if err != nil {
				return false, err
			}

			return resp.StatusCode == 200, nil
		},

		HasBasket: func(name string) (bool, error) {
			req, _ := http.NewRequest(http.MethodGet, PANTRY_URL, nil)
			resp, err := doRateLimitedRequest(req)
			if err != nil {
				return false, err
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return false, err
			}

			var info = PantryInfo{}

			json.Unmarshal(body, &info)
			var names = []string{}
			for _, b := range info.Baskets {
				names = append(names, b.Name)
			}

			return slices.Contains(names, name), nil
		},
	}

	return pantry
}

func CreatePantry(apiKey string) Pantry {
	var PANTRY_URL = "https://getpantry.cloud/apiv1/pantry/" + apiKey

	var client = &http.Client{}

	var pantry = Pantry{
		GetDetails: func() (PantryInfo, error) {
			req, _ := http.NewRequest(http.MethodPut, PANTRY_URL, bytes.NewBuffer(nil))
			resp, err := client.Do(req)
			if err != nil {
				return PantryInfo{}, err
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return PantryInfo{}, err
			}

			var info = PantryInfo{}

			json.Unmarshal(body, &info)

			return info, nil
		},

		UpdateDetails: func(info UpdatedInfo) (PantryInfo, error) {
			reqBody, err := json.Marshal(info)
			if err != nil {
				return PantryInfo{}, err
			}

			req, err := http.NewRequest(http.MethodPut, PANTRY_URL, bytes.NewBuffer(reqBody))
			if err != nil {
				return PantryInfo{}, err
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				return PantryInfo{}, err
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return PantryInfo{}, err
			}

			var newInfo = PantryInfo{}

			json.Unmarshal(body, &newInfo)

			return newInfo, nil
		},

		CreateOrReplaceBasket: func(name string, data any) (bool, error) {
			if reflect.TypeOf(data).Kind() != reflect.Struct {
				return false, errors.New("data must be a struct but got " + reflect.TypeOf(data).Kind().String())
			}

			reqBody, err := json.Marshal(data)
			if err != nil {
				return false, err
			}

			resp, err := http.Post(PANTRY_URL+"/basket/"+name, "application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				return false, err
			}

			return resp.StatusCode == 200, nil
		},

		UpdateBasketContent: func(name string, data any) (any, error) {
			if reflect.TypeOf(data).Kind() != reflect.Struct {
				return nil, errors.New("data must be a struct but got " + reflect.TypeOf(data).Kind().String())
			}

			reqBody, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}

			req, err := http.NewRequest(http.MethodPut, PANTRY_URL+"/basket/"+name, bytes.NewBuffer(reqBody))
			if err != nil {
				return nil, err
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(body, &data)

			return data, err
		},

		GetBasketContent: func(name string, format any) (any, error) {
			resp, err := http.Get(PANTRY_URL + "/basket/" + name)
			if err != nil {
				return nil, err
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(body, &format)

			return format, err
		},

		DeleteBasket: func(name string) (bool, error) {
			req, err := http.NewRequest(http.MethodDelete, PANTRY_URL+"/basket/"+name, nil)
			if err != nil {
				return false, err
			}

			resp, err := client.Do(req)
			if err != nil {
				return false, err
			}

			return resp.StatusCode == 200, nil
		},

		HasBasket: func(name string) (bool, error) {

			resp, err := http.Get(PANTRY_URL)
			if err != nil {
				return false, err
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return false, err
			}

			var info = PantryInfo{}

			json.Unmarshal(body, &info)
			var names = []string{}
			for _, b := range info.Baskets {
				names = append(names, b.Name)
			}

			return slices.Contains(names, name), nil
		},
	}

	return pantry
}

var rl = rate.NewLimiter(rate.Every(time.Second), 2)

func doRateLimitedRequest(req *http.Request) (*http.Response, error) {
	ctx := context.Background()
	err := rl.Wait(ctx)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}
