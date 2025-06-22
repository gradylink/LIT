package main

import "github.com/dave/jennifer/jen"

func Motion(block Block) []jen.Code {
	switch block.Opcode {
	case "motion_movesteps":
		return []jen.Code{
			jen.Id("t").Dot("X").Op("+=").Lit(100).Op("*").Qual("math", "Cos").Call(jen.Parens(jen.Lit(90).Op("-").Id("t").Dot("Direction")).Op("*").Lit(180).Op("/").Qual("math", "Pi")),
			jen.Id("t").Dot("Y").Op("+=").Lit(100).Op("*").Qual("math", "Sin").Call(jen.Parens(jen.Lit(90).Op("-").Id("t").Dot("Direction")).Op("*").Lit(180).Op("/").Qual("math", "Pi")),
			jen.Return(jen.Lit(true)),
		}
	}

	return []jen.Code{}
}
