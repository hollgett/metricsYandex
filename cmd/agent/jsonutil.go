package main

import "encoding/json"

func Marshal(requestdata Metrics) ([]byte, error) {
	return json.Marshal(requestdata)
}
