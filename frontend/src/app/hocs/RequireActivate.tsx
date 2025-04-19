import { Navigate } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { useAuth } from '../../entities/user'
import HocProps from './HocProps'

const RequireActivate: React.FC<HocProps> = observer(({ children }) => {
    const userStore = useAuth()

    if (!userStore.user?.isActivated) {
        return <Navigate to='/activate' />
    }
    
    return (
        <>{ children }</>
    )
})

export default RequireActivate