// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {AnyAction, Dispatch} from 'redux';

import {GetStateFunc} from 'mattermost-redux/types/actions';
import {Client4} from 'mattermost-redux/client';
import {getCurrentChannel} from 'mattermost-redux/selectors/entities/channels';

import {id as pluginId} from './manifest';
import {setTriggerId} from './actions';

let siteURL = '';
let basePath = '';
let apiUrl = `${basePath}/plugins/${pluginId}/api/v0`;

export const setSiteUrl = (url?: string): void => {
    if (url) {
        basePath = new URL(url).pathname.replace(/\/+$/, '');
        siteURL = url;
    } else {
        basePath = '';
        siteURL = '';
    }

    apiUrl = `${basePath}/plugins/${pluginId}/api/v0`;
};

export const getSiteUrl = (): string => {
    return siteURL;
};

export const getApiUrl = (): string => {
    return apiUrl;
};

export const playbookExportProps = (playbook: {id: string, title: string}) => {
    const href = `${apiUrl}/playbooks/${playbook.id}/export`;
    const filename = playbook.title.split(/\s+/).join('_').toLowerCase() + '_playbook.json';
    return [href, filename];
};

export async function clientExecuteCommand(dispatch: Dispatch<AnyAction>, getState: GetStateFunc, command: string, teamId: string) {
    let currentChannel = getCurrentChannel(getState());

    // Default to town square if there is no current channel (i.e., if Mattermost has not yet loaded)
    // or in a different team.
    if (!currentChannel || currentChannel.team_id !== teamId) {
        currentChannel = await Client4.getChannelByName(teamId, 'town-square');
    }

    const args = {
        channel_id: currentChannel?.id,
        team_id: teamId,
    };

    try {
        const data = await Client4.executeCommand(command, args);
        dispatch(setTriggerId(data?.trigger_id));
    } catch (error) {
        console.error(error); //eslint-disable-line no-console
    }
}
