# Required
Manage the complexity of your Product Requirements Document (PRD)

## Installation

```bash
go build
```

## Usage

### Initialize a new project

Create a new requirements project in the current directory:

```bash
reqd init
# or
reqd i
```

This creates a `requirements.yaml` file with your project structure.

### Add requirements

Add new requirements to your project:

```bash
reqd require "Your requirement text"
# or
reqd r "Your requirement text"
```

**Automatic Validation:**
When `OPENAI_API_KEY` is set, requirements are automatically validated using OpenAI's GPT-4o model to ensure they follow best practices (RFC 2119 keywords, clear language, etc.). Without an API key, validation is automatically skipped.

**Parent Proposal:**
When no parent is specified and `OPENAI_API_KEY` is available, the system can suggest an appropriate parent requirement from existing branch requirements (requirements that have children). This helps maintain a well-organized requirement hierarchy.

**Setup OPENAI_API_KEY:**
```bash
export OPENAI_API_KEY="your-api-key-here"
```

**Flags:**
- `--parent` or `-p`: Specify parent requirement ID for nested requirements
- `--no-validate` or `-V`: Skip validation even when API key is configured
- `--no-parent-proposal` or `-P`: Skip parent proposal feature

**Examples:**
```bash
# Basic requirement (validates automatically and may propose parent if API key is set)
reqd require "User must be able to login"

# Child requirement with explicit parent
reqd require "Login form must validate email format" --parent 1.1

# Skip validation and parent proposal
reqd require "Quick requirement" --no-validate --no-parent-proposal

# Skip only parent proposal (still validates if API key available)
reqd require "Another requirement" --no-parent-proposal
```

### Browse requirements

Display your requirements in a flat list format:

```bash
reqd show
# or
reqd s
```

This displays all requirements where hierarchy is indicated by dot-notation IDs (1, 1.1, 1.1.1).

**View specific requirement:**
```bash
reqd show <requirement_id>
```

Shows the specified requirement and all its children in the same flat format.

## File Structure

The tool creates and manages a `requirements.yaml` file with the following structure:

```yaml
name: Your Project Name
requirements:
  - id: "1"
    text: "Main requirement"
    children:
      - id: "1.1"
        text: "Sub-requirement"
```

## Commands

| Command | Alias | Description |
|---------|--------|-------------|
| `init` | `i` | Initialize a new requirements project |
| `require [text]` | `r` | Add a new requirement with optional validation |
| `show [id]` | `s` | Display requirements in flat list format |
