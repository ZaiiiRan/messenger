import { useTranslation } from 'react-i18next'
import { IChat } from '../../../entities/Chat'
import { Dispatch, SetStateAction, useEffect, useState } from 'react'
import { useModal } from '../../modal'
import { apiErrors, ApiErrorsKey } from '../../../shared/api'
import { fetchFriendsAreNotChatting } from '../../../entities/SocialUser'
import { Button } from '../../../shared/ui/Button'
import { Loader } from '../../../shared/ui/Loader'
import IsFetchingStates from '../models/isFetchingStates'
import { addMembersToChat, deleteChat, fetchChat, leaveFromChat, returnToChat } from '../../chatsFetching'
import { Dialog } from '../../../shared/ui/Dialog'
import { UserSelection } from '../../userSelection'
import { observer } from 'mobx-react'

interface ChatManipulationProps {
    chat: IChat,
    isButtonsDisabled: () => boolean,
    isMember: () => boolean,
    isOwner: () => boolean,
    isFetching: IsFetchingStates,
    setIsFetching: Dispatch<SetStateAction<IsFetchingStates>>,
    onDelete: () => void,
}

const ChatManipulation: React.FC<ChatManipulationProps> = ({ chat, isButtonsDisabled, isMember, isOwner, isFetching, setIsFetching, onDelete }) => {
    const { t } = useTranslation('chatProperties')
    const [addMembersShow, setAddMembersShow] = useState<boolean>(false)
    const [selectedUsers, setSelectedUsers] = useState<(number | string)[]>([])
    const { openModal } = useModal()

    const showErrorModal = (e: any) => {
        const errorKey: ApiErrorsKey = e.response?.data?.error
        const errMsg = t(apiErrors[errorKey]) || t('Internal server error')
        openModal(t('Error'), errMsg)
    }

    const deleteChatAction = async () => {
        try {
            setIsFetching({ ...isFetching, delete: true })
            await deleteChat(chat.chat.id)
            onDelete()
        } catch (e: any) {
            showErrorModal(e)
        } finally {
            setIsFetching({ ...isFetching, delete: false })
        }
    }

    const leaveFromChatAction = async () => {
        try {
            setIsFetching({ ...isFetching, leave: true })
            await leaveFromChat(chat.chat.id)
        } catch (e: any) {
            showErrorModal(e)
        } finally {
            setIsFetching({ ...isFetching, leave: false })
        }
    }

    const returnToChatAction = async () => {
        try {
            setIsFetching({ ...isFetching, return: true })
            await returnToChat(chat.chat.id)
            await fetchChat(chat.chat.id)
        } catch (e: any) {
            if (e.response?.data?.error !== 'messages not found') {
                showErrorModal(e)
            }
        } finally {
            setIsFetching({ ...isFetching, return: false })
        }
    }

    const addMembersAction = async () => {
        if (selectedUsers.length === 0) {
            openModal(t('Error'), t('Please select users'))
            return
        }

        try {
            setIsFetching({ ...isFetching, addMembers: true })
            await addMembersToChat(chat.chat.id, selectedUsers)
            setAddMembersShow(false)
        } catch (e: any) {
            showErrorModal(e)
        } finally {
            setIsFetching({ ...isFetching, addMembers: false })
        }
    }

    const selectUser = (id: number | string) => {
        if (!selectedUsers.includes(id)) {
            setSelectedUsers((prevSelectedUsers) => [...prevSelectedUsers, id])
        } else {
            setSelectedUsers((prevSelectedUsers) => prevSelectedUsers.filter((selectedId) => selectedId !== id))
        }
    }

    const fetchUsers = async (search: string, limit: number, offset: number) => {
        return fetchFriendsAreNotChatting(chat.chat.id, search, limit, offset)
    }

    useEffect(() => {
        setSelectedUsers([])
    }, [addMembersShow])

    return (
        <>
            <div className='flex gap-10 mobile:gap-6 2k:gap-12 4k:gap-14 justify-between w-full flex-wrap'>
                <Button
                    className='h-12 flex items-center justify-center 2k:h-16 4k:h-28 w-72 xl:w-60 lg:w-52 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                        rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                    disabled={!isOwner() || chat.you.isLeft || chat.you.isRemoved || isButtonsDisabled()}
                    onClick={deleteChatAction}
                >
                    {
                        isFetching.delete ? (
                            <Loader className='h-3 w-16 2k:h-4 2k:w-24 4k:h-6 4k:w-36'/>
                        ) : t('Delete chat')
                    }
                </Button>
                {
                    chat.you.isLeft ? (
                        <Button
                            className='h-12 flex items-center justify-center 2k:h-16 4k:h-28 w-72 xl:w-60 lg:w-52 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                            disabled={chat.you.isRemoved || isButtonsDisabled()}
                            onClick={returnToChatAction}
                        >
                            {
                                isFetching.return ? (
                                    <Loader className='h-3 w-16 2k:h-4 2k:w-24 4k:h-6 4k:w-36'/>
                                ) : t('Return')
                            }
                        </Button>
                    ) : (
                        <Button
                            className='h-12 flex items-center justify-center 2k:h-16 4k:h-28 w-72 xl:w-60 lg:w-52 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                            disabled={chat.you.isRemoved || isButtonsDisabled()}
                            onClick={leaveFromChatAction}
                        >
                            {
                                isFetching.leave ? (
                                    <Loader className='h-3 w-16 2k:h-4 2k:w-24 4k:h-6 4k:w-36'/>
                                ) : t('Leave')
                            }
                        </Button>
                    )
                }

                <Button
                    className='h-12 flex items-center justify-center 2k:h-16 4k:h-28 w-1/2 xl:w-60 lg:w-52 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                        rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                    disabled={isButtonsDisabled()}
                    onClick={() => setAddMembersShow(true)}
                >
                    { t('Add members') }
                </Button>
            </div>

            <Dialog
                show={addMembersShow}
                setShow={(show) => setAddMembersShow(show)}
                title={t('Add members')}
                id={'add-members'}
            >
                <UserSelection onSelect={selectUser} fetchFunction={fetchUsers} checkSelected={(id) => selectedUsers.includes(id)} />

                <div className='self-end'>
                    <Button
                        className='h-12 flex items-center justify-center 2k:h-16 4k:h-28 w-72 xl:w-60 lg:w-56 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                            rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                        onClick={addMembersAction}
                        disabled={isFetching.addMembers}
                    >
                        { isFetching.addMembers ? (
                            <Loader className='h-3 w-16 2k:h-4 2k:w-24 4k:h-6 4k:w-36'/>
                        ) : t('Add') }
                    </Button>
                </div>
            </Dialog>
        </>
    )
}

export default observer(ChatManipulation)