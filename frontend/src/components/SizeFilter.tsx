'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuLabel,
	DropdownMenuRadioGroup,
	DropdownMenuRadioItem,
	DropdownMenuSeparator,
	DropdownMenuSub,
	DropdownMenuSubContent,
	DropdownMenuSubTrigger,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select';
import { SizeRange } from '@/types/Size';
import { ChevronDown, X, XIcon } from 'lucide-react';
import { cn } from '@/lib/utils';

const UNIT_MULTIPLIERS = { BYTES: 1, KB: 1024, MB: 1024 ** 2, GB: 1024 ** 3 };
type Unit = keyof typeof UNIT_MULTIPLIERS;



interface SizeFilterProps {
	ranges: SizeRange[];
	onApplyFilter: (filter: { min_size: number | null; max_size: number | null }) => void;
}

export function SizeFilter({ ranges, onApplyFilter }: SizeFilterProps) {
	const [activeLabel, setActiveLabel] = useState<string | null>(null);
	const [minVal, setMinVal] = useState('');
	const [minUnit, setMinUnit] = useState<Unit>('MB');
	const [maxVal, setMaxVal] = useState('');
	const [maxUnit, setMaxUnit] = useState<Unit>('MB');

	const handleRangeSelect = (label: string) => {
		const selectedRange = ranges.find((r) => r.label === label);
		if (selectedRange) {
			setActiveLabel(label);
			onApplyFilter({ min_size: selectedRange.minBytes, max_size: selectedRange.maxBytes });
		}
	};

	const handleApplyCustomFilter = () => {
		const minBytes = minVal ? parseInt(minVal) * UNIT_MULTIPLIERS[minUnit] : null;
		const maxBytes = maxVal ? parseInt(maxVal) * UNIT_MULTIPLIERS[maxUnit] : null;
		setActiveLabel('Custom');
		onApplyFilter({ min_size: minBytes, max_size: maxBytes });
	};

	const handleClearFilter = () => {
		setActiveLabel(null);
		setMinVal('');
		setMaxVal('');
		onApplyFilter({ min_size: null, max_size: null });
	};

	return (
		<div className="relative flex items-center">
			<DropdownMenu>
				<DropdownMenuTrigger >
					<Button
						variant="outline"
						className={activeLabel ? "pr-2" : "text-muted-foreground"}>
						<span>Size: {activeLabel || 'Any'}</span>
						<ChevronDown
							className={cn(
								activeLabel ? "mr-4" : "w-3 h-3 ml-0"
							)}
						/>
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent className="w-56">
					<DropdownMenuLabel>Filter by size</DropdownMenuLabel>
					<DropdownMenuSeparator />
					<DropdownMenuRadioGroup value={activeLabel || ''} onValueChange={handleRangeSelect}>
						{ranges.map((range) => (
							<DropdownMenuRadioItem key={range.label} value={range.label}>
								{range.label}
							</DropdownMenuRadioItem>
						))}
					</DropdownMenuRadioGroup>

					<DropdownMenuSeparator />

					<DropdownMenuSub>
						<DropdownMenuSubTrigger>Custom Range</DropdownMenuSubTrigger>
						<DropdownMenuSubContent>
							{/* Stop propagation to prevent the dropdown from closing when interacting with inputs */}
							<div className="p-2 space-y-4" onClick={(e) => e.stopPropagation()}>
								<div className="space-y-2">
									<label className="text-sm font-medium">Minimum</label>
									<div className="flex gap-2">
										<Input type="number" placeholder="Value" value={minVal} onChange={(e) => setMinVal(e.target.value)} />
										<Select value={minUnit} onValueChange={(v: Unit) => setMinUnit(v)}>
											<SelectTrigger className="w-[80px]"><SelectValue /></SelectTrigger>
											<SelectContent>
												<SelectItem value="KB">KB</SelectItem>
												<SelectItem value="MB">MB</SelectItem>
												<SelectItem value="GB">GB</SelectItem>
											</SelectContent>
										</Select>
									</div>
								</div>
								<div className="space-y-2">
									<label className="text-sm font-medium">Maximum</label>
									<div className="flex gap-2">
										<Input type="number" placeholder="Value" value={maxVal} onChange={(e) => setMaxVal(e.target.value)} />
										<Select value={maxUnit} onValueChange={(v: Unit) => setMaxUnit(v)}>
											<SelectTrigger className="w-[80px]"><SelectValue /></SelectTrigger>
											<SelectContent>
												<SelectItem value="KB">KB</SelectItem>
												<SelectItem value="MB">MB</SelectItem>
												<SelectItem value="GB">GB</SelectItem>
											</SelectContent>
										</Select>
									</div>
								</div>
								<Button className="w-full" size="sm" onClick={handleApplyCustomFilter}>
									Apply Custom Filter
								</Button>
							</div>
						</DropdownMenuSubContent>
					</DropdownMenuSub>

				</DropdownMenuContent>
			</DropdownMenu>

			{activeLabel && (
				<Button
					variant="ghost"
					size="icon"
					className="absolute right-1 h-6 w-6"
					onClick={handleClearFilter}
				>
					<X className="w-4 h-4 text-muted-foreground" />
				</Button>
			)}

		</div>
	);
}
