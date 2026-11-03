package feed

import "strconv"

const (
	DiversityWindow      = 10
	MaxPerAuthorInWindow = 2
	MaxPerKindInWindow   = 3
)

func authorKey(c Candidate) string {
	if c.AuthorID > 0 {
		return "u:" + strconv.Itoa(c.AuthorID)
	}
	return "n:" + c.AuthorName
}

func Diversify(items []Scored) []Scored {
	out := make([]Scored, 0, len(items))
	used := make([]bool, len(items))
	authorCount := make(map[string]int)
	kindCount := make(map[Kind]int)

	for pos := 0; pos < DiversityWindow && len(out) < len(items); pos++ {
		pick := -1
		for i, item := range items {
			if used[i] {
				continue
			}
			if authorCount[authorKey(item.Candidate)] >= MaxPerAuthorInWindow {
				continue
			}
			if kindCount[item.Kind] >= MaxPerKindInWindow {
				continue
			}
			pick = i
			break
		}
		if pick == -1 {
			break
		}

		used[pick] = true
		out = append(out, items[pick])
		authorCount[authorKey(items[pick].Candidate)]++
		kindCount[items[pick].Kind]++
	}

	for i, item := range items {
		if !used[i] {
			out = append(out, item)
		}
	}
	return out
}
