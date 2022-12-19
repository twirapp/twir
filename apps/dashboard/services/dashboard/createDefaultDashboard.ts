import type { AuthUser, Dashboard } from '@tsuwari/shared';

/**
 * The user can have access to many dashboards, but his own dashboard
 * is not in the list provided by the api. So we create his dashboard
 * on the client. We do this to avoid unnecessary duplication of data.
 *
 * @returns user dashboard
 */
export const createDefaultDashboard = (user: AuthUser): Dashboard => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const { dashboards: _, ...userData } = user;

  return {
    id: '0',
    channelId: user.id,
    userId: user.id,
    twitchUser: userData,
  };
};
