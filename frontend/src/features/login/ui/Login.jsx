import { Input } from '../../../shared/ui/Input'
import { Button } from '../../../shared/ui/Button'
import { Link } from '../../../shared/ui/Link'
import { Loader } from '../../../shared/ui/Loader'
import { observer } from 'mobx-react'
import { useAuth } from '../../../entities/user'
import { useTranslation } from 'react-i18next'

const Login = observer(({ login, setLogin, loginErr, password, setPassword, passwordErr, onLogin }) => {
    const { t } = useTranslation('loginFeature')
    
    const userStore = useAuth()
    
    return (
        <form className='flex flex-col lg:w-1/3 mobile:w-1/2 mobile:gap-14 md:gap-12 2k:gap-16 4k:gap-24'>
            <h1 
                className='text-center font-extrabold 
                    md:text-3xl mobile:text-2xl 2k:text-5xl 4k:text-7xl'
            >
                { t('Authorization') }
            </h1>
            <Input 
                placeholder={t('Login (username or Email)')}
                className='px-3 py-2 2k:px-4 2k:py-3 4k:px-6 4k:py-5 rounded-lg 
                    md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                value={login}
                onChange={setLogin}
                error={loginErr}
                disabled={userStore.isLoading}
            />
            <Input 
                placeholder={t('Password')}
                className='px-3 py-2 rounded-lg 2k:px-4 2k:py-3 4k:px-6 4k:py-5
                    md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl' 
                password={true}
                value={password}
                onChange={setPassword}
                error={passwordErr}
                disabled={userStore.isLoading}
            />
            <div 
                className='flex md:gap-4 items-center 
                    mobile:gap-2 md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
            >
                <div>{ t('Don\'t have an account?') }</div>
                <Link to='/register'>{ t('Registration') }</Link>
            </div>
            
            <Button 
                className='h-14 2k:h-20 4k:h-32 rounded-3xl font-semibold 
                    md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl flex items-center justify-center'
                onClick={onLogin}
                disabled={userStore.isLoading}
            >
                {
                    userStore.isLoading ? (
                        <Loader className='h-3 w-16 2k:h-4 2k:w-24 4k:h-6 4k:w-36'/>
                    ) : (
                        t('Login') 
                    )
                }
            </Button>
        </form>
    )
})

export default Login
