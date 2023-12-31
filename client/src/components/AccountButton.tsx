import styles from '../styles/Header.module.css';
import { createSignal, createEffect } from 'solid-js';

const AccountButton = () => {
  const [user, setUser] = createSignal(null);

  createEffect(() => {
    const storedUserData = localStorage.getItem('user');
    if (storedUserData) {
      const storedUser = JSON.parse(storedUserData);
      setUser(storedUser);
    }
  });

  return (
    <li class={styles.navItem}>
      <a href={user() ? '/account' : '/login'} class={styles.navListItem}>
        {user() ? 'Account' : 'Login'}
      </a>
    </li>
  );
};

export default AccountButton;
