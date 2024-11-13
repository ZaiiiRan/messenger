
import styles from './ShortUser.module.css'

const ShortUserSkeleton = () => {
    return (
        <div 
            className={`${styles.ShortUser} flex items-center px-5 py-2 
                2k:px-8 2k:py-3 4k:px-12 4k:py-4 rounded-3xl xl:gap-5 mobile:gap-4 2k:gap-8 4k:gap-12`}
        >
            <div className='md:h-2/3 mobile:h-3/5 rounded-3xl aspect-square'>
                <div className={`flex items-center justify-center w-full h-full ${styles.avatarSkeleton} xl:rounded-3xl lg:rounded-2xl mobile:rounded-2xl md:rounded-3xl`}>
                    <div className={`flex items-center justify-center w-1/2 aspect-square ${styles.skeleton}`}>
                        
                    </div>
                </div>
            </div>

            <div className='flex flex-col 4k:gap-2'>
                <div 
                    className={`${styles.skeleton} 2xl:text-xl xl:text-lg lg:text-base 2k:text-2xl 4k:text-3xl
                        md:text-xl sm:text-lg mobile:text-text-base w-48 h-5 2k:h-6 4k:h-8 rounded-lg`}
                >
                </div>

                <div
                    className={`${styles.skeleton} 2xl:text-base xl:text-sm lg:text-sm 2k:text-lg 4k:text-xl
                        md:text-base sm:text-sm mobile:text-sm w-32 h-4 2k:h-5 4k:h-7 mt-2 rounded-md`}
                >
                </div>
            </div>
        </div>
    )
}

export default ShortUserSkeleton
