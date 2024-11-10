/* eslint-disable react/prop-types */
import { NavLink } from 'react-router-dom'
import styles from './NavLink.module.css'


const NavBarLink = ({ className, children, to }) => {
    return (
        <NavLink 
            to={to} 
            className={styles.NavLink + ' ' + className + ' rounded-3xl'}
        >
            {children}
        </NavLink>
    )
}

export default NavBarLink
