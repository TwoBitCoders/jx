# jx: Command-Line JSON Processing with JavaScript Syntax, Powered by Go

Welcome to **jx**, a powerful command-line tool for JSON processing! Developed by **TwoBitCoders**, **jx** fuses the memory-safety and performance of Go with the flexibility of JavaScript syntax. Unlike other tools like jq that require learning a new domain-specific language, **jx** allows you to use your existing JavaScript knowledge to handle JSON with ease.

## Features

- **JavaScript Syntax**: Use familiar JavaScript constructs to filter, transform, and manipulate JSON data.
- **Memory Safe**: Built in Go, ensuring fast and secure processing, even with large datasets.
- **Cross-Platform**: Seamlessly runs on Linux, macOS, and Windows thanks to Go's portability.
- **High Performance**: Comparable to jq, but with the safety and simplicity of Go and JavaScript.
- **Single-binary install**: It's just a single binary, put it somewhere on your path, mark it executable and your good to go

## Why jx?

**jx** is your ultimate tool for JSON manipulation. Whether you're working with APIs, logs, or configuration files, **jx** provides an intuitive, command-line experience without sacrificing performance or security. With the power of Go and the flexibility of JavaScript, it’s ideal for developers who want to avoid learning a new syntax while gaining the safety of a memory-safe language.

**jx** is designed to be fast and user-friendly, eliminating unnecessary complexity from your JSON processing workflows.

### Key Advantages:
- **No Learning Curve**: Use familiar JavaScript syntax—no need to learn jq's DSL.
- **Efficient and Lightweight**: Go's compilation speed and memory safety ensure high performance.
- **No UI Overhead**: Designed to be a lean, CLI-only tool—no unnecessary graphical interface to bloat your workflow.
- **Robust and Portable**: Works out of the box across multiple platforms.

## Latest Release

For standalone installation, you can download the latest release directly from our [Releases Page](https://github.com/TwoBitCoders/jx/releases).

## Installation

To install **jx**, follow these steps:

1. **Clone the Repo**:
    ```sh
    git clone https://github.com/TwoBitCoders/jx
    ```
2. **Build jx**:
    ```sh
    cd jx
    go build -o out/
    ```
3. **Run the Simplest Command**:
    ```sh
    echo '{"foo":42}' | jx x
    ```
    'x' has your parsed JSON in it - so 'jx x' will pretty print your input data.

Alternatively, download pre-built binaries from our [releases page](https://github.com/TwoBitCoders/jx/releases) and add them to your `PATH`.

## Usage

**jx** makes JSON manipulation easy and accessible directly from your terminal.

All examples assume you're using a Bash shell, you are always passed variable 'x', it's value is the result of parsing your JSON.


- **Extract a Field**:
    ```sh
    echo '{"foo":42,"bar":0}' | jx 'x.foo'
    ```
- **Filter an Array**:
    ```sh
    echo '[{"foo":42},{"foo":0}]' | jx 'x.filter(item => item.foo !== 0)'
    ```
- **Transform Data**:
    ```sh
    echo '[{"foo":21}]' | jx 'x.map(item => ({ bar: item.foo * 2 }))'
    ```
- **Multiple Statements**:
    ```sh
    echo '{"foo":21}' | jx '{let op1 = x.foo;let op2 = 2;return op1*op2}'
    ```

Note: If you need multiple statements throw braces around your code, turning it into a "block body".

Note: For PowerShell Users on Windows: 
Set-ExecutionPolicy -ExecutionPolicy Unrestricted

For more detailed examples and advanced usage, check out the [Wiki](https://github.com/TwoBitCoders/jx/wiki/Advanced-Filtering-Techniques-with-jx).

## Contributing

We welcome contributions from the community. Here's how you can get involved:

- **Submit Issues**: Encountered a bug or have a feature request? Open an [issue](https://github.com/TwoBitCoders/jx/issues).
- **Create Pull Requests**: If you want to contribute code, check our contribution guidelines in the `CONTRIBUTING.md` file.
- **Join Discussions**: Discuss features, improvements, or usage questions [here](https://github.com/TwoBitCoders/jx/discussions).

## Support

Love **jx**? Consider supporting us:

- **[Become a Sponsor](https://patreon.com/TwoBitCoders)**: Help us keep improving **jx** with your generous support.
- For more tools and projects, visit the **[TwoBitCoders GitHub](https://github.com/TwoBitCoders)** page.

## License

**jx** is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

**TwoBitCoders**: Empowering developers, one tool at a time.

---

### Sources of Inspiration:
- [jq](https://github.com/stedolan/jq): The original inspiration for JSON processing tools.
- [ripgrep](https://github.com/BurntSushi/ripgrep): We admire the simplicity and speed of ripgrep, a key influence on **jx**'s design philosophy.
