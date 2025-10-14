# Rovers in Mars!

The suggested half-day to complete the exercise has been substantially increased due a couple of design changes mid-development and a desire to incorporate as much production grade testing as possible without using any AI to generate the tests: it was all manually implemented (yes all circa 900 lines of it - table test design allows for plenty of copy pasting but still).

I tried to demonstrate clean architecture, defensive programming, and Test-Driven Development using testify (assert, require and mock were chosen).

The core philosophy was not merely to solve the puzzle, but to engineer a solution as if it were a real-world, mission-critical application: a system that is maintainable, highly testable, and somehow resilient to invalid input.

The exercise gave ample room for decision making, so the following design choices were made:

**Plateau dimensions:**
By default, the plateau cannot be smaller than a 2 * 2 grid (but is configurable via cli switches)

**Rover placement and movement:**
Placing a rover at a location where another rover already exists will not be possible
As a rover moves, if it finds an obstacle (another previously deployed rover) it will still go as far as it can (have to make good use of that solar power in mars man!)

If a rover reaches a boundary, it will stop at the edge preventing from getting lost, crashing against environmental hazards or falling into unexplored terrain:
```Bash
2025/10/14 19:24:56 WARN: Rover 1 ignored move to (0 -1 S): position must be more than 0 and within boundaries
```

**Parsing the inputs:**
As a convenience feature, the parser will accept lowercase values (so n, e, s, w and l, r, m will be accepted)
White spaces (new-line, tabs and spaces) are trimmed


---

## ✨ Features

*   **Robust Domain Model:** The world of the simulation (Rovers, Plateau, Position) is modelled using clean, type-safe structs and "const enums", eliminating magic strings and ensuring correctness at the type level
*   **Decoupled Architecture:** A clear separation of concerns between the core `rover` domain, the `parser` for input handling, and the `app` orchestrator
*   **Comprehensive Error Handling:** Granular, sentinel errors provide clear, contextual feedback for all possible failure modes, from malformed input to in-flight collisions
*   **Extensive Unit & Integration Testing (with testify mocks):** The entire system is validated by a comprehensive suite of table-driven unit tests, proving the correctness of the logic and a generous amount of edge cases
*   **Clean Command-Line Interface:** The application runs as a standard CLI tool, accepting input from either a file (`-file` flag) or a `stdin` pipe, making it flexible and easy to integrate into scripts.

---
### 1. Installation & Setup

Clone the repository and tidy the dependencies:
```bash
git clone https://github.com/while-maybe/mars-rovers.git
cd mars-rovers
go mod tidy
```

### 2. Running the app

The application can be run in two ways: by providing a file path or by piping data to standard input.

#### **From a File**

Use the `-file` flag to specify an input file. An example `data.txt` is included.
```bash
go run ./cmd/cli -file data.txt
```

#### **From Standard Input (stdin)**

Pipe the input data directly to the application. This is useful for scripting and quick tests.
```bash
printf "5 5\n1 2 N\nLMLMLMLMM\n3 3 E\nMMRMMRMRRM" | go run ./cmd/cli
```
*Using `printf` or `echo -e` is recommended for correctly interpreting newline characters.*

#### **Expected Output**
For the proposed standard test case and regardless of the input method chosen, the output will be:

```
info: Mission complete. Final rover positions:
1 3 N
5 1 E
```

---

## 🛠️ Testing Strategy

This project was developed with a Test-Driven Development (TDD) mindset. The comprehensive test suite is a core feature, providing a guarantee of the application's correctness.

The included `Makefile` provides convenient commands for running the tests.

**Run all unit tests:**
```bash
make unit
```
This runs all tests within the `internal` directory.

**Run the integration test:**
```bash
make integration
```
This runs all tests within the `internal` directory.

**Generate and view test coverage:**
```bash
make coverage
```
This runs all tests, generates both a `coverage.out` profile and a detailed `coverage.html` coverage report which can be opened in a browser.

---

## 🏛️ Design & Architectural Decisions

My goal was to build a solution that is not only correct but also robust, maintainable, and demonstrates professional engineering practices without over-engineering. This informed several key architectural decisions.

### 1. Separation of Concerns

The application is split into three distinct packages, each with a single responsibility:

```
internal/
├── app       # Orchestrator
├── parser    # Input Adapter
└── rover     # Core Domain
```

*   **`rover`:** This is the heart of the application. It contains a pure, self-contained domain model with zero dependencies on other packages. It defines the "nouns" (`Rover`, `Plateau`) and "verbs" (`move`, `turnLeft`) of the simulation. Its logic is entirely independent of how the input is provided or how the output is displayed.

*   **`parser`:** This package acts as an "Input adapter." Its sole responsibility is to translate the raw, untrusted string-based input into the clean, validated, type-safe structs required by the `rover` domain. It enforces all formatting and value constraints at the boundary, ensuring that the core engine never receives invalid data.

*   **`app`:** This package holds the high-level application flow. It is decoupled from the concrete implementations of the parser and simulation engine via interfaces, demonstrating the principles of **Dependency Injection**. This makes the application's core workflow fully unit-testable in isolation.

### 2. Why Not Regular Expressions for Parsing?

I made a deliberate choice to use Go's standard library (`strings` and `strconv`) for parsing. This decision was based on several factors:

*   **Readability & Maintainability:** A sequence of `strings.Fields()` and `strconv.Atoi()` calls is a clear, step-by-step recipe that is easier for others to read and debug than  regex patterns.
*   **Error Granularity:** `strconv.Atoi()` returns a specific error for a malformed number, allowing the parser to return a much more helpful error message (e.g., `"invalid plateau width"`) than a generic "format mismatch" from a regex.
*   **Idiomatic Go:** For this type of simple, space-delimited text processing, using the standard library is the most common and idiomatic approach in the Go community.

### 3. Why Use Dependency Injection? Is It Overkill?

For a project of this size, introducing an `App` struct with interfaces might seem like overkill, but it was a conscious decision to showcase a professional approach to building testable software.

By having the `App` struct depend on a `Parser` interface, we can test the entire application flow (`App.Run`) without ever touching the real parser or the filesystem. We can inject a **mock parser** that returns pre-defined data or errors, allowing us to verify that the `App` orchestrates its components correctly under all conditions.

I don't think this is overkill, I see it as a more foundational pattern for building software that can be reliably tested and maintained should it ever need to grow in complexity.

I would be delighted to get feedback in any shape or form and I hope this matches both your expectations and technical needs.

Thanks for reading all of this!