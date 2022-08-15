package app

import (
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-wiki/server/bot"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
)

type wikiDocsService struct {
	store  WikiDocStore
	api    *pluginapi.Client
	logger bot.Logger
}

// WikiDocService is the wikiDoc service for managing playbooks
// userID is the user initiating the event.
type WikiDocService interface {
	// Get retrieves a wikiDoc. Returns ErrNotFound if not found.
	Get(id string) (WikiDoc, error)

	// Create creates a new wikiDoc
	Create(wikiDoc WikiDoc) (string, error)

	// GetPlaybooks retrieves all wikiDocs
	GetWikiDocs(requesterInfo RequesterInfo, opts WikiDocFilterOptions) (*GetWikiDocsResults, error)

	// GetWikiDocsForChannel retrieves all playbooks on the specified channel given the provided options
	//GetWikiDocsForChannel(requesterInfo RequesterInfo, channelID string, opts WikiDocFilterOptions) (GetWikiDocsResults, error)

	// Update updates a wikiDoc
	Update(wikiDoc WikiDoc) error

	// Duplicate duplicates a wikiDoc
	Duplicate(wikiDoc WikiDoc, userID string) (string, error)
}

// DialogFieldPlaybookIDKey is the key for the playbook ID field used in OpenCreatePlaybookRunDialog.
const DialogFieldPlaybookIDKey = "playbookID"

// DialogFieldNameKey is the key for the playbook run name field used in OpenCreatePlaybookRunDialog.
const DialogFieldNameKey = "playbookRunName"

// DialogFieldDescriptionKey is the key for the description textarea field used in UpdatePlaybookRunDialog
const DialogFieldDescriptionKey = "description"

func NewWikiDocService(store WikiDocStore, logger bot.Logger, api *pluginapi.Client) WikiDocService {
	return &wikiDocsService{
		store:  store,
		logger: logger,
		api:    api,
	}
}

func (s *wikiDocsService) Create(playbook WikiDoc) (string, error) {
	playbook.CreateAt = model.GetMillis()
	playbook.UpdateAt = playbook.CreateAt

	newID, err := s.store.Create(playbook)
	if err != nil {
		return "", err
	}
	playbook.ID = newID

	return newID, nil
}

func (s *wikiDocsService) Get(id string) (WikiDoc, error) {
	return s.store.Get(id)
}

func (s *wikiDocsService) GetWikiDocs(requesterInfo RequesterInfo, options WikiDocFilterOptions) (*GetWikiDocsResults, error) {
	results, err := s.store.GetWikiDocs(requesterInfo, options)
	if err != nil {
		return nil, errors.Wrap(err, "can't get playbook runs from the store")
	}
	return &GetWikiDocsResults{
		TotalCount: results.TotalCount,
		PageCount:  results.PageCount,
		HasMore:    results.HasMore,
		Items:      results.Items,
	}, nil
}

func (s *wikiDocsService) Update(playbook WikiDoc) error {
	if playbook.DeleteAt != 0 {
		return errors.New("cannot update a playbook that is archived")
	}

	playbook.UpdateAt = model.GetMillis()

	if err := s.store.Update(playbook); err != nil {
		return err
	}

	return nil
}

func (s *wikiDocsService) Duplicate(wikiDoc WikiDoc, userID string) (string, error) {
	//TODO implement me
	panic("implement me")
}
