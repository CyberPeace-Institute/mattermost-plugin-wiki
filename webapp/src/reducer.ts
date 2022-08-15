// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {combineReducers} from 'redux';

import {
    RECEIVED_TOGGLE_RHS_ACTION,
    ReceivedToggleRHSAction,
    SET_RHS_OPEN,
    SetRHSOpen,
    SET_CLIENT_ID,
    SetClientId,
    SET_RHS_STATE,
    RECEIVED_GLOBAL_SETTINGS,
    ShowPostMenuModal,
    HidePostMenuModal,
    SHOW_POST_MENU_MODAL,
    HIDE_POST_MENU_MODAL,
    ShowChannelActionsModal,
    SHOW_CHANNEL_ACTIONS_MODAL,
    HIDE_CHANNEL_ACTIONS_MODAL,
    ShowRunActionsModal,
    HideRunActionsModal,
    SHOW_RUN_ACTIONS_MODAL,
    HIDE_RUN_ACTIONS_MODAL,
    SetHasViewedChannel,
    SET_HAS_VIEWED_CHANNEL,
    SetRHSAboutCollapsedState,
    SET_RHS_ABOUT_COLLAPSED_STATE,
    SetChecklistCollapsedState,
    SetAllChecklistsCollapsedState,
    SET_CHECKLIST_COLLAPSED_STATE,
    SET_ALL_CHECKLISTS_COLLAPSED_STATE,
    SetEveryChecklistCollapsedState,
    SET_EVERY_CHECKLIST_COLLAPSED_STATE,
    SetRHSState,
} from './types/actions';
import {GlobalSettings} from './types/settings';

import {RHSState} from './types/rhs';
import {ReceivedGlobalSettings} from './actions';

function toggleRHSFunction(state = null, action: ReceivedToggleRHSAction) {
    switch (action.type) {
    case RECEIVED_TOGGLE_RHS_ACTION:
        return action.toggleRHSPluginAction;
    default:
        return state;
    }
}

function rhsOpen(state = false, action: SetRHSOpen) {
    switch (action.type) {
    case SET_RHS_OPEN:
        return action.open || false;
    default:
        return state;
    }
}

function clientId(state = '', action: SetClientId) {
    switch (action.type) {
    case SET_CLIENT_ID:
        return action.clientId || '';
    default:
        return state;
    }
}

function rhsState(state = RHSState.ViewingList, action: SetRHSState) {
    switch (action.type) {
    case SET_RHS_STATE:
        return action.nextState;
    default:
        return state;
    }
}

const globalSettings = (state: GlobalSettings | null = null, action: ReceivedGlobalSettings) => {
    switch (action.type) {
    case RECEIVED_GLOBAL_SETTINGS:
        return action.settings;
    default:
        return state;
    }
};

const postMenuModalVisibility = (state = false, action: ShowPostMenuModal | HidePostMenuModal) => {
    switch (action.type) {
    case SHOW_POST_MENU_MODAL:
        return true;
    case HIDE_POST_MENU_MODAL:
        return false;
    default:
        return state;
    }
};

const channelActionsModalVisibility = (state = false, action: ShowChannelActionsModal) => {
    switch (action.type) {
    case SHOW_CHANNEL_ACTIONS_MODAL:
        return true;
    case HIDE_CHANNEL_ACTIONS_MODAL:
        return false;
    default:
        return state;
    }
};

const runActionsModalVisibility = (state = false, action: ShowRunActionsModal | HideRunActionsModal) => {
    switch (action.type) {
    case SHOW_RUN_ACTIONS_MODAL:
        return true;
    case HIDE_RUN_ACTIONS_MODAL:
        return false;
    default:
        return state;
    }
};

const hasViewedByChannel = (state: Record<string, boolean> = {}, action: SetHasViewedChannel) => {
    switch (action.type) {
    case SET_HAS_VIEWED_CHANNEL:
        return {
            ...state,
            [action.channelId]: action.hasViewed,
        };
    default:
        return state;
    }
};

const rhsAboutCollapsedByChannel = (state: Record<string, boolean> = {}, action: SetRHSAboutCollapsedState) => {
    switch (action.type) {
    case SET_RHS_ABOUT_COLLAPSED_STATE:
        return {
            ...state,
            [action.channelId]: action.collapsed,
        };
    default:
        return state;
    }
};

// checklistCollapsedState keeps a map of channelId -> checklist number -> collapsed
const checklistCollapsedState = (
    state: Record<string, Record<number, boolean>> = {},
    action:
    | SetChecklistCollapsedState
    | SetAllChecklistsCollapsedState
    | SetEveryChecklistCollapsedState,
) => {
    switch (action.type) {
    case SET_CHECKLIST_COLLAPSED_STATE: {
        const setAction = action as SetChecklistCollapsedState;
        return {
            ...state,
            [setAction.key]: {
                ...state[setAction.key],
                [setAction.checklistIndex]: setAction.collapsed,
            },
        };
    }
    case SET_ALL_CHECKLISTS_COLLAPSED_STATE: {
        const setAction = action as SetAllChecklistsCollapsedState;
        const newState: Record<number, boolean> = {};
        for (let i = 0; i < setAction.numOfChecklists; i++) {
            newState[i] = setAction.collapsed;
        }
        return {
            ...state,
            [setAction.key]: newState,
        };
    }
    case SET_EVERY_CHECKLIST_COLLAPSED_STATE: {
        const setAction = action as SetEveryChecklistCollapsedState;
        return {
            ...state,
            [setAction.key]: setAction.state,
        };
    }
    default:
        return state;
    }
};

const reducer = combineReducers({
    toggleRHSFunction,
    rhsOpen,
    clientId,
    rhsState,
    globalSettings,
    postMenuModalVisibility,
    channelActionsModalVisibility,
    runActionsModalVisibility,
    hasViewedByChannel,
    rhsAboutCollapsedByChannel,
    checklistCollapsedState,
});

export default reducer;

export type PlaybooksPluginState = ReturnType<typeof reducer>;
