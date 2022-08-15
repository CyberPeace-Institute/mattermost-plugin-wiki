package sqlstore

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/mattermost-plugin-wiki/server/app"
	"github.com/mattermost/mattermost-plugin-wiki/server/bot"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"
	"math"
)

type sqlWikiDoc struct {
	app.WikiDoc
}

// wikiDocStore is a sql store for playbooks. Use NewPlaybookStore to create it.
type wikiDocStore struct {
	pluginAPI      PluginAPIClient
	log            bot.Logger
	store          *SQLStore
	queryBuilder   sq.StatementBuilderType
	playbookSelect sq.SelectBuilder
}

// Ensure wikiDocStore implements the playbook.Store interface.
var _ app.WikiDocStore = (*wikiDocStore)(nil)

func applyWikiDocFilterOptionsSort(builder sq.SelectBuilder, options app.WikiDocFilterOptions) (sq.SelectBuilder, error) {
	var sort string
	switch options.Sort {
	case app.SortByID:
		sort = "ID"
	case app.SortByCreateAt:
		sort = "CreateAt"
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

// NewWikiDocStore creates a new store for playbook service.
func NewWikiDocStore(pluginAPI PluginAPIClient, log bot.Logger, sqlStore *SQLStore) app.WikiDocStore {
	playbookSelect := sqlStore.builder.
		Select(
			"p.ID",
			"p.TeamID",
			"p.ChannelID",
			"p.CreateAt",
			"p.UpdateAt",
		).
		From("CPI_WikiDocs p")

	newStore := &wikiDocStore{
		pluginAPI:      pluginAPI,
		log:            log,
		store:          sqlStore,
		queryBuilder:   sqlStore.builder,
		playbookSelect: playbookSelect,
	}
	return newStore
}

// Create creates a new playbook
func (p *wikiDocStore) Create(playbook app.WikiDoc) (id string, err error) {
	/*if playbook.ID != "" {
		return "", errors.New("ID should be empty")
	}*/
	playbook.ID = model.NewId()

	rawPlaybook, err := toSQLWikiDoc(playbook)
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
			"ID":          rawPlaybook.ID,
			"Name":        rawPlaybook.Name,
			"Content":     rawPlaybook.Content,
			"Status":      rawPlaybook.Status,
			"OwnerUserID": rawPlaybook.OwnerUserID,
			"TeamID":      rawPlaybook.TeamID,
			"ChannelID":   rawPlaybook.ChannelID,
			"Description": rawPlaybook.Description,
			"CreateAt":    rawPlaybook.CreateAt,
			"UpdateAt":    rawPlaybook.UpdateAt,
			"DeleteAt":    rawPlaybook.DeleteAt,
		}))
	if err != nil {
		return "", errors.Wrap(err, "failed to store new playbook")
	}

	if err = tx.Commit(); err != nil {
		return "", errors.Wrap(err, "could not commit transaction")
	}

	return rawPlaybook.ID, nil
}

// Get retrieves a playbook
func (p *wikiDocStore) Get(id string) (app.WikiDoc, error) {
	if id == "" {
		return app.WikiDoc{}, errors.New("ID cannot be empty")
	}

	tx, err := p.store.db.Beginx()
	if err != nil {
		return app.WikiDoc{}, errors.Wrap(err, "could not begin transaction")
	}
	defer p.store.finalizeTransaction(tx)

	var rawPlaybook sqlWikiDoc
	err = p.store.getBuilder(tx, &rawPlaybook, p.playbookSelect.Where(sq.Eq{"p.ID": id}))
	if err == sql.ErrNoRows {
		return app.WikiDoc{}, errors.Wrapf(app.ErrNotFound, "playbook does not exist for id '%s'", id)
	} else if err != nil {
		return app.WikiDoc{}, errors.Wrapf(err, "failed to get playbook by id '%s'", id)
	}

	playbook, err := toWikiDoc(rawPlaybook)
	if err != nil {
		return app.WikiDoc{}, err
	}

	if err = tx.Commit(); err != nil {
		return app.WikiDoc{}, errors.Wrap(err, "could not commit transaction")
	}

	return playbook, nil
}

// GetWikiDocs retrieves all playbooks that are not deleted.
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
	// 2. Is the playbook open to everyone on the team, or is the user a member of the playbook?
	//    If so, they have permission to view the run.
	return sq.Expr(`
        (
			EXISTS(SELECT 1
					 FROM ChannelMembers as cm
					 WHERE cm.ChannelId = p.ChannelID
					   AND cm.UserId = ?)
		)`, info.UserID)
}

// GetWikiDocsForChannel retrieves all playbooks on the specified team given the provided options.
/*func (p *wikiDocStore) GetWikiDocsForChannel(requesterInfo app.RequesterInfo, channelID string, opts app.WikiDocFilterOptions) (app.GetPlaybooksResults, error) {
	// Check that you are a playbook member or there are no restrictions.
	permissionsAndFilter := p.buildPermissionsExpr(requesterInfo)

	queryForResults := p.store.builder.
		Select(
			"p.ID",
			"p.Name",
			"p.Description",
			"p.TeamID",
			"p.CreateAt",
			"p.DeleteAt",
		).
		From("CPI_WikiDocs AS p").
		//LeftJoin("Channels c ON c.Id = p.ChannelID").
		GroupBy("p.ID").
		Where(permissionsAndFilter)

	queryForResults, err := applyPlaybookFilterOptionsSort(queryForResults, opts)
	if err != nil {
		return app.GetPlaybooksResults{}, errors.Wrap(err, "failed to apply sort options")
	}

	queryForTotal := p.store.builder.
		Select("COUNT(*)").
		From("CPI_WikiDocs AS p").
		Where(permissionsAndFilter)

	if opts.SearchTerm != "" {
		column := "p.Title"
		searchString := opts.SearchTerm

		// Postgres performs a case-sensitive search, so we need to lowercase
		// both the column contents and the search string
		if p.store.db.DriverName() == model.DatabaseDriverPostgres {
			column = "LOWER(p.Name)"
			searchString = strings.ToLower(opts.SearchTerm)
		}

		queryForResults = queryForResults.Where(sq.Like{column: fmt.Sprint("%", searchString, "%")})
		queryForTotal = queryForTotal.Where(sq.Like{column: fmt.Sprint("%", searchString, "%")})
	}

	if !opts.WithArchived {
		queryForResults = queryForResults.Where(sq.Eq{"p.DeleteAt": 0})
		queryForTotal = queryForTotal.Where(sq.Eq{"DeleteAt": 0})
	}

	var playbooks []app.WikiDoc
	err = p.store.selectBuilder(p.store.db, &playbooks, queryForResults)
	if err == sql.ErrNoRows {
		return app.GetPlaybooksResults{}, errors.Wrap(app.ErrNotFound, "no playbooks found")
	} else if err != nil {
		return app.GetPlaybooksResults{}, errors.Wrap(err, "failed to get playbooks")
	}

	var total int
	if err = p.store.getBuilder(p.store.db, &total, queryForTotal); err != nil {
		return app.GetWikiDocsResults{}, errors.Wrap(err, "failed to get total count")
	}

	ids := make([]string, len(playbooks))
	for _, pb := range playbooks {
		ids = append(ids, pb.ID)
	}

	pageCount := 0
	if opts.PerPage > 0 {
		pageCount = int(math.Ceil(float64(total) / float64(opts.PerPage)))
	}
	hasMore := opts.Page+1 < pageCount

	return app.GetWikiDocsResults{
		TotalCount: total,
		PageCount:  pageCount,
		HasMore:    hasMore,
		Items:      playbooks,
	}, nil
}
*/

// Update updates a playbook
func (p *wikiDocStore) Update(playbook app.WikiDoc) (err error) {
	if playbook.ID == "" {
		return errors.New("id should not be empty")
	}

	rawPlaybook, err := toSQLWikiDoc(playbook)
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
			"Name":        rawPlaybook.Name,
			"Content":     rawPlaybook.Content,
			"Status":      rawPlaybook.Status,
			"OwnerUserID": rawPlaybook.OwnerUserID,
			"TeamID":      rawPlaybook.TeamID,
			"ChannelID":   rawPlaybook.ChannelID,
			"Description": rawPlaybook.Description,
			"UpdateAt":    rawPlaybook.UpdateAt,
			"DeleteAt":    rawPlaybook.DeleteAt,
		}).
		Where(sq.Eq{"ID": rawPlaybook.ID}))

	if err != nil {
		return errors.Wrapf(err, "failed to update playbook with id '%s'", rawPlaybook.ID)
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "could not commit transaction")
	}

	return nil
}

// Archive archives a playbook.
func (p *wikiDocStore) Archive(id string) error {
	if id == "" {
		return errors.New("ID cannot be empty")
	}

	_, err := p.store.execBuilder(p.store.db, sq.
		Update("CPI_WikiDocs").
		Set("DeleteAt", model.GetMillis()).
		Where(sq.Eq{"ID": id}))

	if err != nil {
		return errors.Wrapf(err, "failed to delete playbook with id '%s'", id)
	}

	return nil
}

func toSQLWikiDoc(wikiDocs app.WikiDoc) (*sqlWikiDoc, error) {
	return &sqlWikiDoc{
		WikiDoc: wikiDocs,
	}, nil
}

func toWikiDoc(rawPlaybook sqlWikiDoc) (app.WikiDoc, error) {
	p := rawPlaybook.WikiDoc

	return p, nil
}
