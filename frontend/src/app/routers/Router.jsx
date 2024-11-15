import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { AnimatePresence } from 'framer-motion'
import { StartPage } from '../../pages/StartPage'
import { LoginPage } from '../../pages/LoginPage'
import { RegisterPage } from '../../pages/RegisterPage'
import { ActivationPage } from '../../pages/ActivationPage'
import StartRedirect from '../hocs/StartRedirect'
import RequireAuth from '../hocs/RequireAuth'
import RequireActivate from '../hocs/RequireActivate'
import ActivateRedirect from '../hocs/ActivateRedirect'
import MainLayout from '../layouts/MainLayout'
import { OptionsPage } from '../../pages/OptionsPage'
import { FriendsPage } from '../../pages/FriendsPage'
import { MessengingPage } from '../../pages/MessengingPage'

export const Router = () => {
    return(
        <BrowserRouter>
            <AnimatePresence mode="wait">
                <Routes location={location} key={location.pathname}>
                    <Route element={ 
                            <RequireAuth>
                                <RequireActivate>
                                    <MainLayout />
                                </RequireActivate>
                            </RequireAuth> 
                        }
                    >
                        <Route path="/" element={ <MessengingPage /> } />
                        <Route path="/friends" element={ <FriendsPage /> } />
                        <Route path="/options" element={ <OptionsPage /> } />
                    </Route>

                    <Route path='/start' element={
                            <StartRedirect>
                                <StartPage />
                            </StartRedirect>
                        }
                    />

                    <Route path='/login' element={
                            <StartRedirect>
                                <LoginPage />
                            </StartRedirect>
                        }
                    />

                    <Route path='/register' element={
                            <StartRedirect>
                                <RegisterPage />
                            </StartRedirect>
                        }
                    />

                    <Route path='/activate' element={
                            <RequireAuth>
                                <ActivateRedirect>
                                    <ActivationPage />
                                </ActivateRedirect>
                            </RequireAuth>
                        }
                    />
                </Routes>
            </AnimatePresence>
        </BrowserRouter>
    ) 
}