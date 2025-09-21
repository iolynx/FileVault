'use client';

import { create } from 'zustand';
import { User } from '@/types/User';
import api from '@/lib/axios';


interface AuthState {
	user: User | null;
	isBannerDismissed: boolean;
	fetchUser: () => Promise<void>;
	dismissBanner: () => void;
	setUser: (user: User | null) => void;
}

export const useAuthStore = create<AuthState>((set) => ({
	user: null,
	// Check sessionStorage to see if the user has already closed the banner in this session.
	isBannerDismissed: false,
	// typeof window !== 'undefined'
	// 	? window.sessionStorage.getItem('isQuotaBannerDismissed') === 'true'
	// 	: false,

	// Fetches user data from the /auth/me endpoint
	fetchUser: async () => {
		try {
			const response = await api.get('/auth/me');
			set({ user: response.data });
		} catch (error) {
			console.error('Failed to fetch user data', error);
			set({ user: null }); // Clear user on error
		}
	},

	// Sets the banner as dismissed and saves the choice to sessionStorage
	dismissBanner: () => {
		sessionStorage.setItem('isQuotaBannerDismissed', 'true');
		set({ isBannerDismissed: true });
	},

	setUser: (user) => {
		set({ user })
	}
}));
