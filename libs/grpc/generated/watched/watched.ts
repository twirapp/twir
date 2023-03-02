/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "watched";

export interface Request {
  channelsId: string[];
  botId: string;
}

function createBaseRequest(): Request {
  return { channelsId: [], botId: "" };
}

export const Request = {
  encode(message: Request, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.channelsId) {
      writer.uint32(10).string(v!);
    }
    if (message.botId !== "") {
      writer.uint32(18).string(message.botId);
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
          message.channelsId.push(reader.string());
          break;
        case 2:
          message.botId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Request {
    return {
      channelsId: Array.isArray(object?.channelsId) ? object.channelsId.map((e: any) => String(e)) : [],
      botId: isSet(object.botId) ? String(object.botId) : "",
    };
  },

  toJSON(message: Request): unknown {
    const obj: any = {};
    if (message.channelsId) {
      obj.channelsId = message.channelsId.map((e) => e);
    } else {
      obj.channelsId = [];
    }
    message.botId !== undefined && (obj.botId = message.botId);
    return obj;
  },

  create(base?: DeepPartial<Request>): Request {
    return Request.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Request>): Request {
    const message = createBaseRequest();
    message.channelsId = object.channelsId?.map((e) => e) || [];
    message.botId = object.botId ?? "";
    return message;
  },
};

export type WatchedDefinition = typeof WatchedDefinition;
export const WatchedDefinition = {
  name: "Watched",
  fullName: "watched.Watched",
  methods: {
    incrementByChannelId: {
      name: "IncrementByChannelId",
      requestType: Request,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface WatchedServiceImplementation<CallContextExt = {}> {
  incrementByChannelId(request: Request, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
}

export interface WatchedClient<CallOptionsExt = {}> {
  incrementByChannelId(request: DeepPartial<Request>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
