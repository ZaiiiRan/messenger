import { Outlet } from 'react-router-dom'
import { Navigation } from '../../features/navigation/'
import './MainLayout.css'
import { User } from '../../entities/user'

const MainLayout = () => {
    return (
        <div >
            <Navigation />
            <main>
                <Outlet />
            </main>
            <User />
        </div>
    )
}

export default MainLayout
