
# Contributing to pool-maintenance-app

**Note: This is a demo project created to showcase industry best practices and DevOps workflows. Contributions are welcome for demonstration and learning purposes.**

Thank you for your interest in contributing! We welcome all contributions that help improve this project.


## How to Contribute

1. **Fork the repository** and clone your fork locally.
2. **Create a new branch** from `main` for each feature or fix:
   ```sh
   git checkout main
   git pull
   git checkout -b feat/your-feature-name
   ```
3. **Make your changes** with clear, conventional commit messages.

4. **Test your changes locally** before pushing. You can:
   - **Run the app with Go:**
     ```sh
     go run ./cmd/main.go
     ```
   - **Run with Docker (Static musl/Alpine build):**
     ```sh
     docker build -t pool-maintenance-api .
     docker run -p 8080:8080 pool-maintenance-api
     ```
     > The Docker image is now based on Alpine Linux and contains a statically linked Go binary built with musl libc. This eliminates glibc version issues and ensures portability.
   - **Run the full CI pipeline locally with [`act`](https://github.com/nektos/act):**
     ```sh
     act -j build
     ```
     > Note: The artifact upload step is skipped locally, and Trivy/golangci-lint must be installed in the runner image. Security scanning may fail the build if vulnerabilities are found—this is intentional for best practices.
   - Visit [http://localhost:8080/health](http://localhost:8080/health) to verify the health check endpoint.
5. **Push your branch** to your fork and open a Pull Request (PR) to the main repository.
6. **Describe your changes** clearly in the PR description.
<<<<<<< Updated upstream

7. **Request a review** if needed. All PRs require at least one review and must pass CI checks before merging.
=======
7. **Request a review** if needed. All PRs require at least one review and must pass CI checks before merging (including lint, coverage, and security scan steps).
>>>>>>> Stashed changes
8. **After merging**, delete your feature branch if no longer needed.

## Continuous Integration (CI)

Our GitHub Actions pipeline will automatically:
- Lint your code with golangci-lint
- Run tests and measure code coverage
- Build and scan the Docker image and Go dependencies for vulnerabilities (Trivy)
- Run a health check endpoint smoke test

Please ensure your code passes all CI checks before requesting a review.

## Code Style
- Follow Go formatting standards (`go fmt`).
- Use clear, descriptive names and comments.
- Keep functions and files focused and modular.

## Commit Messages
- Use [Conventional Commits](https://www.conventionalcommits.org/):
  - `feat: add new feature`
  - `fix: correct a bug`
  - `docs: update documentation`
  - `ci: update CI configuration`
  - `test: add or update tests`
  - `chore: maintenance or tooling changes`

## Reporting Issues
- Use GitHub Issues to report bugs or request features.
- Provide as much detail as possible (steps to reproduce, expected behavior, etc.).

## Code of Conduct
This project follows a [Code of Conduct](CODE_OF_CONDUCT.md) to foster an open and welcoming environment.

## License
By contributing, you agree that your contributions will be licensed under the MIT License.

---
Thank you for helping make this project better!
