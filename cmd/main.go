package main

import (
	app "github.com/coolorvi/web-calculator/web"
)

func main() {
	app := app.New()
	app.RunServer()
}
