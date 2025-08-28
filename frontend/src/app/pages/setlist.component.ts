import { CommonModule } from '@angular/common';
import { Component, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router, ActivatedRoute } from '@angular/router';
import { SongsService, Song } from '../services/songs.service';
import { ConcertsService } from '../services/concerts.service';

@Component({
  selector: 'app-setlist',
  imports: [CommonModule, FormsModule],
  template: `
    <div class="row" style="align-items:center; margin-bottom: 16px">
      <h1 class="title">{{ concertTitle }} setlist</h1>
      <span class="spacer"></span>
      <button class="btn outline" (click)="goBack()">Back to Concerts</button>
    </div>

    <p *ngIf="error()" style="color:var(--danger)">{{ error() }}</p>

    <div class="card" style="margin-bottom: 16px">
      <h3 class="section-title">Add Song</h3>
      <form class="row" (ngSubmit)="addSong()">
        <input
          class="input"
          placeholder="Song title"
          [(ngModel)]="newSongTitle"
          name="title"
          required
          style="flex: 1;"
        />
        <input
          class="input"
          placeholder="Notes (optional)"
          [(ngModel)]="newSongNotes"
          name="notes"
          style="flex: 1;"
        />
        <button class="btn" type="submit">Add Song</button>
      </form>
    </div>

    <div class="card">
      <h3 class="section-title">Setlist Songs</h3>
      <div *ngIf="!songs().length" class="muted">
        No songs in your setlist yet. Add your first song above.
      </div>
      <div class="list">
        <div class="list-item" *ngFor="let song of songs(); trackBy: trackBySongId">
          <div style="flex: 1;">
            <div>
              <strong>{{ song.title }}</strong>
            </div>
            <div class="meta" *ngIf="song.notes">{{ song.notes }}</div>
          </div>
          <div class="actions-container">
            <div class="actions">
              <button class="btn ghost" (click)="removeSong(song.id)">Remove</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  `,
})
export class SetlistComponent {
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);
  private readonly songsService = inject(SongsService);
  private readonly concertsService = inject(ConcertsService);

  songs = signal<Song[]>([]);
  newSongTitle = '';
  newSongNotes = '';
  error = signal<string | null>(null);
  concertId: number | null = null;
  concertTitle = '';

  constructor() {
    // Get concert ID from route params
    this.route.params.subscribe((params) => {
      this.concertId = +params['id'];
      this.concertTitle = this.loadConcertTitle();
      if (this.concertId) {
        this.loadSetlist();
      }
    });
  }

  loadConcertTitle() {
    if (!this.concertId) return '';
    this.concertsService
      .get(this.concertId)
      .subscribe((concert) => (this.concertTitle = concert.title));
    return this.concertTitle;
  }

  loadSetlist() {
    if (!this.concertId) return;

    this.songsService.list(this.concertId).subscribe({
      next: (songs) => this.songs.set(songs),
      error: (err) => this.error.set(err?.error?.error || 'Failed to load songs'),
    });
  }

  addSong() {
    if (!this.newSongTitle.trim() || !this.concertId) return;

    this.songsService
      .create(this.concertId, {
        title: this.newSongTitle.trim(),
        notes: this.newSongNotes.trim(),
      })
      .subscribe({
        next: (newSong) => {
          this.songs.update((songs) => [...songs, newSong]);
          // Clear form
          this.newSongTitle = '';
          this.newSongNotes = '';
          this.loadSetlist();
        },
        error: (err) => this.error.set(err?.error?.error || 'Failed to add song'),
      });
  }

  removeSong(id: number) {
    if (!this.concertId) return;

    this.songsService.delete(this.concertId, id).subscribe({
      next: () => {
        this.songs.update((songs) => songs.filter((song) => song.id !== id));
      },
      error: (err) => this.error.set(err?.error?.error || 'Failed to remove song'),
    });
  }

  trackBySongId(index: number, song: Song): number {
    return song.id;
  }

  goBack() {
    this.router.navigate(['/concerts']);
  }
}
