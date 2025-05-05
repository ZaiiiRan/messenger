import axios, { AxiosError } from 'axios'
import { ListWidget } from '../../../shared/ui/ListWidget'
import { ChatCard, ChatCardSkeleton } from '../../../features/chatCard'
import { observer } from 'mobx-react'
import { apiErrors, ApiErrorsKey } from '../../../shared/api'
import { useModal } from '../../../features/modal'
import { useTranslation } from 'react-i18next'
import { Dispatch, SetStateAction, useCallback, useEffect, useRef, useState } from 'react'
import chatStore from '../../../entities/Chat/store/ChatStore'
import { fetchPrivateChats, fetchGroupChats } from '../../../features/chatsFetching'

interface ChatListProps {
    open: (chatId: string | number) => void,
    group?: boolean,
    setNewChatModal?: Dispatch<SetStateAction<boolean>>
}

const ChatList: React.FC<ChatListProps> = ({ open, group, setNewChatModal }) => {
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
    }, [])

    useEffect(() => {
        const chatCount = group ? chatStore.getGroupChatCount() : chatStore.getPrivateChatCount()
        if (chatCount > offset) {
            setOffset(chatCount)
        }
    }, [chatStore])

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

    const button = group && setNewChatModal ? (
        <div 
            className='backBtn 2xl:w-10 2xl:h-10 xl:w-9 xl:h-9 lg:w-9 lg:h-8 2k:w-12 2k:h-12 4k:w-14 4k:h-14 
                mobile:w-8 mobile:h-8 md:w-9 md:h-9 
                rounded-3xl flex items-center justify-center'
            onClick={() => setNewChatModal(true)}
        >
            <div className='Icon flex items-center justify-center h-1/2 aspect-square'>
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" strokeWidth="0"></g><g id="SVGRepo_tracerCarrier" strokeLinecap="round" strokeLinejoin="round"></g><g id="SVGRepo_iconCarrier"><path fillRule="evenodd" clipRule="evenodd" d="M19.186 2.09c.521.25 1.136.612 1.625 1.101.49.49.852 1.104 1.1 1.625.313.654.11 1.408-.401 1.92l-7.214 7.213c-.31.31-.688.541-1.105.675l-4.222 1.353a.75.75 0 0 1-.943-.944l1.353-4.221a2.75 2.75 0 0 1 .674-1.105l7.214-7.214c.512-.512 1.266-.714 1.92-.402zm.211 2.516a3.608 3.608 0 0 0-.828-.586l-6.994 6.994a1.002 1.002 0 0 0-.178.241L9.9 14.102l2.846-1.496c.09-.047.171-.107.242-.178l6.994-6.994a3.61 3.61 0 0 0-.586-.828zM4.999 5.5A.5.5 0 0 1 5.47 5l5.53.005a1 1 0 0 0 0-2L5.5 3A2.5 2.5 0 0 0 3 5.5v12.577c0 .76.082 1.185.319 1.627.224.419.558.754.977.978.442.236.866.318 1.627.318h12.154c.76 0 1.185-.082 1.627-.318.42-.224.754-.559.978-.978.236-.442.318-.866.318-1.627V13a1 1 0 1 0-2 0v5.077c0 .459-.021.571-.082.684a.364.364 0 0 1-.157.157c-.113.06-.225.082-.684.082H5.923c-.459 0-.57-.022-.684-.082a.363.363 0 0 1-.157-.157c-.06-.113-.082-.225-.082-.684V5.5z"></path></g></svg>
            </div>
        </div>
    ) : <></>
    
    return (
        <>
            <ListWidget className='h-2/5 w-full flex-grow basis-2/5' title={title} button={button}>
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
        </>
        
    )
}

export default observer(ChatList)