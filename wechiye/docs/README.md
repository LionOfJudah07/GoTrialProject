# WECHIYE Documentation

## Security Architecture
- Database encryption via SQLCipher (AES-256).
- Master password hashed with scrypt.
- Keychain integration for "remember me".
- ECDH for couple pairing.
- AES-GCM for shared data encryption.

## Build Instructions
1. Install Go 1.21+, Wails v2, and a C compiler.
2. Run `go mod tidy`.
3. Run `wails build -tags sqlite_see`.