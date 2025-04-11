import { Navigate } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { useAuth } from '../../entities/user'
import HocProps from './HocProps'

const RequireAuth: React.FC<HocProps> = observer(({ children }) => {
    const userStore = useAuth()

    if (!userStore.isAuth) {
        return <Navigate to='/start' />
    }
    
    return (
        <>{ children }</>
    )
})

export default RequireAuth
