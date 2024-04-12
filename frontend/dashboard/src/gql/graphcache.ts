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
  Upload: { input: any; output: any; }
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
  fileUrl: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  /** IDS of users which has this badge */
  users?: Maybe<Array<Scalars['String']['output']>>;
};

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

export type Mutation = {
  __typename?: 'Mutation';
  badgesAddUser: Scalars['Boolean']['output'];
  badgesCreate: Badge;
  badgesDelete: Scalars['Boolean']['output'];
  badgesRemoveUser: Scalars['Boolean']['output'];
  badgesUpdate: Badge;
  createCommand: Command;
  notificationsCreate: AdminNotification;
  notificationsDelete: Scalars['Boolean']['output'];
  notificationsUpdate: AdminNotification;
  removeCommand: Scalars['Boolean']['output'];
  switchUserAdmin: Scalars['Boolean']['output'];
  switchUserBan: Scalars['Boolean']['output'];
  updateCommand: Command;
};


export type MutationBadgesAddUserArgs = {
  id: Scalars['ID']['input'];
  userId: Scalars['String']['input'];
};


export type MutationBadgesCreateArgs = {
  file: Scalars['Upload']['input'];
  name: Scalars['String']['input'];
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
  notificationsByAdmin: AdminNotificationsResponse;
  notificationsByUser: Array<UserNotification>;
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
  /** `newNotification` will return a stream of `Notification` objects. */
  newNotification: UserNotification;
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

export type TwirBadgeUpdateOpts = {
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
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
  TwirAdminUser?: (data: WithTypename<TwirAdminUser>) => null | string,
  TwirUserTwitchInfo?: (data: WithTypename<TwirUserTwitchInfo>) => null | string,
  TwirUsersResponse?: (data: WithTypename<TwirUsersResponse>) => null | string,
  UserNotification?: (data: WithTypename<UserNotification>) => null | string
}

export type GraphCacheResolvers = {
  Query?: {
    authenticatedUser?: GraphCacheResolver<WithTypename<Query>, Record<string, never>, WithTypename<AuthenticatedUser> | string>,
    commands?: GraphCacheResolver<WithTypename<Query>, Record<string, never>, Array<WithTypename<Command> | string>>,
    notificationsByAdmin?: GraphCacheResolver<WithTypename<Query>, QueryNotificationsByAdminArgs, WithTypename<AdminNotificationsResponse> | string>,
    notificationsByUser?: GraphCacheResolver<WithTypename<Query>, Record<string, never>, Array<WithTypename<UserNotification> | string>>,
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
  notificationsCreate?: GraphCacheOptimisticMutationResolver<MutationNotificationsCreateArgs, WithTypename<AdminNotification>>,
  notificationsDelete?: GraphCacheOptimisticMutationResolver<MutationNotificationsDeleteArgs, Scalars['Boolean']>,
  notificationsUpdate?: GraphCacheOptimisticMutationResolver<MutationNotificationsUpdateArgs, WithTypename<AdminNotification>>,
  removeCommand?: GraphCacheOptimisticMutationResolver<MutationRemoveCommandArgs, Scalars['Boolean']>,
  switchUserAdmin?: GraphCacheOptimisticMutationResolver<MutationSwitchUserAdminArgs, Scalars['Boolean']>,
  switchUserBan?: GraphCacheOptimisticMutationResolver<MutationSwitchUserBanArgs, Scalars['Boolean']>,
  updateCommand?: GraphCacheOptimisticMutationResolver<MutationUpdateCommandArgs, WithTypename<Command>>
};

export type GraphCacheUpdaters = {
  Query?: {
    authenticatedUser?: GraphCacheUpdateResolver<{ authenticatedUser: WithTypename<AuthenticatedUser> }, Record<string, never>>,
    commands?: GraphCacheUpdateResolver<{ commands: Array<WithTypename<Command>> }, Record<string, never>>,
    notificationsByAdmin?: GraphCacheUpdateResolver<{ notificationsByAdmin: WithTypename<AdminNotificationsResponse> }, QueryNotificationsByAdminArgs>,
    notificationsByUser?: GraphCacheUpdateResolver<{ notificationsByUser: Array<WithTypename<UserNotification>> }, Record<string, never>>,
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
    notificationsCreate?: GraphCacheUpdateResolver<{ notificationsCreate: WithTypename<AdminNotification> }, MutationNotificationsCreateArgs>,
    notificationsDelete?: GraphCacheUpdateResolver<{ notificationsDelete: Scalars['Boolean'] }, MutationNotificationsDeleteArgs>,
    notificationsUpdate?: GraphCacheUpdateResolver<{ notificationsUpdate: WithTypename<AdminNotification> }, MutationNotificationsUpdateArgs>,
    removeCommand?: GraphCacheUpdateResolver<{ removeCommand: Scalars['Boolean'] }, MutationRemoveCommandArgs>,
    switchUserAdmin?: GraphCacheUpdateResolver<{ switchUserAdmin: Scalars['Boolean'] }, MutationSwitchUserAdminArgs>,
    switchUserBan?: GraphCacheUpdateResolver<{ switchUserBan: Scalars['Boolean'] }, MutationSwitchUserBanArgs>,
    updateCommand?: GraphCacheUpdateResolver<{ updateCommand: WithTypename<Command> }, MutationUpdateCommandArgs>
  },
  Subscription?: {
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
