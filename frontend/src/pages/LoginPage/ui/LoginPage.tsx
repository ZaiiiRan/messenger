import { motion } from 'framer-motion'
import { useState } from 'react'
import { useModal } from '../../../features/modal'
import { Login } from '../../../features/login'
import { useAuth } from '../../../entities/user' 
import { useTranslation } from 'react-i18next'
import { apiErrors, ApiErrorsKey } from '../../../shared/api'

const LoginPage = () => {
    const { t } = useTranslation('loginPage')

    const [data, setData] = useState<{login: string, password: string}>({ login: '', password: '' })
    const [err, setErr] = useState<{login: boolean, password: boolean}>({ login: false, password: false })

    const { openModal } = useModal()

    const userStore = useAuth()

    const proccessValidateErrors = (errors: {login: boolean, password: boolean}) => {
        let errMsg = ''
        if (errors.login && errors.password) {
            errMsg = t('Enter your login and password')
        } else if (errors.login) {
            errMsg = t('Enter your login')
        } else if (errors.password) {
            errMsg = t('Enter your password')
        } else {
            return false
        }
        openModal(t('Error'), errMsg)
        return true
    }

    const handleLogin = async (e: React.MouseEvent<HTMLButtonElement>) => {
        e.preventDefault()

        let newErrors = { ...err }
        if (data.login.trim() === '') {
            newErrors.login = true
        } else {
            newErrors.login = false
        }
        if (data.password.trim() === '') {
            newErrors.password = true
        } else {
            newErrors.password = false
        }
        setErr(newErrors)
        if (proccessValidateErrors(newErrors)) {
            return
        }

        try {
            userStore.setLoading(true)
            await userStore.login(data.login, data.password)
        } catch (e: any) {
            console.log(e)
            const errorKey: ApiErrorsKey = e.response?.data?.error
            const errMsg = t(apiErrors[errorKey]) || t('Internal server error')
            openModal(t('Error'), errMsg)
        } finally {
            userStore.setLoading(false)
        }
    }

    return (
        <motion.div 
            initial={{ opacity: 0, y: -1000 }}
            animate={{ opacity: 1, y: 0}}
            exit={{ opacity: 0, y: 1000 }}
            transition={{ duration: 0.7 }}
            className='w-full_screen h-full_screen flex flex-col items-center justify-center'
        >
            <Login
                login={data.login}
                setLogin={(e) => setData({ ...data, login: e.target.value })}
                loginErr={err.login}
                password={data.password}
                setPassword={(e) => setData({ ...data, password: e.target.value })}
                passwordErr={err.password}
                onLogin={handleLogin}
            />
        </motion.div>
    )
}

export default LoginPage
