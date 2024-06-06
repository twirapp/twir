/* eslint-disable */
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
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
  availableDashboards: Array<Dashboard>;
  botId?: Maybe<Scalars['ID']['output']>;
  hideOnLandingPage: Scalars['Boolean']['output'];
  id: Scalars['ID']['output'];
  isBanned: Scalars['Boolean']['output'];
  isBotAdmin: Scalars['Boolean']['output'];
  isBotModerator?: Maybe<Scalars['Boolean']['output']>;
  isEnabled?: Maybe<Scalars['Boolean']['output']>;
  selectedDashboardId: Scalars['String']['output'];
  selectedDashboardTwitchUser: TwirUserTwitchInfo;
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

export type BuiltInVariable = {
  __typename?: 'BuiltInVariable';
  canBeUsedInRegistry: Scalars['Boolean']['output'];
  description: Scalars['String']['output'];
  example: Scalars['String']['output'];
  name: Scalars['String']['output'];
  visible: Scalars['Boolean']['output'];
};

export type ChannelAlert = {
  __typename?: 'ChannelAlert';
  audioId?: Maybe<Scalars['ID']['output']>;
  audioVolume?: Maybe<Scalars['Int']['output']>;
  commandIds?: Maybe<Array<Scalars['ID']['output']>>;
  greetingsIds?: Maybe<Array<Scalars['ID']['output']>>;
  id: Scalars['ID']['output'];
  keywordsIds?: Maybe<Array<Scalars['ID']['output']>>;
  name: Scalars['String']['output'];
  rewardIds?: Maybe<Array<Scalars['ID']['output']>>;
};

export type ChannelAlertCreateInput = {
  audioId?: InputMaybe<Scalars['ID']['input']>;
  audioVolume?: InputMaybe<Scalars['Int']['input']>;
  commandIds?: InputMaybe<Array<Scalars['ID']['input']>>;
  greetingsIds?: InputMaybe<Array<Scalars['ID']['input']>>;
  keywordsIds?: InputMaybe<Array<Scalars['ID']['input']>>;
  name: Scalars['String']['input'];
  rewardIds?: InputMaybe<Array<Scalars['ID']['input']>>;
};

export type ChannelAlertUpdateInput = {
  audioId?: InputMaybe<Scalars['ID']['input']>;
  audioVolume?: InputMaybe<Scalars['Int']['input']>;
  commandIds?: InputMaybe<Array<Scalars['ID']['input']>>;
  greetingsIds?: InputMaybe<Array<Scalars['ID']['input']>>;
  keywordsIds?: InputMaybe<Array<Scalars['ID']['input']>>;
  name?: InputMaybe<Scalars['String']['input']>;
  rewardIds?: InputMaybe<Array<Scalars['ID']['input']>>;
};

export enum ChannelRolePermissionEnum {
  CanAccessDashboard = 'CAN_ACCESS_DASHBOARD',
  ManageAlerts = 'MANAGE_ALERTS',
  ManageCommands = 'MANAGE_COMMANDS',
  ManageEvents = 'MANAGE_EVENTS',
  ManageGames = 'MANAGE_GAMES',
  ManageGreetings = 'MANAGE_GREETINGS',
  ManageIntegrations = 'MANAGE_INTEGRATIONS',
  ManageKeywords = 'MANAGE_KEYWORDS',
  ManageModeration = 'MANAGE_MODERATION',
  ManageOverlays = 'MANAGE_OVERLAYS',
  ManageRoles = 'MANAGE_ROLES',
  ManageSongRequests = 'MANAGE_SONG_REQUESTS',
  ManageTimers = 'MANAGE_TIMERS',
  ManageVariables = 'MANAGE_VARIABLES',
  UpdateChannelCategory = 'UPDATE_CHANNEL_CATEGORY',
  UpdateChannelTitle = 'UPDATE_CHANNEL_TITLE',
  ViewAlerts = 'VIEW_ALERTS',
  ViewCommands = 'VIEW_COMMANDS',
  ViewEvents = 'VIEW_EVENTS',
  ViewGames = 'VIEW_GAMES',
  ViewGreetings = 'VIEW_GREETINGS',
  ViewIntegrations = 'VIEW_INTEGRATIONS',
  ViewKeywords = 'VIEW_KEYWORDS',
  ViewModeration = 'VIEW_MODERATION',
  ViewOverlays = 'VIEW_OVERLAYS',
  ViewRoles = 'VIEW_ROLES',
  ViewSongRequests = 'VIEW_SONG_REQUESTS',
  ViewTimers = 'VIEW_TIMERS',
  ViewVariables = 'VIEW_VARIABLES'
}

export type ChatAlerts = {
  __typename?: 'ChatAlerts';
  ban?: Maybe<ChatAlertsBan>;
  chatCleared?: Maybe<ChatAlertsChatCleared>;
  cheers?: Maybe<ChatAlertsCheers>;
  donations?: Maybe<ChatAlertsDonations>;
  firstUserMessage?: Maybe<ChatAlertsFirstUserMessage>;
  followers?: Maybe<ChatAlertsFollowersSettings>;
  messageDelete?: Maybe<ChatAlertsMessageDelete>;
  raids?: Maybe<ChatAlertsRaids>;
  redemptions?: Maybe<ChatAlertsRedemptions>;
  streamOffline?: Maybe<ChatAlertsStreamOffline>;
  streamOnline?: Maybe<ChatAlertsStreamOnline>;
  subscribers?: Maybe<ChatAlertsSubscribers>;
  unbanRequestCreate?: Maybe<ChatAlertsUnbanRequestCreate>;
  unbanRequestResolve?: Maybe<ChatAlertsUnbanRequestResolve>;
};

export type ChatAlertsBan = {
  __typename?: 'ChatAlertsBan';
  cooldown: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  ignoreTimeoutFrom: Array<Scalars['String']['output']>;
  messages: Array<ChatAlertsCountedMessage>;
};

export type ChatAlertsBanInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  ignoreTimeoutFrom?: InputMaybe<Array<InputMaybe<Scalars['String']['input']>>>;
  messages?: InputMaybe<Array<InputMaybe<ChatAlertsCountedMessageInput>>>;
};

export type ChatAlertsChatCleared = {
  __typename?: 'ChatAlertsChatCleared';
  cooldown: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  messages: Array<ChatAlertsMessage>;
};

export type ChatAlertsChatClearedInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  messages?: InputMaybe<Array<InputMaybe<ChatAlertsMessageInput>>>;
};

export type ChatAlertsCheers = {
  __typename?: 'ChatAlertsCheers';
  cooldown: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  messages: Array<ChatAlertsCountedMessage>;
};

export type ChatAlertsCheersInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  messages?: InputMaybe<Array<InputMaybe<ChatAlertsCountedMessageInput>>>;
};

export type ChatAlertsCountedMessage = {
  __typename?: 'ChatAlertsCountedMessage';
  count: Scalars['Int']['output'];
  text: Scalars['String']['output'];
};

export type ChatAlertsCountedMessageInput = {
  count?: InputMaybe<Scalars['Int']['input']>;
  text?: InputMaybe<Scalars['String']['input']>;
};

export type ChatAlertsDonations = {
  __typename?: 'ChatAlertsDonations';
  cooldown: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  messages: Array<ChatAlertsCountedMessage>;
};

export type ChatAlertsDonationsInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  messages?: InputMaybe<Array<InputMaybe<ChatAlertsCountedMessageInput>>>;
};

export type ChatAlertsFirstUserMessage = {
  __typename?: 'ChatAlertsFirstUserMessage';
  cooldown: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  messages: Array<ChatAlertsMessage>;
};

export type ChatAlertsFirstUserMessageInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  messages?: InputMaybe<Array<InputMaybe<ChatAlertsMessageInput>>>;
};

export type ChatAlertsFollowersSettings = {
  __typename?: 'ChatAlertsFollowersSettings';
  cooldown: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  messages: Array<ChatAlertsMessage>;
};

export type ChatAlertsFollowersSettingsInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  messages?: InputMaybe<Array<InputMaybe<ChatAlertsMessageInput>>>;
};

export type ChatAlertsInput = {
  ban?: InputMaybe<ChatAlertsBanInput>;
  chatCleared?: InputMaybe<ChatAlertsChatClearedInput>;
  cheers?: InputMaybe<ChatAlertsCheersInput>;
  donations?: InputMaybe<ChatAlertsDonationsInput>;
  firstUserMessage?: InputMaybe<ChatAlertsFirstUserMessageInput>;
  followers?: InputMaybe<ChatAlertsFollowersSettingsInput>;
  messageDelete?: InputMaybe<ChatAlertsMessageDeleteInput>;
  raids?: InputMaybe<ChatAlertsRaidsInput>;
  redemptions?: InputMaybe<ChatAlertsRedemptionsInput>;
  streamOffline?: InputMaybe<ChatAlertsStreamOfflineInput>;
  streamOnline?: InputMaybe<ChatAlertsStreamOnlineInput>;
  subscribers?: InputMaybe<ChatAlertsSubscribersInput>;
  unbanRequestCreate?: InputMaybe<ChatAlertsUnbanRequestCreateInput>;
  unbanRequestResolve?: InputMaybe<ChatAlertsUnbanRequestResolveInput>;
};

export type ChatAlertsMessage = {
  __typename?: 'ChatAlertsMessage';
  text: Scalars['String']['output'];
};

export type ChatAlertsMessageDelete = {
  __typename?: 'ChatAlertsMessageDelete';
  cooldown: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  messages: Array<ChatAlertsMessage>;
};

export type ChatAlertsMessageDeleteInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  messages?: InputMaybe<Array<InputMaybe<ChatAlertsMessageInput>>>;
};

export type ChatAlertsMessageInput = {
  text?: InputMaybe<Scalars['String']['input']>;
};

export type ChatAlertsRaids = {
  __typename?: 'ChatAlertsRaids';
  cooldown: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  messages: Array<ChatAlertsCountedMessage>;
};

export type ChatAlertsRaidsInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  messages?: InputMaybe<Array<InputMaybe<ChatAlertsCountedMessageInput>>>;
};

export type ChatAlertsRedemptions = {
  __typename?: 'ChatAlertsRedemptions';
  cooldown: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  messages: Array<ChatAlertsMessage>;
};

export type ChatAlertsRedemptionsInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  messages?: InputMaybe<Array<InputMaybe<ChatAlertsMessageInput>>>;
};

export type ChatAlertsSettings = ChatAlertsBan | ChatAlertsChatCleared | ChatAlertsCheers | ChatAlertsDonations | ChatAlertsFirstUserMessage | ChatAlertsFollowersSettings | ChatAlertsMessageDelete | ChatAlertsRaids | ChatAlertsRedemptions | ChatAlertsStreamOffline | ChatAlertsStreamOnline | ChatAlertsSubscribers | ChatAlertsUnbanRequestCreate | ChatAlertsUnbanRequestResolve;

export type ChatAlertsStreamOffline = {
  __typename?: 'ChatAlertsStreamOffline';
  cooldown: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  messages: Array<ChatAlertsMessage>;
};

export type ChatAlertsStreamOfflineInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  messages?: InputMaybe<Array<InputMaybe<ChatAlertsMessageInput>>>;
};

export type ChatAlertsStreamOnline = {
  __typename?: 'ChatAlertsStreamOnline';
  cooldown: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  messages: Array<ChatAlertsMessage>;
};

export type ChatAlertsStreamOnlineInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  messages?: InputMaybe<Array<InputMaybe<ChatAlertsMessageInput>>>;
};

export type ChatAlertsSubscribers = {
  __typename?: 'ChatAlertsSubscribers';
  cooldown: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  messages: Array<ChatAlertsCountedMessage>;
};

export type ChatAlertsSubscribersInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  messages?: InputMaybe<Array<InputMaybe<ChatAlertsCountedMessageInput>>>;
};

export type ChatAlertsUnbanRequestCreate = {
  __typename?: 'ChatAlertsUnbanRequestCreate';
  cooldown: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  messages: Array<ChatAlertsMessage>;
};

export type ChatAlertsUnbanRequestCreateInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  messages?: InputMaybe<Array<InputMaybe<ChatAlertsMessageInput>>>;
};

export type ChatAlertsUnbanRequestResolve = {
  __typename?: 'ChatAlertsUnbanRequestResolve';
  cooldown: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  messages: Array<ChatAlertsMessage>;
};

export type ChatAlertsUnbanRequestResolveInput = {
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  messages?: InputMaybe<Array<InputMaybe<ChatAlertsMessageInput>>>;
};

export type ChatOverlay = {
  __typename?: 'ChatOverlay';
  animation?: Maybe<ChatOverlayAnimation>;
  chatBackgroundColor: Scalars['String']['output'];
  direction: Scalars['String']['output'];
  fontFamily: Scalars['String']['output'];
  fontSize: Scalars['Int']['output'];
  fontStyle: Scalars['String']['output'];
  fontWeight: Scalars['Int']['output'];
  hideBots: Scalars['Boolean']['output'];
  hideCommands: Scalars['Boolean']['output'];
  id: Scalars['String']['output'];
  messageHideTimeout: Scalars['Int']['output'];
  messageShowDelay: Scalars['Int']['output'];
  paddingContainer: Scalars['Int']['output'];
  preset: Scalars['String']['output'];
  showAnnounceBadge: Scalars['Boolean']['output'];
  showBadges: Scalars['Boolean']['output'];
  textShadowColor: Scalars['String']['output'];
  textShadowSize: Scalars['Int']['output'];
};

export enum ChatOverlayAnimation {
  Disabled = 'DISABLED'
}

export type ChatOverlayUpdateOpts = {
  animation?: InputMaybe<ChatOverlayAnimation>;
  chatBackgroundColor?: InputMaybe<Scalars['String']['input']>;
  direction?: InputMaybe<Scalars['String']['input']>;
  fontFamily?: InputMaybe<Scalars['String']['input']>;
  fontSize?: InputMaybe<Scalars['Int']['input']>;
  fontStyle?: InputMaybe<Scalars['String']['input']>;
  fontWeight?: InputMaybe<Scalars['Int']['input']>;
  hideBots?: InputMaybe<Scalars['Boolean']['input']>;
  hideCommands?: InputMaybe<Scalars['Boolean']['input']>;
  id: Scalars['String']['input'];
  messageHideTimeout?: InputMaybe<Scalars['Int']['input']>;
  messageShowDelay?: InputMaybe<Scalars['Int']['input']>;
  paddingContainer?: InputMaybe<Scalars['Int']['input']>;
  preset?: InputMaybe<Scalars['String']['input']>;
  showAnnounceBadge?: InputMaybe<Scalars['Boolean']['input']>;
  showBadges?: InputMaybe<Scalars['Boolean']['input']>;
  textShadowColor?: InputMaybe<Scalars['String']['input']>;
  textShadowSize?: InputMaybe<Scalars['Int']['input']>;
};

export type Command = {
  __typename?: 'Command';
  aliases: Array<Scalars['String']['output']>;
  allowedUsersIds: Array<Scalars['String']['output']>;
  cooldown: Scalars['Int']['output'];
  cooldownRolesIds: Array<Scalars['String']['output']>;
  cooldownType: Scalars['String']['output'];
  default: Scalars['Boolean']['output'];
  defaultName?: Maybe<Scalars['String']['output']>;
  deniedUsersIds: Array<Scalars['String']['output']>;
  description: Scalars['String']['output'];
  enabled: Scalars['Boolean']['output'];
  enabledCategories: Array<Scalars['String']['output']>;
  group?: Maybe<CommandGroup>;
  id: Scalars['ID']['output'];
  isReply: Scalars['Boolean']['output'];
  keepResponsesOrder: Scalars['Boolean']['output'];
  module: Scalars['String']['output'];
  name: Scalars['String']['output'];
  onlineOnly: Scalars['Boolean']['output'];
  requiredMessages: Scalars['Int']['output'];
  requiredUsedChannelPoints: Scalars['Int']['output'];
  requiredWatchTime: Scalars['Int']['output'];
  responses: Array<CommandResponse>;
  rolesIds: Array<Scalars['String']['output']>;
  visible: Scalars['Boolean']['output'];
};

export type CommandGroup = {
  __typename?: 'CommandGroup';
  color: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
};

export type CommandResponse = {
  __typename?: 'CommandResponse';
  commandId: Scalars['ID']['output'];
  id: Scalars['ID']['output'];
  order: Scalars['Int']['output'];
  text: Scalars['String']['output'];
};

export type CommandsCreateOpts = {
  aliases: Array<Scalars['String']['input']>;
  allowedUsersIds: Array<Scalars['String']['input']>;
  cooldown: Scalars['Int']['input'];
  cooldownRolesIds: Array<Scalars['String']['input']>;
  cooldownType: Scalars['String']['input'];
  deniedUsersIds: Array<Scalars['String']['input']>;
  description: Scalars['String']['input'];
  enabled: Scalars['Boolean']['input'];
  enabledCategories: Array<Scalars['String']['input']>;
  groupId?: InputMaybe<Scalars['String']['input']>;
  isReply: Scalars['Boolean']['input'];
  keepResponsesOrder: Scalars['Boolean']['input'];
  name: Scalars['String']['input'];
  onlineOnly: Scalars['Boolean']['input'];
  requiredMessages: Scalars['Int']['input'];
  requiredUsedChannelPoints: Scalars['Int']['input'];
  requiredWatchTime: Scalars['Int']['input'];
  responses: Array<CreateOrUpdateCommandResponseInput>;
  rolesIds: Array<Scalars['String']['input']>;
  visible: Scalars['Boolean']['input'];
};

export type CommandsGroupsCreateOpts = {
  color: Scalars['String']['input'];
  name: Scalars['String']['input'];
};

export type CommandsGroupsUpdateOpts = {
  color?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
};

export type CommandsUpdateOpts = {
  aliases?: InputMaybe<Array<Scalars['String']['input']>>;
  allowedUsersIds?: InputMaybe<Array<Scalars['String']['input']>>;
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  cooldownRolesIds?: InputMaybe<Array<Scalars['String']['input']>>;
  cooldownType?: InputMaybe<Scalars['String']['input']>;
  deniedUsersIds?: InputMaybe<Array<Scalars['String']['input']>>;
  description?: InputMaybe<Scalars['String']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  enabledCategories?: InputMaybe<Array<Scalars['String']['input']>>;
  groupId?: InputMaybe<Scalars['String']['input']>;
  isReply?: InputMaybe<Scalars['Boolean']['input']>;
  keepResponsesOrder?: InputMaybe<Scalars['Boolean']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  onlineOnly?: InputMaybe<Scalars['Boolean']['input']>;
  requiredMessages?: InputMaybe<Scalars['Int']['input']>;
  requiredUsedChannelPoints?: InputMaybe<Scalars['Int']['input']>;
  requiredWatchTime?: InputMaybe<Scalars['Int']['input']>;
  responses?: InputMaybe<Array<CreateOrUpdateCommandResponseInput>>;
  rolesIds?: InputMaybe<Array<Scalars['String']['input']>>;
  visible?: InputMaybe<Scalars['Boolean']['input']>;
};

export type CommunityUser = TwirUser & {
  __typename?: 'CommunityUser';
  id: Scalars['ID']['output'];
  messages: Scalars['Int']['output'];
  twitchProfile: TwirUserTwitchInfo;
  usedChannelPoints: Scalars['Int']['output'];
  usedEmotes: Scalars['Int']['output'];
  watchedMs: Scalars['Int']['output'];
};

export type CommunityUsersOpts = {
  channelId: Scalars['ID']['input'];
  order?: InputMaybe<CommunityUsersOrder>;
  page?: InputMaybe<Scalars['Int']['input']>;
  perPage?: InputMaybe<Scalars['Int']['input']>;
  search?: InputMaybe<Scalars['String']['input']>;
  sortBy?: InputMaybe<CommunityUsersSortBy>;
};

export enum CommunityUsersOrder {
  Asc = 'ASC',
  Desc = 'DESC'
}

export enum CommunityUsersResetType {
  Messages = 'MESSAGES',
  UsedChannelsPoints = 'USED_CHANNELS_POINTS',
  UsedEmotes = 'USED_EMOTES',
  Watched = 'WATCHED'
}

export type CommunityUsersResponse = {
  __typename?: 'CommunityUsersResponse';
  total: Scalars['Int']['output'];
  users: Array<CommunityUser>;
};

export enum CommunityUsersSortBy {
  Messages = 'MESSAGES',
  UsedChannelsPoints = 'USED_CHANNELS_POINTS',
  UsedEmotes = 'USED_EMOTES',
  Watched = 'WATCHED'
}

export type CreateOrUpdateCommandResponseInput = {
  order: Scalars['Int']['input'];
  text: Scalars['String']['input'];
};

export type CreateOrUpdateRoleSettingsInput = {
  requiredMessages: Scalars['Int']['input'];
  requiredUserChannelPoints: Scalars['Int']['input'];
  requiredWatchTime: Scalars['Int']['input'];
};

export type Dashboard = {
  __typename?: 'Dashboard';
  flags: Array<ChannelRolePermissionEnum>;
  id: Scalars['String']['output'];
  twitchProfile: TwirUserTwitchInfo;
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

export type DuelGame = {
  __typename?: 'DuelGame';
  bothDieMessage: Scalars['String']['output'];
  bothDiePercent: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  globalCooldown: Scalars['Int']['output'];
  pointsPerLose: Scalars['Int']['output'];
  pointsPerWin: Scalars['Int']['output'];
  resultMessage: Scalars['String']['output'];
  secondsToAccept: Scalars['Int']['output'];
  startMessage: Scalars['String']['output'];
  timeoutSeconds: Scalars['Int']['output'];
  userCooldown: Scalars['Int']['output'];
};

export type DuelGameOpts = {
  bothDieMessage?: InputMaybe<Scalars['String']['input']>;
  bothDiePercent?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  globalCooldown?: InputMaybe<Scalars['Int']['input']>;
  pointsPerLose?: InputMaybe<Scalars['Int']['input']>;
  pointsPerWin?: InputMaybe<Scalars['Int']['input']>;
  resultMessage?: InputMaybe<Scalars['String']['input']>;
  secondsToAccept?: InputMaybe<Scalars['Int']['input']>;
  startMessage?: InputMaybe<Scalars['String']['input']>;
  timeoutSeconds?: InputMaybe<Scalars['Int']['input']>;
  userCooldown?: InputMaybe<Scalars['Int']['input']>;
};

export type EightBallGame = {
  __typename?: 'EightBallGame';
  answers: Array<Scalars['String']['output']>;
  enabled: Scalars['Boolean']['output'];
};

export type EightBallGameOpts = {
  answers?: InputMaybe<Array<Scalars['String']['input']>>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
};

export enum EmoteStatisticRange {
  LastDay = 'LAST_DAY',
  LastMonth = 'LAST_MONTH',
  LastThreeMonth = 'LAST_THREE_MONTH',
  LastWeek = 'LAST_WEEK',
  LastYear = 'LAST_YEAR'
}

export type EmoteStatisticTopUser = {
  __typename?: 'EmoteStatisticTopUser';
  count: Scalars['Int']['output'];
  twitchProfile: TwirUserTwitchInfo;
  userId: Scalars['String']['output'];
};

export type EmoteStatisticUsage = {
  __typename?: 'EmoteStatisticUsage';
  count: Scalars['Int']['output'];
  timestamp: Scalars['Int']['output'];
};

export type EmoteStatisticUserUsage = {
  __typename?: 'EmoteStatisticUserUsage';
  date: Scalars['Time']['output'];
  twitchProfile: TwirUserTwitchInfo;
  userId: Scalars['String']['output'];
};

export type EmotesStatistic = {
  __typename?: 'EmotesStatistic';
  emoteName: Scalars['String']['output'];
  graphicUsages: Array<EmoteStatisticUsage>;
  lastUsedTimestamp: Scalars['Int']['output'];
  totalUsages: Scalars['Int']['output'];
};

export type EmotesStatisticEmoteDetailedOpts = {
  emoteName: Scalars['String']['input'];
  range: EmoteStatisticRange;
  topUsersPage?: InputMaybe<Scalars['Int']['input']>;
  topUsersPerPage?: InputMaybe<Scalars['Int']['input']>;
  usagesByUsersPage?: InputMaybe<Scalars['Int']['input']>;
  usagesByUsersPerPage?: InputMaybe<Scalars['Int']['input']>;
};

export type EmotesStatisticEmoteDetailedResponse = {
  __typename?: 'EmotesStatisticEmoteDetailedResponse';
  emoteName: Scalars['String']['output'];
  graphicUsages: Array<EmoteStatisticUsage>;
  lastUsedTimestamp: Scalars['Int']['output'];
  topUsers: Array<EmoteStatisticTopUser>;
  topUsersTotal: Scalars['Int']['output'];
  totalUsages: Scalars['Int']['output'];
  usagesByUsersTotal: Scalars['Int']['output'];
  usagesHistory: Array<EmoteStatisticUserUsage>;
};

export type EmotesStatisticResponse = {
  __typename?: 'EmotesStatisticResponse';
  emotes: Array<EmotesStatistic>;
  total: Scalars['Int']['output'];
};

export type EmotesStatisticsOpts = {
  graphicRange?: InputMaybe<EmoteStatisticRange>;
  order?: InputMaybe<EmotesStatisticsOptsOrder>;
  page?: InputMaybe<Scalars['Int']['input']>;
  perPage?: InputMaybe<Scalars['Int']['input']>;
  search?: InputMaybe<Scalars['String']['input']>;
};

export enum EmotesStatisticsOptsOrder {
  Asc = 'ASC',
  Desc = 'DESC'
}

export enum EventsubSubscribeConditionInput {
  Channel = 'CHANNEL',
  ChannelWithBotId = 'CHANNEL_WITH_BOT_ID',
  ChannelWithModeratorId = 'CHANNEL_WITH_MODERATOR_ID',
  User = 'USER'
}

export type EventsubSubscribeInput = {
  condition: EventsubSubscribeConditionInput;
  type: Scalars['String']['input'];
  version: Scalars['String']['input'];
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
  usageCount: Scalars['Int']['output'];
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
  authenticatedUserRegenerateApiKey: Scalars['String']['output'];
  authenticatedUserSelectDashboard: Scalars['Boolean']['output'];
  authenticatedUserUpdatePublicPage: Scalars['Boolean']['output'];
  authenticatedUserUpdateSettings: Scalars['Boolean']['output'];
  badgesAddUser: Scalars['Boolean']['output'];
  badgesCreate: Badge;
  badgesDelete: Scalars['Boolean']['output'];
  badgesRemoveUser: Scalars['Boolean']['output'];
  badgesUpdate: Badge;
  channelAlertsCreate: ChannelAlert;
  channelAlertsDelete: Scalars['Boolean']['output'];
  channelAlertsUpdate: ChannelAlert;
  chatOverlayUpdate: Scalars['Boolean']['output'];
  commandsCreate: Scalars['Boolean']['output'];
  commandsGroupsCreate: Scalars['Boolean']['output'];
  commandsGroupsRemove: Scalars['Boolean']['output'];
  commandsGroupsUpdate: Scalars['Boolean']['output'];
  commandsRemove: Scalars['Boolean']['output'];
  commandsUpdate: Scalars['Boolean']['output'];
  communityResetStats: Scalars['Boolean']['output'];
  dropAllAuthSessions: Scalars['Boolean']['output'];
  eventsubSubscribe: Scalars['Boolean']['output'];
  gamesDuelUpdate: DuelGame;
  gamesEightBallUpdate: EightBallGame;
  gamesRussianRouletteUpdate: RussianRouletteGame;
  gamesSeppukuUpdate: SeppukuGame;
  gamesVotebanUpdate: VotebanGame;
  greetingsCreate: Greeting;
  greetingsRemove: Scalars['Boolean']['output'];
  greetingsUpdate: Greeting;
  keywordCreate: Keyword;
  keywordRemove: Scalars['Boolean']['output'];
  keywordUpdate: Keyword;
  logout: Scalars['Boolean']['output'];
  notificationsCreate: AdminNotification;
  notificationsDelete: Scalars['Boolean']['output'];
  notificationsUpdate: AdminNotification;
  rolesCreate: Scalars['Boolean']['output'];
  rolesRemove: Scalars['Boolean']['output'];
  rolesUpdate: Scalars['Boolean']['output'];
  songRequestsUpdate: Scalars['Boolean']['output'];
  switchUserAdmin: Scalars['Boolean']['output'];
  switchUserBan: Scalars['Boolean']['output'];
  timersCreate: Timer;
  timersRemove: Scalars['Boolean']['output'];
  timersUpdate: Timer;
  updateChatAlerts: ChatAlerts;
  variablesCreate: Variable;
  variablesDelete: Scalars['Boolean']['output'];
  variablesUpdate: Variable;
};


export type MutationAuthenticatedUserSelectDashboardArgs = {
  dashboardId: Scalars['String']['input'];
};


export type MutationAuthenticatedUserUpdatePublicPageArgs = {
  opts: UserUpdatePublicSettingsInput;
};


export type MutationAuthenticatedUserUpdateSettingsArgs = {
  opts: UserUpdateSettingsInput;
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


export type MutationChannelAlertsCreateArgs = {
  input: ChannelAlertCreateInput;
};


export type MutationChannelAlertsDeleteArgs = {
  id: Scalars['ID']['input'];
};


export type MutationChannelAlertsUpdateArgs = {
  id: Scalars['ID']['input'];
  input: ChannelAlertUpdateInput;
};


export type MutationChatOverlayUpdateArgs = {
  opts: ChatOverlayUpdateOpts;
};


export type MutationCommandsCreateArgs = {
  opts: CommandsCreateOpts;
};


export type MutationCommandsGroupsCreateArgs = {
  opts: CommandsGroupsCreateOpts;
};


export type MutationCommandsGroupsRemoveArgs = {
  id: Scalars['ID']['input'];
};


export type MutationCommandsGroupsUpdateArgs = {
  id: Scalars['ID']['input'];
  opts: CommandsGroupsUpdateOpts;
};


export type MutationCommandsRemoveArgs = {
  id: Scalars['ID']['input'];
};


export type MutationCommandsUpdateArgs = {
  id: Scalars['ID']['input'];
  opts: CommandsUpdateOpts;
};


export type MutationCommunityResetStatsArgs = {
  type: CommunityUsersResetType;
};


export type MutationEventsubSubscribeArgs = {
  opts: EventsubSubscribeInput;
};


export type MutationGamesDuelUpdateArgs = {
  opts: DuelGameOpts;
};


export type MutationGamesEightBallUpdateArgs = {
  opts: EightBallGameOpts;
};


export type MutationGamesRussianRouletteUpdateArgs = {
  opts: RussianRouletteGameOpts;
};


export type MutationGamesSeppukuUpdateArgs = {
  opts: SeppukuGameOpts;
};


export type MutationGamesVotebanUpdateArgs = {
  opts: VotebanGameOpts;
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


export type MutationRolesCreateArgs = {
  opts: RolesCreateOrUpdateOpts;
};


export type MutationRolesRemoveArgs = {
  id: Scalars['ID']['input'];
};


export type MutationRolesUpdateArgs = {
  id: Scalars['ID']['input'];
  opts: RolesCreateOrUpdateOpts;
};


export type MutationSongRequestsUpdateArgs = {
  opts: SongRequestsSettingsOpts;
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


export type MutationUpdateChatAlertsArgs = {
  input: ChatAlertsInput;
};


export type MutationVariablesCreateArgs = {
  opts: VariableCreateInput;
};


export type MutationVariablesDeleteArgs = {
  id: Scalars['ID']['input'];
};


export type MutationVariablesUpdateArgs = {
  id: Scalars['ID']['input'];
  opts: VariableUpdateInput;
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

export type PublicCommand = {
  __typename?: 'PublicCommand';
  aliases: Array<Scalars['String']['output']>;
  cooldown: Scalars['Int']['output'];
  cooldownType: Scalars['String']['output'];
  description: Scalars['String']['output'];
  module: Scalars['String']['output'];
  name: Scalars['String']['output'];
  permissions: Array<PublicCommandPermission>;
  responses: Array<Scalars['String']['output']>;
};

export type PublicCommandPermission = {
  __typename?: 'PublicCommandPermission';
  name: Scalars['String']['output'];
  type: Scalars['String']['output'];
};

export type PublicSettings = {
  __typename?: 'PublicSettings';
  description?: Maybe<Scalars['String']['output']>;
  socialLinks: Array<SocialLink>;
};

export type Query = {
  __typename?: 'Query';
  authLink: Scalars['String']['output'];
  authenticatedUser: AuthenticatedUser;
  channelAlerts: Array<ChannelAlert>;
  chatAlerts?: Maybe<ChatAlerts>;
  chatOverlays: Array<ChatOverlay>;
  chatOverlaysById?: Maybe<ChatOverlay>;
  commands: Array<Command>;
  commandsGroups: Array<CommandGroup>;
  commandsPublic: Array<PublicCommand>;
  communityUsers: CommunityUsersResponse;
  emotesStatisticEmoteDetailedInformation: EmotesStatisticEmoteDetailedResponse;
  emotesStatistics: EmotesStatisticResponse;
  gamesDuel: DuelGame;
  gamesEightBall: EightBallGame;
  gamesRussianRoulette: RussianRouletteGame;
  gamesSeppuku: SeppukuGame;
  gamesVoteban: VotebanGame;
  greetings: Array<Greeting>;
  keywords: Array<Keyword>;
  notificationsByAdmin: AdminNotificationsResponse;
  notificationsByUser: Array<UserNotification>;
  rewardsRedemptionsHistory: TwitchRedemptionResponse;
  roles: Array<Role>;
  songRequests?: Maybe<SongRequestsSettings>;
  songRequestsSearchChannelOrVideo: SongRequestsSearchChannelOrVideoResponse;
  timers: Array<Timer>;
  /** Twir badges */
  twirBadges: Array<Badge>;
  /** finding users on twitch with filter does they exists in database */
  twirUsers: TwirUsersResponse;
  /**
   * Get channel badges.
   * If channelId is not provided - selected dashboard/authenticated user channelId is used, depending on context.
   * For example if queried by apiKey - userId belongs to apiKey owner id.
   */
  twitchGetChannelBadges: TwirTwitchChannelBadgeResponse;
  twitchGetChannelRewards: TwirTwitchChannelRewardResponse;
  twitchGetGlobalBadges: TwirTwitchGlobalBadgeResponse;
  twitchGetUserById?: Maybe<TwirUserTwitchInfo>;
  twitchGetUserByName?: Maybe<TwirUserTwitchInfo>;
  /** Channel id is optional */
  twitchRewards: Array<TwitchReward>;
  userPublicSettings: PublicSettings;
  variables: Array<Variable>;
  variablesBuiltIn: Array<BuiltInVariable>;
};


export type QueryAuthLinkArgs = {
  redirectTo: Scalars['String']['input'];
};


export type QueryChatOverlaysByIdArgs = {
  id: Scalars['String']['input'];
};


export type QueryCommandsPublicArgs = {
  channelId: Scalars['ID']['input'];
};


export type QueryCommunityUsersArgs = {
  opts: CommunityUsersOpts;
};


export type QueryEmotesStatisticEmoteDetailedInformationArgs = {
  opts: EmotesStatisticEmoteDetailedOpts;
};


export type QueryEmotesStatisticsArgs = {
  opts: EmotesStatisticsOpts;
};


export type QueryNotificationsByAdminArgs = {
  opts: AdminNotificationsParams;
};


export type QueryRewardsRedemptionsHistoryArgs = {
  opts: TwitchRedemptionsOpts;
};


export type QuerySongRequestsSearchChannelOrVideoArgs = {
  opts: SongRequestsSearchChannelOrVideoOpts;
};


export type QueryTwirUsersArgs = {
  opts: TwirUsersSearchParams;
};


export type QueryTwitchGetChannelBadgesArgs = {
  channelId?: InputMaybe<Scalars['ID']['input']>;
};


export type QueryTwitchGetChannelRewardsArgs = {
  channelId?: InputMaybe<Scalars['ID']['input']>;
};


export type QueryTwitchGetUserByIdArgs = {
  id: Scalars['ID']['input'];
};


export type QueryTwitchGetUserByNameArgs = {
  name: Scalars['String']['input'];
};


export type QueryTwitchRewardsArgs = {
  channelId?: InputMaybe<Scalars['String']['input']>;
};


export type QueryUserPublicSettingsArgs = {
  userId?: InputMaybe<Scalars['String']['input']>;
};

export type Role = {
  __typename?: 'Role';
  channelId: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  permissions: Array<ChannelRolePermissionEnum>;
  settings: RoleSettings;
  type: RoleTypeEnum;
  /** This is a list of user ids */
  users: Array<TwirUserTwitchInfo>;
};

export type RoleSettings = {
  __typename?: 'RoleSettings';
  requiredMessages: Scalars['Int']['output'];
  requiredUserChannelPoints: Scalars['Int']['output'];
  requiredWatchTime: Scalars['Int']['output'];
};

export enum RoleTypeEnum {
  Broadcaster = 'BROADCASTER',
  Custom = 'CUSTOM',
  Moderator = 'MODERATOR',
  Subscriber = 'SUBSCRIBER',
  Viewer = 'VIEWER',
  Vip = 'VIP'
}

export type RolesCreateOrUpdateOpts = {
  name: Scalars['String']['input'];
  permissions: Array<ChannelRolePermissionEnum>;
  settings: CreateOrUpdateRoleSettingsInput;
  /** This is a list of user ids */
  users: Array<Scalars['String']['input']>;
};

export type RussianRouletteGame = {
  __typename?: 'RussianRouletteGame';
  canBeUsedByModerator: Scalars['Boolean']['output'];
  chargedBullets: Scalars['Int']['output'];
  deathMessage: Scalars['String']['output'];
  decisionSeconds: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  initMessage: Scalars['String']['output'];
  surviveMessage: Scalars['String']['output'];
  timeoutSeconds: Scalars['Int']['output'];
  tumberSize: Scalars['Int']['output'];
};

export type RussianRouletteGameOpts = {
  canBeUsedByModerator?: InputMaybe<Scalars['Boolean']['input']>;
  chargedBullets?: InputMaybe<Scalars['Int']['input']>;
  deathMessage?: InputMaybe<Scalars['String']['input']>;
  decisionSeconds?: InputMaybe<Scalars['Int']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  initMessage?: InputMaybe<Scalars['String']['input']>;
  surviveMessage?: InputMaybe<Scalars['String']['input']>;
  timeoutSeconds?: InputMaybe<Scalars['Int']['input']>;
  tumberSize?: InputMaybe<Scalars['Int']['input']>;
};

export type SeppukuGame = {
  __typename?: 'SeppukuGame';
  enabled: Scalars['Boolean']['output'];
  message: Scalars['String']['output'];
  messageModerators: Scalars['String']['output'];
  timeoutModerators: Scalars['Boolean']['output'];
  timeoutSeconds: Scalars['Int']['output'];
};

export type SeppukuGameOpts = {
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  message?: InputMaybe<Scalars['String']['input']>;
  messageModerators?: InputMaybe<Scalars['String']['input']>;
  timeoutModerators?: InputMaybe<Scalars['Boolean']['input']>;
  timeoutSeconds?: InputMaybe<Scalars['Int']['input']>;
};

export type SocialLink = {
  __typename?: 'SocialLink';
  href: Scalars['String']['output'];
  title: Scalars['String']['output'];
};

export type SocialLinkInput = {
  href: Scalars['String']['input'];
  title: Scalars['String']['input'];
};

export type SongRequestsChannelTranslations = {
  __typename?: 'SongRequestsChannelTranslations';
  denied: Scalars['String']['output'];
};

export type SongRequestsChannelTranslationsOpts = {
  denied: Scalars['String']['input'];
};

export type SongRequestsDenyList = {
  __typename?: 'SongRequestsDenyList';
  artistsNames: Array<Scalars['String']['output']>;
  channels: Array<Scalars['String']['output']>;
  songs: Array<Scalars['String']['output']>;
  users: Array<Scalars['String']['output']>;
  words: Array<Scalars['String']['output']>;
};

export type SongRequestsDenyListOpts = {
  artistsNames: Array<Scalars['String']['input']>;
  channels: Array<Scalars['String']['input']>;
  songs: Array<Scalars['String']['input']>;
  users: Array<Scalars['String']['input']>;
  words: Array<Scalars['String']['input']>;
};

export type SongRequestsSearchChannelOrVideoItem = {
  __typename?: 'SongRequestsSearchChannelOrVideoItem';
  id: Scalars['String']['output'];
  thumbnail: Scalars['String']['output'];
  title: Scalars['String']['output'];
};

export type SongRequestsSearchChannelOrVideoOpts = {
  query: Array<Scalars['String']['input']>;
  type: SongRequestsSearchChannelOrVideoOptsType;
};

export enum SongRequestsSearchChannelOrVideoOptsType {
  Channel = 'CHANNEL',
  Video = 'VIDEO'
}

export type SongRequestsSearchChannelOrVideoResponse = {
  __typename?: 'SongRequestsSearchChannelOrVideoResponse';
  items: Array<SongRequestsSearchChannelOrVideoItem>;
};

export type SongRequestsSettings = {
  __typename?: 'SongRequestsSettings';
  acceptOnlyWhenOnline: Scalars['Boolean']['output'];
  announcePlay: Scalars['Boolean']['output'];
  channelPointsRewardId?: Maybe<Scalars['String']['output']>;
  denyList: SongRequestsDenyList;
  enabled: Scalars['Boolean']['output'];
  maxRequests: Scalars['Int']['output'];
  neededVotesForSkip: Scalars['Int']['output'];
  playerNoCookieMode: Scalars['Boolean']['output'];
  song: SongRequestsSongSettings;
  takeSongFromDonationMessages: Scalars['Boolean']['output'];
  translations: SongRequestsTranslations;
  user: SongRequestsUserSettings;
};

export type SongRequestsSettingsOpts = {
  acceptOnlyWhenOnline: Scalars['Boolean']['input'];
  announcePlay: Scalars['Boolean']['input'];
  channelPointsRewardId?: InputMaybe<Scalars['String']['input']>;
  denyList: SongRequestsDenyListOpts;
  enabled: Scalars['Boolean']['input'];
  maxRequests: Scalars['Int']['input'];
  neededVotesForSkip: Scalars['Int']['input'];
  playerNoCookieMode: Scalars['Boolean']['input'];
  song: SongRequestsSongSettingsOpts;
  takeSongFromDonationMessages: Scalars['Boolean']['input'];
  translations: SongRequestsTranslationsOpts;
  user: SongRequestsUserSettingsOpts;
};

export type SongRequestsSongSettings = {
  __typename?: 'SongRequestsSongSettings';
  acceptedCategories: Array<Scalars['String']['output']>;
  maxLength: Scalars['Int']['output'];
  minLength: Scalars['Int']['output'];
  minViews: Scalars['Int']['output'];
};

export type SongRequestsSongSettingsOpts = {
  acceptedCategories: Array<Scalars['String']['input']>;
  maxLength: Scalars['Int']['input'];
  minLength: Scalars['Int']['input'];
  minViews: Scalars['Int']['input'];
};

export type SongRequestsSongTranslations = {
  __typename?: 'SongRequestsSongTranslations';
  ageRestrictions: Scalars['String']['output'];
  alreadyInQueue: Scalars['String']['output'];
  cannotGetInformation: Scalars['String']['output'];
  denied: Scalars['String']['output'];
  live: Scalars['String']['output'];
  maxLength: Scalars['String']['output'];
  maximumOrdered: Scalars['String']['output'];
  minLength: Scalars['String']['output'];
  minViews: Scalars['String']['output'];
  notFound: Scalars['String']['output'];
  requestedMessage: Scalars['String']['output'];
};

export type SongRequestsSongTranslationsOpts = {
  ageRestrictions: Scalars['String']['input'];
  alreadyInQueue: Scalars['String']['input'];
  cannotGetInformation: Scalars['String']['input'];
  denied: Scalars['String']['input'];
  live: Scalars['String']['input'];
  maxLength: Scalars['String']['input'];
  maximumOrdered: Scalars['String']['input'];
  minLength: Scalars['String']['input'];
  minViews: Scalars['String']['input'];
  notFound: Scalars['String']['input'];
  requestedMessage: Scalars['String']['input'];
};

export type SongRequestsTranslations = {
  __typename?: 'SongRequestsTranslations';
  acceptOnlyWhenOnline: Scalars['String']['output'];
  channel: SongRequestsChannelTranslations;
  noText: Scalars['String']['output'];
  notEnabled: Scalars['String']['output'];
  nowPlaying: Scalars['String']['output'];
  song: SongRequestsSongTranslations;
  user: SongRequestsUserTranslations;
};

export type SongRequestsTranslationsOpts = {
  acceptOnlyWhenOnline: Scalars['String']['input'];
  channel: SongRequestsChannelTranslationsOpts;
  noText: Scalars['String']['input'];
  notEnabled: Scalars['String']['input'];
  nowPlaying: Scalars['String']['input'];
  song: SongRequestsSongTranslationsOpts;
  user: SongRequestsUserTranslationsOpts;
};

export type SongRequestsUserSettings = {
  __typename?: 'SongRequestsUserSettings';
  maxRequests: Scalars['Int']['output'];
  minFollowTime: Scalars['Int']['output'];
  minMessages: Scalars['Int']['output'];
  minWatchTime: Scalars['Int']['output'];
};

export type SongRequestsUserSettingsOpts = {
  maxRequests: Scalars['Int']['input'];
  minFollowTime: Scalars['Int']['input'];
  minMessages: Scalars['Int']['input'];
  minWatchTime: Scalars['Int']['input'];
};

export type SongRequestsUserTranslations = {
  __typename?: 'SongRequestsUserTranslations';
  denied: Scalars['String']['output'];
  maxRequests: Scalars['String']['output'];
  minFollow: Scalars['String']['output'];
  minMessages: Scalars['String']['output'];
  minWatched: Scalars['String']['output'];
};

export type SongRequestsUserTranslationsOpts = {
  denied: Scalars['String']['input'];
  maxRequests: Scalars['String']['input'];
  minFollow: Scalars['String']['input'];
  minMessages: Scalars['String']['input'];
  minWatched: Scalars['String']['input'];
};

export type Subscription = {
  __typename?: 'Subscription';
  chatOverlaySettings?: Maybe<ChatOverlay>;
  dashboardStats: DashboardStats;
  /** `newNotification` will return a stream of `Notification` objects. */
  newNotification: UserNotification;
};


export type SubscriptionChatOverlaySettingsArgs = {
  apiKey: Scalars['String']['input'];
  id: Scalars['String']['input'];
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

export type TwirTwitchChannelBadgeResponse = {
  __typename?: 'TwirTwitchChannelBadgeResponse';
  badges: Array<TwitchBadge>;
};

export type TwirTwitchChannelReward = {
  __typename?: 'TwirTwitchChannelReward';
  background_color: Scalars['String']['output'];
  broadcaster_id: Scalars['ID']['output'];
  broadcaster_login: Scalars['String']['output'];
  broadcaster_name: Scalars['String']['output'];
  cooldown_expires_at: Scalars['String']['output'];
  cost: Scalars['Int']['output'];
  global_cooldown_setting: TwirTwitchChannelRewardGlobalCooldownSetting;
  id: Scalars['ID']['output'];
  /** In case of image is not set - default image is used */
  image: TwirTwitchChannelRewardImage;
  is_enabled: Scalars['Boolean']['output'];
  is_in_stock: Scalars['Boolean']['output'];
  is_paused: Scalars['Boolean']['output'];
  is_user_input_required: Scalars['Boolean']['output'];
  max_per_stream_setting: TwirTwitchChannelRewardMaxPerStreamSetting;
  max_per_user_per_stream_setting: TwirTwitchChannelRewardMaxPerUserPerStreamSetting;
  prompt: Scalars['String']['output'];
  redemptions_redeemed_current_stream: Scalars['Int']['output'];
  should_redemptions_skip_request_queue: Scalars['Boolean']['output'];
  title: Scalars['String']['output'];
};

export type TwirTwitchChannelRewardGlobalCooldownSetting = {
  __typename?: 'TwirTwitchChannelRewardGlobalCooldownSetting';
  global_cooldown_seconds: Scalars['Int']['output'];
  is_enabled: Scalars['Boolean']['output'];
};

export type TwirTwitchChannelRewardImage = {
  __typename?: 'TwirTwitchChannelRewardImage';
  url_1x: Scalars['String']['output'];
  url_2x: Scalars['String']['output'];
  url_4x: Scalars['String']['output'];
};

export type TwirTwitchChannelRewardMaxPerStreamSetting = {
  __typename?: 'TwirTwitchChannelRewardMaxPerStreamSetting';
  is_enabled: Scalars['Boolean']['output'];
  max_per_stream: Scalars['Int']['output'];
};

export type TwirTwitchChannelRewardMaxPerUserPerStreamSetting = {
  __typename?: 'TwirTwitchChannelRewardMaxPerUserPerStreamSetting';
  is_enabled: Scalars['Boolean']['output'];
  max_per_user_per_stream: Scalars['Int']['output'];
};

export type TwirTwitchChannelRewardResponse = {
  __typename?: 'TwirTwitchChannelRewardResponse';
  partnerOrAffiliate: Scalars['Boolean']['output'];
  rewards: Array<TwirTwitchChannelReward>;
};

export type TwirTwitchGlobalBadgeResponse = {
  __typename?: 'TwirTwitchGlobalBadgeResponse';
  badges: Array<TwitchBadge>;
};

export type TwirUser = {
  id: Scalars['ID']['output'];
  twitchProfile: TwirUserTwitchInfo;
};

export type TwirUserTwitchInfo = {
  __typename?: 'TwirUserTwitchInfo';
  description: Scalars['String']['output'];
  displayName: Scalars['String']['output'];
  id: Scalars['String']['output'];
  login: Scalars['String']['output'];
  notFound: Scalars['Boolean']['output'];
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

export type TwitchBadge = {
  __typename?: 'TwitchBadge';
  set_id: Scalars['String']['output'];
  versions: Array<TwitchBadgeVersion>;
};

export type TwitchBadgeVersion = {
  __typename?: 'TwitchBadgeVersion';
  id: Scalars['String']['output'];
  image_url_1x: Scalars['String']['output'];
  image_url_2x: Scalars['String']['output'];
  image_url_4x: Scalars['String']['output'];
};

export type TwitchRedemption = {
  __typename?: 'TwitchRedemption';
  channelId: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  prompt?: Maybe<Scalars['String']['output']>;
  redeemedAt: Scalars['Time']['output'];
  reward: TwitchReward;
  user: TwirUserTwitchInfo;
};

export type TwitchRedemptionResponse = {
  __typename?: 'TwitchRedemptionResponse';
  redemptions: Array<TwitchRedemption>;
  total: Scalars['Int']['output'];
};

export type TwitchRedemptionsOpts = {
  byChannelId?: InputMaybe<Scalars['ID']['input']>;
  page?: InputMaybe<Scalars['Int']['input']>;
  perPage?: InputMaybe<Scalars['Int']['input']>;
  rewardsIds?: InputMaybe<Array<Scalars['ID']['input']>>;
  userSearch?: InputMaybe<Scalars['String']['input']>;
};

export type TwitchReward = {
  __typename?: 'TwitchReward';
  backgroundColor: Scalars['String']['output'];
  cost: Scalars['Int']['output'];
  enabled: Scalars['Boolean']['output'];
  id: Scalars['ID']['output'];
  imageUrls?: Maybe<Array<Scalars['String']['output']>>;
  title: Scalars['String']['output'];
  usedTimes: Scalars['Int']['output'];
};

export type UserNotification = Notification & {
  __typename?: 'UserNotification';
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  text: Scalars['String']['output'];
  userId?: Maybe<Scalars['ID']['output']>;
};

export type UserUpdatePublicSettingsInput = {
  description?: InputMaybe<Scalars['String']['input']>;
  socialLinks?: InputMaybe<Array<SocialLinkInput>>;
};

export type UserUpdateSettingsInput = {
  hideOnLandingPage?: InputMaybe<Scalars['Boolean']['input']>;
};

export type Variable = {
  __typename?: 'Variable';
  description?: Maybe<Scalars['String']['output']>;
  evalValue: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  response: Scalars['String']['output'];
  type: VariableType;
};

export type VariableCreateInput = {
  description?: InputMaybe<Scalars['String']['input']>;
  evalValue: Scalars['String']['input'];
  name: Scalars['String']['input'];
  response: Scalars['String']['input'];
  type: VariableType;
};

export enum VariableType {
  Number = 'NUMBER',
  Script = 'SCRIPT',
  Text = 'TEXT'
}

export type VariableUpdateInput = {
  description?: InputMaybe<Scalars['String']['input']>;
  evalValue?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  response?: InputMaybe<Scalars['String']['input']>;
  type?: InputMaybe<VariableType>;
};

export enum VoteBanGameVotingMode {
  Chat = 'CHAT',
  Polls = 'POLLS'
}

export type VotebanGame = {
  __typename?: 'VotebanGame';
  banMessage: Scalars['String']['output'];
  banMessageModerators: Scalars['String']['output'];
  chatVotesWordsNegative: Array<Scalars['String']['output']>;
  chatVotesWordsPositive: Array<Scalars['String']['output']>;
  enabled: Scalars['Boolean']['output'];
  initMessage: Scalars['String']['output'];
  neededVotes: Scalars['Int']['output'];
  surviveMessage: Scalars['String']['output'];
  surviveMessageModerators: Scalars['String']['output'];
  timeoutModerators: Scalars['Boolean']['output'];
  timeoutSeconds: Scalars['Int']['output'];
  voteDuration: Scalars['Int']['output'];
  votingMode: VoteBanGameVotingMode;
};

export type VotebanGameOpts = {
  banMessage?: InputMaybe<Scalars['String']['input']>;
  banMessageModerators?: InputMaybe<Scalars['String']['input']>;
  chatVotesWordsNegative?: InputMaybe<Array<Scalars['String']['input']>>;
  chatVotesWordsPositive?: InputMaybe<Array<Scalars['String']['input']>>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  initMessage?: InputMaybe<Scalars['String']['input']>;
  neededVotes?: InputMaybe<Scalars['Int']['input']>;
  surviveMessage?: InputMaybe<Scalars['String']['input']>;
  surviveMessageModerators?: InputMaybe<Scalars['String']['input']>;
  timeoutModerators?: InputMaybe<Scalars['Boolean']['input']>;
  timeoutSeconds?: InputMaybe<Scalars['Int']['input']>;
  voteDuration?: InputMaybe<Scalars['Int']['input']>;
  votingMode?: InputMaybe<VoteBanGameVotingMode>;
};

export type ChatOverlayWithAdditionalDataQueryVariables = Exact<{
  overlayId: Scalars['String']['input'];
}>;


export type ChatOverlayWithAdditionalDataQuery = { __typename?: 'Query', authenticatedUser: { __typename?: 'AuthenticatedUser', id: string, twitchProfile: { __typename?: 'TwirUserTwitchInfo', login: string, displayName: string, profileImageUrl: string } }, chatOverlaysById?: { __typename?: 'ChatOverlay', id: string, messageHideTimeout: number } | null, twitchGetGlobalBadges: { __typename?: 'TwirTwitchGlobalBadgeResponse', badges: Array<{ __typename?: 'TwitchBadge', set_id: string, versions: Array<{ __typename?: 'TwitchBadgeVersion', image_url_1x: string }> }> }, twitchGetChannelBadges: { __typename?: 'TwirTwitchChannelBadgeResponse', badges: Array<{ __typename?: 'TwitchBadge', set_id: string, versions: Array<{ __typename?: 'TwitchBadgeVersion', image_url_1x: string }> }> } };

export type ChatOverlaySettingsSubscriptionVariables = Exact<{
  id: Scalars['String']['input'];
  apiKey: Scalars['String']['input'];
}>;


export type ChatOverlaySettingsSubscription = { __typename?: 'Subscription', chatOverlaySettings?: { __typename?: 'ChatOverlay', id: string, animation?: ChatOverlayAnimation | null, fontSize: number } | null };


export const ChatOverlayWithAdditionalDataDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"ChatOverlayWithAdditionalData"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"overlayId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"authenticatedUser"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"twitchProfile"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"login"}},{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"profileImageUrl"}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"chatOverlaysById"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"overlayId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"messageHideTimeout"}}]}},{"kind":"Field","name":{"kind":"Name","value":"twitchGetGlobalBadges"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"badges"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"set_id"}},{"kind":"Field","name":{"kind":"Name","value":"versions"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"image_url_1x"}}]}}]}}]}},{"kind":"Field","name":{"kind":"Name","value":"twitchGetChannelBadges"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"badges"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"set_id"}},{"kind":"Field","name":{"kind":"Name","value":"versions"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"image_url_1x"}}]}}]}}]}}]}}]} as unknown as DocumentNode<ChatOverlayWithAdditionalDataQuery, ChatOverlayWithAdditionalDataQueryVariables>;
export const ChatOverlaySettingsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"subscription","name":{"kind":"Name","value":"ChatOverlaySettings"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"apiKey"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"chatOverlaySettings"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"apiKey"},"value":{"kind":"Variable","name":{"kind":"Name","value":"apiKey"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"animation"}},{"kind":"Field","name":{"kind":"Name","value":"fontSize"}}]}}]}}]} as unknown as DocumentNode<ChatOverlaySettingsSubscription, ChatOverlaySettingsSubscriptionVariables>;