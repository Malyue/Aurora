package access_server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DelExitClientResp struct {
}

func DelExitClient(url string, id string) (DelExitClientResp, error) {
	path := url + "/api/v1/exitClient" + fmt.Sprintf("?id=%s", id)
	resp, err := http.Get(path)
	if err != nil {
		return DelExitClientResp{}, err
	}
	defer resp.Body.Close()

	var body DelExitClientResp
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return DelExitClientResp{}, err
	}

	return body, nil
}
