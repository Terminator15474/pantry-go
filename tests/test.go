package main

func main() {
	var apiKey = ""
	var pantry = createPantry(apiKey)
	println(pantry.GetDetails())
}
