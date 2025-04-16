import { Outlet } from 'react-router-dom'
import { Navigation } from '../../features/navigation'
import './MainLayout.css'
import { UserModal } from '../../entities/user'

const MainLayout: React.FC = () => {
    return (
        <div >
            <Navigation />
            <main>
                <Outlet />
            </main>
            <UserModal />
        </div>
    )
}

export default MainLayout
