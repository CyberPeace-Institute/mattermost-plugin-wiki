// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

export interface WikiDoc {
    id: string;
    name: string;
    content: string;
    description: string;
    status: string;
    owner_user_id: string;
    team_id: string;
    channel_id: string;
    create_at: number;
    update_at: number;
    delete_at: number;
}

export interface FetchWikiDocsReturn {
    total_count: number;
    page_count: number;
    has_more: boolean;
    items: WikiDoc[];
}

export enum WikiDocStatus {
    Private = 'Private',
    Published = 'Published',
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function isWikiDoc(arg: any): arg is WikiDoc {
    return Boolean(arg &&
        arg.id && typeof arg.id === 'string' &&
        arg.name && typeof arg.name === 'string' &&
        arg.content && typeof arg.content === 'string' &&
        typeof arg.description === 'string' &&
        typeof arg.status === 'string' &&
        arg.owner_user_id && typeof arg.owner_user_id === 'string' &&
        arg.team_id && typeof arg.team_id === 'string' &&
        arg.channel_id && typeof arg.channel_id === 'string' &&
        typeof arg.create_at === 'number' &&
        typeof arg.update_at === 'number' &&
        typeof arg.delete_at === 'number');
}

function isString(arg: any): arg is string {
    return Boolean(typeof arg === 'string');
}

export function wikiDocIsPublished(wikiDocRun: WikiDoc): boolean {
    return wikiDocRun.status === WikiDocStatus.Published;
}

export interface FetchWikiDocsParams {
    page: number;
    per_page: number;
    team_id?: string;
    channel_id?: string;
    sort?: string;
    direction?: string;
    statuses?: string[];
    owner_user_id?: string;
    search_term?: string;
}
