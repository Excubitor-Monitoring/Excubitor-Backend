package main

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/internal/excubitor"
)

func main() {
	if err := excubitor.Execute(); err != nil {
		panic(err)
	}
}
