import { Navigate } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { useAuth } from '../../entities/user'

const ActivateRedirect = observer(({ children }) => {
    const userStore = useAuth()

    if (userStore.user.is_activated) {
        return <Navigate to='/' />
    }

    return (
        <>{ children }</>
    )
})

export default ActivateRedirect
