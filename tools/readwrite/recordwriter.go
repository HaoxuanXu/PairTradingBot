package readwrite

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func WriteIntSlice(sli *[]int, path string) {
	recordBytes, errMarshal := json.Marshal(sli)
	errWrite := ioutil.WriteFile(path, recordBytes, 0644)
	for errWrite != nil || errMarshal != nil {
		if errMarshal != nil {
			log.Printf("error in marshalling for int slice: %s\n", errMarshal.Error())
		} else if errWrite != nil {
			log.Printf("error in writing to disk for int slice: %s\n", errWrite.Error())
			log.Printf("path: %s\n", path)
		}
		log.Println("Rewriting ...")
		recordBytes, errMarshal = json.Marshal(sli)
		errWrite = ioutil.WriteFile(path, recordBytes, 0644)
	}
}

func WriteFloatSlice(sli *[]float64, path string) {
	recordBytes, errMarshal := json.Marshal(*sli)
	errWrite := ioutil.WriteFile(path, recordBytes, 0644)
	for errMarshal != nil || errWrite != nil {
		if errMarshal != nil {
			log.Printf("error in marshalling for float slice: %s\n", errMarshal.Error())
		} else if errWrite != nil {
			log.Printf("error in writing to disk for float slice: %s\n", errWrite.Error())
			log.Printf("path: %s\n", path)
		}
		log.Println("Rewriting ...")
		recordBytes, errMarshal = json.Marshal(sli)
		errWrite = ioutil.WriteFile(path, recordBytes, 0644)
	}
}
