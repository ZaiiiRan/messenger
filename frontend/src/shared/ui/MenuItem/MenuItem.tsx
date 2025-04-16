import styles from './MenuItem.module.css'

interface MenuItemProps {
    icon?: React.ReactNode,
    text?: React.ReactNode,
    onClick?: (event: React.MouseEvent<HTMLDivElement>) => void
}

const MenuItem: React.FC<MenuItemProps> = ({ icon, text, onClick }) => {
    return (
        <div 
            className={`${styles.MenuItem} flex items-center justify-between px-5 py-2 2k:px-8 2k:py-3 4k:px-12 4k:py-4 rounded-3xl`}
            onClick={onClick}
        >
            <div className='flex gap-4 2k:gap-6 4k:gap-8 items-center h-full'>
                <div className='Icon flex items-center justify-center h-1/3 aspect-square'>
                    { icon }
                </div>
                <div 
                    className='2xl:text-xl xl:text-lg lg:text-base 2k:text-2xl 4k:text-3xl
                        md:text-xl sm:text-lg mobile:text-text-base'
                >
                    { text }
                </div>
            </div>

            <div className='Icon flex items-center justify-center h-1/6 aspect-square'>
                <svg viewBox="0 0 7.425 12.021" xmlns="http://www.w3.org/2000/svg">
                    <defs/>
                    <path id="Vector" d="M7.42 6L1.41 0L0 1.41L4.6 6.01L0 10.6L1.41 12.02L7.42 6Z" fillOpacity="1.000000" fillRule="nonzero"/>
                </svg>
            </div>
        </div>
    )
}

export default MenuItem
