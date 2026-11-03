package feed

import (
	"sort"
	"strings"
	"time"
)

const (
	defaultLimit = 10
	maxLimit     = 30
	excerptRunes = 120
)

// GetFeed recalls, ranks, diversifies and pages the community feed.
func GetFeed(mode, rawCursor string, limit int) (*FeedResult, error) {
	if limit < 1 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	state := CursorState{Mode: mode, Snapshot: time.Now().Unix()}
	if rawCursor != "" {
		decoded, err := DecodeCursor(rawCursor)
		if err != nil {
			return nil, err
		}
		if decoded.Mode != mode {
			return nil, ErrInvalidCursor
		}
		state = decoded
	}

	candidates, err := Recall(state.Snapshot)
	if err != nil {
		return nil, err
	}

	profile := Profile{}
	now := time.Unix(state.Snapshot, 0)
	scored := make([]Scored, 0, len(candidates))
	for _, c := range candidates {
		if !IsVisible(c) {
			continue
		}
		scored = append(scored, Scored{Candidate: c, Score: Score(c, profile, now)})
	}

	if mode == ModeLatest {
		sort.SliceStable(scored, func(i, j int) bool {
			if scored[i].Created != scored[j].Created {
				return scored[i].Created > scored[j].Created
			}
			return rankTieBreak(scored[i], scored[j])
		})
	} else {
		sort.SliceStable(scored, func(i, j int) bool {
			if scored[i].Score != scored[j].Score {
				return scored[i].Score > scored[j].Score
			}
			if scored[i].Created != scored[j].Created {
				return scored[i].Created > scored[j].Created
			}
			return rankTieBreak(scored[i], scored[j])
		})
		scored = Diversify(scored)
	}

	result := &FeedResult{Mode: mode, Items: []Item{}}
	if state.Offset >= len(scored) {
		return result, nil
	}

	end := state.Offset + limit
	if end > len(scored) {
		end = len(scored)
	}
	for _, item := range scored[state.Offset:end] {
		result.Items = append(result.Items, toItem(item.Candidate))
	}
	if end < len(scored) {
		result.NextCursor = EncodeCursor(CursorState{
			Offset:   end,
			Snapshot: state.Snapshot,
			Mode:     mode,
		})
	}
	return result, nil
}

// rankTieBreak keeps ordering deterministic when score and created tie.
func rankTieBreak(a, b Scored) bool {
	if a.Kind != b.Kind {
		return a.Kind < b.Kind
	}
	return a.ID < b.ID
}

func toItem(c Candidate) Item {
	excerpt := c.Text
	if c.Kind == KindContent {
		excerpt = TruncateRunes(c.Text, excerptRunes)
	}

	item := Item{
		Kind:    c.Kind,
		ID:      c.ID,
		Title:   c.Title,
		Excerpt: excerpt,
		Created: c.Created,
		Author:  Author{UID: c.AuthorID, Name: c.AuthorName},
		Stats:   Stats{Views: c.Views, Comments: c.Comments, Likes: c.Likes},
	}
	if c.Kind == KindContent && c.CategoryID > 0 {
		item.Category = &Category{MID: c.CategoryID, Name: c.CategoryName}
	}
	if len(c.Tags) > 0 {
		item.Tags = c.Tags
	}
	return item
}

func TruncateRunes(s string, n int) string {
	s = strings.TrimSpace(s)
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}
