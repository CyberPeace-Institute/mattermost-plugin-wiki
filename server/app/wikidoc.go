package app

import (
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"
	"strings"
)

const (
	StatusPrivate   = "Private"
	StatusPublished = "Published"
)

type WikiDoc struct {
	// ID is the unique identifier of the playbook run.
	ID string `json:"id" export:"-"`

	// Name is the name of the doc.
	Name string `json:"name" export:"name"`

	// Description is field for describing the doc.
	Content string `json:"content" export:"content"`

	// Description is field for describing the doc.
	Description string `json:"description" export:"description"`

	// It can be StatusPrivate ("InProgress") or StatusPublished ("Published")
	Status string `json:"status" export:"status"`

	// OwnerUserID is the user identifier of the playbook run's owner.
	OwnerUserID string `json:"owner_user_id" export:"-"`

	// TeamID is the identifier of the team the playbook run lives in.
	TeamID string `json:"team_id" export:"-"`

	// ChannelID is the identifier of the playbook run's channel.
	ChannelID string `json:"channel_id" export:"-"`

	CreateAt int64 `json:"create_at" export:"-"`
	UpdateAt int64 `json:"update_at" export:"-"`
	DeleteAt int64 `json:"delete_at" export:"-"`
}

// PlaybookStore is an interface for storing playbooks
type WikiDocStore interface {
	// Get retrieves a playbook
	Get(id string) (WikiDoc, error)

	// Create creates a new playbook
	Create(playbook WikiDoc) (string, error)

	// GetWikiDocs retrieves all playbooks
	GetWikiDocs(requesterInfo RequesterInfo, options WikiDocFilterOptions) (*GetWikiDocsResults, error)

	// Update updates a playbook
	Update(playbook WikiDoc) error

	// Archive archives a playbook
	Archive(id string) error
}

const PerPageDefault = 1000

// WikiDocFilterOptions specifies the optional parameters when getting WikiDocs.
type WikiDocFilterOptions struct {
	// Gets all the headers with this TeamID.
	TeamID string `url:"team_id,omitempty"`

	// Gets all the headers with this ChannelId.
	ChannelId string `url:"channel_id,omitempty"`

	// Pagination options.
	Page    int `url:"page,omitempty"`
	PerPage int `url:"per_page,omitempty"`

	// Sort sorts by this header field in json format (eg, "create_at", "end_at", "name", etc.);
	// defaults to "create_at".
	Sort SortField `url:"sort,omitempty"`

	// Direction orders by ascending or descending, defaulting to ascending.
	Direction SortDirection `url:"direction,omitempty"`

	// Statuses filters by all statuses in the list (inclusive)
	Statuses []string

	// OwnerID filters by owner's Mattermost user ID. Defaults to blank (no filter).
	OwnerID string `url:"owner_user_id,omitempty"`

	// SearchTerm returns results of the search term and respecting the other header filter options.
	// The search term acts as a filter and respects the Sort and Direction fields (i.e., results are
	// not returned in relevance order).
	SearchTerm string `url:"search_term,omitempty"`
}

func (o *WikiDocFilterOptions) Clone() WikiDocFilterOptions {
	newPlaybookRunFilterOptions := *o
	if len(o.Statuses) > 0 {
		newPlaybookRunFilterOptions.Statuses = append([]string{}, o.Statuses...)
	}

	return newPlaybookRunFilterOptions
}

// Validate returns a new, validated filter options or returns an error if invalid.
func (o WikiDocFilterOptions) Validate() (WikiDocFilterOptions, error) {
	options := o.Clone()

	if options.PerPage <= 0 {
		options.PerPage = PerPageDefault
	}

	options.Sort = SortField(strings.ToLower(string(options.Sort)))
	switch options.Sort {
	case SortByCreateAt:
	case SortByID:
	case SortByName:
	case SortByOwnerUserID:
	case SortByTeamID:
	case SortByChannelID:
	case SortByStatus:
	case "": // default
		options.Sort = SortByCreateAt
	default:
		return WikiDocFilterOptions{}, errors.Errorf("unsupported sort '%s'", options.Sort)
	}

	options.Direction = SortDirection(strings.ToUpper(string(options.Direction)))
	switch options.Direction {
	case DirectionAsc:
	case DirectionDesc:
	case "": //default
		options.Direction = DirectionAsc
	default:
		return WikiDocFilterOptions{}, errors.Errorf("unsupported direction '%s'", options.Direction)
	}

	if options.TeamID != "" && !model.IsValidId(options.TeamID) {
		return WikiDocFilterOptions{}, errors.New("bad parameter 'team_id': must be 26 characters or blank")
	}

	if options.ChannelId != "" && !model.IsValidId(options.ChannelId) {
		return WikiDocFilterOptions{}, errors.New("bad parameter 'channel_id': must be 26 characters or blank")
	}

	if options.OwnerID != "" && !model.IsValidId(options.OwnerID) {
		return WikiDocFilterOptions{}, errors.New("bad parameter 'owner_id': must be 26 characters or blank")
	}

	for _, s := range options.Statuses {
		if !ValidStatus(s) {
			return WikiDocFilterOptions{}, errors.New("bad parameter in 'statuses': must be InProgress or Finished")
		}
	}

	return options, nil
}

func ValidStatus(status string) bool {
	return status == "" || status == StatusPrivate || status == StatusPublished
}

type GetWikiDocsResults struct {
	TotalCount int       `json:"total_count"`
	PageCount  int       `json:"page_count"`
	HasMore    bool      `json:"has_more"`
	Items      []WikiDoc `json:"items"`
}
