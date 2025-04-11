import { NavLink } from 'react-router-dom'
import styles from './Link.module.css'

interface LinkProps {
    to: string,
    children?: React.ReactNode,
    className?: string
}

const Link: React.FC<LinkProps> = (props: LinkProps) => {
    return (
        <NavLink 
            to={props.to} 
            className={styles.Link + ' ' + props.className}
        >
            {props.children}
        </NavLink>
    )
}

export default Link
