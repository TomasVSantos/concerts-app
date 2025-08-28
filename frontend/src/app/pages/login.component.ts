import { Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../services/auth.service';

@Component({
  standalone: true,
  selector: 'app-login',
  imports: [CommonModule, FormsModule, RouterLink],
  template: `
    <div class="card" style="max-width:480px;margin:0 auto;">
      <h1 class="title">Login</h1>
      <p class="muted" style="margin:6px 0 16px">Access your concerts</p>
      <form (ngSubmit)="onSubmit()" class="row">
        <input
          class="input"
          placeholder="Username"
          name="username"
          [(ngModel)]="username"
          required
        />
        <input
          class="input"
          placeholder="Password"
          name="password"
          [(ngModel)]="password"
          type="password"
          required
        />
        <div class="row" style="justify-content: flex-end; width: 100%">
          <button class="btn" type="submit">Login</button>
        </div>
        <p *ngIf="error()" style="color:var(--danger)">{{ error() }}</p>
      </form>
      <p class="muted" style="margin-top:12px">
        No account? <a routerLink="/register">Register</a>
      </p>
    </div>
  `,
})
export class LoginComponent {
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);

  username = '';
  password = '';
  error = signal<string | null>(null);

  onSubmit() {
    this.error.set(null);
    this.auth.login(this.username, this.password).subscribe({
      next: (res) => {
        this.auth.setAuth(res.token, res.user.username);
        this.router.navigateByUrl('/concerts');
      },
      error: (err) => {
        this.error.set(err?.error?.error || 'Login failed');
      },
    });
  }
}
