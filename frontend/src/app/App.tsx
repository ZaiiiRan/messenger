/* eslint-disable react-refresh/only-export-components */
/* eslint-disable react-hooks/exhaustive-deps */
import { Modal } from '../features/modal'
import { Router } from './routers/Router'
import { useAuth } from '../entities/user'
import { useEffect, useState } from 'react'
import { observer } from 'mobx-react'
import { Loader } from '../shared/ui/Loader'
import { useTranslation } from 'react-i18next'
import '../shared/config/i18n'
import { themeStore } from '../shared/theme'


const App: React.FC = () => {
  useEffect(() => {
    themeStore.applyTheme(themeStore.theme)
  }, [])

  const { t } = useTranslation()
  const [isLandscape, setIsLandscape] = useState(
    window.matchMedia("(orientation: landscape) and (max-height: 599px)").matches
    ||
    window.matchMedia(
      "(orientation: landscape) and (min-width: 1921px) and (max-width: 2560px) and (max-height: 700px)"
    ).matches
    ||
    window.matchMedia(
      "(orientation: landscape) and (min-width: 2561px) and (max-width: 3840px) and (max-height: 1300px)"
    ).matches
  )
  const [isMobile, setIsMobile] = useState(false)
  useEffect(() => {
    const handleOrientationChange = () => {
      const landscapeSmallHeight = window.matchMedia("(orientation: landscape) and (max-height: 599px)").matches
      const landscapeMediumResolution = window.matchMedia(
        "(orientation: landscape) and (min-width: 1921px) and (max-width: 2560px) and (max-height: 700px)"
      ).matches
      const landscapeHighResolution = window.matchMedia(
        "(orientation: landscape) and (min-width: 2561px) and (max-width: 3840px) and (max-height: 1300px)"
      ).matches

      setIsLandscape(landscapeSmallHeight || landscapeMediumResolution || landscapeHighResolution)
    }

    const userAgent = navigator.userAgent || navigator.vendor || (window as any).opera
    const isMobileUserAgent = /android|iphone|ipad|ipod|windows phone|opera mini|iemobile/i.test(userAgent)

    setIsMobile(isMobileUserAgent)

    window.addEventListener("resize", handleOrientationChange)
    return () => window.removeEventListener("resize", handleOrientationChange)
  }, [])


  const userStore = useAuth()
  const [isChecked, setIsChecked] = useState(false)

  useEffect(() => {
    const refresh = async () => {
      try {
        if (localStorage.getItem('token')) {
          await userStore.checkAuth()
        }
      } finally {
        userStore.setBegin(false)
        setIsChecked(true)
      }
    }

    refresh()
  }, [])

  if (!isChecked) {
    return (
      <div className='w-full_screen h-full_screen flex items-center justify-center'>
        <Loader className='h-4 w-1/12 mobile:w-1/4 sm:w-1/6 lg:w-1/12 lg:h-5 xl:h-6 2xl:h-7 2k:h-10 4k:h-14' />
      </div>
    )
  }
  

  return (
    <>
      {
        isLandscape ? (
          <div 
            className='h-full_screen w-full_screen 2xl:text-4xl font-semiboldbold text-center flex items-center justify-center px-10
              xl:text-3xl lg:text-2xl mobile:text-xl'
          >
            {
              isMobile 
              ?
              t('mobile window error')
              :
              t('desktop window error')
            }
          </div>
        ) : (
          <>
            <Router />
            <Modal />
          </>
        )
      }
    </>
    
  )
}

export default observer(App)
