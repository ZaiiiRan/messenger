import { observer } from 'mobx-react'
import modalStore from '../../../app/stores/modalStore/modalStore'
import { Button } from '../../../shared/ui/Button'
import { motion, AnimatePresence } from 'framer-motion'
import styles from './Modal.module.css'

const Modal = observer(() => {
    return (
        <AnimatePresence mode='popLayout'>
            { modalStore.isOpen && 
                <motion.div 
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    exit={{ opacity: 0 }}
                    transition={{ duration: 0.3 }}
                    className={`${styles.ModalBG} fixed top-0 left-0 right-0 bottom-0 flex items-center justify-center bg-black bg-opacity-50 z-50`}
                >
                    <motion.div 
                        initial={{ opacity: 0, transform: 'scale(0)' }}
                        animate={{ opacity: 1, transform: 'scale(1)' }}
                        exit={{ opacity: 0, transform: 'scale(0)' }}
                        transition={{ duration: 0.3 }}
                        className={`${styles.Modal} bg-white
                            flex flex-col rounded-3xl p-12 2k:p-20 4k:p-24 gap-12 2k:gap-16 4k:gap-20 lg:w-1/3 sm:w-96 mobile:w-72`}
                    >
                        <h1 className='text-center font-extrabold 
                            md:text-3xl mobile:text-2xl 2k:text-5xl 4k:text-7xl'
                        >
                            {modalStore.title}
                        </h1>
                        <div className='text-xl mobile:text-xl 2k:text-2xl 4k:text-4xl text-center'>
                            {modalStore.text}
                        </div>
                        <Button 
                            className='h-14 2k:h-20 4k:h-32 rounded-3xl font-semibold 
                                md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                            onClick={() => modalStore.closeModal()}
                        >
                            ОК
                        </Button>
                    </motion.div>
                </motion.div>
            }
        </AnimatePresence>
    )
})

export default Modal
