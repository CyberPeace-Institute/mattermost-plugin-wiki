// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import {AnyAction, Dispatch} from 'redux';

import {IntegrationTypes} from 'mattermost-redux/action_types';
import {GetStateFunc} from 'mattermost-redux/types/actions';

import {
    SetRHSOpen,
    SET_RHS_OPEN,
    SetTriggerId,
    SetClientId,
    SET_CLIENT_ID,
    RECEIVED_GLOBAL_SETTINGS,
    SET_RHS_STATE,
    SetRHSState,
    RECEIVED_TOGGLE_RHS_ACTION,
    ReceivedToggleRHSAction, RECEIVED_CHANNEL_WIKI_DOCS, ReceivedChannelWikiDocs,
} from './types/actions';
import {GlobalSettings} from './types/settings';
import {selectToggleRHS} from './selectors';
import {RHSState} from './types/rhs';
import {WikiDoc} from './types/wikiDoc';
import {modals} from './webapp_globals';
import {makeWikiDocCreateModal, WikiDocCreateModalProps} from './components/modals/create_wikidoc_modal';
import {makeWikiDocViewModal, WikiDocViewModalProps} from './components/modals/view_wikidoc_modal';

export interface ReceivedGlobalSettings {
    type: typeof RECEIVED_GLOBAL_SETTINGS;
    settings: GlobalSettings;
}

export function setRHSViewingSingleWikiDoc(): SetRHSState {
    return {
        type: SET_RHS_STATE,
        nextState: RHSState.ViewingSingleWikiDoc,
    };
}

export function setRHSViewingList(): SetRHSState {
    return {
        type: SET_RHS_STATE,
        nextState: RHSState.ViewingList,
    };
}

export function setRHSOpen(open: boolean): SetRHSOpen {
    return {
        type: SET_RHS_OPEN,
        open,
    };
}

export function displayWikiDocCreateModal(props: WikiDocCreateModalProps) {
    return async (dispatch: Dispatch<AnyAction>) => {
        dispatch(modals.openModal(makeWikiDocCreateModal(props)));
    };
}

export function displayWikiDocViewModal(props: WikiDocViewModalProps) {
    return async (dispatch: Dispatch<AnyAction>) => {
        dispatch(modals.openModal(makeWikiDocViewModal(props)));
    };
}

export function setTriggerId(triggerId: string): SetTriggerId {
    return {
        type: IntegrationTypes.RECEIVED_DIALOG_TRIGGER_ID,
        data: triggerId,
    };
}

export function setToggleRHSAction(toggleRHSPluginAction: () => void): ReceivedToggleRHSAction {
    return {
        type: RECEIVED_TOGGLE_RHS_ACTION,
        toggleRHSPluginAction,
    };
}

export const receivedChannelWikiDocs = (wikiDocs: WikiDoc[]): ReceivedChannelWikiDocs => ({
    type: RECEIVED_CHANNEL_WIKI_DOCS,
    wikiDocs,
});

export function toggleRHS() {
    return (dispatch: Dispatch<AnyAction>, getState: GetStateFunc) => {
        selectToggleRHS(getState())();
    };
}

export function setClientId(clientId: string): SetClientId {
    return {
        type: SET_CLIENT_ID,
        clientId,
    };
}
