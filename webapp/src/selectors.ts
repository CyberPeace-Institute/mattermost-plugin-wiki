
import {GlobalState} from '@mattermost/types/store';

import {haveIChannelPermission} from 'mattermost-redux/selectors/entities/roles';

import {getCurrentChannel} from 'mattermost-redux/selectors/entities/channels';

import Permissions from 'mattermost-redux/constants/permissions';
import {General} from 'mattermost-redux/constants';
import {isCurrentUserSystemAdmin} from 'mattermost-redux/selectors/entities/users';

import {RHSState} from './types/rhs';
import {WikiPluginState} from './reducer';
import {id as pluginId} from './manifest';

const pluginState = (state: GlobalState): WikiPluginState => state['plugins-' + pluginId as keyof GlobalState] as unknown as WikiPluginState || {} as WikiPluginState;

export const currentRHSState = (state: GlobalState): RHSState => pluginState(state).rhsState;

export const isWikiRHSOpen = (state: GlobalState): boolean => pluginState(state).rhsOpen;

export const selectToggleRHS = (state: GlobalState): () => void => pluginState(state).toggleRHSFunction;

export const canUserUpdateWikiDoc = (state: GlobalState) => {
    const channel = getCurrentChannel(state);

    let canManageChannel = false;
    if (channel.type === General.OPEN_CHANNEL) {
        canManageChannel = haveIChannelPermission(state, channel.team_id, channel.id, Permissions.MANAGE_PUBLIC_CHANNEL_PROPERTIES);
    } else if (channel.type === General.PRIVATE_CHANNEL) {
        canManageChannel = haveIChannelPermission(state, channel.team_id, channel.id, Permissions.MANAGE_PRIVATE_CHANNEL_PROPERTIES);
    }

    /*const canManageSystem = haveISystemPermission(state, {
        channel: channelId,
        team: teamId,
        permission: Permissions.MANAGE_SYSTEM,
    });*/

    return canManageChannel || isCurrentUserSystemAdmin(state);
};
