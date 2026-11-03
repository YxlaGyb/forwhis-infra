package feed

import (
	"math"
	"time"
)

const (
	weightInterest  = 30.0
	weightHelpful   = 25.0
	weightFreshness = 20.0
	weightQuality   = 15.0
	weightAffinity  = 10.0
)

const (
	contentHalfLife = 7 * 24 * time.Hour
	wishHalfLife    = 3 * 24 * time.Hour
)

func clamp(v float64) float64 {
	return math.Max(0, math.Min(1, v))
}

func topicMatch(c Candidate, p Profile) float64 {
	if len(p.TopicWeights) == 0 {
		return 0.5
	}
	return clamp(p.TopicWeights[c.CategoryID])
}

func Freshness(c Candidate, now time.Time) float64 {
	halfLife := contentHalfLife
	if c.Kind == KindWish {
		halfLife = wishHalfLife
	}

	age := now.Sub(time.Unix(c.Created, 0))
	if age < 0 {
		age = 0
	}
	return math.Exp(-math.Ln2 * age.Hours() / halfLife.Hours())
}

func helpfulness(_ Candidate) float64 {
	return 0.5
}

func quality(_ Candidate) float64 {
	return 0.5
}

func authorAffinity(c Candidate, p Profile) float64 {
	if _, ok := p.FollowedAuthors[c.AuthorID]; ok {
		return 1
	}
	return 0
}

func Score(c Candidate, p Profile, now time.Time) float64 {
	return weightInterest*topicMatch(c, p) +
		weightHelpful*helpfulness(c) +
		weightFreshness*Freshness(c, now) +
		weightQuality*quality(c) +
		weightAffinity*authorAffinity(c, p)
}
