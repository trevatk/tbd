package resolver

import (
	"fmt"
	"io"
	"net/http"
)

func resolveDID(target string) ([]byte, error) {
	resp, err := http.Get(target)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http get: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode > http.StatusAccepted {
		return nil, fmt.Errorf("unexcepted http response %d %s", resp.StatusCode, string(body))
	}

	return body, nil
}
