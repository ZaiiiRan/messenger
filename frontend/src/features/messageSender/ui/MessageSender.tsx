import { useState } from 'react'
import { Textarea } from '../../../shared/ui/Textarea'
import { useTranslation } from 'react-i18next'
import sendMessage from '../api/sendMessage'

interface MessageSenderProps {
    chatId: string | number
}

const MessageSender: React.FC<MessageSenderProps> = ({ chatId }) => {
    const [message, setMessage] = useState<string>('')
    const { t } = useTranslation('messageSender')

    const send = () => {
        const trimmedMessage = message.trim()
        if (trimmedMessage.length === 0) {
            return
        }
        sendMessage(chatId, trimmedMessage)
        setMessage('')
    }
    
    return (
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
    )
}

export default MessageSender