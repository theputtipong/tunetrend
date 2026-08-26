export interface Song {
  id: string;
  title: string;
  channelTitle: string;
  thumbnailUrl: string;
  viewCount: string;
  countryCode: string;
  categoryId: string;
  publishedAt: string;
  videoType: string;
}

export interface TrendsResponse {
  success: boolean;
  data?: Song[];
  error?: string;
}
