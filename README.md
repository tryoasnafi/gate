# gate

`gate` is a fast and secure CLI tool to manage SSH credentials and connect to remote servers.

## Features

- [x] Master password protected store (Argon2id + ChaCha20 encryption)
- [x] Modular structure using Cobra + Viper
- [x] Native SSH using `golang.org/x/crypto/ssh` (no `sshpass`)
- [x] Import/export support with encryption
- [x] Interactive TTY with color support (`TERM=xterm-256color`)
- [ ] Session-based password cache (like `sudo`)
- [ ] Support identity file
- [ ] Secure copy to clipboard password
- [ ] Support sshd config

## Commands

```sh
# Initialize gate (create vault)
gate init

# Rotate/change master password
gate rotate

# List saved SSH entries
gate list

# Add a new entry
gate new

# Delete an entry
gate delete <label>

# Import from file
gate import ./store.enc

# Connect to a label
gate connect <label>
```

## Structure

Each SSH entry contains:

```json
{
  "user": "root",
  "host": "example.com",
  "port": 22,
  "password": "secret",
  "createdAt": "2025-07-01T12:00:00Z"
}
```

## Requirements

* Go 1.21+
* Linux/macOS terminal (TTY supported)

## Security Notes

* Master password is never stored.
* Gate objects are encrypted using ChaCha20-Poly1305.
* Key is derived via Argon2id(master password).


