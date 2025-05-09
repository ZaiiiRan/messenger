import IShortUser from '../models/IShortUser'
import styles from './ShortUser.module.css'

interface ShortUserProps {
    user: IShortUser,
    onClick?: (event: React.MouseEvent<HTMLDivElement>) => void,
    isClickable?: boolean,
    isSelectable?: boolean,
    isSelected?: boolean
}

const ShortUser: React.FC<ShortUserProps> = ({ user, onClick, isClickable=true, isSelectable, isSelected }) => {
    return (
        <div 
            className={`${styles.ShortUser} flex items-center px-5 py-2 
                2k:px-8 2k:py-3 4k:px-12 4k:py-4 rounded-3xl xl:gap-5 mobile:gap-4 2k:gap-8 4k:gap-12 ${isClickable ? '' : styles.notClickable}`}
            onClick={onClick}
        >
            <div className='md:h-2/3 mobile:h-3/5 rounded-3xl aspect-square'>
                <div className='flex items-center justify-center w-full h-full Avatar-standart xl:rounded-3xl lg:rounded-2xl mobile:rounded-2xl md:rounded-3xl'>
                    <div className='flex items-center justify-center w-1/2 aspect-square'>
                        {
                            user.isActivated && !user.isDeleted && !user.isBanned ? (
                                <svg viewBox="0 0 16 19" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <defs/>
                                    <path id="Vector" d="M8 0C5.23 0 3 2.23 3 5C3 7.76 5.23 10 8 10C10.76 10 13 7.76 13 5C13 2.23 10.76 0 8 0ZM11 5C11 6.65 9.65 8 8 8C6.34 8 5 6.65 5 5C5 3.34 6.34 2 8 2C9.65 2 11 3.34 11 5ZM0 19C0 16.87 0.84 14.84 2.34 13.34C3.84 11.84 5.87 11 8 11C10.12 11 12.15 11.84 13.65 13.34C15.15 14.84 16 16.87 16 19L14 19C14 17.4 13.36 15.88 12.24 14.75C11.11 13.63 9.59 13 8 13C6.4 13 4.88 13.63 3.75 14.75C2.63 15.88 2 17.4 2 19L0 19Z" fill="#0F1828" fillOpacity="1.000000" fillRule="evenodd"/>
                                </svg>
                            ) : (
                                <svg viewBox="-3.5 0 19 19" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" strokeWidth="0"></g><g id="SVGRepo_tracerCarrier" strokeLinecap="round" strokeLinejoin="round"></g><g id="SVGRepo_iconCarrier"><path d="M11.383 13.644A1.03 1.03 0 0 1 9.928 15.1L6 11.172 2.072 15.1a1.03 1.03 0 1 1-1.455-1.456l3.928-3.928L.617 5.79a1.03 1.03 0 1 1 1.455-1.456L6 8.261l3.928-3.928a1.03 1.03 0 0 1 1.455 1.456L7.455 9.716z"></path></g></svg>
                            )
                        }
                    </div>
                </div>
            </div>

            <div className='flex flex-col 4k:gap-2 max-w-[90%]'>
                <div 
                    className='2xl:text-xl xl:text-lg lg:text-base 2k:text-2xl 4k:text-3xl
                        md:text-xl sm:text-lg mobile:text-text-base text-ellipsis overflow-hidden'
                >
                    { user.firstname } { user.lastname }
                </div>

                <div
                    className='2xl:text-base xl:text-sm lg:text-sm 2k:text-lg 4k:text-xl
                        md:text-base sm:text-sm mobile:text-sm  text-ellipsis overflow-hidden'
                >
                    { user.username }
                </div>
            </div>

            {
                isSelectable && (
                    <div className={`${styles.select} ml-auto rounded-full flex items-center justify-cente border-2 w-5 aspect-square`}>
                        {
                            isSelected && (
                                <div className={`${styles.selected} w-full aspect-square rounded-full`}></div>
                            )
                        }
                    </div>
                )
            }
        </div>
    )
}

export default ShortUser
