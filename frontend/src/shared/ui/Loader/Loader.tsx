import styles from './Loader.module.css'

interface LoaderProps {
    className?: string
}

const Loader: React.FC<LoaderProps> = ({ className }) => {
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
