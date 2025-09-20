import React from 'react';

interface StopPropagationWrapperProps {
	children: React.ReactNode;
	className?: string;
}

export const StopPropagationWrapper = ({ children, className }: StopPropagationWrapperProps) => {
	// A single function to handle all events
	const stopPropagation = (e: React.SyntheticEvent) => {
		e.stopPropagation();
	};

	return (
		<div
			onClick={stopPropagation}
			onPointerDown={stopPropagation}
			onMouseDown={stopPropagation}
			onKeyDown={stopPropagation}
			className={className}
		>
			{children}
		</div>
	);
};
