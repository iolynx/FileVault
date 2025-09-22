'use client';
import { useContentStore } from '@/stores/useContentStore';

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
