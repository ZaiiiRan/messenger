import { Button } from '../../../shared/ui/Button'
import { useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import './StartPage.css'
import { useTranslation } from 'react-i18next'

const StartPage = () => {
    const { t } = useTranslation('startPage')
    const navigate = useNavigate()

    const handleClick = () => {
        navigate('/login')
    }

    return (
        <motion.div 
            initial={{ opacity: 0, x: -1000 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: 1000 }}
            transition={{ duration: 0.7 }}
            className='w-full_screen h-full_screen flex flex-col items-center justify-center 
                2xl:gap-24 xl:gap-20 lg:gap-24 md:gap-36 sm:gap-24 mobile:gap-56 2k:gap-36 4k:gap-56'
        >
            <div className='lg:w-1/3 mobile:w-full mobile:p-12 sm:p-0 sm:w-1/2 flex flex-col items-center 
                2xl:gap-9 xl:gap-7 lg:gap-6 md:gap-9 sm:gap-9 mobile:gap-10 2k:gap-14 4k:gap-20'
            >
                <div className='flex justify-center'>
                    <img className='2xl:w-1/3 mobile:w-1/2 4k:w-1/2 shake' src="./message.svg" alt="message" draggable={false}/>
                </div>
                <div className='font-bold text-3xl mobile:text-2xl 2k:text-4xl 4k:text-6xl text-center'>
                    { t('Start messaging your family and friends now') }
                </div>
            </div>
            
            <Button 
                className='lg:w-1/4 mobile:w-1/2 h-14 2k:h-20 4k:h-32 rounded-3xl font-semibold 
                    md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                onClick={handleClick}
            >
                { t('Start chatting') }
            </Button>
        </motion.div>
    )
}

export default StartPage
