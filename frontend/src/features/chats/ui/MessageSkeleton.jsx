/* eslint-disable react/prop-types */
import styles from './Message.module.css'

const MessageSkeleton = ({ className, lines, displayFrom=false, me=false, key, id }) => {
    return (
        <div 
            className={`${me ? 'self-end' : 'self-start'} flex-shrink-0 2xl:max-w-[75%] 2xl:w-[45%] lg:max-w-[90%] lg:w-[70%]
                md:max-w-[70%] md:w-[40%] sm:max-w-[90%] sm:w-[70%] mobile:w-[85%]
                flex items-start md:gap-4 mobile:gap-2 2k:gap-5 4k:gap-6`}
            key={key}
            id={id}
        >
            {
                    displayFrom && !me && (
                        <div 
                            className='md:w-[40px] md:min-w-[40px] sm:w-[35px] sm:min-w-[35px] mobile:w-[30px] mobile:min-w-[30px] 2xl:w-[50px] 2xl:min-w-[50px] 2k:w-[60px] 2k:min-w-[60px] 4k:w-[80px] 4k:min-w-[80px]
                                mobile:rounded-2xl md:rounded-2xl 2k:rounded-xl aspect-square cursor-pointer self-end'
                        >
                            {/* Avatar */}
                            <div className='w-full h-full aspect-square'>
                                <div className={`flex items-center justify-center w-full h-full ${styles.avatarSkeleton} mobile:rounded-2xl md:rounded-2xl 2k:rounded-xl`}>
                                    <div className={`flex items-center justify-center w-1/2 aspect-square ${styles.skeleton}`}>
                                    </div>
                                </div>
                            </div>
                        </div>
                    )
                }
            <div 
                className={`${styles.Message} ${className} 
                    ${me ? 'rounded-bl-3xl' : 'rounded-br-3xl'} rounded-tl-3xl rounded-tr-3xl
                    mobile:px-4 mobile:py-3 md:px-5 md:py-4 2k:px-6 2k:py-5 4k:px-7 4k:py-6
                    break-words flex flex-col md:gap-2 mobile:gap-1 2k:gap-3 4k:gap-4 w-full lg:min-w-[200px] mobile:min-w-[150px]`}
            >
                <div className='flex flex-col md:gap-2 mobile:gap-1 2k:gap-3 4k:gap-4'>
                    {
                        displayFrom && !me && (
                            <div className={`${styles.from} self-start lg:h-3 2k:h-4 4k:h-5
                                md:h-3 sm:h-3 mobile:h-2 my-1 w-[60%] ${styles.skeleton}`}>
                            </div>
                        )
                    }

                    <div className='my-2 flex flex-col gap-3'>
                        {
                            Array(lines).fill(0).map((_, index) => <div key={index} className={`mobile:h-2 md:h-3 2k:h-4 4k:h-5 ${styles.skeleton}`}></div>)
                        }
                    </div>
                </div>
                <div 
                    className='self-end flex gap-2 items-center select-none'
                >
                    <div className={`${styles.Time} ${styles.skeleton} sm:w-11 mobile:w-8 lg:h-3 2k:h-4 4k:h-5 md:h-3 sm:h-3 mobile:h-2 flex-shrink-0`} ></div>
                    {
                        me && (
                            <div className={`flex items-center w-4 h-4 2k:w-6 2k:h-6 4k:w-8 4k:h-8 ${styles.skeleton}`}></div>
                        )
                    }
                </div>
            
            </div>
        </div>
    )
} 

export default MessageSkeleton