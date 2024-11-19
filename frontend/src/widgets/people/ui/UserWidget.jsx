/* eslint-disable react/prop-types */
import { MainWidget } from '../../../shared/ui/MainWidget'
import { SocialUser } from '../../../entities/SocialUser'

const UserWidget = ({ key, selectedUser, checkAfterUpdate = true, setUserManipulation, onLoadError, goBack, className }) => {
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
            <SocialUser id={selectedUser.user_id} onError={onLoadError} setUserManipulation={checkAfterUpdate ? setUserManipulation : () => {}} />
        </MainWidget>
    )
}

export default UserWidget
