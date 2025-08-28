import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Song {
  id: number;
  title: string;
  notes: string;
  concert_id: number;
  order: number;
}

export interface CreateSongRequest {
  title: string;
  notes: string;
}

export interface SongOrderUpdate {
  song_id: number;
  order: number;
}

@Injectable({
  providedIn: 'root',
})
export class SongsService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = 'http://localhost:8080';

  list(concertId: number): Observable<Song[]> {
    return this.http.get<Song[]>(`${this.baseUrl}/concerts/${concertId}/songs`);
  }

  create(concertId: number, song: CreateSongRequest): Observable<Song> {
    return this.http.post<Song>(`${this.baseUrl}/concerts/${concertId}/songs`, song);
  }

  delete(concertId: number, songId: number): Observable<{ deleted: number }> {
    return this.http.delete<{ deleted: number }>(
      `${this.baseUrl}/concerts/${concertId}/songs/${songId}`
    );
  }

  updateOrder(concertId: number, updates: SongOrderUpdate[]): Observable<{ updated: number }> {
    return this.http.put<{ updated: number }>(
      `${this.baseUrl}/concerts/${concertId}/songs/order`,
      updates
    );
  }
}
