package app

// SortField enumerates the available fields we can sort on.
type SortField string

const (
	// SortByCreateAt sorts by the created time of a wikiDoc.
	SortByCreateAt SortField = "create_at"

	// SortByID sorts by the primary key of a wikiDoc.
	SortByID SortField = "id"

	// SortByTeamID sorts by the team id of a wikiDoc.
	SortByTeamID SortField = "team_id"

	// SortByChannelID sorts by the channel id of a wikiDoc.
	SortByChannelID SortField = "channel_id"

	// SortByName sorts by the name of a wikiDoc run.
	SortByName SortField = "name"

	// SortByOwnerUserID sorts by the user id of the owner of a wikiDoc run.
	SortByOwnerUserID SortField = "owner_user_id"

	// SortByStatus sorts by the status of a wikiDoc run.
	SortByStatus SortField = "status"
)

// SortDirection is the type used to specify the ascending or descending order of returned results.
type SortDirection string

const (
	// DirectionDesc is descending order.
	DirectionDesc SortDirection = "DESC"

	// DirectionAsc is ascending order.
	DirectionAsc SortDirection = "ASC"
)

func IsValidDirection(direction SortDirection) bool {
	return direction == DirectionAsc || direction == DirectionDesc
}
