-- Mark emails of existing users verified
UPDATE users SET email_verified = TRUE WHERE email_verified = FALSE;

--Mark all existing users as approved
UPDATE users SET approved = TRUE WHERE approved = FALSE;
