'use client';

import { useEffect } from 'react';
import { useAuthStore } from '@/stores/useAuthStore';
import { User } from '@/types/User';

interface AuthStoreInitializerProps {
	user: User | null;
}

/**
 * A client-side component that initializes the `useAuthStore` with user data
 * fetched from a Server Component.
 *
 * This component acts as a "hydration" bridge, setting the initial state
 * of the client-side store without rendering any UI itself.
 *
 * @param {object} props The component props.
 * @param {User | null} props.user The user object fetched on the server, or null if unauthenticated.
 * @returns {null} This component renders nothing to the DOM.
 */
function AuthStoreInitializer({ user }: AuthStoreInitializerProps): null {
	const setUser = useAuthStore((state) => state.setUser);

	useEffect(() => {
		// This effect runs once on the client, initializing the store.
		setUser(user);
	}, [user, setUser]);

	return null;
}

export default AuthStoreInitializer;
