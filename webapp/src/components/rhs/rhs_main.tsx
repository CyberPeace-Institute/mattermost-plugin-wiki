// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useEffect, useState} from 'react';
import {useDispatch, useSelector} from 'react-redux';

import {GlobalState} from '@mattermost/types/store';
import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/channels';

import {setRHSOpen} from '../../actions';

import RHSHome from './rhs_home';

const RightHandSidebar = () => {
    const dispatch = useDispatch();

    const currentChannelId = useSelector<GlobalState, string>(getCurrentChannelId);

    const [seenChannelId, setSeenChannelId] = useState('');

    useEffect(() => {
        dispatch(setRHSOpen(true));
        return () => {
            dispatch(setRHSOpen(false));
        };
    }, [dispatch]);

    // Update the rhs state when the channel changes
    if (currentChannelId !== seenChannelId) {
        setSeenChannelId(currentChannelId);

        /*if (rhsState) {
            console.log("switching to list view....")
            dispatch(setRHSViewingList());
        }*/
    }

    return <RHSHome />;
};

export default RightHandSidebar;

