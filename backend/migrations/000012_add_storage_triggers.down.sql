-- Drop the trigger for file deletion
DROP TRIGGER IF EXISTS files_after_delete_storage_trigger ON files;

-- Drop the function that handles file deletion logic
DROP FUNCTION IF EXISTS handle_file_deletion();

-- Drop the trigger for file insertion
DROP TRIGGER IF EXISTS files_after_insert_storage_trigger ON files;

-- Drop the function that handles file insertion logic
DROP FUNCTION IF EXISTS update_user_storage_on_insert();
