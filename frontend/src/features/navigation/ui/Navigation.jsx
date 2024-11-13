/* eslint-disable react/prop-types */
import './Navigation.css'
import { NavLink } from '../../../shared/ui/NavLink'
import { Button } from '../../../shared/ui/Button'
import { useAuth } from '../../../entities/user'

const Navigation = ({ className }) => {
    const userStore = useAuth()

    return (
        <div className={`Navigation rounded-3xl flex items-center justify-between ${className}`}>
            <div className='Nav-container'>
                <Button 
                    className='Avatar rounded-3xl'
                    onClick={() => userStore.setOpen(true)}
                >
                    <div className='flex items-center justify-center w-full h-full Avatar-standart rounded-3xl'>
                        <div className='flex items-center justify-center w-1/2 aspect-square'>
                            <svg viewBox="0 0 16 19" fill="none" xmlns="http://www.w3.org/2000/svg">
                                <defs/>
                                <path id="Vector" d="M8 0C5.23 0 3 2.23 3 5C3 7.76 5.23 10 8 10C10.76 10 13 7.76 13 5C13 2.23 10.76 0 8 0ZM11 5C11 6.65 9.65 8 8 8C6.34 8 5 6.65 5 5C5 3.34 6.34 2 8 2C9.65 2 11 3.34 11 5ZM0 19C0 16.87 0.84 14.84 2.34 13.34C3.84 11.84 5.87 11 8 11C10.12 11 12.15 11.84 13.65 13.34C15.15 14.84 16 16.87 16 19L14 19C14 17.4 13.36 15.88 12.24 14.75C11.11 13.63 9.59 13 8 13C6.4 13 4.88 13.63 3.75 14.75C2.63 15.88 2 17.4 2 19L0 19Z" fill="#0F1828" fillOpacity="1.000000" fillRule="evenodd"/>
                            </svg>
                        </div>
                    </div>
                </Button>

                <div className='Nav-container__links'>
                    <NavLink className='flex items-center justify-center NavLink' to='/'>
                        <div className='Link-Icon'>
                            <svg viewBox="0 0 45.2986 43.9343" fill="none" xmlns="http://www.w3.org/2000/svg" >
                                <defs/>
                                <path id="Vector" d="M38.53 6.46C34.83 2.84 29.95 0.59 24.72 0.1C19.5 -0.4 14.26 0.89 9.92 3.75C5.57 6.6 2.38 10.84 0.91 15.73C-0.57 20.62 -0.24 25.86 1.83 30.54C2.05 30.98 2.12 31.47 2.03 31.94L0.05 41.21C-0.03 41.56 -0.01 41.93 0.09 42.28C0.2 42.62 0.39 42.94 0.66 43.2C0.87 43.41 1.13 43.57 1.41 43.68C1.7 43.79 2 43.84 2.3 43.83L2.75 43.83L12.4 41.95C12.89 41.89 13.39 41.96 13.84 42.15C18.66 44.16 24.06 44.48 29.09 43.05C34.13 41.61 38.49 38.52 41.43 34.29C44.37 30.07 45.7 24.98 45.19 19.91C44.68 14.84 42.36 10.09 38.64 6.5L38.53 6.46ZM13.57 24.13C13.13 24.13 12.69 24 12.32 23.76C11.95 23.52 11.66 23.18 11.49 22.78C11.32 22.38 11.28 21.94 11.36 21.51C11.45 21.09 11.66 20.7 11.98 20.39C12.3 20.08 12.7 19.88 13.13 19.79C13.57 19.71 14.02 19.75 14.44 19.92C14.85 20.08 15.2 20.36 15.45 20.72C15.7 21.08 15.83 21.51 15.83 21.94C15.83 22.52 15.59 23.08 15.17 23.49C14.75 23.9 14.17 24.13 13.57 24.13ZM22.59 24.13C22.14 24.13 21.71 24 21.34 23.76C20.97 23.52 20.68 23.18 20.51 22.78C20.34 22.38 20.29 21.94 20.38 21.51C20.47 21.09 20.68 20.7 21 20.39C21.31 20.08 21.71 19.88 22.15 19.79C22.59 19.71 23.04 19.75 23.45 19.92C23.86 20.08 24.22 20.36 24.46 20.72C24.71 21.08 24.84 21.51 24.84 21.94C24.84 22.52 24.61 23.08 24.18 23.49C23.76 23.9 23.19 24.13 22.59 24.13ZM31.61 24.13C31.16 24.13 30.73 24 30.35 23.76C29.98 23.52 29.69 23.18 29.52 22.78C29.35 22.38 29.31 21.94 29.4 21.51C29.48 21.09 29.7 20.7 30.01 20.39C30.33 20.08 30.73 19.88 31.17 19.79C31.6 19.71 32.06 19.75 32.47 19.92C32.88 20.08 33.23 20.36 33.48 20.72C33.73 21.08 33.86 21.51 33.86 21.94C33.86 22.52 33.62 23.08 33.2 23.49C32.78 23.9 32.2 24.13 31.61 24.13Z" fill="#FFFFFF" fillOpacity="1.000000" fillRule="nonzero"/>
                            </svg>
                        </div>
                    </NavLink>

                    <NavLink className='flex items-center justify-center NavLink' to='/friends'>
                        <div className='Link-Icon'>
                            <svg viewBox="0 0 26.6666 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                <defs/>
                                <path id="coolicon" d="M2.66 6.66C2.66 2.98 5.65 0 9.33 0C13.01 0 16 2.98 16 6.66C16 10.34 13.01 13.33 9.33 13.33C5.65 13.33 2.66 10.34 2.66 6.66ZM9.33 2.66C7.12 2.66 5.33 4.45 5.33 6.66C5.33 8.87 7.12 10.66 9.33 10.66C11.54 10.66 13.33 8.87 13.33 6.66C13.33 4.45 11.54 2.66 9.33 2.66ZM18.66 6.66C19.08 6.66 19.5 6.76 19.87 6.95C20.25 7.14 20.57 7.42 20.82 7.76C21.07 8.1 21.23 8.5 21.3 8.91C21.36 9.33 21.33 9.75 21.2 10.16C21.07 10.56 20.84 10.92 20.54 11.22C20.25 11.51 19.88 11.74 19.48 11.87C19.08 12 18.66 12.03 18.24 11.96C17.82 11.89 17.43 11.73 17.09 11.48L15.52 13.64L15.52 13.64C16.2 14.13 16.99 14.46 17.82 14.59C17.91 14.61 18 14.62 18.09 14.63C18.84 14.71 19.59 14.63 20.3 14.4C21.1 14.14 21.83 13.7 22.43 13.1C23.02 12.51 23.47 11.78 23.73 10.98C23.99 10.18 24.06 9.33 23.93 8.5C23.8 7.67 23.47 6.88 22.98 6.2C22.54 5.59 21.98 5.08 21.33 4.71C21.25 4.66 21.17 4.62 21.08 4.58C20.33 4.19 19.5 4 18.66 4L18.66 6.66ZM16 24L18.66 24C18.66 18.84 14.48 14.66 9.33 14.66C4.17 14.66 0 18.84 0 24L2.66 24C2.66 20.31 5.65 17.33 9.33 17.33C13.01 17.33 16 20.31 16 24ZM23.99 24C23.99 23.3 23.85 22.6 23.59 21.96C23.32 21.31 22.93 20.72 22.43 20.23C21.94 19.73 21.35 19.34 20.7 19.07C20.05 18.8 19.36 18.66 18.66 18.66L18.66 16C19.57 16 20.47 16.15 21.33 16.45C21.46 16.5 21.59 16.55 21.72 16.6C22.69 17.01 23.58 17.6 24.32 18.34C25.06 19.08 25.65 19.96 26.05 20.93C26.11 21.06 26.16 21.2 26.2 21.33C26.51 22.18 26.66 23.09 26.66 24L23.99 24Z" fill="#F7F7FC" fillOpacity="1.000000" fillRule="evenodd"/>
                            </svg>
                        </div>
                    </NavLink>

                    <NavLink className='flex items-center justify-center NavLink' to='/options'>
                        <div className='Link-Icon'>
                            <svg viewBox="0 0 42.8492 44.1667" fill="none" xmlns="http://www.w3.org/2000/svg" >
                                <defs/>
                                <path id="Vector" d="M25.44 44.16L17.4 44.16C16.9 44.16 16.41 43.99 16.02 43.67C15.63 43.36 15.35 42.92 15.24 42.43L14.35 38.27C13.15 37.74 12.01 37.08 10.96 36.31L6.9 37.6C6.42 37.75 5.9 37.74 5.43 37.55C4.96 37.37 4.57 37.03 4.32 36.6L0.29 29.64C0.04 29.2 -0.06 28.69 0.02 28.2C0.1 27.7 0.34 27.24 0.71 26.9L3.86 24.03C3.72 22.73 3.72 21.42 3.86 20.12L0.71 17.26C0.34 16.91 0.1 16.46 0.02 15.96C-0.06 15.46 0.04 14.95 0.29 14.52L4.31 7.55C4.56 7.12 4.95 6.78 5.42 6.6C5.89 6.42 6.41 6.4 6.89 6.55L10.95 7.85C11.49 7.45 12.05 7.08 12.63 6.74C13.18 6.43 13.76 6.14 14.35 5.89L15.25 1.73C15.35 1.24 15.63 0.8 16.02 0.48C16.41 0.17 16.9 0 17.4 0L25.44 0C25.94 0 26.43 0.17 26.82 0.48C27.22 0.8 27.49 1.24 27.6 1.73L28.51 5.89C29.7 6.42 30.84 7.07 31.89 7.85L35.95 6.56C36.43 6.41 36.95 6.42 37.42 6.6C37.89 6.79 38.28 7.12 38.53 7.56L42.55 14.52C43.07 15.42 42.89 16.56 42.13 17.26L38.98 20.13C39.12 21.43 39.12 22.74 38.98 24.04L42.13 26.91C42.89 27.61 43.07 28.75 42.55 29.65L38.53 36.61C38.28 37.05 37.89 37.39 37.42 37.57C36.95 37.75 36.43 37.76 35.95 37.61L31.89 36.32C30.84 37.1 29.7 37.75 28.51 38.28L27.6 42.43C27.49 42.92 27.22 43.36 26.82 43.67C26.43 43.99 25.94 44.16 25.44 44.16ZM11.75 31.42L13.56 32.74C13.97 33.04 14.39 33.32 14.83 33.57C15.25 33.81 15.67 34.03 16.11 34.22L18.17 35.13L19.18 39.75L23.67 39.75L24.67 35.13L26.73 34.22C27.63 33.82 28.49 33.33 29.28 32.75L31.09 31.43L35.6 32.86L37.84 28.98L34.35 25.79L34.59 23.56C34.7 22.58 34.7 21.59 34.59 20.62L34.35 18.38L37.84 15.19L35.6 11.3L31.09 12.74L29.28 11.41C28.49 10.83 27.63 10.33 26.73 9.93L24.67 9.03L23.67 4.41L19.18 4.41L18.17 9.03L16.11 9.93C15.21 10.33 14.36 10.82 13.57 11.4L11.75 12.73L7.25 11.29L5 15.19L8.5 18.37L8.25 20.61C8.14 21.59 8.14 22.57 8.25 23.55L8.5 25.78L5 28.97L7.24 32.85L11.75 31.42ZM21.41 30.91C19.07 30.91 16.82 29.98 15.17 28.32C13.51 26.67 12.58 24.42 12.58 22.08C12.58 19.74 13.51 17.49 15.17 15.83C16.82 14.18 19.07 13.25 21.41 13.25C23.76 13.25 26 14.18 27.66 15.83C29.32 17.49 30.25 19.74 30.25 22.08C30.25 24.42 29.32 26.67 27.66 28.32C26 29.98 23.76 30.91 21.41 30.91ZM21.41 17.66C20.55 17.66 19.7 17.92 18.98 18.39C18.26 18.87 17.69 19.55 17.35 20.35C17.02 21.14 16.92 22.02 17.07 22.87C17.23 23.72 17.63 24.51 18.23 25.13C18.82 25.76 19.59 26.19 20.44 26.38C21.28 26.58 22.16 26.52 22.97 26.21C23.78 25.91 24.48 25.37 24.99 24.67C25.5 23.97 25.79 23.14 25.83 22.28L25.83 23.16L25.83 22.08C25.83 20.91 25.36 19.78 24.54 18.96C23.71 18.13 22.58 17.66 21.41 17.66Z" fill="#FFFFFF" fillOpacity="1.000000" fillRule="nonzero"/>
                            </svg>
                        </div>
                    </NavLink>
                </div>
            </div>
        </div>
    )
}

export default Navigation
