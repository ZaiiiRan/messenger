import userStore from '../store/userStore'
import styles from './User.module.css'
import { motion, AnimatePresence } from 'framer-motion'
import { observer } from 'mobx-react'
import { Button } from '../../../shared/ui/Button'
import { useTranslation } from 'react-i18next'

const User = observer(() => {
    const { t } = useTranslation('userCard')

    const logout = async (e) => {
        e.preventDefault()
        await userStore.logout()
    }


    return (
        <AnimatePresence mode='popLayout'>
            {
                userStore.isOpen && (
                <motion.div 
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    exit={{ opacity: 0 }}
                    transition={{ duration: 0.3 }}
                    className={`${styles.UserCardBG} fixed top-0 left-0 right-0 bottom-0 flex items-center justify-center bg-black bg-opacity-50 z-50`}
                >
                    <motion.div 
                        initial={{ opacity: 0, transform: 'scale(0)' }}
                        animate={{ opacity: 1, transform: 'scale(1)' }}
                        exit={{ opacity: 0, transform: 'scale(0)' }}
                        transition={{ duration: 0.3 }}
                        className={`${styles.UserCard} flex flex-col rounded-3xl gap-8 2k:gap-14 4k:gap-24`}
                    >
                        <div className='flex items-center gap-5 2k:gap-7 4k:gap-9'>
                            <div 
                                className='backBtn 2xl:w-10 2xl:h-10 xl:w-9 xl:h-9 lg:w-9 lg:h-8 2k:w-12 2k:h-12 4k:w-14 4k:h-14 
                                    mobile:w-8 mobile:h-8 md:w-9 md:h-9 
                                    rounded-3xl flex items-center justify-center'
                                onClick={() => { userStore.setOpen(false) }}
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
                                { t('Your profile') }
                            </h1>
                        </div>

                        <div className='flex w-full h-auto gap-9 4k:gap-14'>
                            <div className={`${styles.Avatar} aspect-square 2xl:w-1/5 lg:w-1/6 2k:w-1/6 md:w-1/6 sm:w-1/5 mobile:w-1/3 lg:h-full rounded-3xl`}>
                                <div 
                                    className='flex items-center justify-center w-full aspect-square
                                        Avatar-standart xl:rounded-3xl lg:rounded-2xl mobile:rounded-2xl md:rounded-3xl'
                                >
                                    <div className='flex items-center justify-center w-1/2 aspect-square'>
                                        <svg width="100%" height="100%" viewBox="0 0 16 19" fill="none" xmlns="http://www.w3.org/2000/svg">
                                            <defs/>
                                            <path id="Vector" d="M8 0C5.23 0 3 2.23 3 5C3 7.76 5.23 10 8 10C10.76 10 13 7.76 13 5C13 2.23 10.76 0 8 0ZM11 5C11 6.65 9.65 8 8 8C6.34 8 5 6.65 5 5C5 3.34 6.34 2 8 2C9.65 2 11 3.34 11 5ZM0 19C0 16.87 0.84 14.84 2.34 13.34C3.84 11.84 5.87 11 8 11C10.12 11 12.15 11.84 13.65 13.34C15.15 14.84 16 16.87 16 19L14 19C14 17.4 13.36 15.88 12.24 14.75C11.11 13.63 9.59 13 8 13C6.4 13 4.88 13.63 3.75 14.75C2.63 15.88 2 17.4 2 19L0 19Z" fill="#0F1828" fillOpacity="1.000000" fillRule="evenodd"/>
                                        </svg>
                                    </div>
                                </div>
                            </div>

                            <div
                                className='lg:text-base 2k:text-2xl 4k:text-4xl
                                    md:text-xl sm:text-lg mobile:text-sm flex flex-col gap-2 2k:gap-3 4k:gap-5 sm:gap-3'
                            >
                                <div
                                    className='font-bold lg:text-2xl 2k:text-3xl 4k:text-5xl
                                        md:text-3xl sm:text-2xl mobile:text-xl'
                                >
                                    { userStore.user.firstname } { userStore.user.lastname }
                                </div>
                                <div>{ t('Username') }: { userStore.user.username }</div>
                                <div>Email: { userStore.user.email }</div>
                                { 
                                    userStore.user.phone && (
                                        <div>{ t('Phone number') }: { userStore.user.phone }</div>
                                    )    
                                }
                                { 
                                    userStore.user.birthdate && (
                                        <div>{ t('Birthdate') }: { new Date(userStore.user.birthdate).toLocaleDateString() }</div>
                                    )    
                                }
                            </div>
                        </div>

                        <div className={`${styles.LogoutBtn} lg:self-end`}>
                            <div className='2xl:h-16 lg:h-12 xl:h-14 2k:h-20 4k:h-28 sm:h-16 mobile:h-12'>
                                <Button 
                                    className='h-full flex items-center justify-center aspect-square 2xl:rounded-3xl lg:rounded-2xl sm:rounded-3xl mobile:rounded-2xl'
                                    onClick={logout}
                                >
                                    <div className='w-1/2 aspect-square'>
                                        <svg viewBox="0 0 45 46.4018" fill="none" xmlns="http://www.w3.org/2000/svg">
                                            <defs/>
                                            <path id="Vector" d="M18.31 0.28C19.31 -0.02 20.38 -0.09 21.41 0.1C22.45 0.28 23.43 0.7 24.27 1.33C25.12 1.96 25.8 2.77 26.28 3.71C26.75 4.65 26.99 5.69 27 6.74L27 39.65C26.99 40.7 26.75 41.74 26.28 42.68C25.8 43.62 25.12 44.43 24.27 45.06C23.43 45.69 22.45 46.11 21.41 46.29C20.38 46.48 19.31 46.41 18.31 46.11L4.81 42.06C3.42 41.65 2.2 40.79 1.33 39.63C0.46 38.46 0 37.05 0 35.6L0 10.79C0 9.34 0.46 7.93 1.33 6.77C2.2 5.6 3.42 4.75 4.81 4.33L18.31 0.28ZM29.25 5.2C29.25 4.6 29.48 4.03 29.9 3.6C30.33 3.18 30.9 2.95 31.5 2.95L38.25 2.95C40.04 2.95 41.75 3.66 43.02 4.92C44.28 6.19 45 7.91 45 9.7L45 11.95C45 12.54 44.76 13.11 44.34 13.54C43.91 13.96 43.34 14.2 42.75 14.2C42.15 14.2 41.58 13.96 41.15 13.54C40.73 13.11 40.5 12.54 40.5 11.95L40.5 9.7C40.5 9.1 40.26 8.53 39.84 8.1C39.41 7.68 38.84 7.45 38.25 7.45L31.5 7.45C30.9 7.45 30.33 7.21 29.9 6.79C29.48 6.36 29.25 5.79 29.25 5.2ZM42.75 32.2C43.34 32.2 43.91 32.43 44.34 32.85C44.76 33.28 45 33.85 45 34.45L45 36.7C45 38.49 44.28 40.2 43.02 41.47C41.75 42.73 40.04 43.45 38.25 43.45L31.5 43.45C30.9 43.45 30.33 43.21 29.9 42.79C29.48 42.36 29.25 41.79 29.25 41.2C29.25 40.6 29.48 40.03 29.9 39.6C30.33 39.18 30.9 38.95 31.5 38.95L38.25 38.95C38.84 38.95 39.41 38.71 39.84 38.29C40.26 37.86 40.5 37.29 40.5 36.7L40.5 34.45C40.5 33.85 40.73 33.28 41.15 32.85C41.58 32.43 42.15 32.2 42.75 32.2ZM15.75 20.95C15.15 20.95 14.58 21.18 14.15 21.6C13.73 22.03 13.5 22.6 13.5 23.2C13.5 23.79 13.73 24.36 14.15 24.79C14.58 25.21 15.15 25.45 15.75 25.45L15.75 25.45C16.34 25.45 16.92 25.21 17.34 24.79C17.76 24.36 18 23.79 18 23.2C18 22.6 17.76 22.03 17.34 21.6C16.92 21.18 16.34 20.95 15.75 20.95L15.75 20.95Z" fill="#FFFFFF" fillOpacity="1.000000" fillRule="evenodd"/>
                                        </svg>
                                    </div>
                                </Button>
                            </div>
                        </div>
                    </motion.div>
                </motion.div>
                )
            }
        </AnimatePresence>
    )
})

export default User
