package readwrite

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func ReadRecordInt(path string) []int {

	var recordContainer []int
	recordBytes, errReadBytes := ioutil.ReadFile(path)
	errUnmarshal := json.Unmarshal(recordBytes, &recordContainer)

	if errUnmarshal != nil {
		log.Println(errUnmarshal.Error())
	} else if errReadBytes != nil {
		log.Println(errReadBytes.Error())
	}

	return recordContainer
}

func ReadRecordFloat(path string) []float64 {

	var recordContainer []float64
	recordBytes, errReadBytes := ioutil.ReadFile(path)
	errUnmarshal := json.Unmarshal(recordBytes, &recordContainer)

	if errUnmarshal != nil {
		log.Println(errUnmarshal.Error())
	} else if errReadBytes != nil {
		log.Println(errReadBytes.Error())
	}
	return recordContainer
}
