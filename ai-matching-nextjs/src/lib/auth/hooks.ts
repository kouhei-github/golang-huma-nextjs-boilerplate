import useSWR from 'swr';
import { apiClient } from '../../orval/apiClient';
import { saveAuthData, clearAuthData, getUser } from './tokenStorage';

interface LoginRequest {
  email: string;
  password: string;
}

interface AuthResponse {
  accessToken: string;
  idToken: string;
  refreshToken: string;
  tokenType: string;
  expiresAt: string;
  user: {
    id: string;
    email: string;
    firstName: string;
    lastName: string;
    tenants?: string[];
  };
  message?: string;
  requiresConfirmation?: boolean;
}

export function useAuth() {
  const { data: user, error, isLoading, mutate } = useSWR('auth/user', async () => {
    const userData = await getUser();
    if (!userData) {
      throw new Error('No user data found');
    }
    return userData;
  });

  const login = async (credentials: LoginRequest) => {
    try {
      const response = await apiClient.post<AuthResponse>('/auth/login', credentials);
      const { accessToken, idToken, refreshToken, tokenType, expiresAt, user: userData, requiresConfirmation } = response.data;

      if (requiresConfirmation) {
        return response.data;
      }

      await saveAuthData({
        tokens: {
          accessToken,
          idToken,
          refreshToken,
          tokenType,
          expiresAt,
        },
        user: userData
      });

      await mutate(userData);
      return response.data;
    } catch (error) {
      throw error;
    }
  };

  const logout = async () => {
    try {
      await apiClient.post('/auth/logout');
    } catch (error) {
      console.error('Logout API error:', error);
    } finally {
      await clearAuthData();
      await mutate(undefined, false);
    }
  };

  return {
    user,
    isLoading,
    isError: !!error,
    error,
    login,
    logout,
    mutate,
  };
}

export function useSWRWithAuth<T = any>(key: string | null, fetcher?: () => Promise<T>) {
  return useSWR<T>(
    key,
    fetcher || (async () => {
      const response = await apiClient.get<T>(key!);
      return response.data;
    }),
    {
      revalidateOnFocus: false,
      shouldRetryOnError: (error: any) => {
        if (error?.response?.status === 401) {
          return false;
        }
        return true;
      },
    }
  );
}