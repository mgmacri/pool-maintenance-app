
# Contributing to pool-maintenance-app

**Note: This is a demo project created to showcase industry best practices and DevOps workflows. Contributions are welcome for demonstration and learning purposes.**

Thank you for your interest in contributing! We welcome all contributions that help improve this project.


## How to Contribute

1. **Fork the repository** and clone your fork locally (or work directly if you have access).
2. **Create a new branch** from `develop` for each sub-chapter or focused change:
  ```sh
  git checkout develop
  git pull origin develop
  git checkout -b feature/short-descriptor
  ```
  Keep branches short‑lived and narrowly scoped (one sub-chapter / cohesive slice).
3. **Make your changes** with clear [Conventional Commit](https://www.conventionalcommits.org/) messages. A commitlint GitHub Action will validate messages on PRs.

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
   - **View the Swagger API docs:**
     - [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
   - **Regenerate Swagger docs after changing endpoint comments:**
     ```sh
     swag init -g cmd/main.go
     ```
     > This uses [swaggo/swag](https://github.com/swaggo/swag) to generate OpenAPI docs from Go comments. The generated files are in the `docs/` directory and are included in the Docker image.
5. **Push your branch** and open a Pull Request (PR) targeting `develop`.
6. **Describe your changes** clearly in the PR description. Reference any ERS IDs or plan steps (e.g., “Plan 1.1 step 3”).
7. **Request a review**. All PRs require at least one approval and must pass CI checks (lint, tests, security scan, commitlint) before merging.
8. **After merging**, delete your feature branch locally and remotely to keep the branch list clean.

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
Conventional Commits are enforced in CI. Format: `type(scope?): subject`.

Note: Scope is currently OPTIONAL (kept relaxed for early project velocity). We may enforce non-empty scopes once modules/packages stabilize.

Common types:
- `feat`: user-facing feature
- `fix`: bug fix
- `docs`: documentation only
- `test`: add or update tests
- `chore`: tooling / maintenance (no production code behavior change)
- `refactor`: code change that neither fixes a bug nor adds a feature
- `ci`: CI/CD pipeline changes
- `ops`: deployment or infrastructure assets

Rules (initial baseline):
- Use present tense imperative (“add”, not “adds” / “added”).
- Keep subject ≤ 72 chars.
- Body (optional) explains motivation, not just what.
- Reference plan steps or ERS IDs in body when relevant.

Examples:
```
feat: embed build metadata into binary
chore: add CODEOWNERS and PR template
docs: update README with build metadata usage (plan 1.1 step 6)
feat(health): extend /health with uptime_seconds   # optional scoped style
```

## Reporting Issues
- Use GitHub Issues to report bugs or request features.
- Provide as much detail as possible (steps to reproduce, expected behavior, etc.).

## Code of Conduct
This project follows a [Code of Conduct](CODE_OF_CONDUCT.md) to foster an open and welcoming environment.

## License
By contributing, you agree that your contributions will be licensed under the MIT License.

---
Thank you for helping make this project better!
