import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { SourceFeed } from './models/feed-item.model';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class FeedService {
  constructor(private http: HttpClient) { }

  fetchFeeds(): Observable<SourceFeed[]> {
    return this.http.get<SourceFeed[]>('http://localhost:8080/api/feeds');
  }
}
