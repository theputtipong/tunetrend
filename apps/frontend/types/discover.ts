export interface DiscoverItem {
  categoryId: string;
  categoryLabel: string;
  id: string;
  title: string;
  channelTitle: string;
  thumbnailUrl: string;
  viewCount: string;
  countryCode: string;
}

export interface DiscoverResponse {
  success: boolean;
  data?: DiscoverItem[];
  error?: string;
}
