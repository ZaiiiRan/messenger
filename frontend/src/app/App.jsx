/* eslint-disable react-refresh/only-export-components */
/* eslint-disable react-hooks/exhaustive-deps */
import { Modal } from '../features/modal'
import { Router } from './routers/Router'
import { useAuth } from '../entities/user'
import { useEffect, useState } from 'react'
import { observer } from 'mobx-react'
import { Loader } from '../shared/ui/Loader'


function App() {
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
      <Router />
      <Modal />
    </>
    
  )
}

export default observer(App)
