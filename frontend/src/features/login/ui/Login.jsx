/* eslint-disable react/prop-types */
import { Input } from '../../../shared/ui/Input'
import { Button } from '../../../shared/ui/Button'
import { Link } from '../../../shared/ui/Link'

const Login = ({ login, setLogin, loginErr, password, setPassword, passwordErr, onLogin }) => {
    return (
        <form className='flex flex-col lg:w-1/3 mobile:w-1/2 mobile:gap-14 md:gap-12 2k:gap-16 4k:gap-24'>
            <h1 
                className='text-center font-extrabold 
                    md:text-3xl mobile:text-2xl 2k:text-5xl 4k:text-7xl'
            >
                Авторизация
            </h1>
            <Input 
                placeholder='Логин' 
                className='px-3 py-2 2k:px-4 2k:py-3 4k:px-6 4k:py-5 rounded-lg 
                    md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                value={login}
                onChange={setLogin}
                error={loginErr}
            />
            <Input 
                placeholder='Пароль' 
                className='px-3 py-2 rounded-lg 2k:px-4 2k:py-3 4k:px-6 4k:py-5
                    md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl' 
                password={true}
                value={password}
                onChange={setPassword}
                error={passwordErr}
            />
            <div 
                className='flex md:gap-4 items-center 
                    mobile:gap-2 md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
            >
                <div>Нет аккаунта?</div>
                <Link to='/register'>Регистрация</Link>
            </div>
            
            <Button 
                className='h-14 2k:h-20 4k:h-32 rounded-3xl font-semibold 
                    md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
                onClick={onLogin}
            >
                Войти
            </Button>
        </form>
    )
}

export default Login
