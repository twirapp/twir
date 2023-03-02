import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";
export declare const protobufPackage = "emotes_cacher";
export interface Request {
    channelId: string;
}
export declare const Request: {
    encode(message: Request, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Request;
    fromJSON(object: any): Request;
    toJSON(message: Request): unknown;
    create(base?: DeepPartial<Request>): Request;
    fromPartial(object: DeepPartial<Request>): Request;
};
export type EmotesCacherDefinition = typeof EmotesCacherDefinition;
export declare const EmotesCacherDefinition: {
    readonly name: "EmotesCacher";
    readonly fullName: "emotes_cacher.EmotesCacher";
    readonly methods: {
        readonly cacheChannelEmotes: {
            readonly name: "CacheChannelEmotes";
            readonly requestType: {
                encode(message: Request, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): Request;
                fromJSON(object: any): Request;
                toJSON(message: Request): unknown;
                create(base?: DeepPartial<Request>): Request;
                fromPartial(object: DeepPartial<Request>): Request;
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
        readonly cacheGlobalEmotes: {
            readonly name: "CacheGlobalEmotes";
            readonly requestType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
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
export interface EmotesCacherServiceImplementation<CallContextExt = {}> {
    cacheChannelEmotes(request: Request, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    cacheGlobalEmotes(request: Empty, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
}
export interface EmotesCacherClient<CallOptionsExt = {}> {
    cacheChannelEmotes(request: DeepPartial<Request>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    cacheGlobalEmotes(request: DeepPartial<Empty>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
}
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
//# sourceMappingURL=emotes_cacher.d.ts.map