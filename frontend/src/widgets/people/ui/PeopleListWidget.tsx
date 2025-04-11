import { MainWidget } from '../../../shared/ui/MainWidget'
import { useTranslation } from 'react-i18next'
import { Input } from '../../../shared/ui/Input'
import { Dispatch, SetStateAction, useState } from 'react'
import { PeopleList } from '../../../features/people'
import { AxiosResponse } from 'axios'

interface PeopleListWidgetProps {
    key: number | string,
    title: string,
    goBack: () => void,
    fetchFunction: (search: string, limit: number, offset: number) => Promise<AxiosResponse<any, any>>,
    setSelectedUser: Dispatch<SetStateAction<null>>,
    minSearchLength: number,
    className?: string,
    initialAnimation?: { opcaity: number; x: number},
    exitAnimation?: { opcaity: number; x: number},
    animation?: { opcaity: number; x: number},
    userManipulation: boolean,
    setUserManipulation: Dispatch<SetStateAction<boolean>>,
    selectedUser: any
}

const PeopleListWidget: React.FC<PeopleListWidgetProps> = ({
    key, 
    title, 
    goBack, 
    fetchFunction, 
    setSelectedUser, 
    minSearchLength = 0, 
    className,
    initialAnimation={ opacity: 0, x: -500 }, 
    animation={opacity: 1, x: 0 }, 
    exitAnimation={ opacity: 0, x: -500 },
    userManipulation, 
    setUserManipulation,
    selectedUser
}) => {
    const { t } = useTranslation('peopleWidget')
    const [search, setSearch] = useState('')

    return (
        <MainWidget 
            key={key} 
            title={ t(title) } 
            goBack={ goBack } 
            className={className}
            initialAnimation={initialAnimation}
            animation={animation}
            exitAnimation={exitAnimation}
        >
            <div className='flex flex-col items-center'>
                <Input 
                    className='px-3 py-2 2k:px-4 2k:py-3 4k:px-6 4k:py-5 rounded-lg 
                        md:text-lg mobile:text-sm lg:text-sm 2xl:text-lg 2k:text-2xl 4k:text-4xl sm:w-2/3 mobile:w-full lg:w-full 2xl:w-2/3'
                    placeholder={ t('Username or email') }
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                />
            </div>
            <PeopleList 
                search={search} 
                fetchFunction={fetchFunction}
                setSelectedUser={setSelectedUser}
                minSearchLength={minSearchLength}
                selectedUser={selectedUser}
                userManipulation={userManipulation}
                setUserManipulation={setUserManipulation}
            />
        </MainWidget>
    )
}

export default PeopleListWidget
