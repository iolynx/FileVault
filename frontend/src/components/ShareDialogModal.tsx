import { Dispatch, SetStateAction, useEffect, useState } from "react";
import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button";
import { MultiSelect, MultiSelectOption } from "./multi-select";

interface ShareDialogModalProps {
	isOpen: boolean;
	isOpenChange: (open: boolean) => void;
	userOptions: MultiSelectOption[];
	onConfirm: (usersToShare: string[]) => void;
	defaultValue: string[]
}

export function ShareDialogModal({ isOpen, isOpenChange, userOptions, onConfirm, defaultValue }: ShareDialogModalProps) {
	const [usersToShare, setUsersToShare] = useState<string[]>([]);
	return (
		<Dialog open={isOpen} onOpenChange={isOpenChange}>
			<DialogContent className="sm:max-w-md">
				<DialogHeader>
					<DialogTitle>Share To:</DialogTitle>
				</DialogHeader>
				<div className="flex items-center gap-2">
					<div className="grid flex-1 gap-2">
						<MultiSelect
							options={userOptions}
							onValueChange={setUsersToShare}
							defaultValue={defaultValue}
							placeholder="Search For Users"
							animationConfig={{
								badgeAnimation: "none",
							}}
						/>
					</div>
				</div>
				<DialogFooter className="sm:justify-start">
					<DialogClose asChild>
						<Button type="button" variant="secondary" onClick={() => onConfirm(usersToShare)}>
							Done
						</Button>
					</DialogClose>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
};

