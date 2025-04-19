import { MainWidget } from '../../../shared/ui/MainWidget'
import { SocialUser } from '../../../entities/SocialUser'
import { Dispatch, SetStateAction } from 'react'
import { IShortUser } from '../../../entities/ShortUser'

interface UserWidgetProps {
    key?: number | string,
    title?: string,
    goBack: () => void,
    className?: string,
    setUserManipulation: Dispatch<SetStateAction<boolean>>,
    selectedUser: IShortUser,
    checkAfterUpdate?: boolean,
    onLoadError: () => void,
    onMessageClick: (event: React.MouseEvent<HTMLButtonElement>) => void
}

const UserWidget: React.FC<UserWidgetProps> = ({ 
    key, 
    selectedUser, 
    checkAfterUpdate = true, 
    setUserManipulation, 
    onLoadError, 
    goBack, 
    className, 
    onMessageClick 
}) => {
    return (
        <MainWidget 
            key={key}
            title={ selectedUser.username } 
            goBack={goBack}
            initialAnimation={{ opacity: 0 }}
            animation={{ opacity: 1 }}
            exitAnimation={{ opacity: 0 }}
            className={className}
        >
            <SocialUser 
                id={selectedUser.userId} 
                onError={onLoadError} 
                setUserManipulation={checkAfterUpdate ? setUserManipulation : () => {}} 
                onMessageClick={onMessageClick} 
            />
        </MainWidget>
    )
}

export default UserWidget
