# API Client Setup with Orval, SWR, and Token Management

## Overview

This setup provides:
- Automatic API client generation from OpenAPI specs using Orval
- Token storage in IndexedDB (access token, ID token, refresh token)
- User data storage in IndexedDB
- Automatic token refresh on 401 responses
- SWR integration for data fetching

## Usage

### 1. Generate API Client

```bash
npm run generate:api
```

This will fetch the OpenAPI spec from `http://localhost:8080/openapi.json` and generate typed API clients in `src/lib/api/generated/`.

### 2. Login Flow

```typescript
import { useAuth } from '@/lib/auth/hooks';

function LoginComponent() {
  const { login } = useAuth();
  
  const handleLogin = async () => {
    try {
      const response = await login({ email: 'user@example.com', password: 'password' });
      // Tokens and user data are automatically stored in IndexedDB
      // Response includes:
      // - accessToken, idToken, refreshToken
      // - user info (id, email, firstName, lastName, tenants)
      // - requiresConfirmation flag
    } catch (error) {
      console.error('Login failed:', error);
    }
  };
}
```

### 3. Making API Calls

Use the generated hooks from Orval:

```typescript
import { useGetTodos } from '@/lib/api/generated/todos/todos';

function TodoList() {
  const { data, error, isLoading } = useGetTodos();
  
  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;
  
  return (
    <ul>
      {data?.map(todo => <li key={todo.id}>{todo.title}</li>)}
    </ul>
  );
}
```

Or use the custom `useSWRWithAuth` hook:

```typescript
import { useSWRWithAuth } from '@/lib/auth/hooks';

function CustomData() {
  const { data, error, isLoading } = useSWRWithAuth('/api/custom-endpoint');
  // ...
}
```

### 4. Token Management

Tokens and user data are automatically managed:
- Stored in IndexedDB after login (separate stores for tokens and user data)
- Access token added to all API requests via Axios interceptor
- Automatically refreshed when receiving 401 responses
- Both tokens and user data cleared on logout

### 5. Access User Data

```typescript
import { useAuth } from '@/lib/auth/hooks';

function UserComponent() {
  const { user } = useAuth();
  
  if (user) {
    console.log(user.firstName, user.lastName);
    console.log(user.email);
    console.log(user.tenants); // Array of tenant names
  }
}
```

### 6. Environment Variables

Add to `.env.local`:

```
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

## File Structure

```
src/lib/
├── auth/
│   ├── apiClient.ts      # Axios instance with interceptors
│   ├── tokenStorage.ts   # IndexedDB token management
│   ├── hooks.ts          # Auth and SWR hooks
│   └── example-usage.tsx # Usage examples
└── api/
    └── generated/        # Orval generated files
        ├── models/       # TypeScript types
        └── */            # API endpoints grouped by tags
```