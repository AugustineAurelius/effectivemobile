package enrichment

import (
	"encoding/json"
	"io"
	"net/http"
)

func AddGender(name string) string {
	resp, err := http.Get("https://api.genderize.io/?name=" + name)
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
	return data["gender"].(string)
}
