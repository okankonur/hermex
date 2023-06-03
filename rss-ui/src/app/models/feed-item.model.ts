export interface FeedItem {
    Title: string;
    Description: string;
    Link: string;
  }

  export interface SourceFeed {
    Host: string
    Favicon: string;
    Items: FeedItem[];
  }