import { NavLink } from 'react-router-dom'
import styles from './NavLink.module.css'

interface NavBarLinkProps {
    className?: string,
    children?: React.ReactNode,
    to: string
}

const NavBarLink: React.FC<NavBarLinkProps> = ({ className, children, to }) => {
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
