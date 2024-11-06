import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { StartPage } from '../../pages/StartPage'

export const Router = () => {
    return(
        <BrowserRouter>
            <Routes>
                <Route path='/start' element={<StartPage />}></Route>
            </Routes>
        </BrowserRouter>
    ) 
}