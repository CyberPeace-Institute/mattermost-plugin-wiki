// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import Integrations from 'mattermost-redux/action_types/integrations';

import {id as pluginId} from '../manifest';

import {RHSState} from './rhs';
import {WikiDoc} from './wikiDoc';

export const RECEIVED_TOGGLE_RHS_ACTION = pluginId + '_toggle_rhs';
export const SET_RHS_OPEN = pluginId + '_set_rhs_open';
export const SET_CLIENT_ID = pluginId + '_set_client_id';
export const WIKI_DOC_CREATED = pluginId + '_wiki_doc_created';
export const RECEIVED_CHANNEL_WIKI_DOCS = pluginId + '_received_channel_wiki_docs';
export const REMOVED_FROM_CHANNEL = pluginId + '_removed_from_wiki_doc_run_channel';
export const SET_RHS_STATE = pluginId + '_set_rhs_state';
export const RECEIVED_GLOBAL_SETTINGS = pluginId + '_received_global_settings';
export const SHOW_POST_MENU_MODAL = pluginId + '_show_post_menu_modal';
export const HIDE_POST_MENU_MODAL = pluginId + '_hide_post_menu_modal';
export const SHOW_CHANNEL_ACTIONS_MODAL = pluginId + '_show_channel_actions_modal';
export const HIDE_CHANNEL_ACTIONS_MODAL = pluginId + '_hide_channel_actions_modal';
export const SET_HAS_VIEWED_CHANNEL = pluginId + '_set_has_viewed';

export interface ReceivedToggleRHSAction {
    type: typeof RECEIVED_TOGGLE_RHS_ACTION;
    toggleRHSPluginAction: () => void;
}

export interface SetRHSState {
    type: typeof SET_RHS_STATE;
    nextState: RHSState;
}

export interface SetRHSOpen {
    type: typeof SET_RHS_OPEN;
    open: boolean;
}

export interface SetTriggerId {
    type: typeof Integrations.RECEIVED_DIALOG_TRIGGER_ID;
    data: string;
}

export interface SetClientId {
    type: typeof SET_CLIENT_ID;
    clientId: string;
}

export interface ReceivedChannelWikiDocs {
    type: typeof RECEIVED_CHANNEL_WIKI_DOCS;
    wikiDocs: WikiDoc[];
}

export interface WikiDocCreated {
    type: typeof WIKI_DOC_CREATED;
    teamID: string;
}

export interface RemovedFromChannel {
    type: typeof REMOVED_FROM_CHANNEL;
    channelId: string;
}

export interface ShowPostMenuModal {
    type: typeof SHOW_POST_MENU_MODAL;
}

export interface HidePostMenuModal {
    type: typeof HIDE_POST_MENU_MODAL;
}

export interface ShowChannelActionsModal {
    type: typeof SHOW_CHANNEL_ACTIONS_MODAL;
}

export interface HideChannelActionsModal {
    type: typeof HIDE_CHANNEL_ACTIONS_MODAL;
}

export interface SetHasViewedChannel {
    type: typeof SET_HAS_VIEWED_CHANNEL;
    channelId: string;
    hasViewed: boolean;
}
