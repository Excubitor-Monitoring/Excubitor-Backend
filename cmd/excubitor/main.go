package main

import "github.com/Excubitor-Monitoring/Excubitor-Backend/cmd/excubitor/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
