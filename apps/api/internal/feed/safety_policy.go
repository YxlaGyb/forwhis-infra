package feed

import "forwhis-infra/apps/api/internal/models"

// IsVisible is the final safety guard after SQL recall. The repository query
// already filters these states; this keeps a single rule for all sources.
func IsVisible(c Candidate) bool {
	switch c.Kind {
	case KindContent:
		return c.Status == models.StatusPublish && c.Type == models.TypePost && !c.HasPassword
	case KindWish:
		return c.Status == models.WishStatusApproved
	}
	return false
}
