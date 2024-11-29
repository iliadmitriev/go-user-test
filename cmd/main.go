package main

import "github.com/iliadmitriev/go-user-test/internal/app"

func main() {
	application := app.NewApplication()
	_ = application.Run()
}
