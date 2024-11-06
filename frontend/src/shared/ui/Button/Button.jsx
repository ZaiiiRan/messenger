/* eslint-disable react/prop-types */
import styles from './Button.module.css'

const Button = ({ children, className, onClick, disabled=false }) => {
    return (
        <button 
            className={styles.Button + ' ' + className}
            disabled={disabled}
            onClick={onClick}
        >
            {children}
        </button>
    )
}

export default Button
