package main

import (
	"fmt"

	"github.com/Terminator15474/pantry-go"
)

type Test struct {
	Name string
}

func main() {
	// your API key
	apiKey := "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
	pantry_api := pantry.CreatePantry(apiKey)

	fmt.Println(pantry_api.CreateOrReplaceBasket("name", Test{Name: "test"}))
}
