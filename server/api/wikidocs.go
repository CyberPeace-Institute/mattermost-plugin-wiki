package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	pluginapi "github.com/mattermost/mattermost-plugin-api"

	"github.com/mattermost/mattermost-plugin-wiki/server/app"
	"github.com/mattermost/mattermost-plugin-wiki/server/bot"
)

// WikiDocHandler is the API handler.
type WikiDocHandler struct {
	*ErrorHandler
	wikiDocService app.WikiDocService
	permissions    *app.PermissionsService
	pluginAPI      *pluginapi.Client
	log            bot.Logger
}

// NewWikiDocHandler Creates a new Plugin API handler.
func NewWikiDocHandler(
	router *mux.Router,
	wikiDocService app.WikiDocService,
	permissions *app.PermissionsService,
	api *pluginapi.Client,
	log bot.Logger,
) *WikiDocHandler {
	handler := &WikiDocHandler{
		ErrorHandler:   &ErrorHandler{log: log},
		wikiDocService: wikiDocService,
		pluginAPI:      api,
		log:            log,
		permissions:    permissions,
	}

	wikiDocsRouter := router.PathPrefix("/wikiDocs").Subrouter()
	wikiDocsRouter.HandleFunc("", handler.getWikiDocs).Methods(http.MethodGet)

	wikiDocsRouter.HandleFunc("/dialog", handler.createWikiDocFromDialog).Methods(http.MethodPost)

	wikiDocRouter := wikiDocsRouter.PathPrefix("/{id:[A-Za-z0-9]+}").Subrouter()
	wikiDocRouter.HandleFunc("", handler.getWikiDoc).Methods(http.MethodGet)

	wikiDocRouterAuthorized := wikiDocRouter.PathPrefix("").Subrouter()
	wikiDocRouterAuthorized.Use(handler.checkEditPermissions)
	wikiDocRouterAuthorized.HandleFunc("", handler.updateWikiDoc).Methods(http.MethodPatch)
	wikiDocRouterAuthorized.HandleFunc("/status", handler.status).Methods(http.MethodPost)

	//channelRouter := wikiDocsRouter.PathPrefix("/channel").Subrouter()
	//channelRouter.HandleFunc("/{channel_id:[A-Za-z0-9]+}", handler.getWikiDocByChannel).Methods(http.MethodGet)

	return handler
}

func (h *WikiDocHandler) checkEditPermissions(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := r.Header.Get("Mattermost-User-ID")

		wikiDoc, err := h.wikiDocService.Get(vars["id"])
		if err != nil {
			h.HandleError(w, err)
			return
		}

		if !h.PermissionsCheck(w, h.permissions.HasEditPermissionsToWikiDocs(userID, wikiDoc)) {
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Note that this currently does nothing. This is temporary given the removal of stages. Will be used by status.
func (h *WikiDocHandler) updateWikiDoc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wikiDocID := vars["id"]
	userID := r.Header.Get("Mattermost-User-ID")

	var wikiDoc app.WikiDoc
	if err := json.NewDecoder(r.Body).Decode(&wikiDoc); err != nil {
		h.HandleErrorWithCode(w, http.StatusBadRequest, "unable to decode playbook", err)
		return
	}

	wikiDoc.ID = wikiDocID
	oldWikiDoc, err := h.wikiDocService.Get(wikiDocID)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	if !h.PermissionsCheck(w, h.permissions.HasEditPermissionsToWikiDocs(userID, oldWikiDoc)) {
		return
	}

	err = h.wikiDocService.Update(wikiDoc)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	updatedWikiDoc := wikiDoc

	ReturnJSON(w, updatedWikiDoc, http.StatusOK)
}

// createWikiDocFromDialog handles the interactive dialog submission when a user presses confirm on
// the create wiki doc dialog.
func (h *WikiDocHandler) createWikiDocFromDialog(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")

	var request *model.SubmitDialogRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request == nil {
		h.HandleErrorWithCode(w, http.StatusBadRequest, "failed to decode SubmitDialogRequest", err)
		return
	}

	if userID != request.UserId {
		h.HandleErrorWithCode(w, http.StatusBadRequest, "interactive dialog's userID must be the same as the requester's userID", nil)
		return
	}

	var name, description string
	if rawName, ok := request.Submission[app.DialogFieldNameKey].(string); ok {
		name = rawName
	}

	if rawDescription, ok := request.Submission[app.DialogFieldNameKey].(string); ok {
		description = rawDescription
	}

	wikDocId, err := h.createWikiDoc(
		app.WikiDoc{
			OwnerUserID: request.UserId,
			TeamID:      request.TeamId,
			ChannelID:   request.ChannelId,
			Name:        name,
			Description: description,
		},
		request.UserId,
	)
	if err != nil {
		if errors.Is(err, app.ErrMalformedWikiDoc) {
			h.HandleErrorWithCode(w, http.StatusBadRequest, "unable to create wikiDoc", err)
			return
		}

		if errors.Is(err, app.ErrNoPermissions) {
			h.HandleErrorWithCode(w, http.StatusForbidden, "not authorized to make a wikiDoc", err)
			return
		}

		var msg string

		if errors.Is(err, app.ErrChannelDisplayNameInvalid) {
			msg = "The name is invalid or too long. Please use a valid name with fewer than 64 characters."
		}

		if msg != "" {
			resp := &model.SubmitDialogResponse{
				Errors: map[string]string{
					app.DialogFieldNameKey: msg,
				},
			}
			respBytes, _ := json.Marshal(resp)
			_, _ = w.Write(respBytes)
			return
		}

		h.HandleError(w, err)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("/api/v0/wikiDocs/%s", wikDocId))
	w.WriteHeader(http.StatusCreated)
}

func (h *WikiDocHandler) createWikiDoc(wikiDoc app.WikiDoc, userID string) (string, error) {
	if wikiDoc.ID != "" {
		return "", errors.Wrap(app.ErrMalformedWikiDoc, "wikiDoc already has an id")
	}

	if wikiDoc.CreateAt != 0 {
		return "", errors.Wrap(app.ErrMalformedWikiDoc, "wikiDoc already has created at date")
	}

	if wikiDoc.ChannelID == "" {
		return "", errors.Wrap(app.ErrMalformedWikiDoc, "must provide a channel to create a wikiDoc")
	}

	// If a channel is specified, ensure it's from the given team (if one provided), or
	// just grab the team for that channel.
	var channel *model.Channel
	var err error
	if wikiDoc.ChannelID != "" {
		channel, err = h.pluginAPI.Channel.Get(wikiDoc.ChannelID)
		if err != nil {
			return "", errors.Wrapf(err, "failed to get channel")
		}

		if wikiDoc.TeamID == "" {
			wikiDoc.TeamID = channel.TeamId
		} else if channel.TeamId != wikiDoc.TeamID {
			return "", errors.Wrap(app.ErrMalformedWikiDoc, "channel not in given team")
		}
	}

	if wikiDoc.OwnerUserID == "" {
		return "", errors.Wrap(app.ErrMalformedWikiDoc, "missing owner user id of wiki doc")
	}
	if wikiDoc.OwnerUserID != userID {
		return "", errors.Wrap(app.ErrMalformedWikiDoc, "owner user must be the same as the user")
	}

	if strings.TrimSpace(wikiDoc.Name) == "" && wikiDoc.ChannelID == "" {
		return "", errors.Wrap(app.ErrMalformedWikiDoc, "missing name of wiki doc")
	}

	if !app.ValidStatus(wikiDoc.Status) {
		return "", errors.Wrap(app.ErrMalformedWikiDoc, "invalid status provided")
	}

	permission := model.PermissionManagePublicChannelProperties
	permissionMessage := "You are not able to manage public channel properties"
	if channel.Type == model.ChannelTypePrivate {
		permission = model.PermissionManagePrivateChannelProperties
		permissionMessage = "You are not able to manage private channel properties"
	} else if channel.Type == model.ChannelTypeDirect || channel.Type == model.ChannelTypeGroup {
		permission = model.PermissionReadChannel
		permissionMessage = "You do not have access to this channel"
	}

	if !h.pluginAPI.User.HasPermissionToChannel(userID, channel.Id, permission) {
		return "", errors.Wrap(app.ErrNoPermissions, permissionMessage)
	}

	return h.wikiDocService.Create(wikiDoc)
}

func (h *WikiDocHandler) getRequesterInfo(userID string) (app.RequesterInfo, error) {
	return app.GetRequesterInfo(userID, h.pluginAPI)
}

// getWikiDocs handles the GET /runs endpoint.
func (h *WikiDocHandler) getWikiDocs(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")

	filterOptions, err := parseWikiDocsFilterOptions(r.URL, userID)
	if err != nil {
		h.HandleErrorWithCode(w, http.StatusBadRequest, "Bad parameter", err)
		return
	}

	requesterInfo, err := h.getRequesterInfo(userID)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	results, err := h.wikiDocService.GetWikiDocs(requesterInfo, *filterOptions)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	ReturnJSON(w, results, http.StatusOK)
}

// getWikiDoc handles the /doc/{id} endpoint.
func (h *WikiDocHandler) getWikiDoc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wikiDocID := vars["id"]

	playbookRunToGet, err := h.wikiDocService.Get(wikiDocID)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	ReturnJSON(w, playbookRunToGet, http.StatusOK)
}

// updateStatusD handles the POST /doc/{id}/status endpoint, user has edit permissions
func (h *WikiDocHandler) status(w http.ResponseWriter, r *http.Request) {
	wikiDocID := mux.Vars(r)["id"]
	userID := r.Header.Get("Mattermost-User-ID")

	playbookRunToModify, err := h.wikiDocService.Get(wikiDocID)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	isGuest, _ := app.IsGuest(userID, h.pluginAPI)

	if (!app.IsSystemAdmin(userID, h.pluginAPI) && !app.CanPostToChannel(userID, playbookRunToModify.ChannelID, h.pluginAPI)) ||
		isGuest {
		h.HandleErrorWithCode(w, http.StatusForbidden, "Not authorized", fmt.Errorf("user %s cannot post to wiki doc channel %s", userID, playbookRunToModify.ChannelID))
		return
	}

	var options map[string]string

	if err = json.NewDecoder(r.Body).Decode(&options); err != nil {
		h.HandleErrorWithCode(w, http.StatusBadRequest, "unable to decode body into StatusUpdateOptions", err)
		return
	}

	if !app.ValidStatus(options["status"]) {
		h.HandleErrorWithCode(w, http.StatusBadRequest, "invalid status provided", err)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"OK"}`))
}

// parseWikiDocsFilterOptions is only for parsing. Put validation logic in app.validateOptions.
func parseWikiDocsFilterOptions(u *url.URL, currentUserID string) (*app.WikiDocFilterOptions, error) {
	teamId := u.Query().Get("team_id")
	channelId := u.Query().Get("channel_id")

	pageParam := u.Query().Get("page")
	if pageParam == "" {
		pageParam = "0"
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		return nil, errors.Wrapf(err, "bad parameter 'page'")
	}

	perPageParam := u.Query().Get("per_page")
	if perPageParam == "" {
		perPageParam = "0"
	}
	perPage, err := strconv.Atoi(perPageParam)
	if err != nil {
		return nil, errors.Wrapf(err, "bad parameter 'per_page'")
	}

	sort := u.Query().Get("sort")
	direction := u.Query().Get("direction")

	// Parse statuses= query string parameters as an array.
	statuses := u.Query()["statuses"]

	ownerID := u.Query().Get("owner_user_id")

	searchTerm := u.Query().Get("search_term")

	options := app.WikiDocFilterOptions{
		TeamID:     teamId,
		ChannelId:  channelId,
		Page:       page,
		PerPage:    perPage,
		Sort:       app.SortField(sort),
		Direction:  app.SortDirection(direction),
		Statuses:   statuses,
		OwnerID:    ownerID,
		SearchTerm: searchTerm,
	}

	options, err = options.Validate()
	if err != nil {
		return nil, err
	}

	return &options, nil
}
