import { openDB, DBSchema, IDBPDatabase } from 'idb';

interface AuthDB extends DBSchema {
  tokens: {
    key: string;
    value: {
      accessToken: string;
      idToken: string;
      refreshToken: string;
      expiresAt: string;
      tokenType: string;
    };
  };
  user: {
    key: string;
    value: {
      id: string;
      email: string;
      firstName: string;
      lastName: string;
      tenants?: string[];
    };
  };
}

const DB_NAME = 'AuthDatabase';
const DB_VERSION = 1;
const TOKEN_STORE = 'tokens';
const USER_STORE = 'user';
const TOKEN_KEY = 'auth_tokens';
const USER_KEY = 'current_user';

let db: IDBPDatabase<AuthDB> | null = null;

async function getDB(): Promise<IDBPDatabase<AuthDB>> {
  if (!db) {
    db = await openDB<AuthDB>(DB_NAME, DB_VERSION, {
      upgrade(database) {
        if (!database.objectStoreNames.contains(TOKEN_STORE)) {
          database.createObjectStore(TOKEN_STORE);
        }
        if (!database.objectStoreNames.contains(USER_STORE)) {
          database.createObjectStore(USER_STORE);
        }
      },
    });
  }
  return db;
}

export interface TokenData {
  accessToken: string;
  idToken: string;
  refreshToken: string;
  expiresAt: string;
  tokenType: string;
}

export interface UserData {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  tenants?: string[];
}

export interface AuthData {
  tokens: TokenData;
  user: UserData;
}

export async function saveAuthData(authData: AuthData): Promise<void> {
  const database = await getDB();
  await database.put(TOKEN_STORE, authData.tokens, TOKEN_KEY);
  await database.put(USER_STORE, authData.user, USER_KEY);
}

export async function getTokens(): Promise<TokenData | null> {
  const database = await getDB();
  const tokens = await database.get(TOKEN_STORE, TOKEN_KEY);
  return tokens || null;
}

export async function getUser(): Promise<UserData | null> {
  const database = await getDB();
  const user = await database.get(USER_STORE, USER_KEY);
  return user || null;
}

export async function getAccessToken(): Promise<string | null> {
  const tokens = await getTokens();
  return tokens?.accessToken || null;
}

export async function getRefreshToken(): Promise<string | null> {
  const tokens = await getTokens();
  return tokens?.refreshToken || null;
}

export async function getIdToken(): Promise<string | null> {
  const tokens = await getTokens();
  return tokens?.idToken || null;
}

export async function clearAuthData(): Promise<void> {
  const database = await getDB();
  await database.delete(TOKEN_STORE, TOKEN_KEY);
  await database.delete(USER_STORE, USER_KEY);
}

export async function isTokenExpired(): Promise<boolean> {
  const tokens = await getTokens();
  if (!tokens?.expiresAt) {
    return false;
  }
  return new Date() >= new Date(tokens.expiresAt);
}

export async function updateTokens(newTokens: Partial<TokenData>): Promise<void> {
  const currentTokens = await getTokens();
  const currentUser = await getUser();
  if (currentTokens && currentUser) {
    await saveAuthData({ 
      tokens: { ...currentTokens, ...newTokens },
      user: currentUser
    });
  } else {
    throw new Error('No auth data found to update');
  }
}