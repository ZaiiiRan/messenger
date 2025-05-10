import { useTranslation } from 'react-i18next'
import { IChat } from '../../../entities/Chat'
import { useEffect, useState } from 'react'
import { observer } from 'mobx-react'
import { fetchChat } from '../../chatsFetching'
import { Dialog } from '../../../shared/ui/Dialog'
import ChatRenaming from './ChatRenaming'
import IsFetchingStates from '../models/isFetchingStates'
import ChatManipulation from './ChatManipulation'
import MemberManagement from './MemberManagement'

interface ChatPropertiesProps {
    chat: IChat,
    onDelete: () => void,
    show: boolean,
    setShow: (show: boolean) => void,
}

const ChatProperties: React.FC<ChatPropertiesProps> = ({ chat, onDelete, show, setShow }) => {
    const role = chat.you.role
    const { t } = useTranslation('chatProperties')
    const [isFetching, setIsFetching] = useState<IsFetchingStates>({ rename: false, delete: false, leave: false, return: false, addMembers: false, removeMember: false, changeRole: false })

    const isButtonsDisabled = () => {
        return isFetching.rename || isFetching.delete || isFetching.leave || isFetching.return || isFetching.addMembers || isFetching.removeMember || isFetching.changeRole
    }

    const isMember = () => {
        return role !== 'admin' && role !== 'owner'
    }

    const isAdmin = () => {
        return role === 'admin'
    }

    const isOwner = () => {
        return role === 'owner'
    }

    useEffect(() => {
        let isMounted = true

        const updateChatInfo = async () => {
            if (isMounted)
                await fetchChat(chat.chat.id)
        }
        
        updateChatInfo()

        return () => {
            isMounted = false
        }
    }, [])

    return (
        <Dialog
            show={show}
            setShow={setShow}
            title={chat.chat.name ? chat.chat.name : '???'}
        >
            <ChatRenaming 
                chat={chat}
                isButtonsDisabled={isButtonsDisabled}
                isMember={isMember}
                isFetching={isFetching}
                setIsFetching={setIsFetching}
            />

            <ChatManipulation 
                chat={chat}
                isButtonsDisabled={isButtonsDisabled}
                isMember={isMember}
                isOwner={isOwner}
                isFetching={isFetching}
                setIsFetching={setIsFetching}
                onDelete={onDelete}
            />
            
            <MemberManagement
                chat={chat}
                isButtonsDisabled={isButtonsDisabled}
                isAdmin={isAdmin}
                isMember={isMember}
                isOwner={isOwner}
                isFetching={isFetching}
                setIsFetching={setIsFetching}
            />

        </Dialog>
    )
}

export default observer(ChatProperties)
