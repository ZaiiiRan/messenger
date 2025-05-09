import { observer } from 'mobx-react'
import { chatStore, IChat, IMessage } from '../../../entities/Chat'
import Message from './Message'
import MessageSkeleton from './MessageSkeleton'
import getDateLabel from '../../../utils/dateLabel'
import { useTranslation } from 'react-i18next'
import './MessageList.css'
import { useCallback, useEffect, useRef, useState } from 'react'
import fetchMessages from '../api/fetchMessages'
import { AxiosError } from 'axios'
import { useModal } from '../../../features/modal'
import { apiErrors, ApiErrorsKey } from '../../../shared/api'

interface MessageListProps {
    chat?: IChat
}

interface IMessageGroup {
    date: string,
    messages: IMessage[]
}

const MESSAGES_BLOCK_COUNT = 2
const LIMIT = 20

const MessageList: React.FC<MessageListProps> = ({ chat }) => {
    const { t } = useTranslation('messageList')
    const containerRef = useRef<HTMLDivElement | null>(null)
    const skeletonRef = useRef<HTMLDivElement | null>(null)
    const [end, setEnd] = useState<boolean>(false)
    const [loading, setLoading] = useState<boolean>(false)
    const { openModal } = useModal()

    useEffect(() => {
        const container = containerRef.current
        if (!container) return
    
        container.scrollTop = container.scrollHeight
    }, [chat?.chat.id])

    const loadMoreMessages = useCallback(async () => {
        if (!chat || loading || end) return

        const container = containerRef.current
        if (!container) return

        const scrollHeightBefore = container.scrollHeight
        const scrollTopBefore = container.scrollTop

        setLoading(true)
        try {
            const newMessages = await fetchMessages(chat.chat.id, LIMIT, chat.messages.length)
            if (newMessages.length < LIMIT) {
                setEnd(true)
            }

            requestAnimationFrame(() => {
                const scrollHeightAfter = container.scrollHeight
                const heightDiff = scrollHeightAfter - scrollHeightBefore
                container.scrollTop = scrollTopBefore + heightDiff
            })
        } catch (e: any) {
            if (e instanceof AxiosError && e.status === 404) {
                setEnd(true)
            } else {
                const errorKey: ApiErrorsKey = e.response?.data?.error
                const errMsg = t(apiErrors[errorKey]) || t('Internal server error')
                openModal(t('Error'), errMsg)
            }
        } finally {
            setLoading(false)
        }
    }, [chat, loading, end])

    useEffect(() => {
        if (!skeletonRef.current || end) return

        const observer = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting && !loading) {
                    loadMoreMessages()
                }
            },
            { root: containerRef.current, threshold: 0.1 }
        )

        observer.observe(skeletonRef.current)

        return () => {
            if (skeletonRef.current) {
                observer.unobserve(skeletonRef.current)
            }
        }
    }, [loadMoreMessages, end])

    useEffect(() => {
        const container = containerRef.current
        if (!container) return

        const isNearBottom = container.scrollHeight - (container.scrollTop + container.clientHeight) <= 300

        if (isNearBottom && !loading) {
            container.scrollTo({
                top: container.scrollHeight,
                behavior: 'smooth'
            })
        }
    }, [chat?.messages.length, chat?.messages])

    const groupMessagesByDate = () => {
        const chatMessages = chat?.messages || []
        const messages = chatMessages.toSorted((a, b) => new Date(a.sentAt).getTime() - new Date(b.sentAt).getTime())

        const groupedMessages: IMessageGroup[] = []

        let currentGroup: IMessageGroup | null = null

        messages.forEach((msg: IMessage) => {
            const messageDateLabel = getDateLabel(msg.sentAt)

            if (!currentGroup || currentGroup.date !== messageDateLabel) {
                currentGroup = { date: messageDateLabel, messages: [] }
                groupedMessages.push(currentGroup)
            }

            currentGroup.messages.push(msg)
        })

        return groupedMessages
    }

    return (
        <div 
            className='Chat-Main-Area rounded-3xl w-full h-4/5 p-8 flex flex-col gap-5 2k:gap-9 4k:gap-12 overflow-y-scroll relative'
            ref={containerRef}
        >
            {
                !end && (
                    <div 
                        ref={skeletonRef}
                        className='Chat-Main-Area w-full flex flex-col gap-5 2k:gap-9 4k:gap-12'
                    >
                        { Array(MESSAGES_BLOCK_COUNT).fill(0).map((_, index) => (
                            <div 
                                key={index}
                                className='Chat-Main-Area w-full flex flex-col gap-5 2k:gap-9 4k:gap-12'
                            >
                                <MessageSkeleton
                                    displayFrom={chat ? chat.chat.isGroupChat : false}
                                    lines={2}
                                    key={index + 10}
                                />
                                <MessageSkeleton
                                    displayFrom={chat ? chat.chat.isGroupChat : false}
                                    lines={3}
                                    me
                                    key={index + 20}
                                />
                                <MessageSkeleton
                                    displayFrom={chat ? chat.chat.isGroupChat : false}
                                    lines={4}
                                    me
                                    key={index + 30}
                                />
                                <MessageSkeleton
                                    displayFrom={chat ? chat.chat.isGroupChat : false}
                                    lines={1}
                                    key={index + 40}
                                />
                            </div>
                        ))}
                    </div>
            )}
            {
                groupMessagesByDate().map((group: IMessageGroup, index: number) => (
                    <div key={index} className='Chat-Main-Area w-full flex flex-col gap-5 2k:gap-9 4k:gap-12'>
                        <div key={index} className='w-full flex items-center justify-center my-5'>
                            <div
                                className='py-2 px-5 2k:py-3 2k:px-6 rounded-3xl DateBlock lg:text-sm 2k:text-base 4k:text-xl
                                    md:text-sm sm:text-sm mobile:text-xs select-none'
                            >
                                { group.date === 'Yesterday' || group.date === 'Today' ? t(group.date) : group.date }
                            </div>
                        </div>
                        {
                            group.messages.map((msg: IMessage, i: number) => (
                                <Message
                                    isGroupChat={chat ? chat.chat.isGroupChat : false}
                                    key={msg.id}
                                    message={msg}
                                    id={msg.id}
                                />
                            ))
                        }
                    </div>
                ))
            }
        </div>
    )
}

export default observer(MessageList)
