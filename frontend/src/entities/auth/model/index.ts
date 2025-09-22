export interface User {
  id?: string;
  name?: string;
  email?: string;
  picture?: string;
  email_verified?: boolean;
  sub?: string;
}

export interface AuthState {
  isAuthenticated: boolean;
  isLoading: boolean;
  user?: User;
  error?: Error;
}

export interface AuthConfig {
  domain: string;
  clientId: string;
  redirectUri: string;
  audience: string;
  scope: string;
}
