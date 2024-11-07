import { motion, AnimatePresence } from 'framer-motion'
import { useState } from 'react'
import { StepAdditionalInfoRegister, StepEmailUsername, StepNames, StepPassword } from '../../../features/registration/index'
import { validateEmail, validateFirstName, validateLastName, validateUsername, validatePhone, validatePassword } from '../../../entities/user'
import { useModal } from '../../../features/modal'

const RegisterPage = () => {
    const [step, setStep] = useState(1)
    const [data, setData] = useState({ username: '', email: '', firstname: '', lastname: '', password: '', repeatPassword: '', birthdate: '', phone: ''})
    const [err, setErr] = useState({ username: false, email: false, firstname: false, lastname: false, password: false, repeatPassword: false, birthdate: false, phone: false})
    const { openModal, setModalTitle, setModalText } = useModal()

    const namesConfirm = (e) => {
        e.preventDefault()
        let newErr = {}
        let isErr = false
        let valid = validateFirstName(data.firstname)
        if (!valid.valid) {
            newErr["firstname"] = true
            isErr = true
            setModalTitle('Ошибка')
            setModalText(valid.message)
            openModal()
        } else {
            newErr["firstname"] = false
        }
        valid = validateLastName(data.lastname)
        if (!valid.valid) {
            newErr["lastname"] = true
            if (!isErr) {
                isErr = true
                setModalTitle('Ошибка')
                setModalText(valid.message)
                openModal()
            }
        } else {
            newErr["lastname"] = false
        }
        setErr({...err, ...newErr})

        if (!isErr) {
            handleNext(e)
        }
    }

    const handleNext = (e) => {
        e.preventDefault()
        setStep(step + 1)
    }
    const handlePrev = (e) => {
        e.preventDefault()
        setStep(step - 1)
    }

    return (
        <AnimatePresence mode="wait">
        <motion.div 
            initial={{ opacity: 0, x: -1000 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: 1000 }}
            transition={{ duration: 0.7 }}
            className='w-full_screen h-full_screen flex flex-col items-center justify-center'
        >
                {
                    step === 1 && (
                        <StepNames 
                            onNext={namesConfirm}
                            firstname={data.firstname}
                            setFirstname={(e) => setData({ ...data, firstname: e.target.value })}
                            firstnameErr={err.firstname}
                            lastname={data.lastname}
                            setLastname={(e) => setData({ ...data, lastname: e.target.value})}
                            lastnameErr={err.lastname}
                        />
                    )
                }

                {
                    step === 2 && (
                        <StepEmailUsername
                            onNext={handleNext}
                            onPrev={handlePrev}
                        />
                    )
                }

                {
                    step === 3 && (
                        <StepPassword
                            onNext={handleNext}
                            onPrev={handlePrev}
                        />
                    )
                }

                {
                    step === 4 && (
                        <StepAdditionalInfoRegister
                            onPrev={handlePrev}
                        />
                    )
                }
        </motion.div>
        </AnimatePresence>
    )
}

export default RegisterPage
