# TriCrypt

TriCrypt is a custom-built client-server application suite designed for secure communication with full end-to-end encryption (E2EE). It leverages the advanced Three-Layer-Encryption-Protocol (TLEP) to ensure robust security across multiple encryption layers.

## Features

- **Three-Layer-Encryption-Protocol (TLEP):**
  - **Elliptic-Curve Diffie-Hellman (ECDH)**: Uses the `secp521r1` curve for secure key exchange.
  - **Lorenz-Based Chaos System**: Adds dynamic and unpredictable transformations for encryption.
  - **Zero-Trust Pseudo-Language Dictionary**: Obfuscates communication patterns to prevent unauthorized inference.

- **Client Interface**: Built with the Fyne GUI for a user-friendly, cross-platform interface.
- **Future Development**: Planned implementation of a custom terminal-based UI in Go to minimize third-party dependencies and enhance flexibility.

## Getting Started

### Prerequisites

- **Go** (latest version recommended): [Install Go](https://golang.org/doc/install)
- **Fyne GUI Library**: Install via:
  ```bash
  go get fyne.io/fyne/v2
  ```
