package main

import (
	"github.com/dave/jennifer/jen"
)

func main() {
	project, err := Parse("./test-extracted/project.json")
	if err != nil {
		panic(err)
	}

	f := jen.NewFile("main")
	for _, target := range project.Targets {
		if target.IsStage {
			continue
		}

		for _, block := range target.Blocks {
			if block.Opcode != "procedures_definition" {
				continue
			}
			if *target.Blocks[*block.Inputs["custom_block"].BlockID].Mutation.ProcCode == "Main" {
				var statements []jen.Code
				next := block.Next
				for {
					switch target.Blocks[*next].Opcode {
					case "looks_say":
						switch (*target.Blocks[*next].Inputs["MESSAGE"].Value).(type) {
						case string:
							statements = append(statements, jen.Qual("fmt", "Println").Call(jen.Lit(*target.Blocks[*next].Inputs["MESSAGE"].Value)))
						default:
							panic("Unsupported Input Type")
						}
					default:
						panic("Unsupported Opcode: " + target.Blocks[*next].Opcode)
					}
					if target.Blocks[*next].Next == nil {
						break
					}
					next = target.Blocks[*next].Next
				}
				f.Func().Id("main").Params().Block(statements...)
			}
		}
	}

	if err = f.Save("dist/main.go"); err != nil {
		panic(err)
	}
}
