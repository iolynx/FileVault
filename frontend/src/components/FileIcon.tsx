import {
	File,
	FileImage,
	FileText,
	FileVideo,
	FileVolume,
	FileArchive,
	FileCode,
	FileSpreadsheet,
} from 'lucide-react';

// Map mimeType -> File Icons
export function getFileIcon(mimeType: string) {
	if (mimeType.startsWith('image/')) {
		return <FileImage className="h-5 w-5 text-gray-500" />;
	}
	if (mimeType.startsWith('video/')) {
		return <FileVideo className="h-5 w-5 text-gray-500" />;
	}
	if (mimeType.startsWith('audio/')) {
		return <FileVolume className="h-5 w-5 text-gray-500" />;
	}

	switch (mimeType) {
		case 'application/pdf':
			return <FileText className="h-5 w-5 text-red-500" />;
		case 'text/plain':
		case 'text/csv':
			return <FileText className="h-5 w-5 text-gray-500" />;
		case 'application/zip':
		case 'application/x-rar-compressed':
			return <FileArchive className="h-5 w-5 text-yellow-500" />;
		case 'application/vnd.ms-excel':
		case 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet':
			return <FileSpreadsheet className="h-5 w-5 text-green-500" />;
		default:
			return <File className="h-5 w-5 text-gray-500" />;
	}
}
