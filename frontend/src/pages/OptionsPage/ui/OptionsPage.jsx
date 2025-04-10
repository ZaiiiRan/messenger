import { motion, AnimatePresence } from 'framer-motion'
import { OptionsList, AppearanceOptions } from '../../../widgets/options'
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
            className='w-full h-full flex relative lg:gap-10 xl:gap-12 2xl:gap-14 2k:gap-24 4k:gap-36'
        >
            <OptionsList open={open}/>

            <AnimatePresence mode='wait'>
                {
                    selected === 'appearance' && (
                        <AppearanceOptions goBack={ goBack } />
                    )
                }
            </AnimatePresence>
        </motion.div>
    )
}

export default OptionsPage
