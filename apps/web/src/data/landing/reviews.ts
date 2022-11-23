// import TextAvatarUrl from '@/assets/avatar.png';

export interface Review {
  id: number;
  username: string;
  comment: string;
  avatarUrl: string;
  rating: number;
}

export const reviews: Review[] = [
  {
    id: 1,
    username: '7ssk7',
    comment: `I've been using the bot for a few years now, and I'm happy with it. There are convenient integrations with Volaroant, Spotify. It's easy to add and remove commands from the chat. I am pleased with its stability and functionality.`,
    avatarUrl:
      'https://static-cdn.jtvnw.net/jtv_user_pictures/66cb7060-1a8a-4fca-9ccd-f760b70af636-profile_image-70x70.png',
    rating: 5,
  },
];
