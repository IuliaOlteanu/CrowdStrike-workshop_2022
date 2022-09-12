package main

import "lab06/service"

func main() {
	s := service.NewService()
	s.StartWebService()
}