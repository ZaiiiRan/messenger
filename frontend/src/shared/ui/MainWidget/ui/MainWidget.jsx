/* eslint-disable react/prop-types */
import './MainWidget.css'
import { motion } from 'framer-motion'

const MainWidget = ({ goBack, title, children }) => {
    return (
        <motion.div 
            initial={{ opacity: 0, x: -500 }}
            animate={{ opacity: 1, x: 0}}
            exit={{ opacity: 0, x: -500 }}
            transition={{ duration: 0.3 }}
            className='Main-Widget rounded-3xl flex flex-col gap-8 2k:gap-14 4k:gap-24'
        >
            {/* Title */}
            <div className='flex items-center gap-5 2k:gap-7 4k:gap-9'>
                <div 
                    className='backBtn 2xl:w-10 2xl:h-10 xl:w-9 xl:h-9 lg:w-9 lg:h-8 2k:w-12 2k:h-12 4k:w-14 4k:h-14 
                        mobile:w-8 mobile:h-8 md:w-9 md:h-9 
                        rounded-3xl flex items-center justify-center'
                    onClick={goBack}
                >
                    <div className='Icon flex items-center justify-center h-1/4 aspect-square'>
                        <svg viewBox="0 0 7.424 12.02" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <defs/>
                            <path id="Vector" d="M0 6.01L6 12.02L7.42 10.6L2.82 6L7.42 1.4L6 0L0 6.01Z" fillOpacity="1.000000" fillRule="nonzero"/>
                        </svg>
                    </div>
                </div>
                <h1 className='font-extrabold 2xl:text-3xl xl:text-2xl lg:text-xl 2k:text-4xl 4k:text-5xl
                    md:text-3xl sm:text-2xl mobile:text-xl'
                >
                    { title }
                </h1>
            </div>

            { children }
        </motion.div>
    )
}

export default MainWidget
