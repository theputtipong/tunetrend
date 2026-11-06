export interface PlaylistItem {
  videoId: string;
  position: number;
  title: string;
  channelTitle: string;
  thumbnailUrl: string;
  addedAt: string;
}

export interface Playlist {
  id: string;
  name: string;
  createdAt: string;
  updatedAt: string;
  itemCount: number;
  items?: PlaylistItem[];
}

export interface PlaylistVideo {
  videoId: string;
  title: string;
  channelTitle: string;
  thumbnailUrl: string;
}
