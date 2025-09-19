import { cookies } from 'next/headers';
import { User } from '@/types/User';

/**
 * Fetches the current user from the backend API on the server.
 * It forwards the request cookies to authenticate the request.
 * @returns {Promise<User | null>} The user object or null if not authenticated.
 */
export async function getCurrentUser(): Promise<User | null> {
  const cookieStore = await cookies();

  const token = cookieStore.get('jwt');
  if (!token) {
    return null;
  }

  try {
    const response = await fetch('http://localhost:8080/auth/me', {
      headers: {
        'Cookie': `${token.name}=${token.value}`,
      },
      cache: 'no-store',
    });

    if (!response.ok) {
      console.error('Authentication failed:', response.statusText);
      return null;
    }

    const user: User = await response.json();
    return user;

  } catch (error) {
    console.error('Failed to fetch user data:', error);
    return null;
  }
}
