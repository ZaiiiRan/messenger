import { observer } from 'mobx-react'
import modalStore from '../store/modalStore'
import { Button } from '../../../shared/ui/Button'
import { motion, AnimatePresence } from 'framer-motion'
import styles from './Modal.module.css'
import { useTranslation } from 'react-i18next'
import ModalData from '../models/modalData'

const Modal = observer(() => {
    const { t } = useTranslation('modal')

    const confirm = (modal: ModalData) => {
        if (modal.actionFunction) {
            modal.actionFunction()
        }
        modalStore.closeModal(modal.id)
    }

    return (
        <AnimatePresence mode='popLayout'>
            { modalStore.modals.map((modal) => (
                <motion.div 
                    key={modal.id}
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    exit={{ opacity: 0 }}
                    transition={{ duration: 0.3 }}
                    className={`${styles.ModalBG} fixed top-0 left-0 right-0 bottom-0 flex items-center justify-center bg-black bg-opacity-50`}
                >
                    <motion.div 
                        initial={{ opacity: 0, transform: 'scale(0)' }}
                        animate={{ opacity: 1, transform: 'scale(1)' }}
                        exit={{ opacity: 0, transform: 'scale(0)' }}
                        transition={{ duration: 0.3 }}
                        className={`${styles.Modal} bg-white
                            flex flex-col rounded-3xl md:p-12 mobile:p-8 2k:p-20 4k:p-24 gap-12 2k:gap-16 4k:gap-20 lg:w-1/3 sm:w-96 mobile:w-72`}
                    >
                        <h1 className='text-center font-extrabold 
                            md:text-3xl mobile:text-2xl 2k:text-5xl 4k:text-7xl'
                        >
                            {modal.title}
                        </h1>
                        <div className='text-xl mobile:text-lg 2k:text-2xl 4k:text-4xl text-center'>
                            {modal.text}
                        </div>
                        <div className='flex w-full gap-5 2k:gap-8 4k:gap-12'>
                            {
                                modal.actionFunction ? (
                                    <>
                                        <Button 
                                        className='w-full h-14 2k:h-20 4k:h-32 rounded-3xl font-semibold 
                                            md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                                        onClick={() => modalStore.closeModal(modal.id)}
                                        >
                                            { t('No') }
                                        </Button>
                                        <Button 
                                        className='w-full h-14 2k:h-20 4k:h-32 rounded-3xl font-semibold 
                                            md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                                        onClick={() => confirm(modal) }
                                        >
                                            { t('Yes') }
                                        </Button>
                                    </>
                                ) : (
                                    <Button 
                                        className='w-full h-14 2k:h-20 4k:h-32 rounded-3xl font-semibold 
                                            md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                                        onClick={() => modalStore.closeModal(modal.id)}
                                    >
                                        ОК
                                    </Button>
                                )
                            }
                            
                        </div>
                        
                    </motion.div>
                </motion.div>
            ))}
        </AnimatePresence>
    )
})

export default Modal
