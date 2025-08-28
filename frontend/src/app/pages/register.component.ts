import { Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../services/auth.service';

@Component({
  standalone: true,
  selector: 'app-register',
  imports: [CommonModule, FormsModule, RouterLink],
  template: `
    <div class="card" style="max-width:480px;margin:0 auto;">
      <h1 class="title">Register</h1>
      <p class="muted" style="margin:6px 0 16px">Create your account</p>
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
          <button class="btn" type="submit">Create account</button>
        </div>
        <p *ngIf="message()" style="color:#22c55e">{{ message() }}</p>
        <p *ngIf="error()" style="color:var(--danger)">{{ error() }}</p>
      </form>
      <p class="muted" style="margin-top:12px">
        Already have an account? <a routerLink="/login">Login</a>
      </p>
    </div>
  `,
})
export class RegisterComponent {
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);

  username = '';
  password = '';
  error = signal<string | null>(null);
  message = signal<string | null>(null);

  onSubmit() {
    this.error.set(null);
    this.message.set(null);
    this.auth.register(this.username, this.password).subscribe({
      next: () => {
        this.message.set('Registered! You can login now.');
        setTimeout(() => this.router.navigateByUrl('/login'), 500);
      },
      error: (err) => {
        this.error.set(err?.error?.error || 'Register failed');
      },
    });
  }
}
