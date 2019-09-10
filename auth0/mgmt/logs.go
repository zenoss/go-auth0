package mgmt

import (
	"fmt"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/zenoss/go-auth0/auth0/http"
)

type LogsService struct {
	c *http.Client
}

type LogEntry struct {
	UserID       string                 `json:"user_id,omitempty"`
	Date         time.Time              `json:"date,omitempty"`
	Type         string                 `json:"type,omitempty"`
	ClientID     string                 `json:"client_id,omitempty"`
	ClientName   string                 `json:"client_name,omitempty"`
	IP           string                 `json:"ip,omitempty"`
	LocationInfo map[string]interface{} `json:"location_info,omitempty"`
	Details      map[string]interface{} `json:"details,omitempty"`
}

type SearchLogsOpts struct {
	Page          int    `url:"page,omitempty"`
	PerPage       int    `url:"per_page,omitempty"`
	Sort          string `url:"sort,omitempty"`
	Fields        string `url:"fields,omitempty"`
	IncludeFields bool   `url:"include_fields,omitempty"`
	IncludeTotals bool   `url:"include_totals,omitempty"`
	From          string `url:"from,omitempty"`
	Take          int    `url:"take,omitempty"`
	Q             string `url:"q,omitempty"`
}

func (opts *SearchLogsOpts) Encode() (string, error) {
	vals, err := query.Values(opts)
	if err != nil {
		return "", err
	}
	return vals.Encode(), nil
}

func (svc *LogsService) Search(opts SearchLogsOpts) ([]LogEntry, error) {
	var logs []LogEntry
	queryString, err := opts.Encode()
	if err != nil {
		return nil, err
	}
	url := "/logs"
	if queryString != "" {
		url = fmt.Sprintf("/logs?%s", queryString)
	}
	err = svc.c.Get(url, &logs)
	return logs, err
}
