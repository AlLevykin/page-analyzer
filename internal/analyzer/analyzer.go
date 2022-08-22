package analyzer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func Process(ctx context.Context, args interface{}) (interface{}, error) {

	url, ok := args.(string)
	if !ok {
		return "", errors.New("wrong argument type")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if len(body) == 0 {
		return "", errors.New("empty body")
	}

	return fmt.Sprintf("%s: %d", url, bytes.Count(body, []byte{'\n'})+1), nil

}
