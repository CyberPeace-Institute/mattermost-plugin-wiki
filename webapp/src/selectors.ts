
import {GlobalState} from 'mattermost-redux/types/store';

import {id as pluginId} from './manifest';
import {PlaybooksPluginState} from './reducer';
import {RHSState} from './types/rhs';

const pluginState = (state: GlobalState): PlaybooksPluginState => state['plugins-' + pluginId as keyof GlobalState] as unknown as PlaybooksPluginState || {} as PlaybooksPluginState;

export const currentRHSState = (state: GlobalState): RHSState => pluginState(state).rhsState;

export const selectToggleRHS = (state: GlobalState): () => void => pluginState(state).toggleRHSFunction;
