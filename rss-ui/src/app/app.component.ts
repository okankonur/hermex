import { Component, OnInit } from '@angular/core';
import { FeedService } from './feed.service';
import { SourceFeed } from './models/feed-item.model';


@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})

export class AppComponent implements OnInit{
  feeds!: SourceFeed[];
  title = 'rss-frontend';

  constructor(private feedService: FeedService) {} 

  ngOnInit() {
    this.feedService.fetchFeeds().subscribe(feeds => {
      this.feeds = feeds;
    });
  }

}

