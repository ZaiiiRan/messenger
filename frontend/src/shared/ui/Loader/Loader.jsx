/* eslint-disable react/prop-types */
import styles from './Loader.module.css'

const Loader = ({ className }) => {
    return (
        <div className={styles.Loader + ' ' + className}>
        <div
            className={styles.LoaderElement + ' h-full animate-pulse rounded-full'}
        ></div>
        <div
            className={styles.LoaderElement + ' animate-pulse h-full rounded-full'}
        ></div>
        <div
            className={styles.LoaderElement + ' h-full animate-pulse rounded-full'}
        ></div>
    </div>
    )
}

export default Loader
