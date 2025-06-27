# Cognito Configuration Changes

## Summary of Changes

### 1. Auto-Confirm Users on Signup

Modified the Cognito client implementation to automatically confirm users after signup when `COGNITO_AUTO_CONFIRM` environment variable is set to `true`.

**Modified Files:**
- `/ai-matching-golang/src/infrastructure/external/cognito/cognito_client.go`
  - Added auto-confirmation logic in the `SignUp` method
  - Uses `AdminConfirmSignUp` API when `COGNITO_AUTO_CONFIRM=true`

- `/docker-compose.yml`
  - Added `COGNITO_AUTO_CONFIRM: "true"` to golang-app environment variables

- `/ai-matching-golang/.env.example`
  - Added `COGNITO_AUTO_CONFIRM=true` to document the new environment variable

### 2. Cleaned Up Cognito Directory

Removed redundant configuration files from the `/cognito` directory:
- Deleted `config.json`
- Deleted `db/` directory and its contents

The `/cognito/.cognito/` directory is preserved as it's used by cognito-local for runtime data storage.

## How It Works

When a user signs up:
1. The signup request is processed normally through Cognito
2. If `COGNITO_AUTO_CONFIRM=true`, the system automatically confirms the user
3. The signup response will have `UserConfirmed=true`
4. Users can immediately sign in without email confirmation

## Configuration

To enable auto-confirmation:
- Set `COGNITO_AUTO_CONFIRM=true` in your environment variables
- This is already configured in `docker-compose.yml` for local development

To disable auto-confirmation (require email verification):
- Set `COGNITO_AUTO_CONFIRM=false` or remove the environment variable
- Users will need to confirm their email before signing in