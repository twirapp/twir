import { cacheExchange } from '@urql/exchange-graphcache';
import { Resolver as GraphCacheResolver, UpdateResolver as GraphCacheUpdateResolver, OptimisticMutationResolver as GraphCacheOptimisticMutationResolver } from '@urql/exchange-graphcache';

export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Time: { input: any; output: any; }
  Upload: { input: File; output: File; }
};

export type AdminNotification = Notification & {
  __typename?: 'AdminNotification';
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  text: Scalars['String']['output'];
  twitchProfile?: Maybe<TwirUserTwitchInfo>;
  userId?: Maybe<Scalars['ID']['output']>;
};

export type AdminNotificationsParams = {
  page?: InputMaybe<Scalars['Int']['input']>;
  perPage?: InputMaybe<Scalars['Int']['input']>;
  search?: InputMaybe<Scalars['String']['input']>;
  type?: InputMaybe<NotificationType>;
};

export type AdminNotificationsResponse = {
  __typename?: 'AdminNotificationsResponse';
  notifications: Array<AdminNotification>;
  total: Scalars['Int']['output'];
};

export type AuthenticatedUser = TwirUser & {
  __typename?: 'AuthenticatedUser';
  apiKey: Scalars['String']['output'];
  botId?: Maybe<Scalars['ID']['output']>;
  hideOnLandingPage: Scalars['Boolean']['output'];
  id: Scalars['ID']['output'];
  isBanned: Scalars['Boolean']['output'];
  isBotAdmin: Scalars['Boolean']['output'];
  isBotModerator?: Maybe<Scalars['Boolean']['output']>;
  isEnabled?: Maybe<Scalars['Boolean']['output']>;
  twitchProfile: TwirUserTwitchInfo;
};

export type Badge = {
  __typename?: 'Badge';
  createdAt: Scalars['String']['output'];
  enabled: Scalars['Boolean']['output'];
  ffzSlot: Scalars['Int']['output'];
  fileUrl: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  /** IDS of users which has this badge */
  users?: Maybe<Array<Scalars['String']['output']>>;
};

export enum ChannelRolePermissionEnum {
  CanAccessDashboard = 'CAN_ACCESS_DASHBOARD',
  ManageCommands = 'MANAGE_COMMANDS',
  ManageGreetings = 'MANAGE_GREETINGS',
  ManageIntegrations = 'MANAGE_INTEGRATIONS',
  ManageKeywords = 'MANAGE_KEYWORDS',
  ManageModeration = 'MANAGE_MODERATION',
  ManageOverlays = 'MANAGE_OVERLAYS',
  ManageSongRequests = 'MANAGE_SONG_REQUESTS',
  ManageTimers = 'MANAGE_TIMERS',
  ManageVariables = 'MANAGE_VARIABLES',
  UpdateChannelCategory = 'UPDATE_CHANNEL_CATEGORY',
  UpdateChannelTitle = 'UPDATE_CHANNEL_TITLE',
  ViewCommands = 'VIEW_COMMANDS',
  ViewGreetings = 'VIEW_GREETINGS',
  ViewIntegrations = 'VIEW_INTEGRATIONS',
  ViewKeywords = 'VIEW_KEYWORDS',
  ViewModeration = 'VIEW_MODERATION',
  ViewOverlays = 'VIEW_OVERLAYS',
  ViewSongRequests = 'VIEW_SONG_REQUESTS',
  ViewTimers = 'VIEW_TIMERS',
  ViewVariables = 'VIEW_VARIABLES'
}

export type Command = {
  __typename?: 'Command';
  aliases?: Maybe<Array<Scalars['String']['output']>>;
  allowedUsersIds?: Maybe<Array<Scalars['String']['output']>>;
  cooldown?: Maybe<Scalars['Int']['output']>;
  cooldownRolesIds?: Maybe<Array<Scalars['String']['output']>>;
  cooldownType: Scalars['String']['output'];
  default: Scalars['Boolean']['output'];
  defaultName?: Maybe<Scalars['String']['output']>;
  deniedUsersIds?: Maybe<Array<Scalars['String']['output']>>;
  description?: Maybe<Scalars['String']['output']>;
  enabled: Scalars['Boolean']['output'];
  enabledCategories?: Maybe<Array<Scalars['String']['output']>>;
  id: Scalars['ID']['output'];
  isReply: Scalars['Boolean']['output'];
  keepResponsesOrder: Scalars['Boolean']['output'];
  module: Scalars['String']['output'];
  name: Scalars['String']['output'];
  onlineOnly: Scalars['Boolean']['output'];
  requiredMessages: Scalars['Int']['output'];
  requiredUsedChannelPoints: Scalars['Int']['output'];
  requiredWatchTime: Scalars['Int']['output'];
  responses?: Maybe<Array<CommandResponse>>;
  rolesIds?: Maybe<Array<Scalars['String']['output']>>;
  visible: Scalars['Boolean']['output'];
};

export type CommandResponse = {
  __typename?: 'CommandResponse';
  commandId: Scalars['ID']['output'];
  id: Scalars['ID']['output'];
  order: Scalars['Int']['output'];
  text: Scalars['String']['output'];
};

export type CreateCommandInput = {
  aliases?: InputMaybe<Array<Scalars['String']['input']>>;
  description?: InputMaybe<Scalars['String']['input']>;
  name: Scalars['String']['input'];
  responses?: InputMaybe<Array<CreateCommandResponseInput>>;
};

export type CreateCommandResponseInput = {
  order: Scalars['Int']['input'];
  text: Scalars['String']['input'];
};

export type DashboardStats = {
  __typename?: 'DashboardStats';
  categoryId: Scalars['ID']['output'];
  categoryName: Scalars['String']['output'];
  chatMessages: Scalars['Int']['output'];
  followers: Scalars['Int']['output'];
  requestedSongs: Scalars['Int']['output'];
  startedAt?: Maybe<Scalars['Time']['output']>;
  subs: Scalars['Int']['output'];
  title: Scalars['String']['output'];
  usedEmotes: Scalars['Int']['output'];
  viewers?: Maybe<Scalars['Int']['output']>;
};

export type Greeting = {
  __typename?: 'Greeting';
  enabled: Scalars['Boolean']['output'];
  id: Scalars['ID']['output'];
  isReply: Scalars['Boolean']['output'];
  text: Scalars['String']['output'];
  twitchProfile: TwirUserTwitchInfo;
  userId: Scalars['String']['output'];
};

export type GreetingsCreateInput = {
  enabled: Scalars['Boolean']['input'];
  isReply: Scalars['Boolean']['input'];
  text: Scalars['String']['input'];
  userId: Scalars['String']['input'];
};

export type GreetingsUpdateInput = {
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  isReply?: InputMaybe<Scalars['Boolean']['input']>;
  text?: InputMaybe<Scalars['String']['input']>;
  userId?: InputMaybe<Scalars['String']['input']>;
};

export type Keyword = {
  __typename?: 'Keyword';
  cooldown: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  id: Scalars['ID']['output'];
  isRegularExpression: Scalars['Boolean']['output'];
  isReply: Scalars['Boolean']['output'];
  response?: Maybe<Scalars['String']['output']>;
  text: Scalars['String']['output'];
  usages: Scalars['Int']['output'];
};

export type KeywordCreateInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  isRegularExpression?: InputMaybe<Scalars['Boolean']['input']>;
  isReply?: InputMaybe<Scalars['Boolean']['input']>;
  response?: InputMaybe<Scalars['String']['input']>;
  text: Scalars['String']['input'];
  usageCount?: InputMaybe<Scalars['Int']['input']>;
};

export type KeywordUpdateInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  isRegularExpression?: InputMaybe<Scalars['Boolean']['input']>;
  isReply?: InputMaybe<Scalars['Boolean']['input']>;
  response?: InputMaybe<Scalars['String']['input']>;
  text?: InputMaybe<Scalars['String']['input']>;
  usageCount?: InputMaybe<Scalars['Int']['input']>;
};

export type Mutation = {
  __typename?: 'Mutation';
  badgesAddUser: Scalars['Boolean']['output'];
  badgesCreate: Badge;
  badgesDelete: Scalars['Boolean']['output'];
  badgesRemoveUser: Scalars['Boolean']['output'];
  badgesUpdate: Badge;
  createCommand: Command;
  greetingsCreate: Greeting;
  greetingsRemove: Scalars['Boolean']['output'];
  greetingsUpdate: Greeting;
  keywordCreate: Keyword;
  keywordRemove: Scalars['Boolean']['output'];
  keywordUpdate: Keyword;
  notificationsCreate: AdminNotification;
  notificationsDelete: Scalars['Boolean']['output'];
  notificationsUpdate: AdminNotification;
  removeCommand: Scalars['Boolean']['output'];
  switchUserAdmin: Scalars['Boolean']['output'];
  switchUserBan: Scalars['Boolean']['output'];
  timersCreate: Timer;
  timersRemove: Scalars['Boolean']['output'];
  timersUpdate: Timer;
  updateCommand: Command;
};


export type MutationBadgesAddUserArgs = {
  id: Scalars['ID']['input'];
  userId: Scalars['String']['input'];
};


export type MutationBadgesCreateArgs = {
  opts: TwirBadgeCreateOpts;
};


export type MutationBadgesDeleteArgs = {
  id: Scalars['ID']['input'];
};


export type MutationBadgesRemoveUserArgs = {
  id: Scalars['ID']['input'];
  userId: Scalars['String']['input'];
};


export type MutationBadgesUpdateArgs = {
  id: Scalars['ID']['input'];
  opts: TwirBadgeUpdateOpts;
};


export type MutationCreateCommandArgs = {
  opts: CreateCommandInput;
};


export type MutationGreetingsCreateArgs = {
  opts: GreetingsCreateInput;
};


export type MutationGreetingsRemoveArgs = {
  id: Scalars['String']['input'];
};


export type MutationGreetingsUpdateArgs = {
  id: Scalars['String']['input'];
  opts: GreetingsUpdateInput;
};


export type MutationKeywordCreateArgs = {
  opts: KeywordCreateInput;
};


export type MutationKeywordRemoveArgs = {
  id: Scalars['String']['input'];
};


export type MutationKeywordUpdateArgs = {
  id: Scalars['String']['input'];
  opts: KeywordUpdateInput;
};


export type MutationNotificationsCreateArgs = {
  text: Scalars['String']['input'];
  userId?: InputMaybe<Scalars['String']['input']>;
};


export type MutationNotificationsDeleteArgs = {
  id: Scalars['ID']['input'];
};


export type MutationNotificationsUpdateArgs = {
  id: Scalars['ID']['input'];
  opts: NotificationUpdateOpts;
};


export type MutationRemoveCommandArgs = {
  id: Scalars['String']['input'];
};


export type MutationSwitchUserAdminArgs = {
  userId: Scalars['ID']['input'];
};


export type MutationSwitchUserBanArgs = {
  userId: Scalars['ID']['input'];
};


export type MutationTimersCreateArgs = {
  opts: TimerCreateInput;
};


export type MutationTimersRemoveArgs = {
  id: Scalars['String']['input'];
};


export type MutationTimersUpdateArgs = {
  id: Scalars['String']['input'];
  opts: TimerUpdateInput;
};


export type MutationUpdateCommandArgs = {
  id: Scalars['String']['input'];
  opts: UpdateCommandOpts;
};

export type Notification = {
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  text: Scalars['String']['output'];
  userId?: Maybe<Scalars['ID']['output']>;
};

export enum NotificationType {
  Global = 'GLOBAL',
  User = 'USER'
}

export type NotificationUpdateOpts = {
  text?: InputMaybe<Scalars['String']['input']>;
};

export type Query = {
  __typename?: 'Query';
  authenticatedUser: AuthenticatedUser;
  commands: Array<Command>;
  greetings: Array<Greeting>;
  keywords: Array<Keyword>;
  notificationsByAdmin: AdminNotificationsResponse;
  notificationsByUser: Array<UserNotification>;
  timers: Array<Timer>;
  /** Twir badges */
  twirBadges: Array<Badge>;
  /** finding users on twitch with filter does they exists in database */
  twirUsers: TwirUsersResponse;
};


export type QueryNotificationsByAdminArgs = {
  opts: AdminNotificationsParams;
};


export type QueryTwirUsersArgs = {
  opts: TwirUsersSearchParams;
};

export type Subscription = {
  __typename?: 'Subscription';
  dashboardStats: DashboardStats;
  /** `newNotification` will return a stream of `Notification` objects. */
  newNotification: UserNotification;
};

export type Timer = {
  __typename?: 'Timer';
  enabled: Scalars['Boolean']['output'];
  id: Scalars['ID']['output'];
  messageInterval: Scalars['Int']['output'];
  name: Scalars['String']['output'];
  responses: Array<TimerResponse>;
  timeInterval: Scalars['Int']['output'];
};

export type TimerCreateInput = {
  enabled: Scalars['Boolean']['input'];
  messageInterval: Scalars['Int']['input'];
  name: Scalars['String']['input'];
  responses: Array<TimerResponseCreateInput>;
  timeInterval: Scalars['Int']['input'];
};

export type TimerResponse = {
  __typename?: 'TimerResponse';
  id: Scalars['ID']['output'];
  isAnnounce: Scalars['Boolean']['output'];
  text: Scalars['String']['output'];
};

export type TimerResponseCreateInput = {
  isAnnounce: Scalars['Boolean']['input'];
  text: Scalars['String']['input'];
};

export type TimerResponseUpdateInput = {
  isAnnounce: Scalars['Boolean']['input'];
  text: Scalars['String']['input'];
};

export type TimerUpdateInput = {
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  messageInterval?: InputMaybe<Scalars['Int']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  responses?: InputMaybe<Array<TimerResponseUpdateInput>>;
  timeInterval?: InputMaybe<Scalars['Int']['input']>;
};

export type TwirAdminUser = TwirUser & {
  __typename?: 'TwirAdminUser';
  apiKey: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  isBanned: Scalars['Boolean']['output'];
  isBotAdmin: Scalars['Boolean']['output'];
  isBotEnabled: Scalars['Boolean']['output'];
  isBotModerator: Scalars['Boolean']['output'];
  twitchProfile: TwirUserTwitchInfo;
};

export type TwirBadgeCreateOpts = {
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  ffzSlot: Scalars['Int']['input'];
  file: Scalars['Upload']['input'];
  name: Scalars['String']['input'];
};

export type TwirBadgeUpdateOpts = {
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  ffzSlot?: InputMaybe<Scalars['Int']['input']>;
  file?: InputMaybe<Scalars['Upload']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
};

export type TwirUser = {
  id: Scalars['ID']['output'];
  twitchProfile: TwirUserTwitchInfo;
};

export type TwirUserTwitchInfo = {
  __typename?: 'TwirUserTwitchInfo';
  description: Scalars['String']['output'];
  displayName: Scalars['String']['output'];
  login: Scalars['String']['output'];
  profileImageUrl: Scalars['String']['output'];
};

export type TwirUsersResponse = {
  __typename?: 'TwirUsersResponse';
  total: Scalars['Int']['output'];
  users: Array<TwirAdminUser>;
};

export type TwirUsersSearchParams = {
  badges?: InputMaybe<Array<Scalars['String']['input']>>;
  isBanned?: InputMaybe<Scalars['Boolean']['input']>;
  isBotAdmin?: InputMaybe<Scalars['Boolean']['input']>;
  isBotEnabled?: InputMaybe<Scalars['Boolean']['input']>;
  page?: InputMaybe<Scalars['Int']['input']>;
  perPage?: InputMaybe<Scalars['Int']['input']>;
  search?: InputMaybe<Scalars['String']['input']>;
};

export type UpdateCommandOpts = {
  aliases?: InputMaybe<Array<Scalars['String']['input']>>;
  allowedUsersIds?: InputMaybe<Array<Scalars['String']['input']>>;
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  cooldownRolesIds?: InputMaybe<Array<Scalars['String']['input']>>;
  cooldownType?: InputMaybe<Scalars['String']['input']>;
  deniedUsersIds?: InputMaybe<Array<Scalars['String']['input']>>;
  description?: InputMaybe<Scalars['String']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  enabledCategories?: InputMaybe<Array<Scalars['String']['input']>>;
  isReply?: InputMaybe<Scalars['Boolean']['input']>;
  keepResponsesOrder?: InputMaybe<Scalars['Boolean']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  onlineOnly?: InputMaybe<Scalars['Boolean']['input']>;
  requiredMessages?: InputMaybe<Scalars['Int']['input']>;
  requiredUsedChannelPoints?: InputMaybe<Scalars['Int']['input']>;
  requiredWatchTime?: InputMaybe<Scalars['Int']['input']>;
  responses?: InputMaybe<Array<CreateCommandResponseInput>>;
  rolesIds?: InputMaybe<Array<Scalars['String']['input']>>;
  visible?: InputMaybe<Scalars['Boolean']['input']>;
};

export type UserNotification = Notification & {
  __typename?: 'UserNotification';
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  text: Scalars['String']['output'];
  userId?: Maybe<Scalars['ID']['output']>;
};

export type WithTypename<T extends { __typename?: any }> = Partial<T> & { __typename: NonNullable<T['__typename']> };

export type GraphCacheKeysConfig = {
  AdminNotification?: (data: WithTypename<AdminNotification>) => null | string,
  AdminNotificationsResponse?: (data: WithTypename<AdminNotificationsResponse>) => null | string,
  AuthenticatedUser?: (data: WithTypename<AuthenticatedUser>) => null | string,
  Badge?: (data: WithTypename<Badge>) => null | string,
  Command?: (data: WithTypename<Command>) => null | string,
  CommandResponse?: (data: WithTypename<CommandResponse>) => null | string,
  DashboardStats?: (data: WithTypename<DashboardStats>) => null | string,
  Greeting?: (data: WithTypename<Greeting>) => null | string,
  Keyword?: (data: WithTypename<Keyword>) => null | string,
  Timer?: (data: WithTypename<Timer>) => null | string,
  TimerResponse?: (data: WithTypename<TimerResponse>) => null | string,
  TwirAdminUser?: (data: WithTypename<TwirAdminUser>) => null | string,
  TwirUserTwitchInfo?: (data: WithTypename<TwirUserTwitchInfo>) => null | string,
  TwirUsersResponse?: (data: WithTypename<TwirUsersResponse>) => null | string,
  UserNotification?: (data: WithTypename<UserNotification>) => null | string
}

export type GraphCacheResolvers = {
  Query?: {
    authenticatedUser?: GraphCacheResolver<WithTypename<Query>, Record<string, never>, WithTypename<AuthenticatedUser> | string>,
    commands?: GraphCacheResolver<WithTypename<Query>, Record<string, never>, Array<WithTypename<Command> | string>>,
    greetings?: GraphCacheResolver<WithTypename<Query>, Record<string, never>, Array<WithTypename<Greeting> | string>>,
    keywords?: GraphCacheResolver<WithTypename<Query>, Record<string, never>, Array<WithTypename<Keyword> | string>>,
    notificationsByAdmin?: GraphCacheResolver<WithTypename<Query>, QueryNotificationsByAdminArgs, WithTypename<AdminNotificationsResponse> | string>,
    notificationsByUser?: GraphCacheResolver<WithTypename<Query>, Record<string, never>, Array<WithTypename<UserNotification> | string>>,
    timers?: GraphCacheResolver<WithTypename<Query>, Record<string, never>, Array<WithTypename<Timer> | string>>,
    twirBadges?: GraphCacheResolver<WithTypename<Query>, Record<string, never>, Array<WithTypename<Badge> | string>>,
    twirUsers?: GraphCacheResolver<WithTypename<Query>, QueryTwirUsersArgs, WithTypename<TwirUsersResponse> | string>
  },
  AdminNotification?: {
    createdAt?: GraphCacheResolver<WithTypename<AdminNotification>, Record<string, never>, Scalars['Time'] | string>,
    id?: GraphCacheResolver<WithTypename<AdminNotification>, Record<string, never>, Scalars['ID'] | string>,
    text?: GraphCacheResolver<WithTypename<AdminNotification>, Record<string, never>, Scalars['String'] | string>,
    twitchProfile?: GraphCacheResolver<WithTypename<AdminNotification>, Record<string, never>, WithTypename<TwirUserTwitchInfo> | string>,
    userId?: GraphCacheResolver<WithTypename<AdminNotification>, Record<string, never>, Scalars['ID'] | string>
  },
  AdminNotificationsResponse?: {
    notifications?: GraphCacheResolver<WithTypename<AdminNotificationsResponse>, Record<string, never>, Array<WithTypename<AdminNotification> | string>>,
    total?: GraphCacheResolver<WithTypename<AdminNotificationsResponse>, Record<string, never>, Scalars['Int'] | string>
  },
  AuthenticatedUser?: {
    apiKey?: GraphCacheResolver<WithTypename<AuthenticatedUser>, Record<string, never>, Scalars['String'] | string>,
    botId?: GraphCacheResolver<WithTypename<AuthenticatedUser>, Record<string, never>, Scalars['ID'] | string>,
    hideOnLandingPage?: GraphCacheResolver<WithTypename<AuthenticatedUser>, Record<string, never>, Scalars['Boolean'] | string>,
    id?: GraphCacheResolver<WithTypename<AuthenticatedUser>, Record<string, never>, Scalars['ID'] | string>,
    isBanned?: GraphCacheResolver<WithTypename<AuthenticatedUser>, Record<string, never>, Scalars['Boolean'] | string>,
    isBotAdmin?: GraphCacheResolver<WithTypename<AuthenticatedUser>, Record<string, never>, Scalars['Boolean'] | string>,
    isBotModerator?: GraphCacheResolver<WithTypename<AuthenticatedUser>, Record<string, never>, Scalars['Boolean'] | string>,
    isEnabled?: GraphCacheResolver<WithTypename<AuthenticatedUser>, Record<string, never>, Scalars['Boolean'] | string>,
    twitchProfile?: GraphCacheResolver<WithTypename<AuthenticatedUser>, Record<string, never>, WithTypename<TwirUserTwitchInfo> | string>
  },
  Badge?: {
    createdAt?: GraphCacheResolver<WithTypename<Badge>, Record<string, never>, Scalars['String'] | string>,
    enabled?: GraphCacheResolver<WithTypename<Badge>, Record<string, never>, Scalars['Boolean'] | string>,
    ffzSlot?: GraphCacheResolver<WithTypename<Badge>, Record<string, never>, Scalars['Int'] | string>,
    fileUrl?: GraphCacheResolver<WithTypename<Badge>, Record<string, never>, Scalars['String'] | string>,
    id?: GraphCacheResolver<WithTypename<Badge>, Record<string, never>, Scalars['ID'] | string>,
    name?: GraphCacheResolver<WithTypename<Badge>, Record<string, never>, Scalars['String'] | string>,
    users?: GraphCacheResolver<WithTypename<Badge>, Record<string, never>, Array<Scalars['String'] | string>>
  },
  Command?: {
    aliases?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Array<Scalars['String'] | string>>,
    allowedUsersIds?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Array<Scalars['String'] | string>>,
    cooldown?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Scalars['Int'] | string>,
    cooldownRolesIds?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Array<Scalars['String'] | string>>,
    cooldownType?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Scalars['String'] | string>,
    default?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Scalars['Boolean'] | string>,
    defaultName?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Scalars['String'] | string>,
    deniedUsersIds?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Array<Scalars['String'] | string>>,
    description?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Scalars['String'] | string>,
    enabled?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Scalars['Boolean'] | string>,
    enabledCategories?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Array<Scalars['String'] | string>>,
    id?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Scalars['ID'] | string>,
    isReply?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Scalars['Boolean'] | string>,
    keepResponsesOrder?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Scalars['Boolean'] | string>,
    module?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Scalars['String'] | string>,
    name?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Scalars['String'] | string>,
    onlineOnly?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Scalars['Boolean'] | string>,
    requiredMessages?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Scalars['Int'] | string>,
    requiredUsedChannelPoints?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Scalars['Int'] | string>,
    requiredWatchTime?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Scalars['Int'] | string>,
    responses?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Array<WithTypename<CommandResponse> | string>>,
    rolesIds?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Array<Scalars['String'] | string>>,
    visible?: GraphCacheResolver<WithTypename<Command>, Record<string, never>, Scalars['Boolean'] | string>
  },
  CommandResponse?: {
    commandId?: GraphCacheResolver<WithTypename<CommandResponse>, Record<string, never>, Scalars['ID'] | string>,
    id?: GraphCacheResolver<WithTypename<CommandResponse>, Record<string, never>, Scalars['ID'] | string>,
    order?: GraphCacheResolver<WithTypename<CommandResponse>, Record<string, never>, Scalars['Int'] | string>,
    text?: GraphCacheResolver<WithTypename<CommandResponse>, Record<string, never>, Scalars['String'] | string>
  },
  DashboardStats?: {
    categoryId?: GraphCacheResolver<WithTypename<DashboardStats>, Record<string, never>, Scalars['ID'] | string>,
    categoryName?: GraphCacheResolver<WithTypename<DashboardStats>, Record<string, never>, Scalars['String'] | string>,
    chatMessages?: GraphCacheResolver<WithTypename<DashboardStats>, Record<string, never>, Scalars['Int'] | string>,
    followers?: GraphCacheResolver<WithTypename<DashboardStats>, Record<string, never>, Scalars['Int'] | string>,
    requestedSongs?: GraphCacheResolver<WithTypename<DashboardStats>, Record<string, never>, Scalars['Int'] | string>,
    startedAt?: GraphCacheResolver<WithTypename<DashboardStats>, Record<string, never>, Scalars['Time'] | string>,
    subs?: GraphCacheResolver<WithTypename<DashboardStats>, Record<string, never>, Scalars['Int'] | string>,
    title?: GraphCacheResolver<WithTypename<DashboardStats>, Record<string, never>, Scalars['String'] | string>,
    usedEmotes?: GraphCacheResolver<WithTypename<DashboardStats>, Record<string, never>, Scalars['Int'] | string>,
    viewers?: GraphCacheResolver<WithTypename<DashboardStats>, Record<string, never>, Scalars['Int'] | string>
  },
  Greeting?: {
    enabled?: GraphCacheResolver<WithTypename<Greeting>, Record<string, never>, Scalars['Boolean'] | string>,
    id?: GraphCacheResolver<WithTypename<Greeting>, Record<string, never>, Scalars['ID'] | string>,
    isReply?: GraphCacheResolver<WithTypename<Greeting>, Record<string, never>, Scalars['Boolean'] | string>,
    text?: GraphCacheResolver<WithTypename<Greeting>, Record<string, never>, Scalars['String'] | string>,
    twitchProfile?: GraphCacheResolver<WithTypename<Greeting>, Record<string, never>, WithTypename<TwirUserTwitchInfo> | string>,
    userId?: GraphCacheResolver<WithTypename<Greeting>, Record<string, never>, Scalars['String'] | string>
  },
  Keyword?: {
    cooldown?: GraphCacheResolver<WithTypename<Keyword>, Record<string, never>, Scalars['Int'] | string>,
    enabled?: GraphCacheResolver<WithTypename<Keyword>, Record<string, never>, Scalars['Boolean'] | string>,
    id?: GraphCacheResolver<WithTypename<Keyword>, Record<string, never>, Scalars['ID'] | string>,
    isRegularExpression?: GraphCacheResolver<WithTypename<Keyword>, Record<string, never>, Scalars['Boolean'] | string>,
    isReply?: GraphCacheResolver<WithTypename<Keyword>, Record<string, never>, Scalars['Boolean'] | string>,
    response?: GraphCacheResolver<WithTypename<Keyword>, Record<string, never>, Scalars['String'] | string>,
    text?: GraphCacheResolver<WithTypename<Keyword>, Record<string, never>, Scalars['String'] | string>,
    usages?: GraphCacheResolver<WithTypename<Keyword>, Record<string, never>, Scalars['Int'] | string>
  },
  Timer?: {
    enabled?: GraphCacheResolver<WithTypename<Timer>, Record<string, never>, Scalars['Boolean'] | string>,
    id?: GraphCacheResolver<WithTypename<Timer>, Record<string, never>, Scalars['ID'] | string>,
    messageInterval?: GraphCacheResolver<WithTypename<Timer>, Record<string, never>, Scalars['Int'] | string>,
    name?: GraphCacheResolver<WithTypename<Timer>, Record<string, never>, Scalars['String'] | string>,
    responses?: GraphCacheResolver<WithTypename<Timer>, Record<string, never>, Array<WithTypename<TimerResponse> | string>>,
    timeInterval?: GraphCacheResolver<WithTypename<Timer>, Record<string, never>, Scalars['Int'] | string>
  },
  TimerResponse?: {
    id?: GraphCacheResolver<WithTypename<TimerResponse>, Record<string, never>, Scalars['ID'] | string>,
    isAnnounce?: GraphCacheResolver<WithTypename<TimerResponse>, Record<string, never>, Scalars['Boolean'] | string>,
    text?: GraphCacheResolver<WithTypename<TimerResponse>, Record<string, never>, Scalars['String'] | string>
  },
  TwirAdminUser?: {
    apiKey?: GraphCacheResolver<WithTypename<TwirAdminUser>, Record<string, never>, Scalars['String'] | string>,
    id?: GraphCacheResolver<WithTypename<TwirAdminUser>, Record<string, never>, Scalars['ID'] | string>,
    isBanned?: GraphCacheResolver<WithTypename<TwirAdminUser>, Record<string, never>, Scalars['Boolean'] | string>,
    isBotAdmin?: GraphCacheResolver<WithTypename<TwirAdminUser>, Record<string, never>, Scalars['Boolean'] | string>,
    isBotEnabled?: GraphCacheResolver<WithTypename<TwirAdminUser>, Record<string, never>, Scalars['Boolean'] | string>,
    isBotModerator?: GraphCacheResolver<WithTypename<TwirAdminUser>, Record<string, never>, Scalars['Boolean'] | string>,
    twitchProfile?: GraphCacheResolver<WithTypename<TwirAdminUser>, Record<string, never>, WithTypename<TwirUserTwitchInfo> | string>
  },
  TwirUserTwitchInfo?: {
    description?: GraphCacheResolver<WithTypename<TwirUserTwitchInfo>, Record<string, never>, Scalars['String'] | string>,
    displayName?: GraphCacheResolver<WithTypename<TwirUserTwitchInfo>, Record<string, never>, Scalars['String'] | string>,
    login?: GraphCacheResolver<WithTypename<TwirUserTwitchInfo>, Record<string, never>, Scalars['String'] | string>,
    profileImageUrl?: GraphCacheResolver<WithTypename<TwirUserTwitchInfo>, Record<string, never>, Scalars['String'] | string>
  },
  TwirUsersResponse?: {
    total?: GraphCacheResolver<WithTypename<TwirUsersResponse>, Record<string, never>, Scalars['Int'] | string>,
    users?: GraphCacheResolver<WithTypename<TwirUsersResponse>, Record<string, never>, Array<WithTypename<TwirAdminUser> | string>>
  },
  UserNotification?: {
    createdAt?: GraphCacheResolver<WithTypename<UserNotification>, Record<string, never>, Scalars['Time'] | string>,
    id?: GraphCacheResolver<WithTypename<UserNotification>, Record<string, never>, Scalars['ID'] | string>,
    text?: GraphCacheResolver<WithTypename<UserNotification>, Record<string, never>, Scalars['String'] | string>,
    userId?: GraphCacheResolver<WithTypename<UserNotification>, Record<string, never>, Scalars['ID'] | string>
  }
};

export type GraphCacheOptimisticUpdaters = {
  badgesAddUser?: GraphCacheOptimisticMutationResolver<MutationBadgesAddUserArgs, Scalars['Boolean']>,
  badgesCreate?: GraphCacheOptimisticMutationResolver<MutationBadgesCreateArgs, WithTypename<Badge>>,
  badgesDelete?: GraphCacheOptimisticMutationResolver<MutationBadgesDeleteArgs, Scalars['Boolean']>,
  badgesRemoveUser?: GraphCacheOptimisticMutationResolver<MutationBadgesRemoveUserArgs, Scalars['Boolean']>,
  badgesUpdate?: GraphCacheOptimisticMutationResolver<MutationBadgesUpdateArgs, WithTypename<Badge>>,
  createCommand?: GraphCacheOptimisticMutationResolver<MutationCreateCommandArgs, WithTypename<Command>>,
  greetingsCreate?: GraphCacheOptimisticMutationResolver<MutationGreetingsCreateArgs, WithTypename<Greeting>>,
  greetingsRemove?: GraphCacheOptimisticMutationResolver<MutationGreetingsRemoveArgs, Scalars['Boolean']>,
  greetingsUpdate?: GraphCacheOptimisticMutationResolver<MutationGreetingsUpdateArgs, WithTypename<Greeting>>,
  keywordCreate?: GraphCacheOptimisticMutationResolver<MutationKeywordCreateArgs, WithTypename<Keyword>>,
  keywordRemove?: GraphCacheOptimisticMutationResolver<MutationKeywordRemoveArgs, Scalars['Boolean']>,
  keywordUpdate?: GraphCacheOptimisticMutationResolver<MutationKeywordUpdateArgs, WithTypename<Keyword>>,
  notificationsCreate?: GraphCacheOptimisticMutationResolver<MutationNotificationsCreateArgs, WithTypename<AdminNotification>>,
  notificationsDelete?: GraphCacheOptimisticMutationResolver<MutationNotificationsDeleteArgs, Scalars['Boolean']>,
  notificationsUpdate?: GraphCacheOptimisticMutationResolver<MutationNotificationsUpdateArgs, WithTypename<AdminNotification>>,
  removeCommand?: GraphCacheOptimisticMutationResolver<MutationRemoveCommandArgs, Scalars['Boolean']>,
  switchUserAdmin?: GraphCacheOptimisticMutationResolver<MutationSwitchUserAdminArgs, Scalars['Boolean']>,
  switchUserBan?: GraphCacheOptimisticMutationResolver<MutationSwitchUserBanArgs, Scalars['Boolean']>,
  timersCreate?: GraphCacheOptimisticMutationResolver<MutationTimersCreateArgs, WithTypename<Timer>>,
  timersRemove?: GraphCacheOptimisticMutationResolver<MutationTimersRemoveArgs, Scalars['Boolean']>,
  timersUpdate?: GraphCacheOptimisticMutationResolver<MutationTimersUpdateArgs, WithTypename<Timer>>,
  updateCommand?: GraphCacheOptimisticMutationResolver<MutationUpdateCommandArgs, WithTypename<Command>>
};

export type GraphCacheUpdaters = {
  Query?: {
    authenticatedUser?: GraphCacheUpdateResolver<{ authenticatedUser: WithTypename<AuthenticatedUser> }, Record<string, never>>,
    commands?: GraphCacheUpdateResolver<{ commands: Array<WithTypename<Command>> }, Record<string, never>>,
    greetings?: GraphCacheUpdateResolver<{ greetings: Array<WithTypename<Greeting>> }, Record<string, never>>,
    keywords?: GraphCacheUpdateResolver<{ keywords: Array<WithTypename<Keyword>> }, Record<string, never>>,
    notificationsByAdmin?: GraphCacheUpdateResolver<{ notificationsByAdmin: WithTypename<AdminNotificationsResponse> }, QueryNotificationsByAdminArgs>,
    notificationsByUser?: GraphCacheUpdateResolver<{ notificationsByUser: Array<WithTypename<UserNotification>> }, Record<string, never>>,
    timers?: GraphCacheUpdateResolver<{ timers: Array<WithTypename<Timer>> }, Record<string, never>>,
    twirBadges?: GraphCacheUpdateResolver<{ twirBadges: Array<WithTypename<Badge>> }, Record<string, never>>,
    twirUsers?: GraphCacheUpdateResolver<{ twirUsers: WithTypename<TwirUsersResponse> }, QueryTwirUsersArgs>
  },
  Mutation?: {
    badgesAddUser?: GraphCacheUpdateResolver<{ badgesAddUser: Scalars['Boolean'] }, MutationBadgesAddUserArgs>,
    badgesCreate?: GraphCacheUpdateResolver<{ badgesCreate: WithTypename<Badge> }, MutationBadgesCreateArgs>,
    badgesDelete?: GraphCacheUpdateResolver<{ badgesDelete: Scalars['Boolean'] }, MutationBadgesDeleteArgs>,
    badgesRemoveUser?: GraphCacheUpdateResolver<{ badgesRemoveUser: Scalars['Boolean'] }, MutationBadgesRemoveUserArgs>,
    badgesUpdate?: GraphCacheUpdateResolver<{ badgesUpdate: WithTypename<Badge> }, MutationBadgesUpdateArgs>,
    createCommand?: GraphCacheUpdateResolver<{ createCommand: WithTypename<Command> }, MutationCreateCommandArgs>,
    greetingsCreate?: GraphCacheUpdateResolver<{ greetingsCreate: WithTypename<Greeting> }, MutationGreetingsCreateArgs>,
    greetingsRemove?: GraphCacheUpdateResolver<{ greetingsRemove: Scalars['Boolean'] }, MutationGreetingsRemoveArgs>,
    greetingsUpdate?: GraphCacheUpdateResolver<{ greetingsUpdate: WithTypename<Greeting> }, MutationGreetingsUpdateArgs>,
    keywordCreate?: GraphCacheUpdateResolver<{ keywordCreate: WithTypename<Keyword> }, MutationKeywordCreateArgs>,
    keywordRemove?: GraphCacheUpdateResolver<{ keywordRemove: Scalars['Boolean'] }, MutationKeywordRemoveArgs>,
    keywordUpdate?: GraphCacheUpdateResolver<{ keywordUpdate: WithTypename<Keyword> }, MutationKeywordUpdateArgs>,
    notificationsCreate?: GraphCacheUpdateResolver<{ notificationsCreate: WithTypename<AdminNotification> }, MutationNotificationsCreateArgs>,
    notificationsDelete?: GraphCacheUpdateResolver<{ notificationsDelete: Scalars['Boolean'] }, MutationNotificationsDeleteArgs>,
    notificationsUpdate?: GraphCacheUpdateResolver<{ notificationsUpdate: WithTypename<AdminNotification> }, MutationNotificationsUpdateArgs>,
    removeCommand?: GraphCacheUpdateResolver<{ removeCommand: Scalars['Boolean'] }, MutationRemoveCommandArgs>,
    switchUserAdmin?: GraphCacheUpdateResolver<{ switchUserAdmin: Scalars['Boolean'] }, MutationSwitchUserAdminArgs>,
    switchUserBan?: GraphCacheUpdateResolver<{ switchUserBan: Scalars['Boolean'] }, MutationSwitchUserBanArgs>,
    timersCreate?: GraphCacheUpdateResolver<{ timersCreate: WithTypename<Timer> }, MutationTimersCreateArgs>,
    timersRemove?: GraphCacheUpdateResolver<{ timersRemove: Scalars['Boolean'] }, MutationTimersRemoveArgs>,
    timersUpdate?: GraphCacheUpdateResolver<{ timersUpdate: WithTypename<Timer> }, MutationTimersUpdateArgs>,
    updateCommand?: GraphCacheUpdateResolver<{ updateCommand: WithTypename<Command> }, MutationUpdateCommandArgs>
  },
  Subscription?: {
    dashboardStats?: GraphCacheUpdateResolver<{ dashboardStats: WithTypename<DashboardStats> }, Record<string, never>>,
    newNotification?: GraphCacheUpdateResolver<{ newNotification: WithTypename<UserNotification> }, Record<string, never>>
  },
  AdminNotification?: {
    createdAt?: GraphCacheUpdateResolver<Maybe<WithTypename<AdminNotification>>, Record<string, never>>,
    id?: GraphCacheUpdateResolver<Maybe<WithTypename<AdminNotification>>, Record<string, never>>,
    text?: GraphCacheUpdateResolver<Maybe<WithTypename<AdminNotification>>, Record<string, never>>,
    twitchProfile?: GraphCacheUpdateResolver<Maybe<WithTypename<AdminNotification>>, Record<string, never>>,
    userId?: GraphCacheUpdateResolver<Maybe<WithTypename<AdminNotification>>, Record<string, never>>
  },
  AdminNotificationsResponse?: {
    notifications?: GraphCacheUpdateResolver<Maybe<WithTypename<AdminNotificationsResponse>>, Record<string, never>>,
    total?: GraphCacheUpdateResolver<Maybe<WithTypename<AdminNotificationsResponse>>, Record<string, never>>
  },
  AuthenticatedUser?: {
    apiKey?: GraphCacheUpdateResolver<Maybe<WithTypename<AuthenticatedUser>>, Record<string, never>>,
    botId?: GraphCacheUpdateResolver<Maybe<WithTypename<AuthenticatedUser>>, Record<string, never>>,
    hideOnLandingPage?: GraphCacheUpdateResolver<Maybe<WithTypename<AuthenticatedUser>>, Record<string, never>>,
    id?: GraphCacheUpdateResolver<Maybe<WithTypename<AuthenticatedUser>>, Record<string, never>>,
    isBanned?: GraphCacheUpdateResolver<Maybe<WithTypename<AuthenticatedUser>>, Record<string, never>>,
    isBotAdmin?: GraphCacheUpdateResolver<Maybe<WithTypename<AuthenticatedUser>>, Record<string, never>>,
    isBotModerator?: GraphCacheUpdateResolver<Maybe<WithTypename<AuthenticatedUser>>, Record<string, never>>,
    isEnabled?: GraphCacheUpdateResolver<Maybe<WithTypename<AuthenticatedUser>>, Record<string, never>>,
    twitchProfile?: GraphCacheUpdateResolver<Maybe<WithTypename<AuthenticatedUser>>, Record<string, never>>
  },
  Badge?: {
    createdAt?: GraphCacheUpdateResolver<Maybe<WithTypename<Badge>>, Record<string, never>>,
    enabled?: GraphCacheUpdateResolver<Maybe<WithTypename<Badge>>, Record<string, never>>,
    ffzSlot?: GraphCacheUpdateResolver<Maybe<WithTypename<Badge>>, Record<string, never>>,
    fileUrl?: GraphCacheUpdateResolver<Maybe<WithTypename<Badge>>, Record<string, never>>,
    id?: GraphCacheUpdateResolver<Maybe<WithTypename<Badge>>, Record<string, never>>,
    name?: GraphCacheUpdateResolver<Maybe<WithTypename<Badge>>, Record<string, never>>,
    users?: GraphCacheUpdateResolver<Maybe<WithTypename<Badge>>, Record<string, never>>
  },
  Command?: {
    aliases?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    allowedUsersIds?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    cooldown?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    cooldownRolesIds?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    cooldownType?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    default?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    defaultName?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    deniedUsersIds?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    description?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    enabled?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    enabledCategories?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    id?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    isReply?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    keepResponsesOrder?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    module?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    name?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    onlineOnly?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    requiredMessages?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    requiredUsedChannelPoints?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    requiredWatchTime?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    responses?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    rolesIds?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>,
    visible?: GraphCacheUpdateResolver<Maybe<WithTypename<Command>>, Record<string, never>>
  },
  CommandResponse?: {
    commandId?: GraphCacheUpdateResolver<Maybe<WithTypename<CommandResponse>>, Record<string, never>>,
    id?: GraphCacheUpdateResolver<Maybe<WithTypename<CommandResponse>>, Record<string, never>>,
    order?: GraphCacheUpdateResolver<Maybe<WithTypename<CommandResponse>>, Record<string, never>>,
    text?: GraphCacheUpdateResolver<Maybe<WithTypename<CommandResponse>>, Record<string, never>>
  },
  DashboardStats?: {
    categoryId?: GraphCacheUpdateResolver<Maybe<WithTypename<DashboardStats>>, Record<string, never>>,
    categoryName?: GraphCacheUpdateResolver<Maybe<WithTypename<DashboardStats>>, Record<string, never>>,
    chatMessages?: GraphCacheUpdateResolver<Maybe<WithTypename<DashboardStats>>, Record<string, never>>,
    followers?: GraphCacheUpdateResolver<Maybe<WithTypename<DashboardStats>>, Record<string, never>>,
    requestedSongs?: GraphCacheUpdateResolver<Maybe<WithTypename<DashboardStats>>, Record<string, never>>,
    startedAt?: GraphCacheUpdateResolver<Maybe<WithTypename<DashboardStats>>, Record<string, never>>,
    subs?: GraphCacheUpdateResolver<Maybe<WithTypename<DashboardStats>>, Record<string, never>>,
    title?: GraphCacheUpdateResolver<Maybe<WithTypename<DashboardStats>>, Record<string, never>>,
    usedEmotes?: GraphCacheUpdateResolver<Maybe<WithTypename<DashboardStats>>, Record<string, never>>,
    viewers?: GraphCacheUpdateResolver<Maybe<WithTypename<DashboardStats>>, Record<string, never>>
  },
  Greeting?: {
    enabled?: GraphCacheUpdateResolver<Maybe<WithTypename<Greeting>>, Record<string, never>>,
    id?: GraphCacheUpdateResolver<Maybe<WithTypename<Greeting>>, Record<string, never>>,
    isReply?: GraphCacheUpdateResolver<Maybe<WithTypename<Greeting>>, Record<string, never>>,
    text?: GraphCacheUpdateResolver<Maybe<WithTypename<Greeting>>, Record<string, never>>,
    twitchProfile?: GraphCacheUpdateResolver<Maybe<WithTypename<Greeting>>, Record<string, never>>,
    userId?: GraphCacheUpdateResolver<Maybe<WithTypename<Greeting>>, Record<string, never>>
  },
  Keyword?: {
    cooldown?: GraphCacheUpdateResolver<Maybe<WithTypename<Keyword>>, Record<string, never>>,
    enabled?: GraphCacheUpdateResolver<Maybe<WithTypename<Keyword>>, Record<string, never>>,
    id?: GraphCacheUpdateResolver<Maybe<WithTypename<Keyword>>, Record<string, never>>,
    isRegularExpression?: GraphCacheUpdateResolver<Maybe<WithTypename<Keyword>>, Record<string, never>>,
    isReply?: GraphCacheUpdateResolver<Maybe<WithTypename<Keyword>>, Record<string, never>>,
    response?: GraphCacheUpdateResolver<Maybe<WithTypename<Keyword>>, Record<string, never>>,
    text?: GraphCacheUpdateResolver<Maybe<WithTypename<Keyword>>, Record<string, never>>,
    usages?: GraphCacheUpdateResolver<Maybe<WithTypename<Keyword>>, Record<string, never>>
  },
  Timer?: {
    enabled?: GraphCacheUpdateResolver<Maybe<WithTypename<Timer>>, Record<string, never>>,
    id?: GraphCacheUpdateResolver<Maybe<WithTypename<Timer>>, Record<string, never>>,
    messageInterval?: GraphCacheUpdateResolver<Maybe<WithTypename<Timer>>, Record<string, never>>,
    name?: GraphCacheUpdateResolver<Maybe<WithTypename<Timer>>, Record<string, never>>,
    responses?: GraphCacheUpdateResolver<Maybe<WithTypename<Timer>>, Record<string, never>>,
    timeInterval?: GraphCacheUpdateResolver<Maybe<WithTypename<Timer>>, Record<string, never>>
  },
  TimerResponse?: {
    id?: GraphCacheUpdateResolver<Maybe<WithTypename<TimerResponse>>, Record<string, never>>,
    isAnnounce?: GraphCacheUpdateResolver<Maybe<WithTypename<TimerResponse>>, Record<string, never>>,
    text?: GraphCacheUpdateResolver<Maybe<WithTypename<TimerResponse>>, Record<string, never>>
  },
  TwirAdminUser?: {
    apiKey?: GraphCacheUpdateResolver<Maybe<WithTypename<TwirAdminUser>>, Record<string, never>>,
    id?: GraphCacheUpdateResolver<Maybe<WithTypename<TwirAdminUser>>, Record<string, never>>,
    isBanned?: GraphCacheUpdateResolver<Maybe<WithTypename<TwirAdminUser>>, Record<string, never>>,
    isBotAdmin?: GraphCacheUpdateResolver<Maybe<WithTypename<TwirAdminUser>>, Record<string, never>>,
    isBotEnabled?: GraphCacheUpdateResolver<Maybe<WithTypename<TwirAdminUser>>, Record<string, never>>,
    isBotModerator?: GraphCacheUpdateResolver<Maybe<WithTypename<TwirAdminUser>>, Record<string, never>>,
    twitchProfile?: GraphCacheUpdateResolver<Maybe<WithTypename<TwirAdminUser>>, Record<string, never>>
  },
  TwirUserTwitchInfo?: {
    description?: GraphCacheUpdateResolver<Maybe<WithTypename<TwirUserTwitchInfo>>, Record<string, never>>,
    displayName?: GraphCacheUpdateResolver<Maybe<WithTypename<TwirUserTwitchInfo>>, Record<string, never>>,
    login?: GraphCacheUpdateResolver<Maybe<WithTypename<TwirUserTwitchInfo>>, Record<string, never>>,
    profileImageUrl?: GraphCacheUpdateResolver<Maybe<WithTypename<TwirUserTwitchInfo>>, Record<string, never>>
  },
  TwirUsersResponse?: {
    total?: GraphCacheUpdateResolver<Maybe<WithTypename<TwirUsersResponse>>, Record<string, never>>,
    users?: GraphCacheUpdateResolver<Maybe<WithTypename<TwirUsersResponse>>, Record<string, never>>
  },
  UserNotification?: {
    createdAt?: GraphCacheUpdateResolver<Maybe<WithTypename<UserNotification>>, Record<string, never>>,
    id?: GraphCacheUpdateResolver<Maybe<WithTypename<UserNotification>>, Record<string, never>>,
    text?: GraphCacheUpdateResolver<Maybe<WithTypename<UserNotification>>, Record<string, never>>,
    userId?: GraphCacheUpdateResolver<Maybe<WithTypename<UserNotification>>, Record<string, never>>
  },
};

export type GraphCacheConfig = Parameters<typeof cacheExchange>[0] & {
  updates?: GraphCacheUpdaters,
  keys?: GraphCacheKeysConfig,
  optimistic?: GraphCacheOptimisticUpdaters,
  resolvers?: GraphCacheResolvers,
};