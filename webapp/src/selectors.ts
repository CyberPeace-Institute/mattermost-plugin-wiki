
import {GlobalState} from 'mattermost-redux/types/store';

import {id as pluginId} from './manifest';
import {WikiPluginState} from './reducer';
import {RHSState} from './types/rhs';

const pluginState = (state: GlobalState): WikiPluginState => state['plugins-' + pluginId as keyof GlobalState] as unknown as WikiPluginState || {} as WikiPluginState;

export const currentRHSState = (state: GlobalState): RHSState => pluginState(state).rhsState;

export const selectToggleRHS = (state: GlobalState): () => void => pluginState(state).toggleRHSFunction;
