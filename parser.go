package main

import (
	"encoding/json"
	"errors"
	"os"
)

type Target struct {
	IsStage              bool                `json:"isStage"`
	Name                 string              `json:"name"`       // "Stage" if IsStage == true
	Variables            map[string]Variable `json:"variables"`  // Maps IDs to [Variable]s
	Lists                map[string]List     `json:"lists"`      // Maps IDs to [List]s
	Broadcasts           map[string]string   `json:"broadcasts"` // Maps IDs to Names. Normally only present if IsStage == true.
	Blocks               map[string]Block    `json:"blocks"`     // Maps IDs to [Block]s.
	Comments             map[string]Comment  `json:"comments"`   // Maps IDs to [Comment]s
	CurrentCostume       float64             `json:"currentCostume"`
	Costumes             []Asset             `json:"costumes"`
	Sounds               []Asset             `json:"sounds"`
	LayerOrder           float64             `json:"layerOrder"`
	Volume               float64             `json:"volume"`
	Tempo                *float64            `json:"tempo,omitempty"`                // Only present if IsStage == true.
	VideoState           *string             `json:"videoState,omitempty"`           // Either "on", "off", or "on-flipped". Only present if IsStage == true.
	VideoTransparency    *float64            `json:"videoTransparency,omitempty"`    // Only present if IsStage == true.
	TextToSpeechLanguage *string             `json:"textToSpeechLanguage,omitempty"` // Only present if IsStage == true.
	Visible              bool                `json:"visible"`                        // Only present if IsStage == false.
	X                    float64             `json:"x"`                              // Only present if IsStage == false.
	Y                    float64             `json:"y"`                              // Only present if IsStage == false.
	Size                 float64             `json:"size"`                           // Only present if IsStage == false.
	Direction            float64             `json:"direction"`                      // Only present if IsStage == false.
	Draggable            bool                `json:"draggable"`                      // Only present if IsStage == false.
	RotationStyle        string              `json:"rotationStyle"`                  // Either "all around", "left-right", or "don't rotate". Only present if IsStage == false.
}

type Variable struct {
	Name    string
	Value   any // Either float64 or string
	IsCloud bool
}

type List struct {
	Name   string
	Values []any // Array of either float64 or string
}

type Block struct {
	Opcode   string              `json:"opcode"`
	Next     *string             `json:"next,omitempty"`
	Parent   *string             `json:"parent,omitempty"`
	Inputs   map[string][]any    `json:"inputs"` // Maps IDs to Arrays representing Inputs. The first element of each array is 1 if the input is a shadow, 2 if there is no shadow, and 3 if there is a shadow but it is obscured by the input. The second is either the ID of the input or an array representing it as described in the table below. If there is an obscured shadow, the third element is its ID or an array representing it.
	Fields   map[string][]string `json:"fields"`
	Shadow   bool                `json:"shadow"`
	TopLevel bool                `json:"topLevel"`
}

type Comment struct {
	BlockID   float64 `json:"blockId"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Width     float64 `json:"width"`
	Height    float64 `json:"height"`
	Minimized bool    `json:"minimized"`
	Text      string  `json:"text"`
}

type Asset struct {
	AssetID          string   `json:"assetId"`
	Name             string   `json:"name"`
	MD5Ext           string   `json:"md5ext"`
	DataFormat       string   `json:"dataFormat"`
	BitmapResolution *float64 `json:"bitmapResolution,omitempty"` // Only if the asset is a costume.
	RotationCenterX  *float64 `json:"rotationCenterX,omitempty"`  // Only if the asset is a costume.
	RotationCenterY  *float64 `json:"rotationCenterY,omitempty"`  // Only if the asset is a costume.
	Rate             *float64 `json:"rate,omitempty"`             // Only if the asset is a sound.
	SampleCount      *float64 `json:"sampleCount,omitempty"`      // Only if the asset is a sound.
}

type Metadata struct {
	SemVer string `json:"semver"` // Should always be "3.0.0"
	VM     string `json:"vm"`     // Version of scratch-vm
	Agent  string `json:"agent"`  // User Agent of client on last save
}

type Project struct {
	Targets    []Target `json:"targets"`
	Monitors   []string `json:"monitors"`
	Extensions []string `json:"extensions"`
	Meta       Metadata `json:"meta"`
}

func (v *Variable) UnmarshalJSON(data []byte) error {
	var arr []any
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}

	if len(arr) < 2 || len(arr) > 3 {
		return errors.New("Invalid Variable: Variables must have 2 or 3 elements.")
	}

	name, ok := arr[0].(string)
	if !ok {
		return errors.New("Invalid Variable: The first element of a variable must be a string.")
	}
	v.Name = name

	switch value := arr[1].(type) {
	case float64, string:
		v.Value = value
	default:
		return errors.New("Invalid Variable: The second element of a variable must be a string or a number.")
	}

	v.IsCloud = false
	if len(arr) == 3 {
		if cloud, ok := arr[2].(bool); !ok || !cloud {
			return errors.New("Invalid Variable: If it exists, the third element of a variable must be true.")
		}
		v.IsCloud = true
	}

	return nil
}

func (l *List) UnmarshalJSON(data []byte) error {
	var arr []any
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}

	if len(arr) != 2 {
		return errors.New("Invalid List: Lists must have 2 elements.")
	}

	name, ok := arr[0].(string)
	if !ok {
		return errors.New("Invalid List: The first element of a list must be a string.")
	}
	l.Name = name

	values, ok := arr[1].([]any)
	if !ok {
		return errors.New("Invalid List: The second element of a list must be an array.")
	}
	l.Values = values

	return nil
}

func Parse(path string) (*Project, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var project Project
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, err
	}

	return &project, nil
}
