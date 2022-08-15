// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useEffect, useState} from 'react';
import {useDispatch, useSelector} from 'react-redux';

import {GlobalState} from 'mattermost-redux/types/store';
import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/channels';

import {setRHSOpen, setRHSViewingList} from '../../actions';

import {RHSState} from '../../types/rhs';

import {currentRHSState} from '../../selectors';

import RHSHome from './rhs_home';

const RightHandSidebar = () => {
    const dispatch = useDispatch();
    const currentChannelId = useSelector<GlobalState, string>(getCurrentChannelId);

    //const inPlaybookRun = useSelector<GlobalState, boolean>(inPlaybookRunChannel);
    const rhsState = useSelector<GlobalState, RHSState>(currentRHSState);
    const [seenChannelId, setSeenChannelId] = useState('');

    useEffect(() => {
        dispatch(setRHSOpen(true));
        return () => {
            dispatch(setRHSOpen(false));
        };
    }, [dispatch]);

    useEffect(() => {
        console.log('Got channel with id: ' + currentChannelId);
    }, [currentChannelId]);

    // Update the rhs state when the channel changes
    if (currentChannelId !== seenChannelId) {
        setSeenChannelId(currentChannelId);

        if (rhsState) {
            dispatch(setRHSViewingList());
        }
    }

    return <RHSHome />;
};

export default RightHandSidebar;

