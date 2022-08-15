import React, {useRef} from 'react';

import {id as pluginId} from 'src/manifest';

import ReactLogo from './CyberPeace-Logo-White-new.png';

export default function MenuIcon() {
    const myRef = useRef(null);

    return (
        <img
            src={`/static/plugins/${pluginId}/${ReactLogo}`}
            alt='Home'
            height='23px'
            width='23px'
        />
    );
}
