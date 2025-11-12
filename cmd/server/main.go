package main

import (
	"fmt"
	//"little-orm/internal/database"
	. "little-orm/internal/database/querybuilder"
	"little-orm/internal/model"
)

func main() {

	//db := database.GetDB()

	builder := NewSelectBuilder(model.User{})
	query, args := builder.
		Select("ID", "Name").
		Where(Or(
			B(OpEq, C("ID"), L(1)),
			U(OpIsNNull, C("Name")),
			B(OpGte, C("Name"), L("A")),
		)).
		Build()

	fmt.Println(query)
	fmt.Println(args...)
}
