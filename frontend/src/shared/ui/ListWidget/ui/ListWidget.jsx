/* eslint-disable react/prop-types */
import './ListWidget.css'

const ListWidget = ({ title, children }) => {
    return (
        <div className='Widget-List lg:rounded-3xl flex flex-col 
            gap-6 2k:gap-10 4k:gap-14'
        >
            <h1 className='font-extrabold 2xl:text-3xl xl:text-2xl lg:text-xl 2k:text-4xl 4k:text-5xl
                md:text-3xl sm:text-2xl mobile:text-xl'
            >
                { title }
            </h1>

            <div className='scrollContainer Widget-List__container flex flex-col items-center gap-5'>
                { children }
            </div>

        </div>
    )
}

export default ListWidget
