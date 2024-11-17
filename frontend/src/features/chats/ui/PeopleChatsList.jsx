/* eslint-disable react/prop-types */
import { ListWidget } from '../../../shared/ui/ListWidget'
import { useTranslation } from 'react-i18next'
import ChatCard from './ChatCard'

const PeopleChatsList = ({ open }) => {
    const { t } = useTranslation('chatsFeature')

    return (
        <ListWidget className='h-1/2 w-full flex-grow basis-2/5' title={t('People')} >
            <ChatCard name={'Test Test'} lastMessage={{ text: 'Привет 😁😁😁', from:'me', time: new Date() }} unreadCount={5} onClick={() => open(1)} />
            <ChatCard name={'Test Test'} lastMessage={{ text: 'Привет, как дела?', time: new Date(2024, 10, 16, 0, 0, 0, 0) }} unreadCount={1002} onClick={() => open(1)} />
            <ChatCard name={'Test Test'} lastMessage={{ text: 'Привет, как дела?', time: new Date(2011, 0, 1, 2, 3, 4, 567), from: 'me', read:true }} unreadCount={5} onClick={() => open(1)} />
        </ListWidget>
    )
}

export default PeopleChatsList
