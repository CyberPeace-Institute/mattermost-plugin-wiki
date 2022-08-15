// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {useDispatch, useSelector} from 'react-redux';
import Scrollbars from 'react-custom-scrollbars';
import styled, {css} from 'styled-components';

import {getCurrentTeam} from 'mattermost-redux/selectors/entities/teams';

import {FormattedMessage} from 'react-intl';

import {
    renderThumbVertical,
    renderTrackHorizontal,
    renderView,
    RHSContainer,
    RHSContent,
} from './rhs_shared';

const WelcomeBlock = styled.div`
    padding: 4rem 3rem 2rem;
    color: rgba(var(--center-channel-color-rgb), 0.72);
`;

const WelcomeDesc = styled.p`
    font-size: 14px;
    line-height: 21px;
    font-weight: 400;
    margin-bottom: 3rem;
`;

const WelcomeCreateAlt = styled.span`
    display: inline-flex;
    align-items: center;
    vertical-align: top;
    padding: 1rem 0;

    > svg {
        margin-left: 0.5em;
    }
`;

const WelcomeWarn = styled(WelcomeDesc)`
    color: rgba(var(--error-text-color-rgb), 0.72);
`;

const RunDetailMaskSvg = 'data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="100%" height="calc(100% - 15px)" viewBox="0 0 400 137" preserveAspectRatio="none"><path d="M0 0H400V122.629C400 122.629 312 137 200 137C101.5 137 0 122.629 0 122.629V0Z"/></svg>';
type RunDetailProps = { exists: boolean; };

const RunDetail = styled.div<RunDetailProps>`
    display: flex;
    place-content: flex-start;
    place-items: center;
    padding: 2rem 2rem 2rem 4rem;
    background:
        linear-gradient(
            180deg,
            rgba(var(--center-channel-bg-rgb), 0.85) 0%,
            rgba(var(--center-channel-bg-rgb), 0.25) 100%
        ),
        rgba(var(${({exists}) => (exists ? '--button-bg-rgb' : '--center-channel-color-rgb')}), 0.08);
    mask-mode: alpha;
    mask-size: cover;
    mask-repeat: round;
    mask-image: url('${RunDetailMaskSvg}');

    > div {
        margin-left: 2rem;
    }
`;

const RunDetailDesc = styled.span<RunDetailProps>`
    font-weight: 400;
    margin-right: auto;
    display: inline-block;
    margin-right: 2rem;
    ${({exists}) => (exists ? css`
        font-size: 14px;
        line-height: 20px;
        color: var(--button-bg);
    ` : css`
        color: '#6F6F73';
        font-size: 16px;
        line-height: 24px;
    `)}
`;

const Header = styled.div`
    min-height: 13rem;
    margin-bottom: 4rem;
    display: grid;
`;

const Heading = styled.h4`
    font-size: 18px;
    line-height: 24px;
    font-weight: 700;
    color: rgba(var(--center-channel-color-rgb), 0.72);
`;

const PaginationContainer = styled.div`
    position: relative;
    height: 0;
    top: -5rem;
    display: flex;
    justify-content: center;
    padding-top: 1rem;

    button {
        height: 3.25rem;
    }
`;

const ListSection = styled.div`
    margin-top: 1rem;
    margin-bottom: 5rem;
    box-shadow: 0px -1px 0px rgba(var(--center-channel-color-rgb), 0.08);
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
    grid-template-rows: repeat(auto-fill, minmax(100px, 1fr));
    position: relative;

    &::after {
        content: '';
        display: block;
        position: absolute;
        width: 100%;
        height: 1px;
        bottom: 0;
        box-shadow: 0px -1px 0px rgba(var(--center-channel-color-rgb), 0.08);
    }
`;

const RHSHome = () => {
    const dispatch = useDispatch();
    const currentTeam = useSelector(getCurrentTeam);

    //const currentRun = useSelector(currentPlaybookRun);
    //const hasCurrentRun = Boolean(currentRun);
    //const [currentPlaybook, setCurrentPlaybook] = useState<Playbook | null>();

    let headerContent;
    if (!headerContent) {
        headerContent = (
            <WelcomeBlock>
                <Heading>
                    <FormattedMessage defaultMessage='Welcome to the WIKI!' />
                </Heading>
                <WelcomeDesc>
                    <FormattedMessage
                        defaultMessage='Here you will see informative pages that can help you better navigate this channel.'
                        values={{br: <br />}}
                    />
                </WelcomeDesc>

                <WelcomeWarn>
                    <FormattedMessage defaultMessage="There are no wiki pages to view. You don't have permission to create wiki pages in this channel." />
                </WelcomeWarn>
            </WelcomeBlock>
        );
    }

    return (
        <RHSContainer>
            <RHSContent>
                <Scrollbars
                    autoHide={true}
                    autoHideTimeout={500}
                    autoHideDuration={500}
                    renderThumbVertical={renderThumbVertical}
                    renderView={renderView}
                    renderTrackHorizontal={renderTrackHorizontal}
                    style={{position: 'absolute'}}
                >
                    {true && <Header>{headerContent}</Header>}
                </Scrollbars>
            </RHSContent>
        </RHSContainer>
    );
};

export default RHSHome;
