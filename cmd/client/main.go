package main

import "github.com/flohansen/coffee-table/internal/ui/view"

func main() {
	app := view.NewCliView()
	if err := app.Run(); err != nil {
		panic(err)
	}
}
