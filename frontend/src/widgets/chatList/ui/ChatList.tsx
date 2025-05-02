import axios, { AxiosError } from 'axios'
import { ListWidget } from '../../../shared/ui/ListWidget'
import { ChatCard, ChatCardSkeleton } from '../../../features/chatCard'
import { observer } from 'mobx-react'
import { apiErrors, ApiErrorsKey } from '../../../shared/api'
import { useModal } from '../../../features/modal'
import { useTranslation } from 'react-i18next'
import { useCallback, useEffect, useRef, useState } from 'react'
import chatStore from '../../../entities/Chat/store/ChatStore'
import { fetchPrivateChats, fetchGroupChats } from '../../../features/chatsFetching'

interface ChatListProps {
    open: (chatId: string | number) => void,
    group?: boolean
}

const ChatList: React.FC<ChatListProps> = ({ open, group }) => {
    const { t } = useTranslation('chatListWidget')
    const limit: number = 5
    const [offset, setOffset] = useState<number>(0)
    const [isFetching, setFetching] = useState<boolean>(false)
    const [end, setEnd] = useState<boolean>(false)
    const { openModal, setModalTitle, setModalText } = useModal()

    const fetchFunction = group ? fetchGroupChats : fetchPrivateChats
    const title = group ? t('Groups') : t('People')

    const loadChats = async (newOffset = offset, newEnd = end, newLimit = limit) => {
        if (newEnd || isFetching) return
        setFetching(true)

        const source = axios.CancelToken.source()

        try {
            const response = await fetchFunction(newLimit, newOffset)
            const newChats = response.data.chats

            if (newChats.length < limit) setEnd(true)
            
            setOffset((prevOffset) => prevOffset + limit)
        } catch (e: any) {
            console.log(e)
            if (e instanceof AxiosError && e.status === 404) {
                setEnd(true)
            } else {
                setModalTitle(t('Error'))

                const errorKey: ApiErrorsKey = e.response?.data?.error
                setModalText(t(apiErrors[errorKey]) || t('Internal server error'))
                openModal()
            }
        } finally {
            setFetching(false)
        }

        return () => {
            source.cancel("Operation canceled due to new request")
        }
    }

    useEffect(() => {
        loadChats()
    })

    const observerRef = useRef<IntersectionObserver | null>(null)
    const lastChatRef = useCallback((node: HTMLDivElement | null) => {
        if (isFetching) return
        if (observerRef.current) observerRef.current.disconnect()
        observerRef.current = new IntersectionObserver((entries) => {
            if (entries[0].isIntersecting && !end) {
                loadChats(offset, end)
            }
        })
        if (node) observerRef.current.observe(node)
    }, [isFetching])

    const chats = group ? chatStore.getGroupChats() : chatStore.getPrivateChats()
    
    return (
        <ListWidget className='h-2/5 w-full flex-grow basis-2/5' title={title} >
            {
                chats.map((chat, index) => (
                    <ChatCard
                        chat={chat} 
                        onClick={() => open(chat.chat.id)} 
                        key={chat.chat.id}
                        ref={index === chats.length -1 ? lastChatRef : null}
                    />
                ))
            }
            {
                isFetching && (
                    <>
                        {Array.from({ length: 5 }).map((_, index) => (
                            <ChatCardSkeleton key={index} />
                        ))}
                    </>
                )
            }
            {
                end && !isFetching && chats.length === 0 && (
                    <div
                        className='my-4 2xl:text-xl xl:text-lg lg:text-base 2k:text-2xl 4k:text-3xl
                            md:text-xl sm:text-lg mobile:text-text-base text-center'
                    >
                        { t('It\'s empty here for now') }
                    </div>
                )
            }
        </ListWidget>
    )
}

export default observer(ChatList)