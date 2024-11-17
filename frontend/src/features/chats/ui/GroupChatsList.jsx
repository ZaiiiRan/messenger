/* eslint-disable react/prop-types */
import { ListWidget } from '../../../shared/ui/ListWidget'
import { useTranslation } from 'react-i18next'
import ChatCard from './ChatCard'
import ChatCardSkeleton from './ChatCardSkeleton'

const GroupChatsList = ({ open }) => {
    const { t } = useTranslation('chatsFeature')

    return (
        <ListWidget className='h-2/5 w-full flex-grow basis-2/5' title={t('Groups')} >
            <ChatCardSkeleton />
            <ChatCard name={'Test Chat'} type='group' lastMessage={{ text: 'Привет, как дела?', time: new Date(2024, 10, 15, 0, 0, 0, 0), from:'me' }} unreadCount={5} onClick={() => open(1)} />
            <ChatCard name={'Test Test'} type='group' lastMessage={{ text: 'Привет, как дела?', time: new Date(), from: 'Test', read:true }} unreadCount={5} onClick={() => open(1)} />
            <ChatCard name={'Test Test'} type='group' lastMessage={{ text: 'Привет, как дела?', time: new Date(), from:'me', read:true }} unreadCount={5} onClick={() => open(1)} />
        </ListWidget>
    )
}

export default GroupChatsList
