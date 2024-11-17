import { motion, AnimatePresence } from 'framer-motion'
import { useState } from 'react'
import { PeopleChatsList, GroupChatsList, ChatWidget } from '../../../features/chats'
import './MessangingPage.css'

const MessenginPage = () => {
    const [selected, setSelected] = useState(null)

    const open = (chatID) => {
        setSelected(chatID)
    }

    const goBack = () => {
        setSelected(null)
    }

    return (
        <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1}}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.3 }}
            className='w-full h-full flex relative lg:gap-10 xl:gap-12 2xl:gap-14 2k:gap-24 4k:gap-36'
        >
            <div className='chat_lists h-full lg:w-2/5 2k:w-7/20 flex flex-col items-center justify-between lg:gap-10 2k:gap-20 4k:gap-32 mobile:w-full_screen'>
                <PeopleChatsList open={open} />
                <GroupChatsList open={open} />
            </div>
            

            <AnimatePresence mode='wait'>
                {
                    selected && (
                        <ChatWidget key={selected} goBack={goBack} />
                    )
                }
            </AnimatePresence>
        </motion.div>
    )
}

export default MessenginPage
