import { Button } from '../../../shared/ui/Button'
import './StartPage.css'

const StartPage = () => {
    return (
        <div className='w-full_screen h-full_screen flex flex-col items-center justify-center 
            2xl:gap-24 xl:gap-20 lg:gap-24 md:gap-36 sm:gap-48 mobile:gap-56 2k:gap-36 4k:gap-56'
        >
            <div className='md:w-1/3 mobile:w-1/2 flex flex-col items-center 
                2xl:gap-9 xl:gap-7 lg:gap-6 md:gap-9 sm:gap-9 mobile:gap-10 2k:gap-14 4k:gap-20'
            >
                <div className='flex justify-center'>
                    <img className='2xl:w-1/3 mobile:w-1/2 4k:w-1/2 shake' src="./message.svg" alt="message" draggable={false}/>
                </div>
                <div className='font-bold text-3xl mobile:text-2xl 2k:text-4xl 4k:text-6xl text-center'>
                    Начните переписываться со своей семьей и друзьями прямо сейчас
                </div>
            </div>
            
            <Button className='sm:w-1/4 mobile:w-1/2 h-14 2k:h-20 4k:h-32 rounded-3xl font-semibold 
                md:text-lg sm:text-sm 2k:text-2xl 4k:text-4xl'
            >
                Начать общение
            </Button>
        </div>
    )
}

export default StartPage
