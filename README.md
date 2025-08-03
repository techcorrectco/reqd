# reqd
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

**Setup OPENAI_API_KEY:**
```bash
export OPENAI_API_KEY="your-api-key-here"
```

**Flags:**
- `--parent` or `-p`: Specify parent requirement ID for nested requirements
- `--no-validate` or `-n`: Skip validation even when API key is configured

**Examples:**
```bash
# Basic requirement (validates automatically if API key is set)
reqd require "User must be able to login"

# Child requirement
reqd require "Login form must validate email format" --parent 1.1

# Force skip validation even with API key
reqd require "Quick requirement" --no-validate
```

### Browse requirements

View and navigate your requirements interactively:

```bash
reqd show
# or
reqd s
```

**Navigation features:**
- Select requirements with arrow keys and Enter
- Use `..` option to navigate back to parent requirements
- Requirements with children are marked with `+` and show their children when selected
- Requirements without children show Edit/Quit options

**Editing:**
- Select "Edit" to modify a requirement's title inline
- Changes are saved automatically to `requirements.yaml`

### View specific requirement children

Start browsing from a specific requirement:

```bash
reqd show <requirement_id>
```

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
| `require` | `r` | Add a new requirement interactively |
| `show` | `s` | Browse requirements with interactive navigation |
