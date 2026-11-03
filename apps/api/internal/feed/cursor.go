package feed

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

const cursorVersion = 1

var ErrInvalidCursor = errors.New("invalid feed cursor")

// CursorState is the opaque payload returned to clients. The snapshot keeps
// ordering stable across pages: recall is limited to created <= Snapshot and
// the ranker treats Snapshot as "now".
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
