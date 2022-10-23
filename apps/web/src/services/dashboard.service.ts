import type { AuthUser, Dashboard } from '@tsuwari/shared';

/**
 * @param user
 *
 * The user can have access to many dashboards, but his own dashboard
 * is not in the list provided by the api. So we create his dashboard
 * on the client. We do this to avoid unnecessary duplication of data.
 *
 * @returns user dashboard
 */
export const createUserDashboard = (user: AuthUser): Dashboard => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const { dashboards: _dashboards, ...userData } = user;

  return {
    id: '0',
    channelId: user.id,
    userId: user.id,
    twitch: userData,
  };
};
