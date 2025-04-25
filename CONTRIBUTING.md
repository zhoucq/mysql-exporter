# Contributing to MySQL Exporter

Thank you for considering contributing to MySQL Exporter! This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

By participating in this project, you are expected to uphold our [Code of Conduct](CODE_OF_CONDUCT.md).

## How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report. Following these guidelines helps maintainers understand your report, reproduce the behavior, and find related reports.

Before creating bug reports, please check the issue list as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* **Use a clear and descriptive title** for the issue to identify the problem.
* **Describe the exact steps which reproduce the problem** in as many details as possible.
* **Provide specific examples to demonstrate the steps**. Include links to files or GitHub projects, or copy/pasteable snippets, which you use in those examples.
* **Describe the behavior you observed after following the steps** and point out what exactly is the problem with that behavior.
* **Explain which behavior you expected to see instead and why.**
* **Include screenshots and animated GIFs** which show you following the described steps and clearly demonstrate the problem.
* **If the problem wasn't triggered by a specific action**, describe what you were doing before the problem happened.

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion, including completely new features and minor improvements to existing functionality.

* **Use a clear and descriptive title** for the issue to identify the suggestion.
* **Provide a step-by-step description of the suggested enhancement** in as many details as possible.
* **Provide specific examples to demonstrate the steps**. Include copy/pasteable snippets which you use in those examples.
* **Describe the current behavior** and **explain which behavior you expected to see instead** and why.
* **Explain why this enhancement would be useful** to most MySQL Exporter users.

### Pull Requests

* Fill in the required template
* Do not include issue numbers in the PR title
* Include screenshots and animated GIFs in your pull request whenever possible
* Follow the Go style guide
* Include tests when adding features
* End all files with a newline

## Development Process

### Setting Up Development Environment

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/mysql-exporter.git`
3. Add the original repository as upstream: `git remote add upstream https://github.com/zhoucq/mysql-exporter.git`
4. Create a new branch for your changes: `git checkout -b feature/your-feature-name`

### Making Changes

1. Make your changes to the codebase
2. Run tests to ensure your changes don't break existing functionality
3. Commit your changes with a clear commit message
4. Push your changes to your fork
5. Submit a pull request to the main repository

### Code Style

Please follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) for style guidance.

## License

By contributing to MySQL Exporter, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).