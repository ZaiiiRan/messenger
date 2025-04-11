import { Input } from '../../../shared/ui/Input'
import { Button } from '../../../shared/ui/Button'
import { Loader } from '../../../shared/ui/Loader'
import { LinkButton } from '../../../shared/ui/LinkButton'
import { observer } from 'mobx-react'
import { useAuth } from '../../../entities/user'
import { useTranslation } from 'react-i18next'

type Position = 'first' | 'second' | 'third' | 'fourth' | 'fifth' | 'sixth'

interface ActivationAccountProps {
    data: {
        first: string;
        second: string;
        third: string;
        fourth: string;
        fifth: string;
        sixth: string;
    },
    refs: Array<React.Ref<HTMLInputElement>>,
    handleChange: (e: React.ChangeEvent<HTMLInputElement>, position: string) => void,
    handleBackspace: (e: React.KeyboardEvent<HTMLInputElement>, position: string) => void,
    err: {
        first: boolean;
        second: boolean;
        third: boolean;
        fourth: boolean;
        fifth: boolean;
        sixth: boolean;
    },
    submit: (e: React.MouseEvent<HTMLButtonElement>) => void,
    resend: (e: React.MouseEvent<HTMLAnchorElement>) => void
}

const ActivationAccount: React.FC<ActivationAccountProps> = observer(({ refs, data, handleChange, handleBackspace, err, submit, resend }) => {
    const { t } = useTranslation('activationFeature')
    const userStore = useAuth()

    return (
        <form autoComplete='off' className='flex flex-col lg:w-1/2 xl:w-1/3 mobile:w-full sm:px-28 mobile:px-12 md:px-36 lg:px-0 mobile:gap-14 md:gap-12 2k:gap-16 4k:gap-24'>
            <div className='flex flex-col gap-3 2k:gap-6 4k:gap-10'>
                <h1 
                    className='text-center font-extrabold 
                        md:text-3xl mobile:text-2xl 2k:text-5xl 4k:text-7xl'
                >
                    { t('Enter code') }
                </h1>
                <h2 
                    className='text-center font-extrabold 
                        md:text-lg mobile:text-base 2k:text-2xl 4k:text-4xl'
                >
                    { t('We have sent an email with an account activation code to the email you provided') }
                </h2>
            </div>

            <div className='flex 2xl:px-8 2k:px-4 mobile:px-0 w-full items-center justify-between'>
            {(['first', 'second', 'third', 'fourth', 'fifth', 'sixth'] as const).map((position, index) => (
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
                <div>{ t('Didn\'t receive the code?') }</div>
                <LinkButton onClick={resend}>{ t('Send again') }</LinkButton>
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
                        t('Confirm')
                    )
                }
            </Button>
        </form>
    )
})

export default ActivationAccount
