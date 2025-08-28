import { Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';

interface LoginResponse {
  token: string;
  user: { id: number; username: string };
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly baseUrl = 'http://localhost:8080';
  readonly token = signal<string | null>(this.getToken());
  readonly username = signal<string | null>(this.safeGet('username'));

  constructor(private http: HttpClient) {}

  register(username: string, password: string) {
    return this.http.post(`${this.baseUrl}/register`, { username, password });
  }

  login(username: string, password: string) {
    return this.http.post<LoginResponse>(`${this.baseUrl}/login`, { username, password });
  }

  setAuth(token: string, username: string) {
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('jwt', token);
      localStorage.setItem('username', username);
    }
    this.token.set(token);
    this.username.set(username);
  }

  logout() {
    if (typeof localStorage !== 'undefined') {
      localStorage.removeItem('jwt');
      localStorage.removeItem('username');
    }
    this.token.set(null);
    this.username.set(null);
  }

  private getToken(): string | null {
    return this.safeGet('jwt');
  }

  private safeGet(key: string): string | null {
    try {
      return typeof localStorage !== 'undefined' ? localStorage.getItem(key) : null;
    } catch {
      return null;
    }
  }
}
