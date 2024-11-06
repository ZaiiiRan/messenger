import { motion, AnimatePresence } from 'framer-motion'
import { useState } from 'react'
import { StepAdditionalInfoRegister, StepEmailUsername, StepNames, StepPassword } from '../../../features/registration/index'

const RegisterPage = () => {
    const [step, setStep] = useState(1)

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
                            onNext={handleNext}
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
