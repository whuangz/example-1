-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER before_insert_account
BEFORE INSERT ON `account`
FOR EACH ROW
BEGIN
  IF NEW.uid IS NULL THEN
    SET NEW.uid = uuid();
  END IF;
END;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin 
DELETE TRIGGER before_insert_account
-- +goose StatementEnd
