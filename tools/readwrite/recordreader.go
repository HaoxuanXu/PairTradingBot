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

	for errReadBytes != nil || errUnmarshal != nil {
		log.Printf("Error in unmarshalling for int slice: %s\n", errUnmarshal.Error())
		log.Printf("Error in errReadBytes for int slice: %s\n", errReadBytes.Error())
		log.Println("Re-reading ...")
		recordBytes, errReadBytes = ioutil.ReadFile(path)
		errUnmarshal = json.Unmarshal(recordBytes, &recordContainer)
	}

	return recordContainer
}

func ReadRecordFloat(path string) []float64 {

	var recordContainer []float64
	recordBytes, errReadBytes := ioutil.ReadFile(path)
	errUnmarshal := json.Unmarshal(recordBytes, &recordContainer)

	for errReadBytes != nil || errUnmarshal != nil {
		log.Printf("Error in unmarshalling for float slice: %s\n", errUnmarshal.Error())
		log.Printf("Error in errReadBytes for float slice: %s\n", errReadBytes.Error())
		log.Println("Re-reading ...")
		recordBytes, errReadBytes = ioutil.ReadFile(path)
		errUnmarshal = json.Unmarshal(recordBytes, &recordContainer)
	}

	return recordContainer
}
