import { Navigate } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { useAuth } from '../../entities/user'

const RequireAuth = observer(({ children }) => {
    const userStore = useAuth()

    if (!userStore.isAuth) {
        return <Navigate to='/start' />
    }
    
    return (
        <>{ children }</>
    )
})

export default RequireAuth
