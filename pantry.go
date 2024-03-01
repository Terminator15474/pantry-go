package pantry

type Pantry struct {
	GetDetails func() string
}

func createPantry(apiKey string) Pantry {
	var url = "https://getpantry.cloud/apiv1/pantry/" + apiKey
	var pantry = Pantry{
		GetDetails: func() string {
			return url
		},
	}

	return pantry
}
