import React, { forwardRef, useEffect, useState } from 'react'
import { IMessage } from '../../../entities/Chat'
import styles from './Message.module.css'
import { userStore } from '../../../entities/user'
import { observer } from 'mobx-react'
import { shortUserStore } from '../../../entities/SocialUser'

interface MessageProps {
    id: string | number,
    isGroupChat: boolean,
    className?: string,
    message: IMessage
}

const Message: React.FC<MessageProps> = ({ className, message, isGroupChat = false, id }) => {
    const isMessageFromMe = message.memberId === userStore.user?.userId
    const [senderName, setSenderName] = useState<string | null>('')

    useEffect(() => {
        let isMounted = true

        const loadSenderName = async () => {
            if (!isMessageFromMe) {
                const member = await shortUserStore.get(message.memberId)
                if (isMounted) {
                    setSenderName(member ? `${member.firstname} ${member.lastname}` : '???')
                }
            } else {
                if (isMounted) {
                    setSenderName('')
                }
            }
        }

        if (isGroupChat) {
            loadSenderName()
        }

        return () => {
            isMounted = false
        }
    }, [message])
    
    return (
        <div 
            className={`${isMessageFromMe ? 'self-end' : 'self-start'} flex-shrink-0 2xl:max-w-[75%] lg:max-w-[90%] md:max-w-[70%] sm:max-w-[90%] 
                flex items-start md:gap-4 mobile:gap-2 2k:gap-5 4k:gap-6`}
            id={id.toString()}
        >
            {
                    isGroupChat && !isMessageFromMe && (
                        <div 
                            className='md:w-[40px] md:min-w-[40px] sm:w-[35px] sm:min-w-[35px] mobile:w-[30px] mobile:min-w-[30px] 2xl:w-[50px] 2xl:min-w-[50px] 2k:w-[60px] 2k:min-w-[60px] 4k:w-[80px] 4k:min-w-[80px]
                                md:rounded-2xl 2k:rounded-3xl aspect-square cursor-pointer self-end'
                        >
                            <div className='flex items-center justify-center w-full h-full Avatar-standart mobile:rounded-2xl'>
                                <div className='flex items-center justify-center w-1/2 aspect-square'>
                                    <svg viewBox="0 0 16 19" xmlns="http://www.w3.org/2000/svg">
                                        <defs/>
                                        <path id="Vector" d="M8 0C5.23 0 3 2.23 3 5C3 7.76 5.23 10 8 10C10.76 10 13 7.76 13 5C13 2.23 10.76 0 8 0ZM11 5C11 6.65 9.65 8 8 8C6.34 8 5 6.65 5 5C5 3.34 6.34 2 8 2C9.65 2 11 3.34 11 5ZM0 19C0 16.87 0.84 14.84 2.34 13.34C3.84 11.84 5.87 11 8 11C10.12 11 12.15 11.84 13.65 13.34C15.15 14.84 16 16.87 16 19L14 19C14 17.4 13.36 15.88 12.24 14.75C11.11 13.63 9.59 13 8 13C6.4 13 4.88 13.63 3.75 14.75C2.63 15.88 2 17.4 2 19L0 19Z" fill="#0F1828" fillOpacity="1.000000" fillRule="evenodd"/>
                                    </svg>
                                </div>
                            </div>
                        </div>
                    )
                }
            <div 
                className={`${styles.Message} ${isMessageFromMe ? styles.self : ''} ${className} 
                    ${isMessageFromMe ? 'rounded-bl-3xl' : 'rounded-br-3xl'} rounded-tl-3xl rounded-tr-3xl
                    mobile:px-4 mobile:py-3 md:px-5 md:py-4 2k:px-6 2k:py-5 4k:px-7 4k:py-6
                    break-words flex flex-col md:gap-2 mobile:gap-1 2k:gap-3 4k:gap-4 lg:min-w-[200px] mobile:min-w-[150px]`}
            >
                <div className='flex flex-col md:gap-2 mobile:gap-1 2k:gap-3 4k:gap-4'>
                    {
                        isGroupChat && !isMessageFromMe && (
                            <div className={`${styles.from} self-start lg:text-sm 2k:text-base 4k:text-xl
                                md:text-sm sm:text-sm mobile:text-xs cursor-pointer max-w-[70%] w-auto inline-block whitespace-nowrap text-ellipsis overflow-hidden`}>
                                { senderName }
                            </div>
                        )
                    }

                    <div className='mobile:text-sm md:text-base 2k:text-xl 4k:text-2xl whitespace-pre-wrap'>{ message.content }</div>
                </div>
                <div 
                    className='self-end lg:text-sm 2k:text-base 4k:text-xl
                        md:text-sm sm:text-sm mobile:text-xs flex gap-2 items-center select-none'
                >
                    <div className={`${styles.Time} flex-shrink-0`} >{ new Date(message.sentAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: false }) }</div>
                    {/* {
                        isMessageFromMe && (
                            <div className="flex items-center">
                                <svg className={`${!read ? styles.checkMark : styles.checkMarkReaded} w-4 h-4 2k:w-6 2k:h-6 4k:w-8 4k:h-8`} xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" strokeWidth="0"></g><g id="SVGRepo_tracerCarrier" strokeLinecap="round" strokeLinejoin="round"></g><g id="SVGRepo_iconCarrier"> <polygon fillRule="evenodd" points="9.707 14.293 19 5 20.414 6.414 9.707 17.121 4 11.414 5.414 10"></polygon> </g></svg>
                            </div>
                        )
                    } */}
                </div>
            
            </div>
        </div>
    )
}

export default observer(Message)