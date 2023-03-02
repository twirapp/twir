/* eslint-disable */
import Long from "long";
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "dota";

export interface GetPlayerCardRequest {
  accountId: number;
}

export interface GetPlayerCardResponse {
  accountId: string;
  rankTier?: number | undefined;
  leaderboardRank?: number | undefined;
}

function createBaseGetPlayerCardRequest(): GetPlayerCardRequest {
  return { accountId: 0 };
}

export const GetPlayerCardRequest = {
  encode(message: GetPlayerCardRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.accountId !== 0) {
      writer.uint32(8).int64(message.accountId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetPlayerCardRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetPlayerCardRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.accountId = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetPlayerCardRequest {
    return { accountId: isSet(object.accountId) ? Number(object.accountId) : 0 };
  },

  toJSON(message: GetPlayerCardRequest): unknown {
    const obj: any = {};
    message.accountId !== undefined && (obj.accountId = Math.round(message.accountId));
    return obj;
  },

  create(base?: DeepPartial<GetPlayerCardRequest>): GetPlayerCardRequest {
    return GetPlayerCardRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetPlayerCardRequest>): GetPlayerCardRequest {
    const message = createBaseGetPlayerCardRequest();
    message.accountId = object.accountId ?? 0;
    return message;
  },
};

function createBaseGetPlayerCardResponse(): GetPlayerCardResponse {
  return { accountId: "", rankTier: undefined, leaderboardRank: undefined };
}

export const GetPlayerCardResponse = {
  encode(message: GetPlayerCardResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.accountId !== "") {
      writer.uint32(10).string(message.accountId);
    }
    if (message.rankTier !== undefined) {
      writer.uint32(16).int64(message.rankTier);
    }
    if (message.leaderboardRank !== undefined) {
      writer.uint32(24).int64(message.leaderboardRank);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetPlayerCardResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetPlayerCardResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.accountId = reader.string();
          break;
        case 2:
          message.rankTier = longToNumber(reader.int64() as Long);
          break;
        case 3:
          message.leaderboardRank = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetPlayerCardResponse {
    return {
      accountId: isSet(object.accountId) ? String(object.accountId) : "",
      rankTier: isSet(object.rankTier) ? Number(object.rankTier) : undefined,
      leaderboardRank: isSet(object.leaderboardRank) ? Number(object.leaderboardRank) : undefined,
    };
  },

  toJSON(message: GetPlayerCardResponse): unknown {
    const obj: any = {};
    message.accountId !== undefined && (obj.accountId = message.accountId);
    message.rankTier !== undefined && (obj.rankTier = Math.round(message.rankTier));
    message.leaderboardRank !== undefined && (obj.leaderboardRank = Math.round(message.leaderboardRank));
    return obj;
  },

  create(base?: DeepPartial<GetPlayerCardResponse>): GetPlayerCardResponse {
    return GetPlayerCardResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetPlayerCardResponse>): GetPlayerCardResponse {
    const message = createBaseGetPlayerCardResponse();
    message.accountId = object.accountId ?? "";
    message.rankTier = object.rankTier ?? undefined;
    message.leaderboardRank = object.leaderboardRank ?? undefined;
    return message;
  },
};

export type DotaDefinition = typeof DotaDefinition;
export const DotaDefinition = {
  name: "Dota",
  fullName: "dota.Dota",
  methods: {
    getPlayerCard: {
      name: "GetPlayerCard",
      requestType: GetPlayerCardRequest,
      requestStream: false,
      responseType: GetPlayerCardResponse,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface DotaServiceImplementation<CallContextExt = {}> {
  getPlayerCard(
    request: GetPlayerCardRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetPlayerCardResponse>>;
}

export interface DotaClient<CallOptionsExt = {}> {
  getPlayerCard(
    request: DeepPartial<GetPlayerCardRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetPlayerCardResponse>;
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
