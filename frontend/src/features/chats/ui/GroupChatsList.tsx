import { ListWidget } from '../../../shared/ui/ListWidget'
import { useTranslation } from 'react-i18next'
import { ChatCard, ChatCardSkeleton } from '../../../entities/Chat'
import { IChat } from '../../../entities/Chat'

const GroupChatsList = ({ open }) => {
    const { t } = useTranslation('chatsFeature')

    const mock: IChat[] = [
            {
                chat: {
                    name: 'Test Test',
                    id: 1,
                    isGroupChat: true,
                },
                members: [],
                you: {
                    userId: 1,
                    role: 'member',
                    isRemoved: false,
                    isLeft: false,
                    addedBy: 1,
                },
                lastMessage: {
                    id: 1,
                    memberId: 1,
                    chatId: 1,
                    content: 'Привет, как дела?',
                    sentAt: new Date(2025, 3, 30, 0, 0, 0, 0),
                    lastEdited: null
                }
            },
            {
                chat: {
                    name: 'Test Test',
                    id: 1,
                    isGroupChat: true,
                },
                members: [],
                you: {
                    userId: 2,
                    role: 'member',
                    isRemoved: false,
                    isLeft: false,
                    addedBy: 1,
                },
                lastMessage: {
                    id: 2,
                    memberId: 2,
                    chatId: 1,
                    content: 'Привет, как дела?',
                    sentAt: new Date(),
                    lastEdited: null
                }
            },
            {
                chat: {
                    name: 'Test Test',
                    id: 3,
                    isGroupChat: true,
                },
                members: [],
                you: {
                    userId: 1,
                    role: 'member',
                    isRemoved: false,
                    isLeft: false,
                    addedBy: 1,
                },
                lastMessage: {
                    id: 2,
                    memberId: 3,
                    chatId: 1,
                    content: 'Привет, как дела?',
                    sentAt: new Date(),
                    lastEdited: null
                }
            },
            {
                chat: {
                    name: 'Test Test',
                    id: 4,
                    isGroupChat: true,
                },
                members: [],
                you: {
                    userId: 1,
                    role: 'member',
                    isRemoved: false,
                    isLeft: false,
                    addedBy: 1,
                },
                lastMessage: null
            }
        ]

    return (
        <ListWidget className='h-2/5 w-full flex-grow basis-2/5' title={t('Groups')} >
            <ChatCardSkeleton />

            {
                mock.map((chat) => (
                    <ChatCard key={chat.chat.id} chat={chat} onClick={() => open(chat.chat.id)} />
                ))
            }
        </ListWidget>
    )
}

export default GroupChatsList
