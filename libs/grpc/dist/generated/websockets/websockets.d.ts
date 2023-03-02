import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";
export declare const protobufPackage = "websockets";
export interface YoutubeAddSongToQueueRequest {
    channelId: string;
    entityId: string;
}
export interface YoutubeRemoveSongFromQueueRequest {
    channelId: string;
    entityId: string;
}
export interface ObsSetSceneMessage {
    channelId: string;
    sceneName: string;
}
export interface ObsToggleSourceMessage {
    channelId: string;
    sourceName: string;
}
export interface ObsToggleAudioMessage {
    channelId: string;
    audioSourceName: string;
}
export interface ObsAudioSetVolumeMessage {
    channelId: string;
    audioSourceName: string;
    volume: number;
}
export interface ObsAudioIncreaseVolumeMessage {
    channelId: string;
    audioSourceName: string;
    step: number;
}
export interface ObsAudioDecreaseVolumeMessage {
    channelId: string;
    audioSourceName: string;
    step: number;
}
export interface ObsAudioDisableOrEnableMessage {
    channelId: string;
    audioSourceName: string;
}
export interface ObsStopOrStartStream {
    channelId: string;
}
export interface TTSMessage {
    channelId: string;
    text: string;
    voice: string;
    rate: string;
    pitch: string;
    volume: string;
}
export interface TTSSkipMessage {
    channelId: string;
}
export declare const YoutubeAddSongToQueueRequest: {
    encode(message: YoutubeAddSongToQueueRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): YoutubeAddSongToQueueRequest;
    fromJSON(object: any): YoutubeAddSongToQueueRequest;
    toJSON(message: YoutubeAddSongToQueueRequest): unknown;
    create(base?: DeepPartial<YoutubeAddSongToQueueRequest>): YoutubeAddSongToQueueRequest;
    fromPartial(object: DeepPartial<YoutubeAddSongToQueueRequest>): YoutubeAddSongToQueueRequest;
};
export declare const YoutubeRemoveSongFromQueueRequest: {
    encode(message: YoutubeRemoveSongFromQueueRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): YoutubeRemoveSongFromQueueRequest;
    fromJSON(object: any): YoutubeRemoveSongFromQueueRequest;
    toJSON(message: YoutubeRemoveSongFromQueueRequest): unknown;
    create(base?: DeepPartial<YoutubeRemoveSongFromQueueRequest>): YoutubeRemoveSongFromQueueRequest;
    fromPartial(object: DeepPartial<YoutubeRemoveSongFromQueueRequest>): YoutubeRemoveSongFromQueueRequest;
};
export declare const ObsSetSceneMessage: {
    encode(message: ObsSetSceneMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ObsSetSceneMessage;
    fromJSON(object: any): ObsSetSceneMessage;
    toJSON(message: ObsSetSceneMessage): unknown;
    create(base?: DeepPartial<ObsSetSceneMessage>): ObsSetSceneMessage;
    fromPartial(object: DeepPartial<ObsSetSceneMessage>): ObsSetSceneMessage;
};
export declare const ObsToggleSourceMessage: {
    encode(message: ObsToggleSourceMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ObsToggleSourceMessage;
    fromJSON(object: any): ObsToggleSourceMessage;
    toJSON(message: ObsToggleSourceMessage): unknown;
    create(base?: DeepPartial<ObsToggleSourceMessage>): ObsToggleSourceMessage;
    fromPartial(object: DeepPartial<ObsToggleSourceMessage>): ObsToggleSourceMessage;
};
export declare const ObsToggleAudioMessage: {
    encode(message: ObsToggleAudioMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ObsToggleAudioMessage;
    fromJSON(object: any): ObsToggleAudioMessage;
    toJSON(message: ObsToggleAudioMessage): unknown;
    create(base?: DeepPartial<ObsToggleAudioMessage>): ObsToggleAudioMessage;
    fromPartial(object: DeepPartial<ObsToggleAudioMessage>): ObsToggleAudioMessage;
};
export declare const ObsAudioSetVolumeMessage: {
    encode(message: ObsAudioSetVolumeMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ObsAudioSetVolumeMessage;
    fromJSON(object: any): ObsAudioSetVolumeMessage;
    toJSON(message: ObsAudioSetVolumeMessage): unknown;
    create(base?: DeepPartial<ObsAudioSetVolumeMessage>): ObsAudioSetVolumeMessage;
    fromPartial(object: DeepPartial<ObsAudioSetVolumeMessage>): ObsAudioSetVolumeMessage;
};
export declare const ObsAudioIncreaseVolumeMessage: {
    encode(message: ObsAudioIncreaseVolumeMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ObsAudioIncreaseVolumeMessage;
    fromJSON(object: any): ObsAudioIncreaseVolumeMessage;
    toJSON(message: ObsAudioIncreaseVolumeMessage): unknown;
    create(base?: DeepPartial<ObsAudioIncreaseVolumeMessage>): ObsAudioIncreaseVolumeMessage;
    fromPartial(object: DeepPartial<ObsAudioIncreaseVolumeMessage>): ObsAudioIncreaseVolumeMessage;
};
export declare const ObsAudioDecreaseVolumeMessage: {
    encode(message: ObsAudioDecreaseVolumeMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ObsAudioDecreaseVolumeMessage;
    fromJSON(object: any): ObsAudioDecreaseVolumeMessage;
    toJSON(message: ObsAudioDecreaseVolumeMessage): unknown;
    create(base?: DeepPartial<ObsAudioDecreaseVolumeMessage>): ObsAudioDecreaseVolumeMessage;
    fromPartial(object: DeepPartial<ObsAudioDecreaseVolumeMessage>): ObsAudioDecreaseVolumeMessage;
};
export declare const ObsAudioDisableOrEnableMessage: {
    encode(message: ObsAudioDisableOrEnableMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ObsAudioDisableOrEnableMessage;
    fromJSON(object: any): ObsAudioDisableOrEnableMessage;
    toJSON(message: ObsAudioDisableOrEnableMessage): unknown;
    create(base?: DeepPartial<ObsAudioDisableOrEnableMessage>): ObsAudioDisableOrEnableMessage;
    fromPartial(object: DeepPartial<ObsAudioDisableOrEnableMessage>): ObsAudioDisableOrEnableMessage;
};
export declare const ObsStopOrStartStream: {
    encode(message: ObsStopOrStartStream, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ObsStopOrStartStream;
    fromJSON(object: any): ObsStopOrStartStream;
    toJSON(message: ObsStopOrStartStream): unknown;
    create(base?: DeepPartial<ObsStopOrStartStream>): ObsStopOrStartStream;
    fromPartial(object: DeepPartial<ObsStopOrStartStream>): ObsStopOrStartStream;
};
export declare const TTSMessage: {
    encode(message: TTSMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TTSMessage;
    fromJSON(object: any): TTSMessage;
    toJSON(message: TTSMessage): unknown;
    create(base?: DeepPartial<TTSMessage>): TTSMessage;
    fromPartial(object: DeepPartial<TTSMessage>): TTSMessage;
};
export declare const TTSSkipMessage: {
    encode(message: TTSSkipMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TTSSkipMessage;
    fromJSON(object: any): TTSSkipMessage;
    toJSON(message: TTSSkipMessage): unknown;
    create(base?: DeepPartial<TTSSkipMessage>): TTSSkipMessage;
    fromPartial(object: DeepPartial<TTSSkipMessage>): TTSSkipMessage;
};
export type WebsocketDefinition = typeof WebsocketDefinition;
export declare const WebsocketDefinition: {
    readonly name: "Websocket";
    readonly fullName: "websockets.Websocket";
    readonly methods: {
        readonly youtubeAddSongToQueue: {
            readonly name: "YoutubeAddSongToQueue";
            readonly requestType: {
                encode(message: YoutubeAddSongToQueueRequest, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): YoutubeAddSongToQueueRequest;
                fromJSON(object: any): YoutubeAddSongToQueueRequest;
                toJSON(message: YoutubeAddSongToQueueRequest): unknown;
                create(base?: DeepPartial<YoutubeAddSongToQueueRequest>): YoutubeAddSongToQueueRequest;
                fromPartial(object: DeepPartial<YoutubeAddSongToQueueRequest>): YoutubeAddSongToQueueRequest;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly youtubeRemoveSongToQueue: {
            readonly name: "YoutubeRemoveSongToQueue";
            readonly requestType: {
                encode(message: YoutubeRemoveSongFromQueueRequest, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): YoutubeRemoveSongFromQueueRequest;
                fromJSON(object: any): YoutubeRemoveSongFromQueueRequest;
                toJSON(message: YoutubeRemoveSongFromQueueRequest): unknown;
                create(base?: DeepPartial<YoutubeRemoveSongFromQueueRequest>): YoutubeRemoveSongFromQueueRequest;
                fromPartial(object: DeepPartial<YoutubeRemoveSongFromQueueRequest>): YoutubeRemoveSongFromQueueRequest;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly obsSetScene: {
            readonly name: "ObsSetScene";
            readonly requestType: {
                encode(message: ObsSetSceneMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): ObsSetSceneMessage;
                fromJSON(object: any): ObsSetSceneMessage;
                toJSON(message: ObsSetSceneMessage): unknown;
                create(base?: DeepPartial<ObsSetSceneMessage>): ObsSetSceneMessage;
                fromPartial(object: DeepPartial<ObsSetSceneMessage>): ObsSetSceneMessage;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly obsToggleSource: {
            readonly name: "ObsToggleSource";
            readonly requestType: {
                encode(message: ObsToggleSourceMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): ObsToggleSourceMessage;
                fromJSON(object: any): ObsToggleSourceMessage;
                toJSON(message: ObsToggleSourceMessage): unknown;
                create(base?: DeepPartial<ObsToggleSourceMessage>): ObsToggleSourceMessage;
                fromPartial(object: DeepPartial<ObsToggleSourceMessage>): ObsToggleSourceMessage;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly obsToggleAudio: {
            readonly name: "ObsToggleAudio";
            readonly requestType: {
                encode(message: ObsToggleAudioMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): ObsToggleAudioMessage;
                fromJSON(object: any): ObsToggleAudioMessage;
                toJSON(message: ObsToggleAudioMessage): unknown;
                create(base?: DeepPartial<ObsToggleAudioMessage>): ObsToggleAudioMessage;
                fromPartial(object: DeepPartial<ObsToggleAudioMessage>): ObsToggleAudioMessage;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly obsAudioSetVolume: {
            readonly name: "ObsAudioSetVolume";
            readonly requestType: {
                encode(message: ObsAudioSetVolumeMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): ObsAudioSetVolumeMessage;
                fromJSON(object: any): ObsAudioSetVolumeMessage;
                toJSON(message: ObsAudioSetVolumeMessage): unknown;
                create(base?: DeepPartial<ObsAudioSetVolumeMessage>): ObsAudioSetVolumeMessage;
                fromPartial(object: DeepPartial<ObsAudioSetVolumeMessage>): ObsAudioSetVolumeMessage;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly obsAudioIncreaseVolume: {
            readonly name: "ObsAudioIncreaseVolume";
            readonly requestType: {
                encode(message: ObsAudioIncreaseVolumeMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): ObsAudioIncreaseVolumeMessage;
                fromJSON(object: any): ObsAudioIncreaseVolumeMessage;
                toJSON(message: ObsAudioIncreaseVolumeMessage): unknown;
                create(base?: DeepPartial<ObsAudioIncreaseVolumeMessage>): ObsAudioIncreaseVolumeMessage;
                fromPartial(object: DeepPartial<ObsAudioIncreaseVolumeMessage>): ObsAudioIncreaseVolumeMessage;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly obsAudioDecreaseVolume: {
            readonly name: "ObsAudioDecreaseVolume";
            readonly requestType: {
                encode(message: ObsAudioDecreaseVolumeMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): ObsAudioDecreaseVolumeMessage;
                fromJSON(object: any): ObsAudioDecreaseVolumeMessage;
                toJSON(message: ObsAudioDecreaseVolumeMessage): unknown;
                create(base?: DeepPartial<ObsAudioDecreaseVolumeMessage>): ObsAudioDecreaseVolumeMessage;
                fromPartial(object: DeepPartial<ObsAudioDecreaseVolumeMessage>): ObsAudioDecreaseVolumeMessage;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly obsAudioEnable: {
            readonly name: "ObsAudioEnable";
            readonly requestType: {
                encode(message: ObsAudioDisableOrEnableMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): ObsAudioDisableOrEnableMessage;
                fromJSON(object: any): ObsAudioDisableOrEnableMessage;
                toJSON(message: ObsAudioDisableOrEnableMessage): unknown;
                create(base?: DeepPartial<ObsAudioDisableOrEnableMessage>): ObsAudioDisableOrEnableMessage;
                fromPartial(object: DeepPartial<ObsAudioDisableOrEnableMessage>): ObsAudioDisableOrEnableMessage;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly obsAudioDisable: {
            readonly name: "ObsAudioDisable";
            readonly requestType: {
                encode(message: ObsAudioDisableOrEnableMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): ObsAudioDisableOrEnableMessage;
                fromJSON(object: any): ObsAudioDisableOrEnableMessage;
                toJSON(message: ObsAudioDisableOrEnableMessage): unknown;
                create(base?: DeepPartial<ObsAudioDisableOrEnableMessage>): ObsAudioDisableOrEnableMessage;
                fromPartial(object: DeepPartial<ObsAudioDisableOrEnableMessage>): ObsAudioDisableOrEnableMessage;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly obsStopStream: {
            readonly name: "ObsStopStream";
            readonly requestType: {
                encode(message: ObsStopOrStartStream, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): ObsStopOrStartStream;
                fromJSON(object: any): ObsStopOrStartStream;
                toJSON(message: ObsStopOrStartStream): unknown;
                create(base?: DeepPartial<ObsStopOrStartStream>): ObsStopOrStartStream;
                fromPartial(object: DeepPartial<ObsStopOrStartStream>): ObsStopOrStartStream;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly obsStartStream: {
            readonly name: "ObsStartStream";
            readonly requestType: {
                encode(message: ObsStopOrStartStream, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): ObsStopOrStartStream;
                fromJSON(object: any): ObsStopOrStartStream;
                toJSON(message: ObsStopOrStartStream): unknown;
                create(base?: DeepPartial<ObsStopOrStartStream>): ObsStopOrStartStream;
                fromPartial(object: DeepPartial<ObsStopOrStartStream>): ObsStopOrStartStream;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly textToSpeechSay: {
            readonly name: "TextToSpeechSay";
            readonly requestType: {
                encode(message: TTSMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): TTSMessage;
                fromJSON(object: any): TTSMessage;
                toJSON(message: TTSMessage): unknown;
                create(base?: DeepPartial<TTSMessage>): TTSMessage;
                fromPartial(object: DeepPartial<TTSMessage>): TTSMessage;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly textToSpeechSkip: {
            readonly name: "TextToSpeechSkip";
            readonly requestType: {
                encode(message: TTSSkipMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): TTSSkipMessage;
                fromJSON(object: any): TTSSkipMessage;
                toJSON(message: TTSSkipMessage): unknown;
                create(base?: DeepPartial<TTSSkipMessage>): TTSSkipMessage;
                fromPartial(object: DeepPartial<TTSSkipMessage>): TTSSkipMessage;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly responseStream: false;
            readonly options: {};
        };
    };
};
export interface WebsocketServiceImplementation<CallContextExt = {}> {
    youtubeAddSongToQueue(request: YoutubeAddSongToQueueRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    youtubeRemoveSongToQueue(request: YoutubeRemoveSongFromQueueRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    obsSetScene(request: ObsSetSceneMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    obsToggleSource(request: ObsToggleSourceMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    obsToggleAudio(request: ObsToggleAudioMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    obsAudioSetVolume(request: ObsAudioSetVolumeMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    obsAudioIncreaseVolume(request: ObsAudioIncreaseVolumeMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    obsAudioDecreaseVolume(request: ObsAudioDecreaseVolumeMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    obsAudioEnable(request: ObsAudioDisableOrEnableMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    obsAudioDisable(request: ObsAudioDisableOrEnableMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    obsStopStream(request: ObsStopOrStartStream, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    obsStartStream(request: ObsStopOrStartStream, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    textToSpeechSay(request: TTSMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    textToSpeechSkip(request: TTSSkipMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
}
export interface WebsocketClient<CallOptionsExt = {}> {
    youtubeAddSongToQueue(request: DeepPartial<YoutubeAddSongToQueueRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    youtubeRemoveSongToQueue(request: DeepPartial<YoutubeRemoveSongFromQueueRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    obsSetScene(request: DeepPartial<ObsSetSceneMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    obsToggleSource(request: DeepPartial<ObsToggleSourceMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    obsToggleAudio(request: DeepPartial<ObsToggleAudioMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    obsAudioSetVolume(request: DeepPartial<ObsAudioSetVolumeMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    obsAudioIncreaseVolume(request: DeepPartial<ObsAudioIncreaseVolumeMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    obsAudioDecreaseVolume(request: DeepPartial<ObsAudioDecreaseVolumeMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    obsAudioEnable(request: DeepPartial<ObsAudioDisableOrEnableMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    obsAudioDisable(request: DeepPartial<ObsAudioDisableOrEnableMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    obsStopStream(request: DeepPartial<ObsStopOrStartStream>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    obsStartStream(request: DeepPartial<ObsStopOrStartStream>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    textToSpeechSay(request: DeepPartial<TTSMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    textToSpeechSkip(request: DeepPartial<TTSSkipMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
}
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
//# sourceMappingURL=websockets.d.ts.map