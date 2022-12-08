import React, {ComponentProps, useState} from 'react';

import {useIntl} from 'react-intl';

import styled, {css} from 'styled-components';

import Select from 'react-select';

import MarkdownTextbox from '../markdown/markdown_textbox';
import GenericModal, {InlineLabel} from '../widgets/generic_modal';

const ID = 'wikiDoc_create';

export const makeWikiDocCreateModal = (props: WikiDocCreateModalProps) => ({
    modalId: ID,
    dialogType: WikiDocCreateModal,
    dialogProps: props,
});

export type WikiDocCreateModalProps = {
    createFunc: (name: string, description: string, status: string, content: string) => Promise<void>
} & Partial<ComponentProps<typeof GenericModal>>;

const BaseInput = styled.input`
    transition: border-color ease-in-out .15s, box-shadow ease-in-out .15s, -webkit-box-shadow ease-in-out .15s;
    background-color: rgb(var(--center-channel-bg-rgb));
    border: none;
    box-shadow: inset 0 0 0 1px rgba(var(--center-channel-color-rgb), 0.16);
    border-radius: 4px;
    min-height: 40px;
    line-height: 40px;
    padding: 0 16px;
    font-size: 14px;

    &:focus {
        box-shadow: inset 0 0 0 2px var(--button-bg);
    }
`;

const commonSelectStyle = css`
    flex-grow: 1;
    background-color: var(--center-channel-bg);
`;

const StyledSelect = styled(Select).attrs((props) => {
    return {
        classNamePrefix: 'wikiDocs-rselect',
        ...props,
    };
})`
    ${commonSelectStyle}
`;

const SizedGenericModal = styled(GenericModal)`
    width: calc(80vw);
    max-width: 1200px;
    @media (max-width: 800px) {
      .width: calc(80vw);
    }
`;

const Body = styled.div`
	display: flex;
	flex-direction: column;

	& > div, & > input {
		margin-bottom: 24px;
	}
`;

const WikiDocCreateModal = ({createFunc, ...modalProps}: WikiDocCreateModalProps) => {
    const {formatMessage} = useIntl();
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [status, setStatus] = useState('');
    const [content, setContent] = useState('');

    const create = createFunc;

    const handleStatusSet = (option: {value: string}) => {
        setStatus(option.value);
    };

    const requirementsMet = (name !== '');

    return (
        <SizedGenericModal
            id={ID}
            modalHeaderText={formatMessage({defaultMessage: 'View WikiDoc'})}
            {...modalProps}
            confirmButtonText={formatMessage({defaultMessage: 'Create a doc'})}
            cancelButtonText={formatMessage({defaultMessage: 'Cancel'})}
            isConfirmDisabled={!requirementsMet}
            handleConfirm={() => create(name, description, status, content)}
            showCancel={true}
            autoCloseOnCancelButton={true}
            autoCloseOnConfirmButton={true}
        >
            <Body>
                <InlineLabel>{formatMessage({defaultMessage: 'Wiki name'})}</InlineLabel>
                <BaseInput
                    autoFocus={true}
                    type={'text'}
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                />
                <InlineLabel>{formatMessage({defaultMessage: 'Description'})}</InlineLabel>
                {/*<BaseInput
                    autoFocus={false}
                    type={'text'}
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                />
                <StyledSelect
                    filterOption={null}
                    isMulti={false}
                    placeholder={formatMessage({defaultMessage: 'Set status'})}
                    onChange={handleStatusSet}
                    options={getWikiDocStatuses()}
                    value={getWikiDocStatuses().find((val) => val.value === status)}
                    isClearable={false}
                    maxMenuHeight={380}
                />*/}
                <MarkdownTextbox
                    value={content}
                    setValue={setContent}
                    placeholder={formatMessage({defaultMessage: 'You can add the content of the doc here or edit it later'})}
                />
            </Body>
        </SizedGenericModal>
    );
};
