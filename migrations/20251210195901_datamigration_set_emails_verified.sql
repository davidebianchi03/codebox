-- Mark emails of existing users verified
UPDATE users SET email_verified = TRUE WHERE email_verified = FALSE;
