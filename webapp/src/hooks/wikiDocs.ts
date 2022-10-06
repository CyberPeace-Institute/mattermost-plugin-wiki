import {useEffect, useState} from 'react';
import debounce from 'debounce';
import {useIntl} from 'react-intl';

import {useSelector} from 'react-redux';

import {getCurrentTeamId} from 'mattermost-redux/selectors/entities/teams';
import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/channels';

import {
    fetchWikiDoc as clientFetchWikiDoc,
    fetchWikiDocs as clientFetchWikiDocs,
    saveWikiDoc,
} from 'src/client';

import {FetchWikiDocsParams, WikiDoc} from '../types/wikiDoc';

type ParamsState = Required<Omit<FetchWikiDocsParams, 'team_id' | 'channel_id' | 'owner_user_id' | 'statuses'>>;

const searchDebounceDelayMilliseconds = 300;

export async function getWikiDocOrFetch(id: string, wikiDocs: WikiDoc[] | null) {
    return wikiDocs?.find((p) => p.id === id) ?? clientFetchWikiDoc(id);
}

type EditWikiDocReturn = [WikiDoc | undefined, (update: Partial<WikiDoc>) => void]

export function useEditWikiDoc(id: WikiDoc['id']): EditWikiDocReturn {
    const [wikiDoc, setWikiDoc] = useState<WikiDoc | undefined>();
    useEffect(() => {
        clientFetchWikiDoc(id).then(setWikiDoc);
    }, [id]);

    const updateWikiDoc = (update: Partial<WikiDoc>) => {
        if (wikiDoc) {
            const updatedWikiDoc: WikiDoc = {...wikiDoc, ...update};
            setWikiDoc(updatedWikiDoc);
            saveWikiDoc(updatedWikiDoc);
        }
    };

    return [wikiDoc, updateWikiDoc];
}

export function useWikiDocsCrud(
    defaultParams: Partial<FetchWikiDocsParams>,
    {infinitePaging} = {infinitePaging: false},
) {
    const {formatMessage} = useIntl();
    const teamId = useSelector(getCurrentTeamId);
    const channelId = useSelector(getCurrentChannelId);
    const [wikiDocs, setWikiDocs] = useState<WikiDoc[] | null>(null);
    const [isLoading, setLoading] = useState(true);
    const [hasMore, setHasMore] = useState(false);
    const [totalCount, setTotalCount] = useState(0);
    const [selectedWikiDoc, setSelectedWikiDocState] = useState<WikiDoc | null>();
    const [params, setParamsState] = useState<ParamsState>({
        sort: 'name',
        direction: 'asc',
        page: 0,
        per_page: 10,
        search_term: '',
        ...defaultParams,
    });

    const setParams = (newParams: Partial<ParamsState>) => {
        setParamsState({...params, ...newParams});
    };

    useEffect(() => {
        fetchWikiDocs();
    }, [params, teamId, channelId]);

    const setSelectedWikiDoc = async (nextSelected: WikiDoc | string | null) => {
        if (typeof nextSelected !== 'string') {
            return setSelectedWikiDocState(nextSelected);
        }

        if (!nextSelected) {
            return setSelectedWikiDocState(null);
        }

        return setSelectedWikiDocState(await getWikiDocOrFetch(nextSelected, wikiDocs) ?? null);
    };

    /**
     * Go to specific or next page
     * @param page - defaults to next page if there is one
     */
    const setPage = (page = (hasMore && params.page + 1) || 0) => {
        setParams({page});
    };

    const fetchWikiDocs = async () => {
        setLoading(true);
        const result = await clientFetchWikiDocs(teamId, channelId, params);
        if (result) {
            setWikiDocs(infinitePaging && wikiDocs ? [...wikiDocs, ...result.items] : result.items);
            setTotalCount(result.total_count);
            setHasMore(result.has_more);
        }
        setLoading(false);
    };

    const sortBy = (colName: FetchWikiDocsParams['sort']) => {
        if (params.sort === colName) {
            // we're already sorting on this column; reverse the direction
            const newSortDirection = params.direction === 'asc' ? 'desc' : 'asc';
            setParams({direction: newSortDirection});
            return;
        }

        setParams({sort: colName, direction: 'desc'});
    };

    const setSearchTerm = (term: string) => {
        setLoading(true);
        setParams({search_term: term, page: 0});
    };
    const setSearchTermDebounced = debounce(setSearchTerm, searchDebounceDelayMilliseconds);

    const isFiltering = (params?.search_term?.length ?? 0) > 0;

    return [
        wikiDocs,
        {isLoading, totalCount, hasMore, params, selectedWikiDoc},
        {
            setPage,
            setParams,
            sortBy,
            setSelectedWikiDoc,
            setSearchTerm: setSearchTermDebounced,
            isFiltering,
            fetchWikiDocs,
        },
    ] as const;
}
