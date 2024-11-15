import { motion, AnimatePresence } from 'framer-motion'
import { useState } from 'react'
import { ListWidget } from '../../../shared/ui/ListWidget'
import { useTranslation } from 'react-i18next'


const MessenginPage = () => {
    const [selected, setSelected] = useState(null)
    const { t } = useTranslation('messengingPage')

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
            <div className='h-full lg:w-2/5 2k:w-7/20 flex flex-col items-center justify-between lg:gap-10 2k:gap-20 4k:gap-32 mobile:w-full_screen'>
                <ListWidget className='h-2/5 w-full flex-grow basis-2/5' title={t('People')} >
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                </ListWidget>
                <ListWidget className='h-1/2 w-full flex-grow basis-2/5' title={t('Groups')} >
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                    <div>test</div>
                </ListWidget>
            </div>
            

            <AnimatePresence mode='wait'>
                
            </AnimatePresence>
        </motion.div>
    )
}

export default MessenginPage
