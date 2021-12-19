package readwrite

import (
	"encoding/json"
	"io/ioutil"
)

func WriteIntSlice(sli []int, path string) {
	recordBytes, _ := json.Marshal(sli)
	ioutil.WriteFile(path, recordBytes, 0644)
}

func WriteFloatSlice(sli []float64, path string) {
	recordBytes, _ := json.Marshal(sli)
	ioutil.WriteFile(path, recordBytes, 0644)
}
