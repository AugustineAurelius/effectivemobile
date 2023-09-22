package enrichment

import (
	"encoding/json"
	"io"
	"net/http"
)

func AddAge(name string) float64 {
	resp, err := http.Get("https://api.agify.io/?name=" + name)
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)

	var data map[string]interface{}

	err = json.Unmarshal(body, &data)
	return data["age"].(float64)
}
