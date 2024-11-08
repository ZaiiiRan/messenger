/* eslint-disable react/prop-types */
import { Input } from '../../../shared/ui/Input'
import { Button } from '../../../shared/ui/Button'
import { Link } from '../../../shared/ui/Link'
import { motion } from 'framer-motion'
import { useTranslation } from 'react-i18next'

const StepAdditionalInfoRegister = ({ onNext, onPrev, phone, setPhone, phoneErr, birthdate, setBirthdate, birthdateErr }) => {
    const { t } = useTranslation('registerFeature')
    
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
            className='flex flex-col lg:w-1/2 xl:w-1/3 mobile:w-full sm:px-28 mobile:px-12 md:px-36 lg:px-0 mobile:gap-14 md:gap-12 2k:gap-16 4k:gap-24'
            onKeyDown={handleFormKeyDown}
        >
            <div className='flex flex-col gap-3 2k:gap-6 4k:gap-10'>
                <h1 
                    className='text-center font-extrabold 
                        md:text-3xl mobile:text-2xl 2k:text-5xl 4k:text-7xl'
                >
                    { t('Registration') }
                </h1>
                <h2 
                    className='text-center font-extrabold 
                        md:text-lg mobile:text-base 2k:text-2xl 4k:text-4xl'
                >
                    { t('You can provide additional information') }
                </h2>
            </div>

            <Input 
                placeholder={t('Phone number (optional)')} 
                className='px-3 py-2 2k:px-4 2k:py-3 4k:px-6 4k:py-5 rounded-lg 
                    md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                value={phone}
                onChange={setPhone}
                error={phoneErr}
                phone={true}
            />
            <Input 
                placeholder={t('Date of birth (optional)')}
                className='px-3 py-2 rounded-lg 2k:px-4 2k:py-3 4k:px-6 4k:py-5
                    md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl' 
                value={birthdate}
                onChange={setBirthdate}
                error={birthdateErr}
                date={true}
            />
            <div 
                className='flex md:gap-4 items-center 
                    mobile:gap-2 md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
            >
                <div>{ t('Already have an account?') }</div>
                <Link to='/login'>{ t('Login') }</Link>
            </div>

            <div className='flex items-center justify-between mobile:gap-3 md:gap-9 gap-9 flex-wrap'>
                <Button 
                    className='flex-grow h-14 2k:h-20 4k:h-32 rounded-3xl font-semibold 
                        md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                    onClick={onPrev}
                >
                    { t('Back') }
                </Button>
                <Button 
                    className='flex-grow h-14 2k:h-20 4k:h-32 rounded-3xl font-semibold 
                        md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                    onClick={onNext}
                    id='submitBtn'
                >
                    { t('Register') } 
                </Button>
            </div>
        </motion.form>
    )
}

export default StepAdditionalInfoRegister
