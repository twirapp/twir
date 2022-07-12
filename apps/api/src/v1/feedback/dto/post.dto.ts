class FeedBackFile {
  name: string;
  id: string;
}

export interface FeedBackPostDto {
  email?: string;
  text: string;
  files?: FeedBackFile[];
}