/* eslint-disable react-hooks/exhaustive-deps */
/* eslint-disable react/prop-types */
import { motion } from 'framer-motion'
import './ChatWidget.css'
import { Textarea } from '../../../shared/ui/Textarea'
import { useState, useRef, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { Loader } from '../../../shared/ui/Loader'
import Message from './Message'
import getDateLabel from '../../../utils/dateLabel'
import getTime from '../../../utils/timeLabel'
import MessageSkeleton from './MessageSkeleton'

const MESSAGES_BLOCK_COUNT = 2

const ChatWidget = ({ key, goBack }) => {
    const { t } = useTranslation('chatsFeature')
    const isGroupChat = false
    const [message, setMessage] = useState('')
    const [loading, setLoading] = useState(false)

    const [messages, setMessages] = useState([])
    const chatContainerRef = useRef(null)
    const [hasNewMessages, setHasNewMessages] = useState(false)
    const [loadingMessages, setLoadingMessages] = useState(false)
    const [isLoadingOldMessages, setIsLoadingOldMessages] = useState(false)

    const [scrollBtnShow, setScrollBtnShow] = useState(false)
    const [hasScrolledToUnreadMessage, setHasScrolledToUnreadMessage] = useState(false)

    const [initialize, setInitialize] = useState(false)

    const send = () => {
        if (message.length === 0) {
            return
        }
        console.log(message)
    }

    // mock
    useEffect(() => {
        const initialMessages = [
            { id:1, text: 'Привет', time: new Date(2024, 10, 17, 0, 0, 0), from:'Миша Петров', read: true },
            { id:2, text: 'Как дела?', time: new Date(2024, 10, 17, 0, 0, 0), from:'Миша Петров', read: true },
            { id:3, text: 'Норм', time: new Date(2024, 10, 17, 0, 0, 0), from:'me', read: true },
            { id:4, text: 'Привет', time: new Date(2024, 10, 17, 0, 0, 0), from:'Миша Петров', read: true },
            { id:5, text: 'Как дела?', time: new Date(2024, 10, 17, 0, 0, 0), from:'Миша Петров', read: true },
            { id:6, text: 'Норм', time: new Date(2024, 10, 17, 0, 0, 0), from:'me', read: true },
            { id:7, text: 'Привет', time: new Date(2024, 10, 17, 0, 0, 0), from:'Миша Петров', read: true },
            { id:8, text: 'Как дела?', time: new Date(2024, 10, 17, 0, 0, 0), from:'Миша Петров', read: true },
            { id:9, text: 'Норм', time: new Date(2024, 10, 18, 0, 0, 0), from:'me', read: true },
            { id:10, text: 'Привет', time: new Date(2024, 10, 18, 0, 0, 0), from:'Миша Петров', read: true },
            { id:11, text: 'Как дела?', time: new Date(2024, 10, 18, 0, 0, 0), from:'Миша Петров', read: true},
            { id:12, text: 'Норм', time: new Date(2024, 10, 18, 0, 0, 0), from:'me', read: true },
            { id:13, text: 'Очень длинное сообщение', from:'me', time: new Date(2024, 10, 18, 0, 0, 0), read: true },
            { id:14, text: 'Очень длинное сообщение', from:'Миша Петров', time: new Date(2024, 10, 18, 0, 0, 0), read: false },
            { id:15, text: 'Очень длинное сообщение', from:'Миша Петров', time: new Date(2024, 10, 18, 0, 0, 0), read: false },
            { id:16, text: 'Очень длинное сообщение', from:'Миша Петров', time: new Date(2024, 10, 18, 0, 0, 0), read: false },
            { id:17, text: 'Очень длинное сообщение', from:'Миша Петров', time: new Date(2024, 10, 18, 0, 0, 0), read: false },
            { id:18, text: 'Очень длинное сообщение', from:'Миша Петров', time: new Date(2024, 10, 18, 0, 0, 0), read: false },
            { id:19, text: 'Очень длинное сообщение', from:'Миша Петров', time: new Date(2024, 10, 18, 0, 0, 0), read: false },
            { id:20, text: 'Очень длинное сообщение', from:'Миша Петров', time: new Date(2024, 10, 18, 0, 0, 0), read: false },
        ]
        setMessages(initialMessages)
    }, [])

    const loadMoreMessages = () => {
        if (loadingMessages) return

        setLoadingMessages(true)
        setIsLoadingOldMessages(true)
        setTimeout(() => {
            const newMessages = [
                { id:21, text: 'Очень старое сообщениее', from:'Миша Петров', time: new Date(2023, 10, 18, 0, 0, 0), read: true },
                { id:22, text: 'Очень старое сообщение', from:'Миша Петров', time: new Date(2023, 10, 18, 0, 0, 0), read: true },
            ]
            setMessages((prevMessages) => [...newMessages, ...prevMessages])
            setLoadingMessages(false)
            setIsLoadingOldMessages(false)
        }, 1000)
    }

    const handleScroll = () => {
        const { scrollTop, scrollHeight, clientHeight } = chatContainerRef.current
        if (scrollTop === 0) {
            loadMoreMessages()
        }

        messages.forEach((msg) => {
            const messageElement = document.getElementById(msg.id)
            if (messageElement) {
                const rect = messageElement.getBoundingClientRect()
                if (rect.top >= 0 && rect.bottom <= clientHeight) {
                    if (!msg.read) {
                        setMessages((prevMessages) => prevMessages.map((message) =>
                            message.id === msg.id ? { ...message, read: true } : message
                        ))
                    }
                }
            }
        })
        if (scrollTop + clientHeight > scrollHeight - 300) {
            setScrollBtnShow(false)
        } else {
            setScrollBtnShow(true)
        }

        if (scrollTop + clientHeight === scrollHeight) {
            setHasNewMessages(false)
        }
    }

    const handleScrollButtonClick = () => {
        if (hasScrolledToUnreadMessage) {
            scrollToBottom()
        } else {
            scrollToFirstUnreadMessage()
        }
        setHasScrolledToUnreadMessage(!hasScrolledToUnreadMessage)
    }

    const scrollToBottom = () => {
        if (chatContainerRef.current) {
            chatContainerRef.current.scrollTo({
                top: chatContainerRef.current.scrollHeight,
                behavior: 'smooth',
            })
        }
    }

    const scrollToFirstUnreadMessage = () => {
        const firstUnreadMessageFromOtherUser = messages.find(msg => !msg.read && msg.from !== 'me')
        if (firstUnreadMessageFromOtherUser) {
            const messageElement = document.getElementById(firstUnreadMessageFromOtherUser.id)
            if (messageElement) {
                messageElement.scrollIntoView({ behavior: 'smooth', block: 'center' })
            }
        } else {
            scrollToBottom()
        }
    }

    useEffect(() => {
        if (messages.length > 0 && !initialize) {
            scrollToFirstUnreadMessage()
            setInitialize(true)
        }        
    }, [messages])



    const groupMessagesByDate = (messages) => {
        const groupedMessages = []
    
        let currentGroup = null
    
        messages.forEach(msg => {
            const messageDateLabel = getDateLabel(msg.time)
    
            if (!currentGroup || currentGroup.date !== messageDateLabel) {
                currentGroup = { date: messageDateLabel, messages: [] }
                groupedMessages.push(currentGroup)
            }
    
            currentGroup.messages.push(msg)
        })
    
        return groupedMessages
    }

    return (
        <motion.div 
            key={key}
            initial={{ opacity: 0, x: -500 }}
            animate={{ opacity: 1, x: 0}}
            exit={{ opacity: 0, x: -500 }}
            transition={{ duration: 0.3 }}
            className='Chat-Widget rounded-3xl flex flex-col gap-8 2k:gap-14 4k:gap-24'
        >
            {/* Back btn */}
            <div className='flex items-center gap-5 2k:gap-7 4k:gap-9'>
                <div 
                    className='backBtn 2xl:w-10 2xl:h-10 xl:w-9 xl:h-9 lg:w-9 lg:h-8 2k:w-12 2k:h-12 4k:w-14 4k:h-14 
                        mobile:w-8 mobile:h-8 md:w-9 md:h-9 
                        rounded-3xl flex items-center justify-center'
                    onClick={goBack}
                >
                    <div className='Icon flex items-center justify-center h-1/4 aspect-square'>
                        <svg viewBox="0 0 7.424 12.02" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <defs/>
                            <path id="Vector" d="M0 6.01L6 12.02L7.42 10.6L2.82 6L7.42 1.4L6 0L0 6.01Z" fillOpacity="1.000000" fillRule="nonzero"/>
                        </svg>
                    </div>
                </div>
                <div className='ChatHeader flex items-center h-full gap-5 py-3 px-4 2k:py-5 2k:px-6 2k:h-24 4k:py-7 4k:px-8 4k:h-32 rounded-3xl max-w-[90%]'>
                    {/* Avatar */}
                    <div className='h-full rounded-3xl aspect-square'>
                        <div className='flex items-center justify-center w-full h-full Avatar-standart xl:rounded-3xl lg:rounded-2xl mobile:rounded-2xl md:rounded-3xl'>
                            <div className='flex items-center justify-center w-1/2 aspect-square'>
                                {
                                    isGroupChat ? (
                                        <svg viewBox="0 0 26.6666 24" xmlns="http://www.w3.org/2000/svg">
                                            <defs/>
                                            <path id="coolicon" d="M2.66 6.66C2.66 2.98 5.65 0 9.33 0C13.01 0 16 2.98 16 6.66C16 10.34 13.01 13.33 9.33 13.33C5.65 13.33 2.66 10.34 2.66 6.66ZM9.33 2.66C7.12 2.66 5.33 4.45 5.33 6.66C5.33 8.87 7.12 10.66 9.33 10.66C11.54 10.66 13.33 8.87 13.33 6.66C13.33 4.45 11.54 2.66 9.33 2.66ZM18.66 6.66C19.08 6.66 19.5 6.76 19.87 6.95C20.25 7.14 20.57 7.42 20.82 7.76C21.07 8.1 21.23 8.5 21.3 8.91C21.36 9.33 21.33 9.75 21.2 10.16C21.07 10.56 20.84 10.92 20.54 11.22C20.25 11.51 19.88 11.74 19.48 11.87C19.08 12 18.66 12.03 18.24 11.96C17.82 11.89 17.43 11.73 17.09 11.48L15.52 13.64L15.52 13.64C16.2 14.13 16.99 14.46 17.82 14.59C17.91 14.61 18 14.62 18.09 14.63C18.84 14.71 19.59 14.63 20.3 14.4C21.1 14.14 21.83 13.7 22.43 13.1C23.02 12.51 23.47 11.78 23.73 10.98C23.99 10.18 24.06 9.33 23.93 8.5C23.8 7.67 23.47 6.88 22.98 6.2C22.54 5.59 21.98 5.08 21.33 4.71C21.25 4.66 21.17 4.62 21.08 4.58C20.33 4.19 19.5 4 18.66 4L18.66 6.66ZM16 24L18.66 24C18.66 18.84 14.48 14.66 9.33 14.66C4.17 14.66 0 18.84 0 24L2.66 24C2.66 20.31 5.65 17.33 9.33 17.33C13.01 17.33 16 20.31 16 24ZM23.99 24C23.99 23.3 23.85 22.6 23.59 21.96C23.32 21.31 22.93 20.72 22.43 20.23C21.94 19.73 21.35 19.34 20.7 19.07C20.05 18.8 19.36 18.66 18.66 18.66L18.66 16C19.57 16 20.47 16.15 21.33 16.45C21.46 16.5 21.59 16.55 21.72 16.6C22.69 17.01 23.58 17.6 24.32 18.34C25.06 19.08 25.65 19.96 26.05 20.93C26.11 21.06 26.16 21.2 26.2 21.33C26.51 22.18 26.66 23.09 26.66 24L23.99 24Z" fill="#0F1828" fillOpacity="1.000000" fillRule="evenodd"/>
                                        </svg>
                                    ) : (
                                        <svg viewBox="0 0 16 19" xmlns="http://www.w3.org/2000/svg">
                                            <defs/>
                                            <path id="Vector" d="M8 0C5.23 0 3 2.23 3 5C3 7.76 5.23 10 8 10C10.76 10 13 7.76 13 5C13 2.23 10.76 0 8 0ZM11 5C11 6.65 9.65 8 8 8C6.34 8 5 6.65 5 5C5 3.34 6.34 2 8 2C9.65 2 11 3.34 11 5ZM0 19C0 16.87 0.84 14.84 2.34 13.34C3.84 11.84 5.87 11 8 11C10.12 11 12.15 11.84 13.65 13.34C15.15 14.84 16 16.87 16 19L14 19C14 17.4 13.36 15.88 12.24 14.75C11.11 13.63 9.59 13 8 13C6.4 13 4.88 13.63 3.75 14.75C2.63 15.88 2 17.4 2 19L0 19Z" fill="#0F1828" fillOpacity="1.000000" fillRule="evenodd"/>
                                        </svg>
                                    )
                                }
                            </div>
                        </div>
                    </div>
                    {/* Title */}
                    <h1 className='font-extrabold 2xl:text-3xl xl:text-2xl lg:text-xl 2k:text-4xl 4k:text-5xl
                        md:text-3xl sm:text-2xl mobile:text-xl whitespace-nowrap text-ellipsis overflow-hidden max-w-[80%]'
                    >
                        { 'ChatName' }
                    </h1>
                </div>
            </div>

            {/* Main */}
            <div 
                className='Chat-Main-Area rounded-3xl w-full h-4/5 p-8 flex flex-col gap-5 2k:gap-9 4k:gap-12 overflow-y-scroll relative'
                ref={chatContainerRef}
                onScroll={handleScroll}
            >
                {
                    loading ? (
                        <div className='w-full h-full flex items-center justify-center'>
                            <Loader 
                                className='2xl:w-1/6 xl:w-1/4 lg:w-1/4 md:w-1/5 md:h-5 sm:w-1/4 sm:h-4 mobile:w-1/2 mobile:h-4 lg:h-4 xl:h-5 2xl:h-5 
                                    2k:h-8 4k:h-9' 
                            />
                        </div>
                    ) : (

                        <>
                            {
                                Array(MESSAGES_BLOCK_COUNT).fill(0).map(() => (
                                    <>
                                        <MessageSkeleton 
                                            displayFrom={true}
                                            lines={2}
                                        />
                                        <MessageSkeleton 
                                            displayFrom={true}
                                            lines={3}
                                            me={true}
                                        />
                                        <MessageSkeleton 
                                            displayFrom={true}
                                            lines={4}
                                            me={true}
                                        />
                                        <MessageSkeleton 
                                            displayFrom={true}
                                            lines={1}
                                        />
                                    </>
                                ))
                            }
                            {
                                groupMessagesByDate(messages).map((group, index) => 
                                    <>
                                        <div key={index} className='w-full flex items-center justify-center my-5'>
                                            <div
                                                className='py-2 px-5 2k:py-3 2k:px-6 rounded-3xl DateBlock lg:text-sm 2k:text-base 4k:text-xl
                                                    md:text-sm sm:text-sm mobile:text-xs select-none'
                                            >
                                                { group.date === 'Yesterday' || group.date === 'Today' ? t(group.date) : group.date }
                                            </div>
                                        </div>
                                        {
                                            group.messages.map((msg) => (
                                                <Message
                                                    displayFrom={true}
                                                    key={msg.id}
                                                    text={msg.text}
                                                    time={getTime(msg.time)}
                                                    from={msg.from}
                                                    read={msg.read}
                                                    id={msg.id}
                                                />
                                            ))
                                        }
                                    </>
                                )
                            }
                        </>
                    )
                }
                {
                    scrollBtnShow && (
                        <div 
                            className={`sticky z-50 bottom-0 right-0 Scroll-Btn 
                                2xl:h-14 2xl:w-14 xl:h-12 xl:w-12 mobile:w-10 mobile:h-10 2k:w-16 2k:h-16 4k:w-20 4k:h-20 aspect-square 
                                flex items-center justify-center 2xl:rounded-3xl rounded-2xl self-end`}
                            onClick={handleScrollButtonClick}
                        >
                            <div className='w-2/3 aspect-square'>
                                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" strokeWidth="0"></g><g id="SVGRepo_tracerCarrier" strokeLinecap="round" strokeLinejoin="round"></g><g id="SVGRepo_iconCarrier"> <path d="M5.70711 9.71069C5.31658 10.1012 5.31658 10.7344 5.70711 11.1249L10.5993 16.0123C11.3805 16.7927 12.6463 16.7924 13.4271 16.0117L18.3174 11.1213C18.708 10.7308 18.708 10.0976 18.3174 9.70708C17.9269 9.31655 17.2937 9.31655 16.9032 9.70708L12.7176 13.8927C12.3271 14.2833 11.6939 14.2832 11.3034 13.8927L7.12132 9.71069C6.7308 9.32016 6.09763 9.32016 5.70711 9.71069Z"></path> </g></svg>
                            </div>
                        </div>
                    )
                }
                
            </div>

            <div className='Chat-Control-Area flex items-center justify-between h-[10%] px-3 gap-5'>
                <Textarea 
                    className='w-11/12 px-5 py-3 2k:px-6 2k:py-4 4k:px-7 4k:py-5 
                        rounded-lg md:text-lg mobile:text-sm lg:text-sm 2xl:text-lg 2k:text-2xl 4k:text-4xl 2k:h-24 4k:h-40' 
                    placeholder={t('Enter message')}
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                />

                <div className='Send-Button aspect-square max-h-[80%] flex items-center justify-center rounded-3xl'
                    onClick={send}
                >
                    <svg className='w-1/2 h-1/2 aspect-square' viewBox="0 0 20 15"  xmlns="http://www.w3.org/2000/svg">
                        <defs/>
                        <path id="Vector" d="M17.77 0.21C17.67 0.11 17.54 0.04 17.4 0.01C17.26 -0.02 17.11 -0.01 16.98 0.04L0.48 6.04C0.34 6.09 0.21 6.19 0.13 6.32C0.04 6.44 0 6.59 0 6.74C0 6.89 0.04 7.04 0.13 7.17C0.21 7.29 0.34 7.39 0.48 7.44L6.92 10.02L11.68 5.25L12.73 6.3L7.96 11.08L10.54 17.52C10.59 17.66 10.69 17.78 10.81 17.87C10.94 17.95 11.08 18 11.23 18C11.39 17.99 11.53 17.94 11.66 17.86C11.78 17.77 11.87 17.64 11.92 17.5L17.92 1C17.98 0.87 17.99 0.72 17.96 0.58C17.93 0.44 17.86 0.32 17.77 0.21Z" fillOpacity="1.000000" fillRule="nonzero"/>
                    </svg>
                </div>
            </div>
        </motion.div>
    )
}

export default ChatWidget
