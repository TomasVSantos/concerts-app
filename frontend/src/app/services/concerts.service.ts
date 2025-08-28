import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';

export interface Concert {
  id: number;
  title: string;
  date: string;
  location: string;
  user_id: number;
}

@Injectable({ providedIn: 'root' })
export class ConcertsService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = 'http://localhost:8080/concerts';

  list() {
    return this.http.get<Concert[]>(this.baseUrl);
  }

  add(payload: { title: string; date: string; location: string }) {
    return this.http.post<Concert>(this.baseUrl, payload);
  }

  remove(id: number) {
    return this.http.delete<{ deleted: number }>(`${this.baseUrl}/${id}`);
  }
}
