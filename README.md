# pantry-go
---
Helper library to interface with the [Pantry](https://github.com/imRohan/pantry) API in [Go](https://go.dev).

## Installation
You can download this repo with **go get**.
```sh
go get github.com/Terminator15474/pantry-go/
```

## Usage
```go
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

```