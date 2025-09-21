'use client';

import { useState, useEffect } from 'react';

/**
 * A custom hook to debounce a value.
 * @param value The value to debounce.
 * @param delay The debounce delay in milliseconds.
 * @returns The debounced value.
 */
export function useDebounce<T>(value: T, delay: number): T {
	const [debouncedValue, setDebouncedValue] = useState<T>(value);

	useEffect(
		() => {
			// Setting up a timer to update the debounced value after the specified delay
			const handler = setTimeout(() => {
				setDebouncedValue(value);
			}, delay);

			// This function cancels the previous timer, so only the latest one will execute
			return () => {
				clearTimeout(handler);
			};
		},
		[value, delay]
	);

	return debouncedValue;
}
