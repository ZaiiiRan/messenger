import styles from './ChatCard.module.css'
import formatChatTime from '../../../utils/formatChatTime'
import { useTranslation } from 'react-i18next'
import { observer } from 'mobx-react'
import IChat from '../../../entities/Chat/models/IChat'
import { forwardRef, useEffect, useState } from 'react'
import { shortUserStore } from '../../../entities/SocialUser'
import { userStore } from '../../../entities/user'

interface ChatCardProps {
    onClick: () => void,
    chat: IChat,
    key: any,
}

const ChatCard = forwardRef<HTMLDivElement, ChatCardProps>(({ chat, onClick, key }, ref) => {
    const isGroupChat = chat.chat.isGroupChat
    const { t } = useTranslation('chatCard')

    const [senderName, setSenderName] = useState<string | null>(null)
    const [name, setName] = useState<string | null>(null)

    useEffect(() => {
        let isMounted = true

        const loadSenderName = async () => {
            if (chat.lastMessage && chat.lastMessage.memberId !== userStore.user?.userId && isGroupChat) {
                const member = await shortUserStore.get(chat.lastMessage.memberId)
                if (isMounted) {
                    setSenderName(member ? `${member.firstname} ${member.lastname[0]}` : '???')
                }
            } else {
                if (isMounted) {
                    setSenderName('')
                }
            }
        }

        loadSenderName()

        return () => {
            isMounted = false
        }
    }, [chat.lastMessage])

    useEffect(() => {
        let isMounted = true

        const loadPartnerName = async () => {
            if (!isGroupChat) {
                const partner = await shortUserStore.get(chat.members[0].userId)
                if (isMounted) {
                    setName(partner ? `${partner.firstname} ${partner.lastname}` : '???')
                } else {
                    if (isMounted) {
                        setName('')
                    }
                }
            }
        }

        if (chat.chat.name) {
            setName(chat.chat.name)
        } else if (!isGroupChat) {
            loadPartnerName()
        }

        return () => {
            isMounted = false
        }
    }, [])

    return (
        <div 
            className={`${styles.ChatCard} flex items-center px-5 py-2 
                2k:px-8 2k:py-3 4k:px-12 4k:py-4 rounded-3xl xl:gap-5 mobile:gap-4 2k:gap-8 4k:gap-12`}
            onClick={onClick}
            key={key}
            ref={ref}
        >
            {/* Avatar */}
            <div className='md:h-2/3 mobile:h-3/5 rounded-3xl aspect-square'>
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

            {/* Chat Info */}
            <div className='w-4/5 overflow-hidden flex flex-col 4k:gap-2'>
                <div 
                    className='2xl:text-xl xl:text-lg lg:text-base 2k:text-2xl 4k:text-3xl
                        md:text-xl sm:text-lg mobile:text-text-base whitespace-nowrap text-ellipsis overflow-hidden'
                >
                    { name }
                </div>

                {/* Last Message */}
                {
                    chat.lastMessage && (
                        <div
                            className='max-w-[90%] 2xl:text-base xl:text-sm lg:text-sm 2k:text-lg 4k:text-xl
                                md:text-base sm:text-sm mobile:text-sm whitespace-nowrap text-ellipsis overflow-hidden'
                        >
                            {
                                isGroupChat && (chat.lastMessage.memberId !== userStore.user?.userId ) && (
                                    <span className='font-semibold'>{senderName}: </span>
                                )
                            }
                            { chat.lastMessage.content }
                        </div>
                    )
                }
            </div>

            {/* Status */}
            {
                chat.lastMessage && (
                    <div className='flex flex-col items-end justify-between gap-2 2k:gap-3 4k:gap-4 ml-auto'>
                        <div
                            className={`${styles.Time} lg:text-sm 2k:text-base 4k:text-xl
                                md:text-sm sm:text-sm mobile:text-xs`}
                        >
                            {formatChatTime(chat.lastMessage.sentAt) === 'Yesterday' ? t('Yesterday') : formatChatTime(chat.lastMessage.sentAt)}
                        </div>
                        {/* {
                            isMessageFromMe ? (
                                <div className="flex items-center gap-1">
                                    <svg className={`${!read ? styles.checkMark : styles.checkMarkReaded} w-4 h-4 2k:w-6 2k:h-6 4k:w-8 4k:h-8`} xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" strokeWidth="0"></g><g id="SVGRepo_tracerCarrier" strokeLinecap="round" strokeLinejoin="round"></g><g id="SVGRepo_iconCarrier"> <polygon fillRule="evenodd" points="9.707 14.293 19 5 20.414 6.414 9.707 17.121 4 11.414 5.414 10"></polygon> </g></svg>
                                </div>
                            ) : (
                                unreadCount && (
                                    <div 
                                        className={`${styles.Unread} rounded-3xl text-center p-1 2k:p-1.5 4k:p-2 2xl:text-sm lg:text-xs 2k:text-base 4k:text-xl
                                            md:text-sm sm:text-sm mobile:text-xs`}
                                    >
                                        { unreadCount < 1000 ? unreadCount : '999+' }
                                    </div>
                                )
                            )
                        } */}
                    </div>
                )
            }
        </div>
    )
})

export default observer(ChatCard)