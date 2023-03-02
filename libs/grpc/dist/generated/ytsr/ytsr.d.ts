import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
export declare const protobufPackage = "ytsr";
export interface SearchRequest {
    search: string;
}
export interface SongAuthor {
    name: string;
    channelId: string;
    avatarUrl?: string | undefined;
}
export interface Song {
    title: string;
    id: string;
    views: number;
    duration: number;
    thumbnailUrl?: string | undefined;
    isLive: boolean;
    author?: SongAuthor | undefined;
}
export interface SearchResponse {
    songs: Song[];
}
export declare const SearchRequest: {
    encode(message: SearchRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SearchRequest;
    fromJSON(object: any): SearchRequest;
    toJSON(message: SearchRequest): unknown;
    create(base?: DeepPartial<SearchRequest>): SearchRequest;
    fromPartial(object: DeepPartial<SearchRequest>): SearchRequest;
};
export declare const SongAuthor: {
    encode(message: SongAuthor, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SongAuthor;
    fromJSON(object: any): SongAuthor;
    toJSON(message: SongAuthor): unknown;
    create(base?: DeepPartial<SongAuthor>): SongAuthor;
    fromPartial(object: DeepPartial<SongAuthor>): SongAuthor;
};
export declare const Song: {
    encode(message: Song, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Song;
    fromJSON(object: any): Song;
    toJSON(message: Song): unknown;
    create(base?: DeepPartial<Song>): Song;
    fromPartial(object: DeepPartial<Song>): Song;
};
export declare const SearchResponse: {
    encode(message: SearchResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SearchResponse;
    fromJSON(object: any): SearchResponse;
    toJSON(message: SearchResponse): unknown;
    create(base?: DeepPartial<SearchResponse>): SearchResponse;
    fromPartial(object: DeepPartial<SearchResponse>): SearchResponse;
};
export type YtsrDefinition = typeof YtsrDefinition;
export declare const YtsrDefinition: {
    readonly name: "Ytsr";
    readonly fullName: "ytsr.Ytsr";
    readonly methods: {
        readonly search: {
            readonly name: "Search";
            readonly requestType: {
                encode(message: SearchRequest, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): SearchRequest;
                fromJSON(object: any): SearchRequest;
                toJSON(message: SearchRequest): unknown;
                create(base?: DeepPartial<SearchRequest>): SearchRequest;
                fromPartial(object: DeepPartial<SearchRequest>): SearchRequest;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(message: SearchResponse, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): SearchResponse;
                fromJSON(object: any): SearchResponse;
                toJSON(message: SearchResponse): unknown;
                create(base?: DeepPartial<SearchResponse>): SearchResponse;
                fromPartial(object: DeepPartial<SearchResponse>): SearchResponse;
            };
            readonly responseStream: false;
            readonly options: {};
        };
    };
};
export interface YtsrServiceImplementation<CallContextExt = {}> {
    search(request: SearchRequest, context: CallContext & CallContextExt): Promise<DeepPartial<SearchResponse>>;
}
export interface YtsrClient<CallOptionsExt = {}> {
    search(request: DeepPartial<SearchRequest>, options?: CallOptions & CallOptionsExt): Promise<SearchResponse>;
}
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
//# sourceMappingURL=ytsr.d.ts.map