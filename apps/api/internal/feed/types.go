package feed

const (
	ModeRecommend = "recommend"
	ModeLatest    = "latest"
)

type Kind string

const (
	KindContent Kind = "content"
	KindWish    Kind = "wish"
)

// Candidate is a recalled feed entry before ranking.
type Candidate struct {
	Kind         Kind
	ID           int
	AuthorID     int
	AuthorName   string
	Title        string
	Text         string
	CategoryID   int
	CategoryName string
	Tags         []string
	Created      int64
	Views        int
	Comments     int
	Likes        int

	// Safety fields, checked again by IsVisible.
	Status      string
	Type        string
	HasPassword bool
}

// Profile holds personalization signals. V1 leaves it empty, so the ranker
// degrades to neutral values; V2 only needs to fill these maps.
type Profile struct {
	TopicWeights    map[int]float64
	FollowedAuthors map[int]struct{}
}

type Scored struct {
	Candidate
	Score float64
}

type Author struct {
	UID  int    `json:"uid"`
	Name string `json:"name"`
}

type Category struct {
	MID  int    `json:"mid"`
	Name string `json:"name"`
}

type Stats struct {
	Views    int `json:"views"`
	Comments int `json:"comments"`
	Likes    int `json:"likes"`
}

type Item struct {
	Kind     Kind      `json:"kind"`
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Excerpt  string    `json:"excerpt"`
	Created  int64     `json:"created"`
	Author   Author    `json:"author"`
	Category *Category `json:"category,omitempty"`
	Tags     []string  `json:"tags,omitempty"`
	Stats    Stats     `json:"stats"`
}

type FeedResult struct {
	Mode       string `json:"mode"`
	Items      []Item `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
}
