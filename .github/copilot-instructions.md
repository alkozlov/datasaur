# Copilot Instructions for Go App

This is a Go 1.24.3 application using VS Code and Copilot Agent mode.

- Use Go idioms and standard library when feasible.
- Always pass `context.Context` in public functions.
- Follow project’s lint/style rules (e.g. gofmt, goimports).
- Write clear comments and use go doc conventions.
- Include error handling with wrapped errors (`fmt.Errorf("...: %w", err)`).
- Development environment is Windows 11 with Go 1.24.3.
- Use cmd for running commands, not PowerShell.
- Do not run application. You can build it only and suggest how to run it.

# Copilot instructions for React App

This is a React 18.2.0 application using VS Code and Copilot Agent mode.

- Use React idioms and standard library when feasible.
- Follow project’s lint/style rules (e.g. ESLint, Prettier).
- Write clear comments and use JSDoc conventions.
- Include error handling with descriptive messages.
- Use functional components and hooks.
- Use TypeScript for type safety.
- Use `npm` for package management.
- Do not run application. You can build it only and suggest how to run it.