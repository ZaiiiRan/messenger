import { Outlet } from 'react-router-dom'
import { Navigation } from '../../features/navigation/'
import './MainLayout.css'

const MainLayout = () => {
    return (
        <div >
            <Navigation />
            <main>
                <Outlet />
            </main>
        </div>
    )
}

export default MainLayout
