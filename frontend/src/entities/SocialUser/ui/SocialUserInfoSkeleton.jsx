import styles from './SocialUser.module.css'

const SocialUserInfoSkeleton = () => {
    return (
        <>
            <div className='flex w-full h-auto gap-9 4k:gap-14'>
                <div className={`aspect-square 2xl:w-1/5 lg:w-1/6 2k:w-1/6 md:w-1/6 sm:w-1/5 mobile:w-1/3 lg:h-full rounded-3xl`}>
                    <div 
                        className={`flex items-center justify-center w-full aspect-square
                            Avatar-standart xl:rounded-3xl lg:rounded-2xl mobile:rounded-2xl md:rounded-3xl ${styles.avatarSkeleton}`}
                    >
                        <div className='flex items-center justify-center w-1/2 aspect-square'>

                        </div>
                    </div>
                </div>

                <div
                    className='flex flex-col gap-2 2k:gap-3 4k:gap-5 sm:gap-3'
                >
                    <div>
                        <div className={`${styles.skeleton} 4k:w-96 2k:w-80 2xl:w-72 xl:w-64 lg:w-56 md:w-64 sm:w-48 mobile:w-36 lg:h-7 2k:h-8 4k:h-9 md:h-8 sm:h-7 mobile:h-6`}></div>
                    </div>
                    <div className={`${styles.skeleton} 4k:w-80 2k:w-64 2xl:w-48 xl:w-40 lg:w-40 md:w-40 sm:w-56 mobile:w-40 lg:h-5 2k:h-7 4k:h-9 md:h-6 sm:h-6 mobile:h-4`}></div>
                    <div className={`${styles.skeleton} 4k:w-80 2k:w-72 2xl:w-64 xl:w-56 lg:w-48 md:w-56 sm:w-64 mobile:w-44 lg:h-5 2k:h-7 4k:h-9 md:h-6 sm:h-6 mobile:h-4`}></div>
                    <div className={`${styles.skeleton} 4k:w-80 2k:w-72 2xl:w-60 xl:w-52 lg:w-52 md:w-52 sm:w-60 mobile:w-40 lg:h-5 2k:h-7 4k:h-9 md:h-6 sm:h-6 mobile:h-4`}></div>
                    <div className={`${styles.skeleton} 4k:w-72 2k:w-60 2xl:w-56 xl:w-48 lg:w-48 md:w-48 sm:w-52 mobile:w-36 lg:h-5 2k:h-7 4k:h-9 md:h-6 sm:h-6 mobile:h-4`}></div>
                </div>
            </div>

            <div className='flex flex-col gap-7 mt-2 2k:mt-4 4k:mt-6 2k:gap-10 4k:gap-14'>
                <div className={`${styles.buttonSkeleton} h-14 2k:h-20 4k:h-32 w-80 xl:w-72 lg:w-64 md:w-60 sm:w-56 mobile:w-56 2k:w-96 rounded-3xl`}></div>
                <div className={`${styles.buttonSkeleton} h-14 2k:h-20 4k:h-32 w-80 xl:w-72 lg:w-64 md:w-60 sm:w-56 mobile:w-56 2k:w-96 rounded-3xl`}></div>
                
            </div>
        </>
    )
}

export default SocialUserInfoSkeleton
