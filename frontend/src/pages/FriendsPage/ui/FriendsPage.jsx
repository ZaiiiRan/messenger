import { motion, AnimatePresence } from 'framer-motion'
import { useState } from 'react'
import { FriendsMenu } from '../../../features/friends'


const FriendsPage = () => {
    const [selected, setSelected] = useState(null)

    const open = (optionGroup) => {
        setSelected(optionGroup)
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
        <FriendsMenu open={open}/>

            <AnimatePresence mode='wait'>
                
            </AnimatePresence>
        </motion.div>
    )
}

export default FriendsPage
