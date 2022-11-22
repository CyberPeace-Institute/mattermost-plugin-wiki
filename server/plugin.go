package main

import (
	"github.com/CyberPeace-Institute/mattermost-plugin-wiki/server/app"
	"github.com/pkg/errors"
	"net/http"
	"sync"

	"github.com/mattermost/mattermost-server/v6/plugin"

	"github.com/CyberPeace-Institute/mattermost-plugin-wiki/server/api"
	"github.com/CyberPeace-Institute/mattermost-plugin-wiki/server/bot"
	"github.com/CyberPeace-Institute/mattermost-plugin-wiki/server/sqlstore"

	"github.com/sirupsen/logrus"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-plugin-api/cluster"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	handler         *api.Handler
	wikiDocsService app.WikiDocService
	permissions     *app.PermissionsService

	bot       *bot.Bot
	pluginAPI *pluginapi.Client
}

// ServeHTTP routes incoming HTTP requests to the plugin's REST API.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.handler.ServeHTTP(w, r)
}

// OnActivate Called when this plugin is activated.
func (p *Plugin) OnActivate() error {
	pluginAPIClient := pluginapi.NewClient(p.API, p.Driver)
	p.pluginAPI = pluginAPIClient

	logger := logrus.StandardLogger()
	pluginapi.ConfigureLogrus(logger, pluginAPIClient)

	apiClient := sqlstore.NewClient(pluginAPIClient)
	p.bot = bot.New(pluginAPIClient, "no id")

	sqlStore, err := sqlstore.New(apiClient, p.bot)
	if err != nil {
		return errors.Wrapf(err, "failed creating the SQL store")
	}

	wikiDocStore := sqlstore.NewWikiDocStore(apiClient, p.bot, sqlStore)

	p.wikiDocsService = app.NewWikiDocService(wikiDocStore, p.bot, pluginAPIClient)

	p.permissions = app.NewPermissionsService(p.wikiDocsService, pluginAPIClient)

	mutex, err := cluster.NewMutex(p.API, "CPI_dbMutex")
	if err != nil {
		return errors.Wrapf(err, "failed creating cluster mutex")
	}
	mutex.Lock()
	if err = sqlStore.RunMigrations(); err != nil {
		mutex.Unlock()
		return errors.Wrapf(err, "failed to run migrations")
	}
	mutex.Unlock()

	p.handler = api.NewHandler(pluginAPIClient, p.bot)

	api.NewWikiDocHandler(
		p.handler.APIRouter,
		p.wikiDocsService,
		p.permissions,
		pluginAPIClient,
		p.bot,
	)
	return nil
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
