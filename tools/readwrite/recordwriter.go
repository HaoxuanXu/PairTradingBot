package readwrite

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func WriteIntSlice(sli []int, path string) {
	recordBytes, errMarshal := json.Marshal(sli)
	errWrite := ioutil.WriteFile(path, recordBytes, 0644)
	if errWrite != nil || errMarshal != nil {
		log.Printf("error in serializing: %s\n", errMarshal)
		log.Printf("error in writing to disk: %s\n", errWrite)
		log.Println("Rewriting ...")
		WriteIntSlice(sli, path)
	}
}

func WriteFloatSlice(sli []float64, path string) {
	recordBytes, errMarshal := json.Marshal(sli)
	errWrite := ioutil.WriteFile(path, recordBytes, 0644)
	if errMarshal != nil || errWrite != nil {
		log.Printf("error in serializing: %s\n", errMarshal)
		log.Printf("error in writing to disk: %s\n", errWrite)
		log.Println("Rewriting ...")
		WriteFloatSlice(sli, path)
	}
}
