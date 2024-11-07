/* eslint-disable react/prop-types */
import { Input } from '../../../shared/ui/Input'
import { Button } from '../../../shared/ui/Button'
import { Link } from '../../../shared/ui/Link'
import { motion } from 'framer-motion'

const StepPassword = ({ onNext, onPrev, password, setPassword, passwordErr, repeatPassword,  setRepeatPassword, repeatPasswordErr }) => {
    const handleFormKeyDown = (e) => {
        if (e.key === 'Enter') {
            e.preventDefault()
            document.getElementById('submitBtn').click()
        }
    }

    return (
        <motion.form 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1}}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.3 }}
            className='flex flex-col lg:w-1/3 mobile:w-1/2 mobile:gap-14 md:gap-12 2k:gap-16 4k:gap-24'
            onKeyDown={handleFormKeyDown}
        >
            <div className='flex flex-col gap-3 2k:gap-6 4k:gap-10'>
                <h1 
                    className='text-center font-extrabold 
                        md:text-3xl mobile:text-2xl 2k:text-5xl 4k:text-7xl'
                >
                    Регистрация
                </h1>
                <h2 
                    className='text-center font-extrabold 
                        md:text-lg mobile:text-base 2k:text-2xl 4k:text-4xl'
                >
                    Придумайте пароль
                </h2>
            </div>

            <Input 
                placeholder='Пароль' 
                className='px-3 py-2 2k:px-4 2k:py-3 4k:px-6 4k:py-5 rounded-lg 
                    md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                password={true}
                value={password}
                onChange={setPassword}
                error={passwordErr}
            />
            <Input 
                placeholder='Повторите пароль' 
                className='px-3 py-2 rounded-lg 2k:px-4 2k:py-3 4k:px-6 4k:py-5
                    md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl' 
                    password={true}
                value={repeatPassword}
                onChange={setRepeatPassword}
                error={repeatPasswordErr}
            />
            <div 
                className='flex md:gap-4 items-center 
                    mobile:gap-2 md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
            >
                <div>Уже есть аккаунт?</div>
                <Link to='/login'>Войти</Link>
            </div>

            <div className='flex items-center justify-between mobile:gap-3 md:gap-9 gap-9 flex-wrap'>
                <Button 
                    className='flex-grow h-14 2k:h-20 4k:h-32 rounded-3xl font-semibold 
                        md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                    onClick={onPrev}
                >
                    Назад
                </Button>
                <Button 
                    className='flex-grow h-14 2k:h-20 4k:h-32 rounded-3xl font-semibold 
                        md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                    onClick={onNext}
                    id='submitBtn'
                >
                    Далее
                </Button>
            </div>
        </motion.form>
    )
}

export default StepPassword
