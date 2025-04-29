import { Button } from '../../../shared/ui/Button'
import { useTranslation } from 'react-i18next'
import { motion, AnimatePresence } from 'framer-motion'
import styles from './SendMessageModal.module.css'
import { ShortUser } from '../../../entities/SocialUser'
import { Textarea } from '../../../shared/ui/Textarea'
import { useEffect, useState } from 'react'
import { useModal } from '../../modal'
import { createPortal } from 'react-dom'

const SendMessageModal = ({ recipient, show, setShow }) => {
    const { t } = useTranslation('chatsFeature')
    const [message, setMessage] = useState('')

    const { openModal, setModalText, setModalTitle } = useModal()

    const send = () => {
        if (message.trim().length === 0) {
            setModalTitle(t('Error'))
            setModalText(t('Message is empty'))
            openModal()
        }
    }

    useEffect(() => {
        setMessage('')
    }, [show, recipient])

    return ( 
        createPortal(
            <AnimatePresence mode='popLayout'>
                {
                    show && (
                        <motion.div 
                        layoutId="modal"
                            initial={{ opacity: 0, x: 0, y: 0 }}
                            animate={{ opacity: 1, x: 0, y: 0 }}
                            exit={{ opacity: 0, x: 0, y: 0 }}
                            transition={{ duration: 0.3 }}
                            className={`${styles.SendMessageModalBG} fixed top-0 left-0 right-0 bottom-0 flex items-center justify-center bg-black bg-opacity-50`}
                        >
                            <motion.div 
                            layoutId="modal-content"
                                initial={{ opacity: 0, transform: 'scale(0)' }}
                                animate={{ opacity: 1, transform: 'scale(1)' }}
                                exit={{ opacity: 0, transform: 'scale(0)' }}
                                transition={{ duration: 0.3 }}
                                className={`${styles.SendMessageModal} flex flex-col rounded-3xl gap-8 2k:gap-14 4k:gap-24`}
                            >
                                <div className='flex items-center gap-5 2k:gap-7 4k:gap-9'>
                                    <div 
                                        className='backBtn 2xl:w-10 2xl:h-10 xl:w-9 xl:h-9 lg:w-9 lg:h-8 2k:w-12 2k:h-12 4k:w-14 4k:h-14 
                                            mobile:w-8 mobile:h-8 md:w-9 md:h-9 
                                            rounded-3xl flex items-center justify-center'
                                        onClick={() => { setShow(false) }}
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
                                        { t('New message') }
                                    </h1>
                                </div>
                                <ShortUser 
                                    isClickable={false}
                                    user={recipient}
                                />

                                <Textarea 
                                    className='w-full px-5 py-3 2k:px-6 2k:py-4 4k:px-7 4k:py-5 
                                        rounded-lg md:text-lg mobile:text-sm lg:text-sm 2xl:text-lg 2k:text-2xl 4k:text-4xl 2k:h-24 4k:h-40' 
                                    placeholder={t('Enter message')}
                                    value={message}
                                    onChange={(e) => setMessage(e.target.value)}
                                />

                                <div className='self-end'>
                                    <Button
                                        className='h-12 flex items-center justify-center 2k:h-16 4k:h-28 w-72 xl:w-60 lg:w-56 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                                            rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                                        onClick={send}
                                    >
                                        { t('Send') }
                                    </Button>
                                </div>
                            </motion.div>
                        </motion.div>
                    )
                }
            </AnimatePresence>, document.body
        )  
    )
}

export default SendMessageModal
