package shared

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

func ExtractBase64ImageData(imageURL string) (mimeType string, data []byte, err error) {
	u, err := url.Parse(imageURL)
	if err != nil {
		return "", nil, fmt.Errorf("invalid URL: %w", err)
	}

	if u.Scheme != "data" {
		return "", nil, fmt.Errorf("not a data URL")
	}

	// 使用strings.SplitN而不是strings.Split避免潜在问题，如果base64编码的数据中包含逗号
	parts := strings.SplitN(u.Opaque, ",", 2)
	if len(parts) != 2 {
		return "", nil, fmt.Errorf("invalid data URL format")
	}

	mimeType = parts[0]
	if !strings.Contains(mimeType, ";base64") {
		return "", nil, fmt.Errorf("invalid data URL format: missing base64 encoding")
	}
	mimeType = strings.TrimSuffix(mimeType, ";base64")

	data, err = base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", nil, fmt.Errorf("base64 decode error: %w", err)
	}

	return mimeType, data, nil
}
