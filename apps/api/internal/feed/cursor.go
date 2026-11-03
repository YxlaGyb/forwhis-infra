package feed

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

const cursorVersion = 1

var ErrInvalidCursor = errors.New("invalid feed cursor")

// Snapshot keeps ordering stable across pages.
type CursorState struct {
	Version  int    `json:"v"`
	Offset   int    `json:"o"`
	Snapshot int64  `json:"t"`
	Mode     string `json:"m"`
}

func EncodeCursor(c CursorState) string {
	c.Version = cursorVersion
	raw, _ := json.Marshal(c)
	return base64.RawURLEncoding.EncodeToString(raw)
}

func DecodeCursor(s string) (CursorState, error) {
	raw, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return CursorState{}, ErrInvalidCursor
	}

	var c CursorState
	if err := json.Unmarshal(raw, &c); err != nil {
		return CursorState{}, ErrInvalidCursor
	}
	if c.Version != cursorVersion || c.Offset < 0 || c.Snapshot <= 0 {
		return CursorState{}, ErrInvalidCursor
	}
	if c.Mode != ModeRecommend && c.Mode != ModeLatest {
		return CursorState{}, ErrInvalidCursor
	}
	return c, nil
}
