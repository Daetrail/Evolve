# General Endpoints

## Auth
/auth/sign-up: 
- username: string, 4 <= length <= 15, unique
- email: string, valid email, unique
- password: string, needs 1 uppercase, 1 lowercase, 1 symbol, all ASCII characters, 8 <= length <= 32

/auth/sign-in:
- email: string, valid email
- password: string

/auth/verify:
- No fields
- Status 200 if bearer token is valid