package app

import (
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"

	"github.com/CyberPeace-Institute/mattermost-plugin-wiki/server/bot"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
)

type wikiDocsService struct {
	store  WikiDocStore
	api    *pluginapi.Client
	logger bot.Logger
}

// WikiDocService is the wikiDoc service for managing wikiDocs
// userID is the user initiating the event.
type WikiDocService interface {
	// Get retrieves a wikiDoc. Returns ErrNotFound if not found.
	Get(id string) (WikiDoc, error)

	// Create creates a new wikiDoc
	Create(wikiDoc WikiDoc) (string, error)

	// GetWikiDocs retrieves all wikiDocs
	GetWikiDocs(requesterInfo RequesterInfo, opts WikiDocFilterOptions) (*GetWikiDocsResults, error)

	// GetWikiDocsForChannel retrieves all wikiDocs on the specified channel given the provided options
	//GetWikiDocsForChannel(requesterInfo RequesterInfo, channelID string, opts WikiDocFilterOptions) (GetWikiDocsResults, error)

	// Update updates a wikiDoc
	Update(wikiDoc WikiDoc) error

	// Duplicate duplicates a wikiDoc
	Duplicate(wikiDoc WikiDoc, userID string) (string, error)
}

// DialogFieldWikiDocIDKey is the key for the wikiDoc ID field used in OpenCreateWikiDocRunDialog.
const DialogFieldWikiDocIDKey = "wikiDocID"

// DialogFieldNameKey is the key for the wikiDoc run name field used in OpenCreateWikiDocRunDialog.
const DialogFieldNameKey = "wikiDocName"

// DialogFieldDescriptionKey is the key for the description textarea field used in UpdateWikiDocRunDialog
const DialogFieldDescriptionKey = "description"

func NewWikiDocService(store WikiDocStore, logger bot.Logger, api *pluginapi.Client) WikiDocService {
	return &wikiDocsService{
		store:  store,
		logger: logger,
		api:    api,
	}
}

func (s *wikiDocsService) Create(wikiDoc WikiDoc) (string, error) {
	wikiDoc.CreateAt = model.GetMillis()
	wikiDoc.UpdateAt = wikiDoc.CreateAt

	newID, err := s.store.Create(wikiDoc)
	if err != nil {
		return "", err
	}
	wikiDoc.ID = newID

	return newID, nil
}

func (s *wikiDocsService) Get(id string) (WikiDoc, error) {
	return s.store.Get(id)
}

func (s *wikiDocsService) GetWikiDocs(requesterInfo RequesterInfo, options WikiDocFilterOptions) (*GetWikiDocsResults, error) {
	results, err := s.store.GetWikiDocs(requesterInfo, options)
	if err != nil {
		return nil, errors.Wrap(err, "can't get wikiDoc runs from the store")
	}
	return &GetWikiDocsResults{
		TotalCount: results.TotalCount,
		PageCount:  results.PageCount,
		HasMore:    results.HasMore,
		Items:      results.Items,
	}, nil
}

func (s *wikiDocsService) Update(wikiDoc WikiDoc) error {
	if wikiDoc.DeleteAt != 0 {
		return errors.New("cannot update a wikiDoc that is archived")
	}

	wikiDoc.UpdateAt = model.GetMillis()

	if err := s.store.Update(wikiDoc); err != nil {
		return err
	}

	return nil
}

func (s *wikiDocsService) Duplicate(wikiDoc WikiDoc, userID string) (string, error) {
	//TODO implement me
	panic("implement me")
}
