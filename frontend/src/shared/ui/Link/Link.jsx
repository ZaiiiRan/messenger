/* eslint-disable react/prop-types */
import { NavLink } from 'react-router-dom'
import styles from './Link.module.css'

const Link = ({ to, children, className }) => {
    return (
        <NavLink 
            to={to} 
            className={styles.Link + ' ' + className}
        >
            {children}
        </NavLink>
    )
}

export default Link
