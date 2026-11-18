package dto

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/adocoder12/social_blog/internal/model"
)

// FeedQuery is the feed's URL query (?limit=&offset=&sort=&tags=&search=&since=&until=).
type FeedQuery struct {
	Limit  int        `json:"limit" validate:"gte=1,lte=20"`
	Offset int        `json:"offset" validate:"gte=0"`
	Sort   string     `json:"sort" validate:"oneof=asc desc"`
	Tags   []string   `json:"tags" validate:"max=5,dive,max=30"`
	Search string     `json:"search" validate:"max=100"`
	Since  *time.Time `json:"since"`
	Until  *time.Time `json:"until"`
}

// ParseFeedQuery reads the query string, applying defaults for missing params.
// It only checks that values parse; ranges are checked by the validator.
func ParseFeedQuery(r *http.Request) (FeedQuery, error) {
	q := FeedQuery{Limit: 20, Offset: 0, Sort: "desc"}
	qs := r.URL.Query()

	var err error
	if q.Limit, err = intParam(qs, "limit", q.Limit); err != nil {
		return q, err
	}
	if q.Offset, err = intParam(qs, "offset", q.Offset); err != nil {
		return q, err
	}
	if v := qs.Get("sort"); v != "" {
		q.Sort = v
	}

	// "go, ,rust" -> ["go", "rust"]: trim spaces and drop empty tags.
	if v := qs.Get("tags"); v != "" {
		for t := range strings.SplitSeq(v, ",") {
			if t = strings.TrimSpace(t); t != "" {
				q.Tags = append(q.Tags, t)
			}
		}
	}

	q.Search = strings.TrimSpace(qs.Get("search"))

	if q.Since, err = timeParam(qs, "since"); err != nil {
		return q, err
	}
	if q.Until, err = timeParam(qs, "until"); err != nil {
		return q, err
	}

	return q, nil
}

func (q FeedQuery) ToModel() model.FeedQuery {
	return model.FeedQuery{
		Limit:  q.Limit,
		Offset: q.Offset,
		Sort:   q.Sort,
		Tags:   q.Tags,
		Search: q.Search,
		Since:  q.Since,
		Until:  q.Until,
	}
}

// intParam returns def when the param is missing, an error when it isn't a number.
func intParam(qs url.Values, key string, def int) (int, error) {
	v := qs.Get(key)
	if v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid %s %q: must be a number", key, v)
	}
	return n, nil
}

// timeParam parses an RFC 3339 time (e.g. 2026-01-02T15:04:05Z). Missing = nil.
func timeParam(qs url.Values, key string) (*time.Time, error) {
	v := qs.Get(key)
	if v == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return nil, fmt.Errorf("invalid %s %q: use RFC 3339, e.g. 2026-01-02T15:04:05Z", key, v)
	}
	return &t, nil
}
