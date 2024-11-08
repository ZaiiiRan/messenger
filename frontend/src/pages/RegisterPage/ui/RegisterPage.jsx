import { motion, AnimatePresence } from 'framer-motion'
import { useState } from 'react'
import { StepAdditionalInfoRegister, StepEmailUsername, StepNames, StepPassword } from '../../../features/registration/index'
import { validateEmail, validateFirstName, validateLastName, validateUsername, validatePhone, validatePassword, validateBirthdate } from '../../../entities/user'
import { useModal } from '../../../features/modal'
import { useAuth } from '../../../entities/user' 

const RegisterPage = () => {
    const [step, setStep] = useState(1)
    const [data, setData] = useState({ username: '', email: '', firstname: '', lastname: '', password: '', repeatPassword: '', birthdate: '', phone: ''})
    const [err, setErr] = useState({ username: false, email: false, firstname: false, lastname: false, password: false, repeatPassword: false, birthdate: false, phone: false})
    const { openModal, setModalTitle, setModalText } = useModal()
    const userStore = useAuth()

    const validateStep = (stepData, validationFunctions) => {
        let isValid = true
        let newErrors = {}
        
        validationFunctions.forEach(({ field, validate }) => {
            const result = validate(stepData[field].trim())
            if (!result.valid) {
                newErrors[field] = true
                isValid = false
                setModalTitle('Ошибка')
                setModalText(result.message)
                openModal()
            } else {
                newErrors[field] = false
            }
        })
        
        setErr((prevErr) => ({ ...prevErr, ...newErrors }))
        return isValid
    }

    const handleNext = async (e, validators) => {
        e.preventDefault()
        if (validators && !validateStep(data, validators)) return
        if (step < 4) setStep(step + 1)
        else await handleRegister()
    }

    const handlePrev = (e) => {
        e.preventDefault()
        setStep(step - 1)
    }

    const handleRegister = async () => {
        try {
            userStore.setLoading(true)
            await userStore.register(data.username, data.email, data.password, data.firstname, data.lastname, data.phone, data.birthdate)
        } catch (e) {
            console.log(e)
            setModalTitle('Ошибка')
            setModalText(e.response?.data?.error || 'Внутренняя ошибка сервера')
            openModal()
        } finally {
            userStore.setLoading(false)
        }
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
                            onNext={(e) => handleNext(e, [
                                { field: 'lastname', validate: validateLastName },
                                { field: 'firstname', validate: validateFirstName },
                            ])}
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
                            onNext={(e) => handleNext(e, [
                                { field: 'username', validate: validateUsername },
                                { field: 'email', validate: validateEmail }
                            ])}
                            onPrev={handlePrev}
                            email={data.email}
                            setEmail={(e) => setData({ ...data, email: e.target.value })}
                            emailErr={err.email}
                            username={data.username}
                            setUsername={(e) => setData({ ...data, username: e.target.value })}
                            usernameErr={err.username}
                        />
                    )
                }

                {
                    step === 3 && (
                        <StepPassword
                            onNext={(e) => handleNext(e, [
                                { field: 'repeatPassword', validate: (value) => ({
                                    valid: value === data.password,
                                    message: 'Пароли не совпадают'
                                }) },
                                { field: 'password', validate: validatePassword }
                            ])}
                            onPrev={handlePrev}
                            password={data.password}
                            setPassword={(e) => setData({ ...data, password: e.target.value })}
                            passwordErr={err.password}
                            repeatPassword={data.repeatPassword}
                            setRepeatPassword={(e) => setData({ ...data, repeatPassword: e.target.value })}
                            repeatPasswordErr={err.repeatPassword}
                        />
                    )
                }

                {
                    step === 4 && (
                        <StepAdditionalInfoRegister
                            onPrev={handlePrev}
                            onNext={ (e) => handleNext(e, [
                                { field: 'birthdate', validate: validateBirthdate },
                                { field: 'phone', validate: validatePhone },
                            ])}
                            phone={data.phone}
                            setPhone={(e) => setData({ ...data, phone: e.target.value })}
                            phoneErr={err.phone}
                            birthdate={data.birthdate}
                            setBirthdate={(e) => setData({ ...data, birthdate: e.target.value })}
                            birthdateErr={err.birthdate}
                        />
                    )
                }
        </motion.div>
        </AnimatePresence>
    )
}

export default RegisterPage
