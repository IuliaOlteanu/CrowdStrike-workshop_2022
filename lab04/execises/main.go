package main

import "example.com/rest/service"

func main() {
	s := service.NewService()
	s.StartWebService()
}