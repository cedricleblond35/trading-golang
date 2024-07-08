package main

// https://github.com/go-gota/gota
// https://www.freecodecamp.org/news/exploratory-data-analysis-in-go-with-gota/

import (
	"fmt"

	"github.com/go-gota/gota/dataframe"
	// "github.com/go-gota/gota/series"
)

func main() {
	type Dog struct {
		Name       string
		Color      string
		Height     int
		Vaccinated bool
	}

	dogs := []Dog{
		{"Buster", "Black", 56, false},
		{"Jake", "White", 61, false},
		{"Bingo", "Brown", 50, true},
		{"Gray", "Cream", 68, false},
	}

	dogsDf := dataframe.LoadStructs(dogs)

	fmt.Println(dogs)

	fmt.Println(dogsDf.Elem(0, 2))
	fmt.Println(dogsDf.Col("Color"))
	fmt.Println(dogsDf.Col("Color").Records())

	// nous voulons sélectionner des colonnes spécifiques, par un index ou un nom de colonne:
	fmt.Println(dogsDf.Select([]int{0, 2}))
	fmt.Println(dogsDf.Select([]string{"Name", "Color"}))

	fmt.Println(dogsDf.Dims())
}
