/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "emotes_cacher";

export interface Request {
  channelId: string;
}

function createBaseRequest(): Request {
  return { channelId: "" };
}

export const Request = {
  encode(message: Request, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Request {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.channelId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Request {
    return { channelId: isSet(object.channelId) ? String(object.channelId) : "" };
  },

  toJSON(message: Request): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    return obj;
  },

  create(base?: DeepPartial<Request>): Request {
    return Request.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Request>): Request {
    const message = createBaseRequest();
    message.channelId = object.channelId ?? "";
    return message;
  },
};

export type EmotesCacherDefinition = typeof EmotesCacherDefinition;
export const EmotesCacherDefinition = {
  name: "EmotesCacher",
  fullName: "emotes_cacher.EmotesCacher",
  methods: {
    cacheChannelEmotes: {
      name: "CacheChannelEmotes",
      requestType: Request,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    cacheGlobalEmotes: {
      name: "CacheGlobalEmotes",
      requestType: Empty,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface EmotesCacherServiceImplementation<CallContextExt = {}> {
  cacheChannelEmotes(request: Request, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  cacheGlobalEmotes(request: Empty, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
}

export interface EmotesCacherClient<CallOptionsExt = {}> {
  cacheChannelEmotes(request: DeepPartial<Request>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  cacheGlobalEmotes(request: DeepPartial<Empty>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
