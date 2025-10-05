// Commitlint configuration enforcing Conventional Commits.
// Docs: https://www.conventionalcommits.org/
// This keeps history machine-readable for future changelog automation and ERS traceability.
module.exports = {
  extends: ['@commitlint/config-conventional'],
  // Initial phase: scopes are OPTIONAL to keep friction low.
  // Later (after codebase modularizes) we can re-enable strict scopes.
  rules: {
    'scope-empty': [0, 'always'] // do not error if scope is omitted
  }
};
