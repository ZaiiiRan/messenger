import { Select } from '../../../shared/ui/Select'
import { useTranslation } from 'react-i18next'
import { useMemo } from 'react'
import { observer } from 'mobx-react-lite'
import { themeStore } from '../../../shared/theme'
import OptionsItem from '../models/OptionsItem'

const ThemeChanging = observer(() => {
    const { t } = useTranslation('optionsFeature')

    const themeOptions = useMemo(() => [
        { key: 'light', label: t('Light') },
        { key: 'dark', label: t('Dark') },
        { key: 'system', label: t('System') }
    ], [t])

    const handleThemeChange = (selectedTheme: OptionsItem) => {
        themeStore.setTheme(selectedTheme.key)
    }

    return (
        <div className='flex 2xl:w-1/2 xl:w-4/6 justify-between'>
            <div className='flex items-center gap-3 2k:gap-5 4k:gap-8'>
                <div className='Option-Icon'>
                    <svg viewBox="0 -3 24 24" xmlns="http://www.w3.org/2000/svg">
                        <defs/>
                        <path id="Vector" d="M11 20L9 20L9 17L11 17L11 20ZM16.36 17.77L14.24 15.65L15.65 14.24L17.77 16.36L16.36 17.77L16.36 17.77ZM3.63 17.77L2.22 16.36L4.34 14.24L5.75 15.65L3.63 17.77L3.63 17.77ZM10 15C7.23 15 4.99 12.76 4.99 9.99C4.99 7.23 7.23 4.99 10 4.99C12.76 4.99 15 7.23 15 10C15 12.76 12.76 15 10 15ZM10 6.99C8.33 6.99 6.99 8.34 6.99 10C6.99 11.66 8.34 13 10 13C11.66 13 13 11.66 13 10C13 8.33 11.66 6.99 10 6.99ZM20 11L17 11L17 9L20 9L20 11ZM3 11L0 11L0 9L3 9L3 11ZM15.65 5.75L14.24 4.34L16.36 2.22L17.77 3.63L15.65 5.75L15.65 5.75ZM4.34 5.75L2.22 3.63L3.63 2.22L5.75 4.34L4.34 5.75L4.34 5.75ZM11 3L9 3L9 0L11 0L11 3Z" fillOpacity="1.000000" fillRule="nonzero"/>
                    </svg>
                </div>
                <div
                    className='2xl:text-lg xl:text-base lg:text-sm 2k:text-xl 4k:text-2xl
                        mobile:text-base md:text-lg'
                > 
                    { t('Theme') }
                </div>
            </div>
            
            <Select 
                className='sm:w-3/5 mobile:w-6/12 2xl:text-lg xl:text-base lg:text-sm 2k:text-xl 4k:text-2xl
                    mobile:text-base md:text-lg z-40' 
                options={themeOptions} 
                onChange={handleThemeChange} 
                defaultValue={themeOptions.find(opt => opt.key === themeStore.theme)} 
            />
        </div>
    )
})

export default ThemeChanging