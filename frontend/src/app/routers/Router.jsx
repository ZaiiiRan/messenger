import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { AnimatePresence } from 'framer-motion'
import { StartPage } from '../../pages/StartPage'
import { LoginPage } from '../../pages/LoginPage'
import { RegisterPage } from '../../pages/RegisterPage'
import StartRedirect from '../hocs/StartRedirect'
import RequireAuth from '../hocs/RequireAuth'

export const Router = () => {
    return(
        <BrowserRouter>
            <AnimatePresence mode="wait">
                <Routes location={location} key={location.pathname}>
                    <Route path='/' index element={<RequireAuth><div>че то</div></RequireAuth>}></Route>
                    <Route path='/start' element={<StartRedirect><StartPage /></StartRedirect>}></Route>
                    <Route path='/login' element={<StartRedirect><LoginPage /></StartRedirect>}></Route>
                    <Route path='/register' element={<StartRedirect><RegisterPage /></StartRedirect>}></Route>
                </Routes>
            </AnimatePresence>
        </BrowserRouter>
    ) 
}