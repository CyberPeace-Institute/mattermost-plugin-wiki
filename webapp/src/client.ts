// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {AnyAction, Dispatch} from 'redux';

import {GetStateFunc} from 'mattermost-redux/types/actions';
import {Client4} from 'mattermost-redux/client';
import {ClientError} from '@mattermost/client';
import {getCurrentChannel} from 'mattermost-redux/selectors/entities/channels';
import qs from 'qs';

import {id as pluginId} from './manifest';
import {setTriggerId} from './actions';
import {FetchWikiDocsParams, FetchWikiDocsReturn, isWikiDoc, WikiDoc} from './types/wikiDoc';

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

export const doGet = async <TData = any>(url: string) => {
    const {data} = await doFetchWithResponse<TData>(url, {method: 'get'});

    return data;
};

export const doPost = async <TData = any>(url: string, body = {}) => {
    const {data} = await doFetchWithResponse<TData>(url, {
        method: 'POST',
        body,
    });

    return data;
};

export const doDelete = async <TData = any>(url: string, body = {}) => {
    const {data} = await doFetchWithResponse<TData>(url, {
        method: 'DELETE',
        body,
    });

    return data;
};

export const doPut = async <TData = any>(url: string, body = {}) => {
    const {data} = await doFetchWithResponse<TData>(url, {
        method: 'PUT',
        body,
    });

    return data;
};

export const doPatch = async <TData = any>(url: string, body = {}) => {
    const {data} = await doFetchWithResponse<TData>(url, {
        method: 'PATCH',
        body,
    });

    return data;
};

export const doFetchWithResponse = async <TData = any>(url: string, options = {}) => {
    const response = await fetch(url, Client4.getOptions(options));
    let data;
    if (response.ok) {
        const contentType = response.headers.get('content-type');
        if (contentType === 'application/json') {
            data = await response.json() as TData;
        }

        return {
            response,
            data,
        };
    }

    data = await response.text();

    throw new ClientError(Client4.url, {
        message: data || '',
        status_code: response.status,
        url,
    });
};

export const doFetchWithTextResponse = async <TData extends string>(url: string, options = {}) => {
    const response = await fetch(url, Client4.getOptions(options));

    let data;
    if (response.ok) {
        data = await response.text() as TData;

        return {
            response,
            data,
        };
    }

    data = await response.text();

    throw new ClientError(Client4.url, {
        message: data || '',
        status_code: response.status,
        url,
    });
};

export const doFetchWithoutResponse = async (url: string, options = {}) => {
    const response = await fetch(url, Client4.getOptions(options));

    if (response.ok) {
        return;
    }

    throw new ClientError(Client4.url, {
        message: '',
        status_code: response.status,
        url,
    });
};

export async function fetchWikiDocs(teamId: string, channelId: string, params: FetchWikiDocsParams) {
    const queryParams = qs.stringify({...params, team_id: teamId, channel_id: channelId}, {addQueryPrefix: true, indices: false});

    let data = await doGet(`${apiUrl}/wikiDocs${queryParams}`);
    if (!data) {
        data = {items: [], total_count: 0, page_count: 0, has_more: false} as FetchWikiDocsReturn;
    }

    return data as FetchWikiDocsReturn;
}

export async function fetchWikiDoc(id: string) {
    const data = await doGet(`${apiUrl}/wikiDocs/${id}`);
    // eslint-disable-next-line no-process-env
    if (!isWikiDoc(data)) {
        // eslint-disable-next-line no-console
        console.error('expected a WikiDoc in fetchWikiDoc, received:', data);
    }

    return data;
}

export async function createWikiDoc(channel_id: string, user_id: string, team_id: string, name: string, description: string, status: string, content: string) {
    const run = await doPost(`${apiUrl}/wikiDocs/dialog`, JSON.stringify({
        user_id,
        channel_id,
        team_id,
        submission: {
            name,
            description,
            status,
            content,
        },
    }));
    return run as WikiDoc;
}

export async function saveWikiDoc(wikiDoc: WikiDoc) {
    if (!wikiDoc.id) {
        console.error('No wikiDoc id provided');
        return {};
    }

    const wiki = await doPatch(`${apiUrl}/wikiDocs/${wikiDoc.id}`, JSON.stringify(wikiDoc));
    return wiki as WikiDoc;
}

export async function updateWikiDocContent(wikiId: string, content: string) {
    const run = await doPost(`${apiUrl}/wikiDocs/${wikiId}/content`, JSON.stringify({
        content,
    }));
    return run as WikiDoc;
}

export async function updateWikiDocStatus(wikiId: string, status: string) {
    const run = await doPost(`${apiUrl}/wikiDocs/${wikiId}/status`, JSON.stringify({
        status,
    }));
    return run as WikiDoc;
}

export async function deleteWikiDoc(wikiId: string) {
    const run = await doDelete(`${apiUrl}/wikiDocs/${wikiId}`);
    return run as WikiDoc;
}

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
