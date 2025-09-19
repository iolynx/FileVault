// The shape of a single filter option for the dropdown
export type FilterOption = {
	value: string;
	label: string;
};

// The shape of the currently active filter state
export type ActiveFilter = {
	column: string;
	value: string;
};
