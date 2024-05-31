/* Do not change, this code is generated from Golang structs */


export enum ChannelOverlayNowPlayingPreset {
    TRANSPARENT = "TRANSPARENT",
    AIDEN_REDESIGN = "AIDEN_REDESIGN",
    SIMPLE_LINE = "SIMPLE_LINE",
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
export interface Modules {
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