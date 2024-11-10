/* eslint-disable react/prop-types */
import './OptionsList.css'
import { useTranslation } from 'react-i18next'

const OptionsList = ({ open }) => {
    const { t } = useTranslation('optionsFeature')

    return (
        <div className='Options-List rounded-3xl flex flex-col gap-6 2k:gap-10 4k:gap-14'>
            <h1 className='font-extrabold 2xl:text-3xl xl:text-2xl lg:text-xl 2k:text-4xl 4k:text-5xl'>
                { t('Settings') }
            </h1>

            <div className='Options-List__container flex flex-col items-center gap-5'>
                <div 
                    className='Options-List__element flex items-center justify-between px-5 py-2 2k:px-8 2k:py-3 4k:px-12 4k:py-4 rounded-3xl'
                    onClick={() => open('appearance')}
                >
                    <div className='flex gap-4 2k:gap-6 4k:gap-8 items-center h-full'>
                        <div className='Option-Icon flex items-center justify-center h-1/3 aspect-square'>
                            <svg viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                                <defs/>
                                <path id="Vector" d="M11 20L9 20L9 17L11 17L11 20ZM16.36 17.77L14.24 15.65L15.65 14.24L17.77 16.36L16.36 17.77L16.36 17.77ZM3.63 17.77L2.22 16.36L4.34 14.24L5.75 15.65L3.63 17.77L3.63 17.77ZM10 15C7.23 15 4.99 12.76 4.99 9.99C4.99 7.23 7.23 4.99 10 4.99C12.76 4.99 15 7.23 15 10C15 12.76 12.76 15 10 15ZM10 6.99C8.33 6.99 6.99 8.34 6.99 10C6.99 11.66 8.34 13 10 13C11.66 13 13 11.66 13 10C13 8.33 11.66 6.99 10 6.99ZM20 11L17 11L17 9L20 9L20 11ZM3 11L0 11L0 9L3 9L3 11ZM15.65 5.75L14.24 4.34L16.36 2.22L17.77 3.63L15.65 5.75L15.65 5.75ZM4.34 5.75L2.22 3.63L3.63 2.22L5.75 4.34L4.34 5.75L4.34 5.75ZM11 3L9 3L9 0L11 0L11 3Z" fillOpacity="1.000000" fillRule="nonzero"/>
                            </svg>
                        </div>
                        <div className='2xl:text-xl xl:text-lg lg:text-base 2k:text-2xl 4k:text-3xl'>{ t('Appearance') }</div>
                    </div>

                    <div className='Option-Icon flex items-center justify-center h-1/6 aspect-square'>
                        <svg viewBox="0 0 7.425 12.021" xmlns="http://www.w3.org/2000/svg">
                            <defs/>
                            <path id="Vector" d="M7.42 6L1.41 0L0 1.41L4.6 6.01L0 10.6L1.41 12.02L7.42 6Z" fillOpacity="1.000000" fillRule="nonzero"/>
                        </svg>
                    </div>
                </div>

            </div>

        </div>
    )
}

export default OptionsList
