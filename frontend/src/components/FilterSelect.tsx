'use client';

import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select';
import { Button } from './ui/button';
import { X } from 'lucide-react';

type SelectOption = {
	value: string;
	label: string;
};

interface FilterSelectProps {
	options: SelectOption[];
	value: string;
	onChange: (newValue: string) => void;
	placeholder: string;
}

export function FilterSelect({
	options,
	value,
	onChange,
	placeholder,
}: FilterSelectProps) {
	return (
		<div className="relative flex items-center">
			<Select
				value={value}
				onValueChange={onChange}
			>
				<SelectTrigger className={`w-[px] ` + (value && "pr-8")}>
					<SelectValue placeholder={placeholder} />
				</SelectTrigger>
				<SelectContent>
					{options.map((option) => (
						<SelectItem key={option.value} value={option.value}>
							{option.label}
						</SelectItem>
					))}
				</SelectContent>
			</Select>

			{value && (
				<Button
					variant="ghost"
					size="icon"
					className="absolute right-1 h-6 w-6"
					onClick={() => onChange('')}
				>
					<X className="h-4 w-4 text-muted-foreground" />
				</Button>
			)}
		</div>
	);
}
