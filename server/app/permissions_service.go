package app

import (
	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"
)

// ErrNoPermissions if the error is caused by the user not having permissions
var ErrNoPermissions = errors.New("does not have permissions")

type PermissionsService struct {
	wikiDocsService WikiDocService
	pluginAPI       *pluginapi.Client
}

func NewPermissionsService(
	wikiDocsService WikiDocService,
	pluginAPI *pluginapi.Client,
) *PermissionsService {
	return &PermissionsService{
		wikiDocsService,
		pluginAPI,
	}
}

func (p *PermissionsService) WikiDocIsPublic(wikiDoc WikiDoc) bool {
	return wikiDoc.Status == StatusPublished
}

func (p *PermissionsService) HasEditPermissionsToWikiDocs(userID string, wikiDoc WikiDoc) error {
	if IsSystemAdmin(userID, p.pluginAPI) || CanManageChannelProperties(userID, wikiDoc.ChannelID, p.pluginAPI) {
		return nil
	}

	return ErrNoPermissions
}

func (p *PermissionsService) canReadChannel(userID string, channelID string) bool {
	if channelID == "" || userID == "" {
		return false
	}

	return p.pluginAPI.User.HasPermissionToTeam(userID, channelID, model.PermissionReadChannel)
}

func (p *PermissionsService) WikiDocCreate(wikiDoc WikiDoc) error {
	if IsSystemAdmin(wikiDoc.OwnerUserID, p.pluginAPI) || CanManageChannelProperties(wikiDoc.OwnerUserID, wikiDoc.ChannelID, p.pluginAPI) {
		return nil
	}

	return ErrNoPermissions
}

func (p *PermissionsService) DeleteWikiDoc(userID string, wikiDoc WikiDoc) error {
	if IsSystemAdmin(userID, p.pluginAPI) || CanManageChannelProperties(userID, wikiDoc.ChannelID, p.pluginAPI) {
		return nil
	}

	return ErrNoPermissions
}

func (p *PermissionsService) WikiDocView(userID string, wikiDocID string) error {
	wikiDoc, err := p.wikiDocsService.Get(wikiDocID)
	if err != nil {
		return errors.Wrapf(err, "Unable to get wikidoc to determine permissions, wikiDoc id `%s`", wikiDocID)
	}

	if p.canReadChannel(userID, wikiDoc.ChannelID) {
		return nil
	}

	return ErrNoPermissions
}

func (p *PermissionsService) WikiDocList(userID string, channelID string) error {
	// Can list wikiDocs if you are on the team
	if p.canReadChannel(userID, channelID) {
		return nil
	}

	return ErrNoPermissions
}

func (p *PermissionsService) WikiDocMakePrivate(userID string, wikiDoc WikiDoc) error {
	if IsSystemAdmin(userID, p.pluginAPI) || CanManageChannelProperties(userID, wikiDoc.ChannelID, p.pluginAPI) {
		return nil
	}

	return ErrNoPermissions
}

func (p *PermissionsService) WikiDocMakePublic(userID string, wikiDoc WikiDoc) error {
	if IsSystemAdmin(userID, p.pluginAPI) || CanManageChannelProperties(userID, wikiDoc.ChannelID, p.pluginAPI) {
		return nil
	}

	return ErrNoPermissions
}

// IsSystemAdmin returns true if the userID is a system admin
func IsSystemAdmin(userID string, pluginAPI *pluginapi.Client) bool {
	return pluginAPI.User.HasPermissionTo(userID, model.PermissionManageSystem)
}

// IsGuest returns true if the userID is a system guest
func IsGuest(userID string, pluginAPI *pluginapi.Client) (bool, error) {
	user, err := pluginAPI.User.Get(userID)
	if err != nil {
		return false, errors.Wrapf(err, "Unable to get user to determine permissions, user id `%s`", userID)
	}

	return user.IsGuest(), nil
}

// CanManageChannelProperties returns true if the userID is allowed to manage the properties of channelID
func CanManageChannelProperties(userID, channelID string, pluginAPI *pluginapi.Client) bool {
	channel, err := pluginAPI.Channel.Get(channelID)
	if err != nil {
		return false
	}

	permission := model.PermissionManagePublicChannelProperties
	if channel.Type == model.ChannelTypePrivate {
		permission = model.PermissionManagePrivateChannelProperties
	}

	return pluginAPI.User.HasPermissionToChannel(userID, channelID, permission)
}

func CanPostToChannel(userID, channelID string, pluginAPI *pluginapi.Client) bool {
	return pluginAPI.User.HasPermissionToChannel(userID, channelID, model.PermissionCreatePost)
}

func IsMemberOfTeam(userID, teamID string, pluginAPI *pluginapi.Client) bool {
	teamMember, err := pluginAPI.Team.GetMember(teamID, userID)
	if err != nil {
		return false
	}

	return teamMember.DeleteAt == 0
}

// RequesterInfo holds the userID and teamID that this request is regarding, and permissions
// for the user making the request
type RequesterInfo struct {
	UserID  string
	TeamID  string
	IsAdmin bool
	IsGuest bool
}

func GetRequesterInfo(userID string, pluginAPI *pluginapi.Client) (RequesterInfo, error) {
	isAdmin := IsSystemAdmin(userID, pluginAPI)

	isGuest, err := IsGuest(userID, pluginAPI)
	if err != nil {
		return RequesterInfo{}, err
	}

	return RequesterInfo{
		UserID:  userID,
		IsAdmin: isAdmin,
		IsGuest: isGuest,
	}, nil
}
