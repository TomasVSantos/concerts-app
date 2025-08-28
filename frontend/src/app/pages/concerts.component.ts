import { Component, inject, signal, PLATFORM_ID } from '@angular/core';
import { CommonModule, isPlatformBrowser } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { ConcertsService, Concert } from '../services/concerts.service';
import { AuthService } from '../services/auth.service';

@Component({
  standalone: true,
  selector: 'app-concerts',
  imports: [CommonModule, FormsModule],
  template: `
    <div class="row" style="align-items:center; margin-bottom: 16px">
      <h1 class="title">Your concerts</h1>
      <span class="spacer"></span>
    </div>

    <div class="card" *ngIf="auth.token()" style="margin-bottom: 16px">
      <h3 class="section-title">Add concert</h3>
      <form class="row" (ngSubmit)="add()">
        <input class="input" placeholder="Title" [(ngModel)]="title" name="title" required />
        <input
          class="input"
          placeholder="Date (YYYY-MM-DD)"
          [(ngModel)]="date"
          name="date"
          required
        />
        <input
          class="input"
          placeholder="Location"
          [(ngModel)]="location"
          name="location"
          required
        />
        <span class="spacer"></span>
        <button class="btn" type="submit">Add</button>
      </form>
    </div>

    <p *ngIf="error()" style="color:var(--danger)">{{ error() }}</p>

    <div class="card">
      <h3 class="section-title">List</h3>
      <div *ngIf="!concerts().length" class="muted">No concerts yet. Add your first one.</div>
      <div class="list">
        <div class="list-item" *ngFor="let c of concerts()">
          <div>
            <div>
              <strong>{{ c.title }}</strong>
            </div>
            <div class="meta">{{ c.date }} • {{ c.location }}</div>
          </div>
          <div class="actions-container">
            <div class="actions">
              <button class="btn outline" (click)="editSetlist(c.id)">Edit Setlist</button>
            </div>
            <div class="actions">
              <button class="btn ghost" (click)="remove(c.id)">Delete</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  `,
})
export class ConcertsComponent {
  private readonly api = inject(ConcertsService);
  readonly auth = inject(AuthService);
  private readonly router = inject(Router);
  private readonly platformId = inject(PLATFORM_ID);

  concerts = signal<Concert[]>([]);
  error = signal<string | null>(null);

  title = '';
  date = '';
  location = '';

  constructor() {
    if (!this.auth.token()) {
      this.router.navigateByUrl('/login');
    }
    if (isPlatformBrowser(this.platformId)) {
      this.load();
    }
  }

  load() {
    this.api.list().subscribe({
      next: (list) => this.concerts.set(list),
      error: (err) => this.error.set(err?.error?.error || 'Failed to load'),
    });
  }

  add() {
    this.api.add({ title: this.title, date: this.date, location: this.location }).subscribe({
      next: (c) => {
        this.concerts.update((arr) => [c, ...arr]);
        this.title = this.date = this.location = '';
        this.load(); // Reload the list
      },
      error: (err) => this.error.set(err?.error?.error || 'Failed to add'),
    });
  }

  remove(id: number) {
    this.api.remove(id).subscribe({
      next: () => this.concerts.update((arr) => arr.filter((c) => c.id !== id)),
      error: (err) => this.error.set(err?.error?.error || 'Failed to delete'),
    });
  }

  editSetlist(id: number) {
    this.router.navigate(['/concerts', id, 'setlist']);
  }
}
