import api from "@/lib/axios";
import { toast } from "sonner";

// Helper function to make the API call to fetch files and folders
export const getContent = async (folderId: string | null, filters: any) => {
  try {
    const params = new URLSearchParams();

    // Add the folderId if it exists
    if (folderId) {
      params.append('folder_id', folderId);
    }

    params.append('limit', '10')
    params.append('offset', '0')

    for (const key in filters) {
      if (filters[key]) {
        params.append(key, filters[key]);
      }
    }

    const response = await api.get('/files', { params });
    if (typeof response.data === 'string') {
      if (response.data === "") {
        return [];
      }
      return JSON.parse(response.data);
    }
    return response.data || [];
  } catch (error) {
    throw error;
  }
};
