-- === TRIGGER FOR WHEN A NEW FILE IS CREATED ===

-- This function increments the storage
CREATE OR REPLACE FUNCTION update_user_storage_on_insert()
RETURNS TRIGGER AS $$
BEGIN
    -- Increment the storage_used for the user who owns the new file.
    UPDATE users
    SET storage_used = storage_used + NEW.size
    WHERE id = NEW.owner_id;

    -- Increment teh blob refcount
    UPDATE blobs
    SET refcount = refcount + 1
    WHERE id = NEW.blob_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- This trigger calls the function after a file is inserted
CREATE TRIGGER files_after_insert_storage_trigger
AFTER INSERT ON files
FOR EACH ROW
EXECUTE FUNCTION update_user_storage_on_insert();


-- === TRIGGER FOR WHEN A FILE IS DELETED ===

-- This functino handles all deletion logic
CREATE OR REPLACE FUNCTION handle_file_deletion()
RETURNS TRIGGER AS $$
BEGIN
    --  decrement the user's storage usage
    UPDATE users
    SET storage_used = storage_used - OLD.size
    WHERE id = OLD.owner_id;

    -- decrement the blob refcount
    UPDATE blobs
    SET refcount = refcount - 1
    WHERE id = OLD.blob_id;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- This trigger calls the function after a file is deleted
CREATE TRIGGER files_after_delete_storage_trigger
AFTER DELETE ON files
FOR EACH ROW
EXECUTE FUNCTION handle_file_deletion();
