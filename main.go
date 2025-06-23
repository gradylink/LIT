package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/spf13/pflag"
)

func main() {
	input := pflag.StringP("input", "i", "", "The Scratch project to transpile.")
	output := pflag.StringP("output", "o", "output.go", "The file to put the transpiled project.")
	projectId := pflag.Int("project-id", 0, "The URL of a Scratch project to transpile.")

	pflag.Parse()

	if !pflag.Lookup("input").Changed && !pflag.Lookup("project-id").Changed {
		fmt.Fprintln(os.Stderr, "Error: either --input (-i) or --project-id must be present.")
		pflag.Usage()
		os.Exit(1)
	}
	if pflag.Lookup("input").Changed {
		if _, err := os.Stat(*input); errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "Error: %s does not exist.\n", *input)
			pflag.Usage()
			os.Exit(1)
		} else if err != nil {
			panic(err)
		}
	} else {
		resp, err := http.Get(fmt.Sprintf("https://api.scratch.mit.edu/projects/%d", *projectId))
		if err != nil {
			panic(err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		var projectData map[string]any
		err = json.Unmarshal(body, &projectData)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()

		f, err := os.CreateTemp("", "*.sb3")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		resp, err = http.Get(fmt.Sprintf("https://projects.scratch.mit.edu/%d?token=%s", *projectId, projectData["project_token"]))
		if err != nil {
			panic(err)
		}
		if _, err = io.Copy(f, resp.Body); err != nil {
			panic(err)
		}
		resp.Body.Close()
		*input = f.Name()
	}

	// Extract SB3
	reader, err := zip.OpenReader(*input)
	if err != nil {
		panic(err)
	}
	extractPath, err := os.MkdirTemp("", "sb3-")
	if err != nil {
		panic(err)
	}
	for _, f := range reader.File {
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filepath.Join(extractPath, f.Name), os.ModePerm); err != nil {
				panic(err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(filepath.Join(extractPath, f.Name)), os.ModePerm); err != nil {
			panic(err)
		}
		extractedFile, err := os.OpenFile(filepath.Join(extractPath, f.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}
		defer extractedFile.Close()
		zippedFile, err := f.Open()
		if err != nil {
			panic(err)
		}
		defer zippedFile.Close()
		if _, err := io.Copy(extractedFile, zippedFile); err != nil {
			panic(err)
		}
	}

	project, err := Parse(filepath.Join(extractPath, "project.json"))
	if err != nil {
		panic(err)
	}

	f := jen.NewFile("main")
	f.ImportName("github.com/hajimehoshi/ebiten/v2", "ebiten")

	f.Type().Id("Block").Struct(
		jen.Id("Opcode").String(),
		jen.Id("Blocks").Index().Id("Block"),
		jen.Id("Callback").Func().Params(jen.Id("t").Op("*").Id("Target")).Params(jen.Bool()),
	).Line()
	f.Type().Id("Stack").Struct(
		jen.Id("Opcode").String(),
		jen.Id("Blocks").Index().Id("Block"),
		jen.Id("CurrentBlock").Uint64(),
		jen.Id("Running").Bool(),
	).Line()
	f.Type().Id("Target").Struct(
		jen.Id("Name").String(),
		jen.Id("IsStage").Bool(),
		jen.Id("CurrentCostume").Uint64(),
		jen.Id("Costumes").Index().Qual("github.com/hajimehoshi/ebiten/v2", "Image"),
		jen.Id("Layer").Int64(),
		jen.Id("Volume").Uint8(),
		jen.Id("Visible").Bool(),
		jen.Id("X").Float64(),
		jen.Id("Y").Float64(),
		jen.Id("Size").Uint8(),
		jen.Id("Direction").Int16(),
		jen.Id("RotationStyle").String(),
		jen.Id("Stacks").Index().Id("Stack"),
	).Line()
	f.Type().Id("Game").Struct(
		jen.Id("Targets").Index().Id("Target"),
	).Line()

	f.Func().Params(jen.Id("g").Op("*").Id("Game")).Id("Update").Params().Params(jen.Error()).Block(jen.Return(jen.Nil())).Line()                                                                                                               // Update Function
	f.Func().Params(jen.Id("g").Op("*").Id("Game")).Id("Draw").Params(jen.Id("screen").Op("*").Qual("github.com/hajimehoshi/ebiten/v2", "Image")).Block().Line()                                                                                // Draw Function
	f.Func().Params(jen.Id("g").Op("*").Id("Game")).Id("Layout").Params(jen.Id("outsideWidth"), jen.Id("outsideHeight").Int()).Params(jen.Id("screenWidth"), jen.Id("screenHeight").Int()).Block(jen.Return(jen.Lit(480), jen.Lit(360))).Line() // Layout Function

	var targets []jen.Code

	for _, target := range project.Targets {
		var stacks []jen.Code
		for _, block := range target.Blocks {
			if block.Opcode == "event_whenflagclicked" || block.Opcode == "event_whenkeypressed" || block.Opcode == "event_whenthisspriteclicked" || block.Opcode == "event_whenstageclicked" || block.Opcode == "event_whenbackdropswitchesto" || block.Opcode == "event_whengreaterthan" || block.Opcode == "event_whenbroadcastreceived" || block.Opcode == "control_start_as_clone" || block.Opcode == "procedures_definition" {
				var blocks []jen.Code
				next := block.Next
				for {
					var block []jen.Code
					switch strings.Split(target.Blocks[*next].Opcode, "_")[0] {
					case "looks":
						block = Looks(target.Blocks[*next])
					case "motion":
						block = Motion(target.Blocks[*next])
					default:
						//panic("Unsupported Opcode: " + target.Blocks[*next].Opcode)
					}
					blocks = append(blocks, jen.Values(jen.Dict{
						jen.Id("Callback"): jen.Func().Params(jen.Id("t").Op("*").Id("Target")).Params(jen.Bool()).Block(block...),
						jen.Id("Blocks"):   jen.Index().Id("Block").Values(),
						jen.Id("Opcode"):   jen.Lit(target.Blocks[*next].Opcode),
					}))
					if target.Blocks[*next].Next == nil {
						break
					}
					next = target.Blocks[*next].Next
				}

				stacks = append(stacks, jen.Values(jen.Dict{
					jen.Id("Opcode"):       jen.Lit(block.Opcode),
					jen.Id("Running"):      jen.Lit(false),
					jen.Id("CurrentBlock"): jen.Lit(0),
					jen.Id("Blocks"):       jen.Index().Id("Block").Values(blocks...),
				}))
			}
		}

		targets = append(targets, jen.Values(jen.Dict{
			jen.Id("Name"):           jen.Lit(target.Name),
			jen.Id("IsStage"):        jen.Lit(target.IsStage),
			jen.Id("CurrentCostume"): jen.Lit(target.CurrentCostume),
			jen.Id("Costumes"):       jen.Index().Qual("github.com/hajimehoshi/ebiten/v2", "Image").Values(),
			jen.Id("Layer"):          jen.Lit(target.LayerOrder),
			jen.Id("Volume"):         jen.Lit(target.Volume),
			jen.Id("Visible"):        jen.Lit(target.Visible),
			jen.Id("X"):              jen.Lit(target.X),
			jen.Id("Y"):              jen.Lit(target.Y),
			jen.Id("Size"):           jen.Lit(target.Size),
			jen.Id("Direction"):      jen.Lit(target.Direction),
			jen.Id("RotationStyle"):  jen.Lit(target.RotationStyle),
			jen.Id("Stacks"):         jen.Index().Id("Stack").Values(stacks...),
		}))
	}

	f.Func().Id("main").Params().Block(
		jen.Qual("github.com/hajimehoshi/ebiten/v2", "SetWindowSize").Call(jen.Lit(480), jen.Lit(360)),
		jen.Qual("github.com/hajimehoshi/ebiten/v2", "SetWindowTitle").Call(jen.Lit("LIT Project")),
		jen.Err().Op(":=").Qual("github.com/hajimehoshi/ebiten/v2", "RunGame").Call(jen.Op("&").Id("Game").Values(jen.Dict{
			jen.Id("Targets"): jen.Index().Id("Target").Values(targets...),
		})),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Qual("log", "Fatal").Call(jen.Err()),
		),
	)

	if err = f.Save(*output); err != nil {
		panic(err)
	}
}
