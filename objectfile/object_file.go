package objectfile

import (
	"encoding/json"
	"os"

	"patito/quadruples"
)

type ObjectConstant struct {
	Literal string `json:"literal"`
	Type    string `json:"type"`
	Address int    `json:"address"`
}

type ObjectFile struct {
	Quadruples []quadruples.Quadruple `json:"quadruples"`
	Constants  []ObjectConstant       `json:"constants"`
}

func Save(path string, obj ObjectFile) error {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func Load(path string) (*ObjectFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var obj ObjectFile
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
