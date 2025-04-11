import { motion } from 'framer-motion'
import { useAuth } from '../../../entities/user'
import { useState, useRef } from 'react'
import { useModal } from '../../../features/modal'
import { ActivationAccount } from '../../../features/activation/'
import { useTranslation } from 'react-i18next'
import { apiErrors, ApiErrorsKey, apiMessages, ApiMessagesKey } from '../../../shared/api'

type DataKeys = 'first' | 'second' | 'third' | 'fourth' | 'fifth' | 'sixth'

const ActivationPage = () => {
    const { t } = useTranslation('activationPage')
    const userStore = useAuth()
    
    const [data, setData] = useState<{ [key in DataKeys]: string }>({ first: '', second: '', third: '', fourth: '', fifth: '', sixth: '' })
    const [err, setErr] = useState({ first: false, second: false, third: false, fourth: false, fifth: false, sixth: false })

    const { openModal, setModalTitle, setModalText } = useModal()

    const firstRef = useRef<HTMLInputElement>(null)
    const secondRef = useRef<HTMLInputElement>(null)
    const thirdRef = useRef<HTMLInputElement>(null)
    const fourthRef = useRef<HTMLInputElement>(null)
    const fifthRef = useRef<HTMLInputElement>(null)
    const sixthRef = useRef<HTMLInputElement>(null)

    const refs = [firstRef, secondRef, thirdRef, fourthRef, fifthRef, sixthRef]

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>, position: string) => {
        const value = e.target.value.slice(-1)
        setData((prev) => ({ ...prev, [position]: value }))

        if (value && position !== 'sixth') {
            const nextIndex = refs.findIndex(ref => ref.current?.name === position) + 1
            refs[nextIndex].current?.focus()
        }
    }

    const handleBackspace = (e: React.KeyboardEvent<HTMLInputElement>, position: string) => {
        const pos = position as DataKeys
        if (e.key === 'Backspace' && !data[pos] && position !== 'first') {
            const prevIndex = refs.findIndex(ref => ref.current?.name === position) - 1
            refs[prevIndex].current?.focus()
        }
    }

    const validate = () => {
        let newErr = { first: false, second: false, third: false, fourth: false, fifth: false, sixth: false }
        let hasErr = false
        if (data.first === '') {
            newErr.first = true
            hasErr = true
        }
        if (data.second === '') {
            newErr.second = true
            hasErr = true
        }
        if (data.third === '') {
            newErr.third = true
            hasErr = true
        }
        if (data.fourth === '') {
            newErr.fourth = true
            hasErr = true
        }
        if (data.fifth === '') {
            newErr.fifth = true
            hasErr = true
        }
        if (data.sixth === '') {
            newErr.sixth = true
            hasErr = true
        }
        setErr(newErr)
        if (hasErr) {
            setModalTitle(t('Error'))
            setModalText(t('The code is not completely filled out'))
            openModal()
        }
        return !hasErr
    }

    const submit = async (e: React.MouseEvent<HTMLButtonElement>) => {
        e.preventDefault()
        if (!validate()) return

        const code = data.first + data.second + data.third + data.fourth + data.fifth + data.sixth
        try {
            userStore.setLoading(true)
            await userStore.activate(code)
        } catch (e: any) {
            setModalTitle(t('Error'))

            const errorKey: ApiErrorsKey = e.response?.data?.error
            setModalText(t(apiErrors[errorKey]) || t('Internal server error'))
            openModal()
        } finally {
            userStore.setLoading(false)
        }
    }

    const resend = async (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault()
        try {
            const response = await userStore.resend()
            console.log
            setModalTitle(t('Activation code'))

            const messageKey: ApiMessagesKey = response?.message
            setModalText(t(apiMessages[messageKey]))
            openModal()
        } catch (e: any) {
            setModalTitle(t('Error'))

            const errorKey: ApiErrorsKey = e.response?.data?.error
            setModalText(t(apiErrors[errorKey]) || t('Internal server error'))
            openModal()
        }
    }

    return (
        <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.7 }}
            className='w-full_screen h-full_screen flex flex-col items-center justify-center'
        >
            <ActivationAccount 
                refs={refs}
                data={data}
                handleChange={handleChange}
                handleBackspace={handleBackspace}
                err={err}
                submit={submit}
                resend={resend}
            />
        </motion.div>
    )
}

export default ActivationPage
