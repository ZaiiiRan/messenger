import { Navigate } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { useAuth } from '../../entities/user'

const StartRedirect = observer(({ children }) => {
    const userStore = useAuth()

    if (userStore.isAuth && !userStore.isBegin) {
        return <Navigate to='/' />
    }

    return (
        <>{ children }</>
    )
})

export default StartRedirect
