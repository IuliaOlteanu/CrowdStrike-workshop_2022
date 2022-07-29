package main

import "fmt"


type Square struct {
	side int
}

type Rectangle struct {
	width int
	length int
}

type Circle struct {
	radius int
}

type Shape interface {
	getName()
	accept(ShapeVisitor) string
}

type ShapeVisitor interface {
	visitForSquare(square *Square)
	visitForCircle(circle *Circle)
	visitForrectangle(rectangle *Rectangle)
}

func (f *Square) visitForSquare(shape *Shape) {
	
}
func main() {
	fmt.Println("ex3")
}