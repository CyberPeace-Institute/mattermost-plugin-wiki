package sqlstore

import (
	"database/sql"
	"fmt"
	"github.com/CyberPeace-Institute/mattermost-plugin-wiki/server/app"
	"github.com/CyberPeace-Institute/mattermost-plugin-wiki/server/bot"
	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"
	"math"
)

type sqlWikiDoc struct {
	app.WikiDoc
}

// wikiDocStore is a sql store for wikiDocs. Use NewWikiDocStore to create it.
type wikiDocStore struct {
	pluginAPI     PluginAPIClient
	log           bot.Logger
	store         *SQLStore
	queryBuilder  sq.StatementBuilderType
	wikiDocSelect sq.SelectBuilder
}

// Ensure wikiDocStore implements the wikiDoc.Store interface.
var _ app.WikiDocStore = (*wikiDocStore)(nil)

func applyWikiDocFilterOptionsSort(builder sq.SelectBuilder, options app.WikiDocFilterOptions) (sq.SelectBuilder, error) {
	var sort string
	switch options.Sort {
	case app.SortByID:
		sort = "ID"
	case app.SortByCreateAt:
		sort = "CreateAt"
	case app.SortByName:
		sort = "Name"
	case app.SortByStatus:
		sort = "Status"
	case "":
		// Default to a stable sort if none explicitly provided.
		sort = "ID"
	default:
		return sq.SelectBuilder{}, errors.Errorf("unsupported sort parameter '%s'", options.Sort)
	}

	var direction string
	switch options.Direction {
	case app.DirectionAsc:
		direction = "ASC"
	case app.DirectionDesc:
		direction = "DESC"
	case "":
		// Default to an ascending sort if none explicitly provided.
		direction = "ASC"
	default:
		return sq.SelectBuilder{}, errors.Errorf("unsupported direction parameter '%s'", options.Direction)
	}

	builder = builder.OrderByClause(fmt.Sprintf("%s %s", sort, direction))

	page := options.Page
	perPage := options.PerPage
	if page < 0 {
		page = 0
	}
	if perPage < 0 {
		perPage = 0
	}

	builder = builder.
		Offset(uint64(page * perPage)).
		Limit(uint64(perPage))

	return builder, nil
}

// NewWikiDocStore creates a new store for wikiDoc service.
func NewWikiDocStore(pluginAPI PluginAPIClient, log bot.Logger, sqlStore *SQLStore) app.WikiDocStore {
	wikiDocSelect := sqlStore.builder.
		Select(
			"p.ID",
			"p.TeamID",
			"p.ChannelID",
			"p.CreateAt",
			"p.UpdateAt",
		).
		From("CPI_WikiDocs p")

	newStore := &wikiDocStore{
		pluginAPI:     pluginAPI,
		log:           log,
		store:         sqlStore,
		queryBuilder:  sqlStore.builder,
		wikiDocSelect: wikiDocSelect,
	}
	return newStore
}

// Create creates a new wikiDoc
func (p *wikiDocStore) Create(wikiDoc app.WikiDoc) (id string, err error) {
	/*if wikiDoc.ID != "" {
		return "", errors.New("ID should be empty")
	}*/
	wikiDoc.ID = model.NewId()

	rawWikiDoc, err := toSQLWikiDoc(wikiDoc)
	if err != nil {
		return "", err
	}

	tx, err := p.store.db.Beginx()
	if err != nil {
		return "", errors.Wrap(err, "could not begin transaction")
	}
	defer p.store.finalizeTransaction(tx)

	_, err = p.store.execBuilder(tx, sq.
		Insert("CPI_WikiDocs").
		SetMap(map[string]interface{}{
			"ID":          rawWikiDoc.ID,
			"Name":        rawWikiDoc.Name,
			"Content":     rawWikiDoc.Content,
			"Status":      rawWikiDoc.Status,
			"OwnerUserID": rawWikiDoc.OwnerUserID,
			"TeamID":      rawWikiDoc.TeamID,
			"ChannelID":   rawWikiDoc.ChannelID,
			"Description": rawWikiDoc.Description,
			"CreateAt":    rawWikiDoc.CreateAt,
			"UpdateAt":    rawWikiDoc.UpdateAt,
			"DeleteAt":    rawWikiDoc.DeleteAt,
		}))
	if err != nil {
		return "", errors.Wrap(err, "failed to store new wikiDoc")
	}

	if err = tx.Commit(); err != nil {
		return "", errors.Wrap(err, "could not commit transaction")
	}

	return rawWikiDoc.ID, nil
}

// Get retrieves a wikiDoc
func (p *wikiDocStore) Get(id string) (app.WikiDoc, error) {
	if id == "" {
		return app.WikiDoc{}, errors.New("ID cannot be empty")
	}

	tx, err := p.store.db.Beginx()
	if err != nil {
		return app.WikiDoc{}, errors.Wrap(err, "could not begin transaction")
	}
	defer p.store.finalizeTransaction(tx)

	var rawWikiDoc sqlWikiDoc
	err = p.store.getBuilder(tx, &rawWikiDoc, p.wikiDocSelect.Where(sq.Eq{"p.ID": id}))
	if err == sql.ErrNoRows {
		return app.WikiDoc{}, errors.Wrapf(app.ErrNotFound, "wikiDoc does not exist for id '%s'", id)
	} else if err != nil {
		return app.WikiDoc{}, errors.Wrapf(err, "failed to get wikiDoc by id '%s'", id)
	}

	wikiDoc, err := toWikiDoc(rawWikiDoc)
	if err != nil {
		return app.WikiDoc{}, err
	}

	if err = tx.Commit(); err != nil {
		return app.WikiDoc{}, errors.Wrap(err, "could not commit transaction")
	}

	return wikiDoc, nil
}

// GetWikiDocs retrieves all wikiDocs that are not deleted.
// Members are not retrieved for this as the query would be large and we don't need it for this for now.
// This is only used for the keywords feature
func (p *wikiDocStore) GetWikiDocs(requesterInfo app.RequesterInfo, options app.WikiDocFilterOptions) (*app.GetWikiDocsResults, error) {
	queryForTotal := p.store.builder.
		Select("COUNT(*)").
		From("CPI_WikiDocs AS w") /*.
		Join("Channels AS c ON (c.Id = w.ChannelId)")*/

	queryForResults := p.store.builder.
		Select(
			"p.ID",
			"p.Name",
			"p.Content",
			"p.Description",
			"p.Status",
			"p.OwnerUserID",
			"p.TeamID",
			"p.ChannelID",
			"p.CreateAt",
			"p.UpdateAt",
			"p.DeleteAt",
		).
		From("CPI_WikiDocs AS p").
		Where(sq.Eq{"p.DeleteAt": 0})

	if options.OwnerID != "" {
		queryForResults = queryForResults.Where(sq.Eq{"w.OwnerUserID": options.OwnerID})
		queryForTotal = queryForTotal.Where(sq.Eq{"w.OwnerUserID": options.OwnerID})
	}

	if options.TeamID != "" {
		queryForResults = queryForResults.Where(sq.Eq{"w.TeamID": options.TeamID})
		queryForTotal = queryForTotal.Where(sq.Eq{"w.TeamID": options.TeamID})
	}

	if options.ChannelId != "" {
		queryForResults = queryForResults.Where(sq.Eq{"w.ChannelID": options.ChannelId})
		queryForTotal = queryForTotal.Where(sq.Eq{"w.ChannelID": options.ChannelId})
	}

	queryForResults, err := applyWikiDocFilterOptionsSort(queryForResults, options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply sort options")
	}

	tx, err := p.store.db.Beginx()
	if err != nil {
		return nil, errors.Wrap(err, "could not begin transaction")
	}
	defer p.store.finalizeTransaction(tx)

	var wikiDocs []app.WikiDoc
	err = p.store.selectBuilder(tx, &wikiDocs, queryForResults)

	if err == sql.ErrNoRows {
		return nil, errors.Wrap(app.ErrNotFound, "no wikiDocs found")
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to get wikiDocs")
	}

	var total int
	if err = p.store.getBuilder(tx, &total, queryForTotal); err != nil {
		return nil, errors.Wrap(err, "failed to get total count")
	}
	pageCount := 0
	if options.PerPage > 0 {
		pageCount = int(math.Ceil(float64(total) / float64(options.PerPage)))
	}
	hasMore := options.Page+1 < pageCount

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "could not commit transaction")
	}

	return &app.GetWikiDocsResults{
		TotalCount: total,
		PageCount:  pageCount,
		HasMore:    hasMore,
		Items:      wikiDocs,
	}, nil
}

func (p *wikiDocStore) buildPermissionsExpr(info app.RequesterInfo) sq.Sqlizer {
	if info.IsAdmin {
		return nil
	}

	// Guests must be channel members
	if info.IsGuest {
		return sq.Expr(`(
			p.Status = "Published"
			AND
			  EXISTS(SELECT 1
						 FROM ChannelMembers as cm
						 WHERE cm.ChannelId = p.ChannelID
						   AND cm.UserId = ?)
		)`, info.UserID)
	}

	// 1. Is the user a channel member? If so, they have permission to view the run.
	// 2. Is the wikiDoc open to everyone on the team, or is the user a member of the wikiDoc?
	//    If so, they have permission to view the run.
	return sq.Expr(`
        (
			EXISTS(SELECT 1
					 FROM ChannelMembers as cm
					 WHERE cm.ChannelId = p.ChannelID
					   AND cm.UserId = ?)
		)`, info.UserID)
}

// Update updates a wikidoc
func (p *wikiDocStore) Update(wikiDoc app.WikiDoc) (err error) {
	if wikiDoc.ID == "" {
		return errors.New("id should not be empty")
	}

	rawWikiDoc, err := toSQLWikiDoc(wikiDoc)
	if err != nil {
		return err
	}

	tx, err := p.store.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "could not begin transaction")
	}
	defer p.store.finalizeTransaction(tx)

	_, err = p.store.execBuilder(tx, sq.
		Update("CPI_WikiDocs").
		SetMap(map[string]interface{}{
			"Name":        rawWikiDoc.Name,
			"Content":     rawWikiDoc.Content,
			"Status":      rawWikiDoc.Status,
			"OwnerUserID": rawWikiDoc.OwnerUserID,
			"TeamID":      rawWikiDoc.TeamID,
			"ChannelID":   rawWikiDoc.ChannelID,
			"Description": rawWikiDoc.Description,
			"UpdateAt":    rawWikiDoc.UpdateAt,
			"DeleteAt":    rawWikiDoc.DeleteAt,
		}).
		Where(sq.Eq{"ID": rawWikiDoc.ID}))

	if err != nil {
		return errors.Wrapf(err, "failed to update wikiDoc with id '%s'", rawWikiDoc.ID)
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "could not commit transaction")
	}

	return nil
}

// Archive archives a wikiDoc.
func (p *wikiDocStore) Archive(id string) error {
	if id == "" {
		return errors.New("ID cannot be empty")
	}

	_, err := p.store.execBuilder(p.store.db, sq.
		Update("CPI_WikiDocs").
		Set("DeleteAt", model.GetMillis()).
		Where(sq.Eq{"ID": id}))

	if err != nil {
		return errors.Wrapf(err, "failed to delete wikiDoc with id '%s'", id)
	}

	return nil
}

func toSQLWikiDoc(wikiDocs app.WikiDoc) (*sqlWikiDoc, error) {
	return &sqlWikiDoc{
		WikiDoc: wikiDocs,
	}, nil
}

func toWikiDoc(rawWikiDoc sqlWikiDoc) (app.WikiDoc, error) {
	p := rawWikiDoc.WikiDoc

	return p, nil
}
