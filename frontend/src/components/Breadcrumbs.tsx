'use client';
import { useContentStore } from '@/stores/useContentStore';

/**
 * Breadcrumbs component
 * 
 * Renders a navigable breadcrumb trail for the current folder path.
 * Allows users to quickly navigate to any parent folder by clicking on its name.
 * 
 * Uses the `useContentStore` hook to:
 * - Access the current folder path
 * - Navigate to a specific folder in the path
 * 
 * @returns {JSX.Element} JSX element rendering the breadcrumb navigation
 * 
 * @component
 */
export function Breadcrumbs() {
	const { path, navigateToPathIndex } = useContentStore();

	return (
		<nav className="flex items-center text-md">
			{path.map((segment, index) => (
				<div key={segment.id || 'home'} className="flex items-center cursor-pointer">
					<button
						onClick={() => navigateToPathIndex(index)}
						className="text-gray-300 hover:text-gray-700 disabled:hover:text-gray-500 disabled:cursor-pointer"
						disabled={index === path.length - 1} // Disable the last item
					>
						{segment.name}
					</button>
					{index < path.length - 1 && <span className="mx-2">/</span>}
				</div>
			))}
		</nav>
	);
}
