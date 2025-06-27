import React, { useState } from 'react';
import { useAuth, useSWRWithAuth } from './hooks';

export function LoginForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const { login, isLoading } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await login({ email, password });
    } catch (error) {
      console.error('Login failed:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="Email"
        required
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Password"
        required
      />
      <button type="submit" disabled={isLoading}>
        {isLoading ? 'Loading...' : 'Login'}
      </button>
    </form>
  );
}

export function UserProfile() {
  const { user, logout, isLoading } = useAuth();

  if (isLoading) return <div>Loading...</div>;
  if (!user) return <div>Not authenticated</div>;

  return (
    <div>
      <h2>Welcome {user.firstName} {user.lastName}</h2>
      <p>Email: {user.email}</p>
      {user.tenants && user.tenants.length > 0 && (
        <div>
          <h3>Tenants:</h3>
          <ul>
            {user.tenants.map((tenant, index) => (
              <li key={index}>{tenant}</li>
            ))}
          </ul>
        </div>
      )}
      <button onClick={logout}>Logout</button>
    </div>
  );
}

interface Todo {
  id: string;
  title: string;
  completed: boolean;
}

export function TodoList() {
  const { data: todos, error, isLoading, mutate } = useSWRWithAuth<Todo[]>('/api/todos');

  if (isLoading) return <div>Loading todos...</div>;
  if (error) return <div>Error loading todos</div>;

  return (
    <ul>
      {todos?.map((todo) => (
        <li key={todo.id}>
          {todo.title} - {todo.completed ? 'Done' : 'Pending'}
        </li>
      ))}
    </ul>
  );
}