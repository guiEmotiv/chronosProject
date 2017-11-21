package parser

import (
	"../core"
	"encoding/json"
	"io/ioutil"
)

func GetPlacesFromJson(path string) []*core.Place {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}
	var places []*core.Place

	err = json.Unmarshal(data, &places)
	if err != nil {
		panic(err.Error())
	}
	return places

}
