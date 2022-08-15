// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import {AnyAction, Dispatch} from 'redux';

import {generateId} from 'mattermost-redux/utils/helpers';
import {IntegrationTypes} from 'mattermost-redux/action_types';
import {getCurrentTeamId} from 'mattermost-redux/selectors/entities/teams';
import {addChannelMember} from 'mattermost-redux/actions/channels';
import {DispatchFunc, GetStateFunc} from 'mattermost-redux/types/actions';

import {getCurrentChannelId} from 'mattermost-webapp/packages/mattermost-redux/src/selectors/entities/common';

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
    ReceivedToggleRHSAction,
} from './types/actions';
import {clientExecuteCommand} from './client';
import {GlobalSettings} from './types/settings';
import {selectToggleRHS} from './selectors';
import {RHSState} from './types/rhs';

export function startPlaybookRunById(teamId: string, playbookId: string, timeout = 0) {
    return async (dispatch: Dispatch<AnyAction>, getState: GetStateFunc) => {
        // Add unique id
        const clientId = generateId();
        dispatch(setClientId(clientId));

        const command = `/playbook run-playbook ${playbookId} ${clientId}`;

        // When dispatching from the playbooks product, the switch to channels resets the websocket
        // connection, losing the event that opens this dialog. Allow the caller to specify a
        // timeout as a gross workaround.
        await new Promise((resolve) => setTimeout(() => {
            clientExecuteCommand(dispatch, getState, command, teamId);
            // eslint-disable-next-line no-undefined
            resolve(undefined);
        }, timeout));
    };
}

export function finishRun(teamId: string) {
    return async (dispatch: Dispatch, getState: GetStateFunc) => {
        await clientExecuteCommand(dispatch, getState, '/playbook finish', teamId);
    };
}

export function addToTimeline(postId: string) {
    return async (dispatch: Dispatch, getState: GetStateFunc) => {
        const currentTeamId = getCurrentTeamId(getState());

        await clientExecuteCommand(dispatch, getState, `/playbook add ${postId}`, currentTeamId);
    };
}

export function addNewTask(checklist: number) {
    return async (dispatch: Dispatch<AnyAction>, getState: GetStateFunc) => {
        const currentTeamId = getCurrentTeamId(getState());

        await clientExecuteCommand(dispatch, getState, `/playbook checkadd ${checklist}`, currentTeamId);
    };
}

export function addToCurrentChannel(userId: string) {
    return async (dispatch: DispatchFunc, getState: GetStateFunc) => {
        const currentChannelId = getCurrentChannelId(getState());

        dispatch(addChannelMember(currentChannelId, userId));
    };
}

export interface ReceivedGlobalSettings {
    type: typeof RECEIVED_GLOBAL_SETTINGS;
    settings: GlobalSettings;
}

export function setRHSViewingWelcome(): SetRHSState {
    return {
        type: SET_RHS_STATE,
        nextState: RHSState.ViewingWelcome,
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
