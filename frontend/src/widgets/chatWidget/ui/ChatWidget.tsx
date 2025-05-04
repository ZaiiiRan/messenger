import { motion } from 'framer-motion'
import { Textarea } from '../../../shared/ui/Textarea'
import React, { useState, useRef, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { chatStore } from '../../../entities/Chat'
import { observer } from 'mobx-react'
import './ChatWidget.css'
import { shortUserStore } from '../../../entities/SocialUser'
import ChatWidgetHeader from './ChatWidgetHeader'
import { MessageSender } from '../../../features/messageSender'

interface IChatWidgetProps {
    goBack: () => void,
    selected: string | number
}

const ChatWidget: React.FC<IChatWidgetProps> = ({ goBack, selected }) => {
    const { t } = useTranslation('chatWidget')
    const chat = chatStore.get(selected)
    const [message, setMessage] = useState<string>('')
    const [chatName, setChatName] = useState<string>(chat?.chat.name || '')

    const isGroupChat = chat?.chat.isGroupChat

    useEffect(() => {
        let isMounted = true
        
        const loadPartnerName = async () => {
            if (!isGroupChat && chat) {
                const partner = await shortUserStore.get(chat.members[0].userId)
                if (isMounted) {
                    setChatName(partner ? `${partner.firstname} ${partner.lastname}` : '???')
                } else {
                    if (isMounted) {
                        setChatName('')
                    }
                }
            }
        }

        if (chat?.chat.name) {
            setChatName(chat.chat.name)
        } else if (!isGroupChat) {
            loadPartnerName()
        }

        return () => {
            isMounted = false
        }
    }, [selected])

    return (
        <motion.div 
            initial={{ opacity: 0, x: -500 }}
            animate={{ opacity: 1, x: 0}}
            exit={{ opacity: 0, x: -500 }}
            transition={{ duration: 0.3 }}
            className='Chat-Widget rounded-3xl flex flex-col gap-8 2k:gap-14 4k:gap-24'
        >
            {/* Header */}
            <ChatWidgetHeader goBack={goBack} chatName={chatName} isGroupChat/>

            {/* Main */}
            <div 
                className='Chat-Main-Area rounded-3xl w-full h-4/5 p-8 flex flex-col gap-5 2k:gap-9 4k:gap-12 overflow-y-scroll relative'
            >
            
            </div>

            {/* Message Sender */}
            <MessageSender chatId={selected} />
        </motion.div>
    )
}

export default observer(ChatWidget)