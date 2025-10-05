// Commitlint configuration enforcing Conventional Commits.
// Docs: https://www.conventionalcommits.org/
// This keeps history machine-readable for future changelog automation and ERS traceability.
module.exports = {
  extends: ['@commitlint/config-conventional'],
  // You can loosen or tighten rules later; start with defaults.
  rules: {
    // Example: allow chore scope optional; adjust when modules emerge.
    'scope-empty': [2, 'never'],
    // Temporarily allow no scope (set to 2 'never' later). Comment out above line if scopes not desired yet.
    // 'scope-empty': [0, 'always'],
  }
};
