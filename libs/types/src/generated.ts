/* Do not change, this code is generated from Golang structs */


export interface YoutubeBlacklistSettingsChannels {
    id: string;
    title: string;
    thumbNail: string;
}
export interface YoutubeBlacklistSettingsSongs {
    id: string;
    title: string;
    thumbNail: string;
}
export interface YoutubeBlacklistSettingsUsers {
    userId: string;
    userName: string;
}
export interface YoutubeBlacklistSettings {
    users: YoutubeBlacklistSettingsUsers[];
    songs: YoutubeBlacklistSettingsSongs[];
    channels: YoutubeBlacklistSettingsChannels[];
    artistsNames: string[];
}
export interface YotubeSongSettings {
    maxLength?: number;
    minViews?: number;
    acceptedCategories: string[];
}
export interface YoutubeUserSettings {
    maxRequests?: number;
    minWatchTime?: number;
    minMessages?: number;
    minFollowTime?: number;
}
export interface YoutubeSettings {
    maxRequests?: number;
    acceptOnlyWhenOnline?: boolean;
    channelPointsRewardName?: string;
    user?: YoutubeUserSettings;
    song?: YotubeSongSettings;
    blacklist?: YoutubeBlacklistSettings;
}