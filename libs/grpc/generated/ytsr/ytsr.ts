/* eslint-disable */
import Long from "long";
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "ytsr";

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

function createBaseSearchRequest(): SearchRequest {
  return { search: "" };
}

export const SearchRequest = {
  encode(message: SearchRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.search !== "") {
      writer.uint32(10).string(message.search);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SearchRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSearchRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.search = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SearchRequest {
    return { search: isSet(object.search) ? String(object.search) : "" };
  },

  toJSON(message: SearchRequest): unknown {
    const obj: any = {};
    message.search !== undefined && (obj.search = message.search);
    return obj;
  },

  create(base?: DeepPartial<SearchRequest>): SearchRequest {
    return SearchRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SearchRequest>): SearchRequest {
    const message = createBaseSearchRequest();
    message.search = object.search ?? "";
    return message;
  },
};

function createBaseSongAuthor(): SongAuthor {
  return { name: "", channelId: "", avatarUrl: undefined };
}

export const SongAuthor = {
  encode(message: SongAuthor, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.channelId !== "") {
      writer.uint32(18).string(message.channelId);
    }
    if (message.avatarUrl !== undefined) {
      writer.uint32(26).string(message.avatarUrl);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SongAuthor {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSongAuthor();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.channelId = reader.string();
          break;
        case 3:
          message.avatarUrl = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SongAuthor {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      channelId: isSet(object.channelId) ? String(object.channelId) : "",
      avatarUrl: isSet(object.avatarUrl) ? String(object.avatarUrl) : undefined,
    };
  },

  toJSON(message: SongAuthor): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.channelId !== undefined && (obj.channelId = message.channelId);
    message.avatarUrl !== undefined && (obj.avatarUrl = message.avatarUrl);
    return obj;
  },

  create(base?: DeepPartial<SongAuthor>): SongAuthor {
    return SongAuthor.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SongAuthor>): SongAuthor {
    const message = createBaseSongAuthor();
    message.name = object.name ?? "";
    message.channelId = object.channelId ?? "";
    message.avatarUrl = object.avatarUrl ?? undefined;
    return message;
  },
};

function createBaseSong(): Song {
  return { title: "", id: "", views: 0, duration: 0, thumbnailUrl: undefined, isLive: false, author: undefined };
}

export const Song = {
  encode(message: Song, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.title !== "") {
      writer.uint32(10).string(message.title);
    }
    if (message.id !== "") {
      writer.uint32(18).string(message.id);
    }
    if (message.views !== 0) {
      writer.uint32(24).uint64(message.views);
    }
    if (message.duration !== 0) {
      writer.uint32(32).uint64(message.duration);
    }
    if (message.thumbnailUrl !== undefined) {
      writer.uint32(42).string(message.thumbnailUrl);
    }
    if (message.isLive === true) {
      writer.uint32(48).bool(message.isLive);
    }
    if (message.author !== undefined) {
      SongAuthor.encode(message.author, writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Song {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSong();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.title = reader.string();
          break;
        case 2:
          message.id = reader.string();
          break;
        case 3:
          message.views = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.duration = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.thumbnailUrl = reader.string();
          break;
        case 6:
          message.isLive = reader.bool();
          break;
        case 7:
          message.author = SongAuthor.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Song {
    return {
      title: isSet(object.title) ? String(object.title) : "",
      id: isSet(object.id) ? String(object.id) : "",
      views: isSet(object.views) ? Number(object.views) : 0,
      duration: isSet(object.duration) ? Number(object.duration) : 0,
      thumbnailUrl: isSet(object.thumbnailUrl) ? String(object.thumbnailUrl) : undefined,
      isLive: isSet(object.isLive) ? Boolean(object.isLive) : false,
      author: isSet(object.author) ? SongAuthor.fromJSON(object.author) : undefined,
    };
  },

  toJSON(message: Song): unknown {
    const obj: any = {};
    message.title !== undefined && (obj.title = message.title);
    message.id !== undefined && (obj.id = message.id);
    message.views !== undefined && (obj.views = Math.round(message.views));
    message.duration !== undefined && (obj.duration = Math.round(message.duration));
    message.thumbnailUrl !== undefined && (obj.thumbnailUrl = message.thumbnailUrl);
    message.isLive !== undefined && (obj.isLive = message.isLive);
    message.author !== undefined && (obj.author = message.author ? SongAuthor.toJSON(message.author) : undefined);
    return obj;
  },

  create(base?: DeepPartial<Song>): Song {
    return Song.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Song>): Song {
    const message = createBaseSong();
    message.title = object.title ?? "";
    message.id = object.id ?? "";
    message.views = object.views ?? 0;
    message.duration = object.duration ?? 0;
    message.thumbnailUrl = object.thumbnailUrl ?? undefined;
    message.isLive = object.isLive ?? false;
    message.author = (object.author !== undefined && object.author !== null)
      ? SongAuthor.fromPartial(object.author)
      : undefined;
    return message;
  },
};

function createBaseSearchResponse(): SearchResponse {
  return { songs: [] };
}

export const SearchResponse = {
  encode(message: SearchResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.songs) {
      Song.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SearchResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSearchResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.songs.push(Song.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SearchResponse {
    return { songs: Array.isArray(object?.songs) ? object.songs.map((e: any) => Song.fromJSON(e)) : [] };
  },

  toJSON(message: SearchResponse): unknown {
    const obj: any = {};
    if (message.songs) {
      obj.songs = message.songs.map((e) => e ? Song.toJSON(e) : undefined);
    } else {
      obj.songs = [];
    }
    return obj;
  },

  create(base?: DeepPartial<SearchResponse>): SearchResponse {
    return SearchResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SearchResponse>): SearchResponse {
    const message = createBaseSearchResponse();
    message.songs = object.songs?.map((e) => Song.fromPartial(e)) || [];
    return message;
  },
};

export type YtsrDefinition = typeof YtsrDefinition;
export const YtsrDefinition = {
  name: "Ytsr",
  fullName: "ytsr.Ytsr",
  methods: {
    search: {
      name: "Search",
      requestType: SearchRequest,
      requestStream: false,
      responseType: SearchResponse,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface YtsrServiceImplementation<CallContextExt = {}> {
  search(request: SearchRequest, context: CallContext & CallContextExt): Promise<DeepPartial<SearchResponse>>;
}

export interface YtsrClient<CallOptionsExt = {}> {
  search(request: DeepPartial<SearchRequest>, options?: CallOptions & CallOptionsExt): Promise<SearchResponse>;
}

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var tsProtoGlobalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new tsProtoGlobalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
