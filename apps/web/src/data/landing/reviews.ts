import TextAvatarUrl from '@/assets/avatar.png';

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
    username: 'random_usergsdagdsagsadgsda',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: TextAvatarUrl,
    rating: 4,
  },
  {
    id: 2,
    username: 'random_user',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: TextAvatarUrl,
    rating: 4,
  },
  {
    id: 3,
    username: 'random_user',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: TextAvatarUrl,
    rating: 2,
  },
  {
    id: 4,
    username: 'random_user',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: TextAvatarUrl,
    rating: 1,
  },
  {
    id: 5,
    username: 'random_user',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: TextAvatarUrl,
    rating: 4,
  },
  {
    id: 6,
    username: 'random_user',
    comment:
      'Praesent dolor quis aliquam nulla id in orci. Mi sit pulvinar nunc blandit egestas cras. Sed habitant amet ultrices vitae. At volutpat enim vel quam dignissim ut justo.',
    avatarUrl: TextAvatarUrl,
    rating: 4,
  },
];
