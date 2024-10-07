# Fishy

Welcome to the Fishy CLI Tool! This application consists of a compiler and virtual machine for **FishyASM** and **Fishy Bytecode**, designed to streamline the development process within the Fishy ecosystem. Fishy is a toy and **should not** be used in production.

## Table of Contents

- [Fishy](#fishy)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Examples](#examples)
    - [Hello World](#hello-world)
  - [Contributing](#contributing)
  - [License](#license)

## Features

- **Compiler**: Convert FishyASM source code into Fishy Bytecode for execution on the virtual machine.
- **Virtual Machine**: Execute Fishy Bytecode with support for standard library functions and custom scripts.
- **Cross-Platform**: Works on Windows, macOS, and Linux.

## Installation

To install the Fishy CLI Tool, follow these steps:

1. Clone the repository:

    ```bash
    git clone https://github.com/ciathefed/fishy
    ```

2. Navigate to the project directory:

    ```bash
    cd fishy
    ```

3. Build the project:

   ```bash
   go build
   ```

## Usage

Add the `-h` option for more details

```bash
fishy <command> [options]
```

## Examples

For more examples, please check the [examples folder](https://github.com/ciathefed/fishy/tree/main/examples).

### Hello World

```asm
.section data
message:
    db "Hello, World!\n", 0

.section text
_start:
    mov x0, $1
    mov x1, message
    mov x2, $14
    mov x15, $4
    syscall

    mov x0, $0
    mov x15, $1
    syscall
```

To compile and run this program:

```bash
fishy build hello.fi
fishy run out.fbc
```

## Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/YourFeature`).
3. Make your changes and commit them (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/YourFeature`).
5. Open a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/ciathefed/fishy/blob/main/LICENSE) file for details.
