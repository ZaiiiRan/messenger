import styles from './ChatCard.module.css'

const ChatCardSkeleton = () => {
    return (
        <div 
            className={`${styles.ChatCard} flex items-center px-5 py-2 
                2k:px-8 2k:py-3 4k:px-12 4k:py-4 rounded-3xl xl:gap-5 mobile:gap-4 2k:gap-8 4k:gap-12`}
        >
            {/* Avatar */}
            <div className='md:h-2/3 mobile:h-3/5 rounded-3xl aspect-square'>
                <div className={`flex items-center justify-center w-full h-full ${styles.avatarSkeleton} xl:rounded-3xl lg:rounded-2xl mobile:rounded-2xl md:rounded-3xl`}>
                    <div className={`flex items-center justify-center w-1/2 aspect-square ${styles.skeleton}`}>
                        
                    </div>
                </div>
            </div>

            {/* Chat Info */}
            <div className='w-4/5 overflow-hidden flex flex-col 4k:gap-2'>
                <div 
                    className={`${styles.skeleton} w-40 h-5 2k:h-6 4k:h-8 rounded-lg`}
                >
                </div>

                {/* Last Message */}
                <div
                    className={`${styles.skeleton} w-32 h-4 2k:h-5 4k:h-7 mt-2 rounded-md`}
                >
                </div>
            </div>

            {/* Status */}
            <div className='flex flex-col items-end justify-between gap-2 2k:gap-3 4k:gap-4 ml-auto'>
                <div
                    className={`${styles.skeleton} w-10 h-4 2k:h-5 4k:h-7 rounded-lg`}
                >
                </div>
                <div
                    className={`${styles.skeleton} p-1 2k:p-1.5 4k:p-2 w-5 h-4 2k:h-5 4k:h-7 rounded-lg`}
                >
                </div>
            </div>
        </div>
    )
}

export default ChatCardSkeleton