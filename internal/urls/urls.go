package urls

import (
	"errors"
	"net/url"
	"strings"
)

func Parse(s string, sep string) ([]string, error) {
	ss := strings.Split(s, sep)
	us := make([]string, 0, len(ss))
	for _, rawURL := range ss {
		u, err := url.Parse(strings.TrimSpace(rawURL))
		if err != nil {
			return nil, err
		}
		if u.Host == "" && u.Path == "" {
			return nil, errors.New("invalid URI for request")
		}
		if u.Scheme == "" {
			u.Scheme = "https"
		}
		u, err = url.Parse(u.String())
		if err != nil {
			return nil, err
		}
		us = append(us, u.String())
	}
	return us, nil
}
