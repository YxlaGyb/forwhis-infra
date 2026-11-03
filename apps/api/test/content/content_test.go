package content_test

import (
	"testing"

	"forwhis-infra/apps/api/internal/models"
	"forwhis-infra/apps/api/pkg/utils"
)

func TestTextLimit(t *testing.T) {
	cases := []struct {
		contentType string
		want        int
	}{
		{models.TypePost, 1000},
		{models.TypeArticle, 30000},
		{models.TypePage, 30000},
	}
	for _, tc := range cases {
		if got := models.TextLimit(tc.contentType); got != tc.want {
			t.Errorf("TextLimit(%q) = %d, want %d", tc.contentType, got, tc.want)
		}
	}
}

func TestCountRunes(t *testing.T) {
	if got := utils.CountRunes(" 你好 世界\n"); got != 4 {
		t.Errorf("CountRunes = %d, want 4", got)
	}
}
