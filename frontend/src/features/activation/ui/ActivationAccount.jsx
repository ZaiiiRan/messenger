import { Input } from '../../../shared/ui/Input'
import { Button } from '../../../shared/ui/Button'
import { Loader } from '../../../shared/ui/Loader'
import { LinkButton } from '../../../shared/ui/LinkButton'
import { observer } from 'mobx-react'
import { useAuth } from '../../../entities/user'

const ActivationAccount = observer(({ refs, data, handleChange, handleBackspace, err, submit, resend }) => {
    const userStore = useAuth()

    return (
        <form autoComplete='off' className='flex flex-col lg:w-1/2 xl:w-1/3 mobile:w-full sm:px-28 mobile:px-12 md:px-36 lg:px-0 mobile:gap-14 md:gap-12 2k:gap-16 4k:gap-24'>
            <div className='flex flex-col gap-3 2k:gap-6 4k:gap-10'>
                <h1 
                    className='text-center font-extrabold 
                        md:text-3xl mobile:text-2xl 2k:text-5xl 4k:text-7xl'
                >
                    Введите код
                </h1>
                <h2 
                    className='text-center font-extrabold 
                        md:text-lg mobile:text-base 2k:text-2xl 4k:text-4xl'
                >
                    Мы отправили письмо с кодом активации аккаунта на указанный вами Email
                </h2>
            </div>

            <div className='flex 2xl:px-8 2k:px-4 mobile:px-0 w-full items-center justify-between'>
            {['first', 'second', 'third', 'fourth', 'fifth', 'sixth'].map((position, index) => (
                    <Input 
                        disabled = { userStore.isLoading }
                        key={position}
                        ref={refs[index]}
                        name={position}
                        className='md:w-16 mobile:w-10 sm:w-16 2k:w-24 4k:w-32 text-center font-extrabold sm:p-2 mobile:p-1 rounded-lg 2k:p-4 4k:p-6
                            md:text-lg mobile:text-sm sm:text-base 2k:text-2xl 4k:text-4xl' 
                        value={data[position]}
                        onChange={(e) => handleChange(e, position)}
                        onKeyDown={(e) => handleBackspace(e, position)}
                        error={err[position]}
                        oneDigit
                    />
                ))}
            </div>
            
            <div 
                className='flex md:gap-4 items-center 
                    mobile:gap-2 md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl'
            >
                <div>Не пришел код?</div>
                <LinkButton onClick={resend}>Отправить снова</LinkButton>
            </div>
            
            <Button 
                className='h-14 2k:h-20 4k:h-32 rounded-3xl font-semibold 
                    md:text-lg mobile:text-sm 2k:text-2xl 4k:text-4xl flex items-center justify-center'
                onClick={submit}
                disabled={userStore.isLoading}
            >
                {
                    userStore.isLoading ? (
                        <Loader className='h-3 w-16 2k:h-4 2k:w-24 4k:h-6 4k:w-36'/>
                    ) : (
                        'Подтвердить'
                    )
                }
            </Button>
        </form>
    )
})

export default ActivationAccount
