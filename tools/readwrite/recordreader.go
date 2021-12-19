package readwrite

import (
	"encoding/json"
	"io/ioutil"
)

func ReadRecordInt(path string) []int {

	var recordContainer []int
	recordBytes, _ := ioutil.ReadFile(path)
	json.Unmarshal(recordBytes, &recordContainer)

	return recordContainer
}

func ReadRecordFloat(path string) []float64 {

	var recordContainer []float64
	recordBytes, _ := ioutil.ReadFile(path)
	json.Unmarshal(recordBytes, &recordContainer)

	return recordContainer
}