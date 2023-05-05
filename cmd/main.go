package main

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/cmd/excubitor"
)

func main() {
	if err := excubitor.Execute(); err != nil {
		panic(err)
	}
}
