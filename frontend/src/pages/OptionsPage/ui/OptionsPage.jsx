import './OptionsPage.css'
import { motion, AnimatePresence } from 'framer-motion'
import { OptionsList, AppearanceOptions } from '../../../features/options'
import { useState } from 'react'


const OptionsPage = () => {
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
            className='w-full h-full flex Options-Page'
        >
            <OptionsList open={open}/>

            <AnimatePresence mode='wait'>
                {
                    selected === 'appearance' && (
                        <AppearanceOptions goBack={goBack} />
                    )
                }
            </AnimatePresence>
        </motion.div>
    )
}

export default OptionsPage
