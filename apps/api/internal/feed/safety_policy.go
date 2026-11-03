package feed

import "forwhis-infra/apps/api/internal/models"

func IsVisible(c Candidate) bool {
	switch c.Kind {
	case KindContent:
		return c.Status == models.StatusPublish && c.Type == models.TypePost && !c.HasPassword
	case KindWish:
		return c.Status == models.WishStatusApproved
	}
	return false
}
