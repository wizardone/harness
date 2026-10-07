package main
import (
	"fmt"
	"harness/providers"
)
func printDefaults(client providers.Client){
	fmt.Println(client.Describe())
}