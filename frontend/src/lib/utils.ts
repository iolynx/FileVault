import { MultiSelectOption } from "@/components/multi-select";
import { User } from "@/types/User";
import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatBytes(bytes: number, decimals = 2): string {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
}

export const mapUsersToOptions = (users: User[]): MultiSelectOption[] => {
  return users.map((u) => ({
    label: `${u.name} (${u.email})`,
    value: u.id,

    icon: undefined,
    disabled: false,
  }));
};

export const toSentenceCase = (s: string) => {
  if (!s) {
    return "";
  }
  return s.charAt(0).toUpperCase() + s.slice(1).toLowerCase();
}
