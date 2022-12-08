import React from 'react';

import {Store, Action} from 'redux';

import {GlobalState} from '@mattermost/types/store';

import {PluginRegistry} from 'mattermost-webapp/plugins/registry';

import {Client4} from 'mattermost-redux/client';

import manifest from './manifest';

// eslint-disable-next-line import/no-unresolved
import MenuIcon from './components/header/MenuIcon';

import {setSiteUrl} from './client';
import RightHandSidebar from './components/rhs/rhs_main';
import {setToggleRHSAction} from './actions';

type WindowObject = {
    location: {
        origin: string;
        protocol: string;
        hostname: string;
        port: string;
    };
    basename?: string;
}

// From mattermost-webapp/utils
function getSiteURLFromWindowObject(obj: WindowObject): string {
    let siteURL = '';
    if (obj.location.origin) {
        siteURL = obj.location.origin;
    } else {
        siteURL = obj.location.protocol + '//' + obj.location.hostname + (obj.location.port ? ':' + obj.location.port : '');
    }

    if (siteURL[siteURL.length - 1] === '/') {
        siteURL = siteURL.substring(0, siteURL.length - 1);
    }

    if (obj.basename) {
        siteURL += obj.basename;
    }

    if (siteURL[siteURL.length - 1] === '/') {
        siteURL = siteURL.substring(0, siteURL.length - 1);
    }

    return siteURL;
}

function getSiteURL(): string {
    return getSiteURLFromWindowObject(window);
}

export default class Plugin {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    public async initialize(registry: PluginRegistry, store: Store<GlobalState, Action<Record<string, unknown>>>) {
        // Consume the SiteURL so that the client is subpath aware. We also do this for Client4
        // in our version of the mattermost-redux, since webapp only does it in its copy.
        const siteUrl = getSiteURL();
        setSiteUrl(siteUrl);
        Client4.setUrl(siteUrl);

        const RHSWrapped = () => (
            <RightHandSidebar />
        );

        const {toggleRHSPlugin} = registry.registerRightHandSidebarComponent(RHSWrapped, 'Wiki');
        const boundToggleRHSAction = (): void => store.dispatch(toggleRHSPlugin);

        // Store the toggleRHS action to use later
        store.dispatch(setToggleRHSAction(boundToggleRHSAction));

        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
        registry.registerChannelHeaderButtonAction(
            MenuIcon,
            boundToggleRHSAction,
            'Plugin Menu Item',
            'Toggle Wiki',
        );
    }
}

declare global {
    interface Window {
        registerPlugin(id: string, plugin: Plugin): void;
    }
}

window.registerPlugin(manifest.id, new Plugin());
