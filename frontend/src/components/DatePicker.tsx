"use client"

import * as React from "react"
import { Calendar } from "@/components/ui/calendar"
import { Button } from "@/components/ui/button"
import { X } from "lucide-react"

interface DatePickerProps {
	date: Date | undefined;
	onDateChange: (date: Date | undefined) => void;
}

export function DatePicker({ date, onDateChange }: DatePickerProps) {
	const [open, setOpen] = React.useState(false)

	const handleClear = () => {
		onDateChange(undefined);
		setOpen(false);
	}

	return (
		<div className="flex flex-col">
			<Calendar
				mode="single"
				selected={date}
				onSelect={(newDate) => {
					onDateChange(newDate)
					setOpen(false)
				}}
			/>

			{date && (
				<div className="p-2 border-t flex justify-center">
					<Button
						variant="outline"
						size="sm"
						onClick={handleClear}
					>
						Clear
					</Button>
				</div>
			)}
		</div>
	)
}
