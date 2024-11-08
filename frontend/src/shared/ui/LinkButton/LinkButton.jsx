/* eslint-disable react/prop-types */
import styles from './LinkButton.module.css'

const LinkButton = ({ onClick, children, className }) => {
    return (
        <a 
            className={styles.LinkButton + ' ' + className}
            onClick={onClick}
        >
            {children}
        </a>
    )
}

export default LinkButton
