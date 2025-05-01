import { ListWidget } from '../../../shared/ui/ListWidget'
import { useTranslation } from 'react-i18next'
import { ChatCard, ChatCardSkeleton } from '../../../entities/Chat'
import { IChat } from '../../../entities/Chat'

const PeopleChatsList = ({ open }) => {
    const { t } = useTranslation('chatsFeature')

    const mock: IChat[] = [
        {
            chat: {
                name: 'Test Test',
                id: 1,
                isGroupChat: false,
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
                id: 4,
                memberId: 1,
                chatId: 1,
                content: 'Привет 😁😁😁',
                sentAt: new Date(),
                lastEdited: null
            }
        },
        {
            chat: {
                name: 'Test Test',
                id: 1,
                isGroupChat: false,
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
                id: 5,
                memberId: 2,
                chatId: 1,
                content: 'Привет, как дела?',
                sentAt: new Date(2024, 10, 16, 0, 0, 0, 0),
                lastEdited: null
            }
        },
        {
            chat: {
                name: 'Test Test',
                id: 3,
                isGroupChat: false,
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
                id: 6,
                memberId: 1,
                chatId: 1,
                content: 'Привет, как дела?',
                sentAt: new Date(2011, 0, 1, 2, 3, 4, 567),
                lastEdited: null
            }
        }
    ]

    return (
        <ListWidget className='h-1/2 w-full flex-grow basis-2/5' title={t('People')} >
            {
                mock.map((chat) => (
                    <ChatCard key={chat.chat.id}  chat={chat} onClick={() => open(chat.chat.id)} />
                ))
            }
        </ListWidget>
    )
}

export default PeopleChatsList
