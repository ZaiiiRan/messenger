/* eslint-disable react/prop-types */
/* eslint-disable react-hooks/exhaustive-deps */
import { Select } from '../../../shared/ui/Select'
import 'flag-icons/css/flag-icons.min.css'
import { useTranslation } from 'react-i18next'
import { useState, useEffect, useMemo } from 'react'
import { motion } from 'framer-motion'

const AppearanceOptions = ({ goBack }) => {
    const { i18n } = useTranslation()
    const { t } = useTranslation('optionsFeature')

    const langs = useMemo(() =>  [
        {
            key: 'ru',
            label: (
                <div className='flex items-center gap-5'>
                    <span className="fi fi-ru"></span>
                    <div>Русский</div>
                </div>
            )
        },
        {
            key: 'en',
            label: (
                <div className='flex items-center gap-5'>
                    <span className="fi fi-us"></span>
                    <div>English (US)</div>
                </div>
            )
        }
    ], [t])

    const defaultLang = langs.find(lang => lang.key === i18n.language) || langs[0]

    const handleLanguageChange = (selectedLang) => {
        i18n.changeLanguage(selectedLang.key)
    }

    const themeOptions = useMemo(() => [
        { key: 'light', label: t('Light') },
        { key: 'dark', label: t('Dark') },
        { key: 'system', label: t('System') }
    ], [t])

    const [theme, setTheme] = useState('system')

    const handleThemeChange = (selectedTheme) => {
        setTheme(selectedTheme.key)
        updateThemeClass(selectedTheme.key)
    }

    const updateThemeClass = (selectedTheme) => {
        document.body.classList.remove('light-theme', 'dark-theme')
        if (selectedTheme === 'light') {
            document.body.classList.add('light-theme')
        } else if (selectedTheme === 'dark') {
            document.body.classList.add('dark-theme')
        }
    }

    useEffect(() => {
        if (theme === 'system') {
            const mediaQuery = window.matchMedia('(prefers-color-scheme: light)')
            const systemThemeListener = (e) => {
                document.body.classList.toggle('light-theme', e.matches)
                document.body.classList.toggle('dark-theme', !e.matches)
            };

            mediaQuery.addEventListener('change', systemThemeListener);
            systemThemeListener(mediaQuery)

            return () => mediaQuery.removeEventListener('change', systemThemeListener)
        }
    }, [theme])

    return (
        <motion.div 
            initial={{ opacity: 0, x: -500 }}
            animate={{ opacity: 1, x: 0}}
            exit={{ opacity: 0, x: -500 }}
            transition={{ duration: 0.3 }}
            className='Options-Widget rounded-3xl flex flex-col gap-8 2k:gap-14 4k:gap-24'
        >
            {/* Title */}
            <div className='flex items-center gap-5 2k:gap-7 4k:gap-9'>
                <div 
                    className='backBtn 2xl:w-10 2xl:h-10 xl:w-9 xl:h-9 lg:w-9 lg:h-8 2k:w-12 2k:h-12 4k:w-14 4k:h-14 
                        mobile:w-8 mobile:h-8 md:w-9 md:h-9 
                        rounded-3xl flex items-center justify-center'
                    onClick={goBack}
                >
                    <div className='Option-Icon flex items-center justify-center h-1/4 aspect-square'>
                        <svg viewBox="0 0 7.424 12.02" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <defs/>
                            <path id="Vector" d="M0 6.01L6 12.02L7.42 10.6L2.82 6L7.42 1.4L6 0L0 6.01Z" fillOpacity="1.000000" fillRule="nonzero"/>
                        </svg>
                    </div>
                </div>
                <h1 className='font-extrabold 2xl:text-3xl xl:text-2xl lg:text-xl 2k:text-4xl 4k:text-5xl
                    md:text-3xl sm:text-2xl mobile:text-xl'
                >
                    { t('Appearance') }
                </h1>
            </div>
            
            {/* Options */}
            <div className='flex flex-col gap-8 2k:gap-12 4k:gap-20'>
                <div className='flex 2xl:w-1/2 xl:w-4/6 justify-between'>
                    <div className='flex items-center gap-3'>
                        <div className='Option-Icon h-1/2'>
                            <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" strokeWidth="0"></g><g id="SVGRepo_tracerCarrier" strokeLinecap="round" strokeLinejoin="round"></g><g id="SVGRepo_iconCarrier"> <path fillRule="evenodd" clipRule="evenodd" d="M12.6247 5.21914C11.0416 3.95267 9.18319 3.99214 7.85649 4.27145C7.04755 4.44175 6.22551 4.72012 5.50386 5.13176C5.19228 5.30981 5 5.64115 5 6.00001V20C5 20.5523 5.44772 21 6 21C6.55228 21 7 20.5523 7 20V14.6294C8.37617 14.0493 10.124 13.7799 11.3753 14.7809C12.9584 16.0473 14.8168 16.0079 16.1435 15.7286C17.1559 15.5154 17.9441 15.1521 18.2954 14.9747C18.6869 14.7771 19 14.4734 19 14V6.00001C19 5.64353 18.8102 5.31402 18.5019 5.13509C18.1938 4.95629 17.8144 4.95462 17.505 5.13114L17.5041 5.13162C16.0661 5.91734 14.0013 6.32045 12.6247 5.21914ZM7 6.62938V12.499C8.88136 11.8968 11.021 11.9362 12.6247 13.2191C13.5416 13.9527 14.6832 13.9922 15.7315 13.7715C16.2336 13.6657 16.6769 13.5068 17 13.3706V7.50105C16.739 7.5846 16.4511 7.6638 16.1435 7.72856C14.8168 8.00787 12.9584 8.04734 11.3753 6.78087C10.124 5.77986 8.37617 6.04932 7 6.62938Z"></path> </g></svg>
                        </div>
                        <div 
                            className='2xl:text-lg xl:text-base lg:text-sm 2k:text-xl 4k:text-2xl
                                mobile:text-base md:text-lg'
                        >
                            { t('Language') }
                        </div>
                    </div>
                    
                    <Select 
                        className='sm:w-3/5 mobile:w-6/12 2xl:text-lg xl:text-base lg:text-sm 2k:text-xl 4k:text-2xl
                            mobile:text-base md:text-lg' 
                        options={langs} 
                        onChange={handleLanguageChange} 
                        defaultValue={defaultLang} 
                    />
                </div>
                <div className='flex 2xl:w-1/2 xl:w-4/6 justify-between'>
                    <div className='flex items-center gap-3'>
                        <div className='Option-Icon h-1/2'>
                            <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
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
                            mobile:text-base md:text-lg' 
                        options={themeOptions} 
                        onChange={handleThemeChange} 
                        defaultValue={themeOptions.find(opt => opt.key === theme)} 
                    />
                </div>
                
            </div>
                
        </motion.div>
    )
}

export default AppearanceOptions
