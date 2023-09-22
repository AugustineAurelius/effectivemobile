package enrichment

import (
	"encoding/json"
	"io"
	"net/http"
)

func AddNation(name string) string {
	resp, err := http.Get("https://api.nationalize.io/?name=" + name)
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

	return data["country"].([]interface{})[0].(map[string]interface{})["country_id"].(string)
}
