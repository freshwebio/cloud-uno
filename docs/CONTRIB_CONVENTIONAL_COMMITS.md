# Conventional Commits in Cloud Uno

## Commit types

**The following commit types are supported in the Cloud Uno project:**

- `fix:` - Should be used for any bug fixes.
- `feat:` - Should be used for any new features added, regardless of the size of the feature.
- `chore:` - Should be used for tasks such as releases or patching dependencies.
- `ci:` - Should be used for any work on GitHub Action workflows or scripts used in CI.
- `docs:`- Should be used for adding or modifying documentation.
- `style:` - Should be used for code formatting commits.
- `refactor:` - Should be used for any type of refactoring work that is not a part of a feature or bug fix.
- `perf:` - Should be used for a commit that represents performance improvements.
- `test:` - Should be used for commits that are purely for automated tests.

## Commit scopes

**The following commit scopes are supported:**

- `gcloud` - This commit scope should be used for a commit that represents work that pertains to the Google Cloud service emulators.
- `aws` - This commit scope should be used for a commit that represents work that pertains to the AWS service emulators.
- `azure` - This commit scope should be used for a commit that represents work that pertains to the Azure service emulators.
- `hostagent` - This commit scope should be used for commits that represent work that only affects the host agent.
- `server` - This commit scope should be used for commits that represent work that only affects the Cloud Uno server but is not specific to any specific cloud provider emulation.

The commit scope can be omitted for changes that cut across these scopes.
However, it's best to check in commits that map to a specific scope where possible.


## Commit footers

**The following custom footers are supported:**

- `Services: service1,service2` - This footer must be provided when a commit pertains to some work on emulators for one or more specific cloud services. 
  This helps with tooling and contributes to a rich, easy to navigate commit history making it fast to filter down to work on specific emulators.
  ***The service name should match the service names used in configuration and documented in the main README!***


## Example commit

### With commit scope

```bash
git commit -m 'feat(aws): add functionality to list secrets

The list secrets endpoint was excluded
in early versions of the project and was forgotten
about until recently.

Services: secretsmanager
'
```

### Without commit scope

```bash
git commit -m 'fix: correct proxy configuration for all services'
```
