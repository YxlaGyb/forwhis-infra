package feed_test

import (
	"math"
	"testing"
	"time"

	"forwhis-infra/apps/api/internal/feed"
	"forwhis-infra/apps/api/internal/models"
)

func TestCursorRoundTrip(t *testing.T) {
	original := feed.CursorState{Offset: 10, Snapshot: 1762300000, Mode: feed.ModeRecommend}
	decoded, err := feed.DecodeCursor(feed.EncodeCursor(original))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if decoded.Offset != original.Offset || decoded.Snapshot != original.Snapshot || decoded.Mode != original.Mode {
		t.Fatalf("cursor round trip mismatch: %+v", decoded)
	}
	if _, err := feed.DecodeCursor("not-a-cursor"); err == nil {
		t.Fatal("expected invalid cursor error")
	}
}

func TestFreshness(t *testing.T) {
	now := time.Unix(1762300000, 0)
	if got := feed.Freshness(feed.Candidate{Kind: feed.KindContent, Created: now.Unix()}, now); got != 1 {
		t.Fatalf("fresh content = %v, want 1", got)
	}

	weekOld := feed.Candidate{Kind: feed.KindContent, Created: now.Add(-7 * 24 * time.Hour).Unix()}
	if got := feed.Freshness(weekOld, now); math.Abs(got-0.5) > 1e-9 {
		t.Fatalf("7-day-old content = %v, want 0.5", got)
	}

	sameAge := now.Add(-3 * 24 * time.Hour).Unix()
	content := feed.Freshness(feed.Candidate{Kind: feed.KindContent, Created: sameAge}, now)
	wish := feed.Freshness(feed.Candidate{Kind: feed.KindWish, Created: sameAge}, now)
	if wish >= content {
		t.Fatalf("wish freshness %v should be below content %v at same age", wish, content)
	}
}

func TestDiversifyCapsWithFourKinds(t *testing.T) {
	kinds := []feed.Kind{feed.KindContent, feed.KindWish, "share", "question"}
	var items []feed.Scored
	score := 100.0
	for _, kind := range kinds {
		for i := 0; i < 5; i++ {
			items = append(items, feed.Scored{
				Candidate: feed.Candidate{Kind: kind, ID: len(items) + 1, AuthorID: len(items) + 1},
				Score:     score,
			})
			score--
		}
	}

	out := feed.Diversify(items)
	if len(out) != len(items) {
		t.Fatalf("Diversify dropped items: %d != %d", len(out), len(items))
	}

	authorCount := make(map[int]int)
	kindCount := make(map[feed.Kind]int)
	for pos, item := range out[:feed.DiversityWindow] {
		authorCount[item.AuthorID]++
		kindCount[item.Kind]++
		if authorCount[item.AuthorID] > feed.MaxPerAuthorInWindow {
			t.Fatalf("author %d exceeds cap at position %d", item.AuthorID, pos)
		}
		if kindCount[item.Kind] > feed.MaxPerKindInWindow {
			t.Fatalf("kind %s exceeds cap at position %d", item.Kind, pos)
		}
	}
}

func TestDiversifyRelaxesWhenExhausted(t *testing.T) {
	var items []feed.Scored
	for i := 0; i < 12; i++ {
		items = append(items, feed.Scored{
			Candidate: feed.Candidate{Kind: feed.KindContent, ID: i + 1, AuthorID: 1},
			Score:     float64(100 - i),
		})
	}

	out := feed.Diversify(items)
	if len(out) != 12 {
		t.Fatalf("Diversify dropped items: %d", len(out))
	}
	for i, item := range out {
		if item.ID != i+1 {
			t.Fatalf("relaxed order changed at %d: got id %d", i, item.ID)
		}
	}
}

func TestIsVisible(t *testing.T) {
	cases := []struct {
		name string
		c    feed.Candidate
		want bool
	}{
		{"published content", feed.Candidate{Kind: feed.KindContent, Status: models.StatusPublish, Type: models.TypePost}, true},
		{"draft content", feed.Candidate{Kind: feed.KindContent, Status: models.StatusDraft, Type: models.TypePost}, false},
		{"locked content", feed.Candidate{Kind: feed.KindContent, Status: models.StatusPublish, Type: models.TypePost, HasPassword: true}, false},
		{"page", feed.Candidate{Kind: feed.KindContent, Status: models.StatusPublish, Type: models.TypePage}, false},
		{"approved wish", feed.Candidate{Kind: feed.KindWish, Status: models.WishStatusApproved}, true},
		{"hidden wish", feed.Candidate{Kind: feed.KindWish, Status: models.WishStatusHidden}, false},
	}
	for _, tc := range cases {
		if got := feed.IsVisible(tc.c); got != tc.want {
			t.Errorf("%s: IsVisible = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestTruncateRunes(t *testing.T) {
	if got := feed.TruncateRunes("  你好世界  ", 2); got != "你好…" {
		t.Errorf("TruncateRunes = %q, want %q", got, "你好…")
	}
	if got := feed.TruncateRunes("短文本", 10); got != "短文本" {
		t.Errorf("TruncateRunes = %q, want %q", got, "短文本")
	}
}
