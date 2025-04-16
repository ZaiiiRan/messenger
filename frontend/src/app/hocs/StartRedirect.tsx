import { Navigate } from 'react-router-dom'
import { observer } from 'mobx-react-lite'
import { useAuth } from '../../entities/user'
import HocProps from './HocProps'

const StartRedirect: React.FC<HocProps> = observer(({ children }) => {
    const userStore = useAuth()

    if (userStore.isAuth && !userStore.isBegin) {
        return <Navigate to='/' />
    }

    return (
        <>{ children }</>
    )
})

export default StartRedirect
