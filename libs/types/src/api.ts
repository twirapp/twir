/* Do not change, this code is generated from Golang structs */


export enum ChannelOverlayNowPlayingPreset {
    TRANSPARENT = "TRANSPARENT",
    AIDEN_REDESIGN = "AIDEN_REDESIGN",
}
export interface BotInfo {
    isMod: boolean;
    botId: string;
    botName: string;
    enabled: boolean;
}
export interface Bot {
    GET: BotInfo;
}
export interface TTSSettings {
    enabled?: boolean;
    rate: number;
    volume: number;
    pitch: number;
    voice: string;
    allow_users_choose_voice_in_main_command: boolean;
    max_symbols: number;
    disallowed_voices: string[];
    do_not_read_emoji: boolean;
    do_not_read_twitch_emotes: boolean;
    do_not_read_links: boolean;
    read_chat_messages: boolean;
    read_chat_messages_nicknames: boolean;
}
export interface TTS {
    GET: TTSSettings;
    POST: TTSSettings;
}
export interface OBSWebSocketSettings {
    serverPort: number;
    serverAddress: string;
    serverPassword: string;
}
export interface OBS {
    GET: OBSWebSocketSettings;
    POST: OBSWebSocketSettings;
}
export interface SearchResult {
    id: string;
    title: string;
    thumbNail: string;
}
export interface YouTubeChannelTranslations {
    denied: string;
}
export interface YouTubeSongTranslations {
    denied: string;
    notFound: string;
    alreadyInQueue: string;
    ageRestrictions: string;
    cannotGetInformation: string;
    live: string;
    maxLength: string;
    minLength: string;
    requestedMessage: string;
    maximumOrdered: string;
    minViews: string;
}
export interface YouTubeUserTranslations {
    denied: string;
    maxRequests: string;
    minMessages: string;
    minWatched: string;
    minFollow: string;
}
export interface YouTubeTranslations {
    nowPlaying: string;
    notEnabled: string;
    noText: string;
    acceptOnlyWhenOnline: string;
    user: YouTubeUserTranslations;
    song: YouTubeSongTranslations;
    channel: YouTubeChannelTranslations;
}
export interface YouTubeDenyList {
    users: string[];
    songs: string[];
    channels: string[];
    artistsNames: string[];
    words: string[];
}
export interface YouTubeSongSettings {
    minLength: number;
    maxLength: number;
    minViews: number;
    acceptedCategories: string[];
    wordsDenyList: string[];
}
export interface YouTubeUserSettings {
    maxRequests: number;
    minWatchTime: number;
    minMessages: number;
    minFollowTime: number;
}
export interface YouTubeSettings {
    enabled?: boolean;
    acceptOnlyWhenOnline?: boolean;
    playerNoCookieMode?: boolean;
    maxRequests: number;
    channelPointsRewardId: string;
    announcePlay?: boolean;
    neededVotesVorSkip: number;
    user: YouTubeUserSettings;
    song: YouTubeSongSettings;
    denyList: YouTubeDenyList;
    translations: YouTubeTranslations;
}
export interface YouTube {
    POST: YouTubeSettings;
    GET: YouTubeSettings;
    SEARCH: SearchResult[];
}
export interface Modules {
    YouTube: YouTube;
    OBS: OBS;
    TTS: TTS;
}
export interface Channels {
    MODULES: Modules;
    BOT: Bot;
}
export interface V1 {
    CHANNELS: Channels;
}