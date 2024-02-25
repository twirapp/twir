/* Do not change, this code is generated from Golang structs */


export interface DudesGrowRequest {
    channelId: string;
    userId: string;
    userName: string;
    userLogin: string;
    color: string;
}
export interface DudesChangeColorRequest {
    channelId: string;
    userId: string;
    userName: string;
    userLogin: string;
    color: string;
}
export interface ChannelRoleUser {
    id: string;
    userId: string;
}
export interface UsersStats {
    id: string;
    userId: string;
    channelId: string;
    messages: number;
    watched: number;
    usedChannelPoints: number;
    isMod: boolean;
    isVip: boolean;
    isSubscriber: boolean;
    reputation: number;
    emotes: number;
}
export interface Tokens {
    id: string;
    accessToken: string;
    refreshToken: string;
    expiresIn: number;
    obtainmentTimestamp: Time;
    scopes: string[];
}
export interface NullString {
    String: string;
    Valid: boolean;
}
export interface Users {
    id: string;
    tokenId: NullString;
    isTester: boolean;
    isBotAdmin: boolean;
    apiKey: string;
    channel?: Channels;
    token?: Tokens;
    stats?: UsersStats;
    hide_on_landing_page: boolean;
    roles: ChannelRoleUser[];
}
export interface Channels {
    id: string;
    isEnabled: boolean;
    isTwitchBanned: boolean;
    isBanned: boolean;
    isBotMod: boolean;
    botId: string;
}
export interface Time {

}
export interface DudesUserSettings {
    id: number[];
    channelId: string;
    userId: string;
    dudeColor: string;
    createdAt: Time;
    updatedAt: Time;
    channel?: Channels;
    user?: Users;
}