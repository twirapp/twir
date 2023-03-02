/* eslint-disable */
import Long from "long";
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "tokens";

export interface Token {
  accessToken: string;
  scopes: string[];
}

export interface GetUserTokenRequest {
  userId: string;
}

export interface GetBotTokenRequest {
  botId: string;
}

export interface UpdateToken {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
  scopes: string[];
}

function createBaseToken(): Token {
  return { accessToken: "", scopes: [] };
}

export const Token = {
  encode(message: Token, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.accessToken !== "") {
      writer.uint32(10).string(message.accessToken);
    }
    for (const v of message.scopes) {
      writer.uint32(18).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Token {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseToken();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.accessToken = reader.string();
          break;
        case 2:
          message.scopes.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Token {
    return {
      accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: Token): unknown {
    const obj: any = {};
    message.accessToken !== undefined && (obj.accessToken = message.accessToken);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    return obj;
  },

  create(base?: DeepPartial<Token>): Token {
    return Token.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Token>): Token {
    const message = createBaseToken();
    message.accessToken = object.accessToken ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    return message;
  },
};

function createBaseGetUserTokenRequest(): GetUserTokenRequest {
  return { userId: "" };
}

export const GetUserTokenRequest = {
  encode(message: GetUserTokenRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetUserTokenRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetUserTokenRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.userId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetUserTokenRequest {
    return { userId: isSet(object.userId) ? String(object.userId) : "" };
  },

  toJSON(message: GetUserTokenRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<GetUserTokenRequest>): GetUserTokenRequest {
    return GetUserTokenRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetUserTokenRequest>): GetUserTokenRequest {
    const message = createBaseGetUserTokenRequest();
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseGetBotTokenRequest(): GetBotTokenRequest {
  return { botId: "" };
}

export const GetBotTokenRequest = {
  encode(message: GetBotTokenRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.botId !== "") {
      writer.uint32(10).string(message.botId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetBotTokenRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetBotTokenRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.botId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetBotTokenRequest {
    return { botId: isSet(object.botId) ? String(object.botId) : "" };
  },

  toJSON(message: GetBotTokenRequest): unknown {
    const obj: any = {};
    message.botId !== undefined && (obj.botId = message.botId);
    return obj;
  },

  create(base?: DeepPartial<GetBotTokenRequest>): GetBotTokenRequest {
    return GetBotTokenRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetBotTokenRequest>): GetBotTokenRequest {
    const message = createBaseGetBotTokenRequest();
    message.botId = object.botId ?? "";
    return message;
  },
};

function createBaseUpdateToken(): UpdateToken {
  return { accessToken: "", refreshToken: "", expiresIn: 0, scopes: [] };
}

export const UpdateToken = {
  encode(message: UpdateToken, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.accessToken !== "") {
      writer.uint32(10).string(message.accessToken);
    }
    if (message.refreshToken !== "") {
      writer.uint32(18).string(message.refreshToken);
    }
    if (message.expiresIn !== 0) {
      writer.uint32(24).int64(message.expiresIn);
    }
    for (const v of message.scopes) {
      writer.uint32(34).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateToken {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateToken();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.accessToken = reader.string();
          break;
        case 2:
          message.refreshToken = reader.string();
          break;
        case 3:
          message.expiresIn = longToNumber(reader.int64() as Long);
          break;
        case 4:
          message.scopes.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateToken {
    return {
      accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
      refreshToken: isSet(object.refreshToken) ? String(object.refreshToken) : "",
      expiresIn: isSet(object.expiresIn) ? Number(object.expiresIn) : 0,
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: UpdateToken): unknown {
    const obj: any = {};
    message.accessToken !== undefined && (obj.accessToken = message.accessToken);
    message.refreshToken !== undefined && (obj.refreshToken = message.refreshToken);
    message.expiresIn !== undefined && (obj.expiresIn = Math.round(message.expiresIn));
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => e);
    } else {
      obj.scopes = [];
    }
    return obj;
  },

  create(base?: DeepPartial<UpdateToken>): UpdateToken {
    return UpdateToken.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateToken>): UpdateToken {
    const message = createBaseUpdateToken();
    message.accessToken = object.accessToken ?? "";
    message.refreshToken = object.refreshToken ?? "";
    message.expiresIn = object.expiresIn ?? 0;
    message.scopes = object.scopes?.map((e) => e) || [];
    return message;
  },
};

export type TokensDefinition = typeof TokensDefinition;
export const TokensDefinition = {
  name: "Tokens",
  fullName: "tokens.Tokens",
  methods: {
    requestAppToken: {
      name: "RequestAppToken",
      requestType: Empty,
      requestStream: false,
      responseType: Token,
      responseStream: false,
      options: {},
    },
    requestUserToken: {
      name: "RequestUserToken",
      requestType: GetUserTokenRequest,
      requestStream: false,
      responseType: Token,
      responseStream: false,
      options: {},
    },
    requestBotToken: {
      name: "RequestBotToken",
      requestType: GetBotTokenRequest,
      requestStream: false,
      responseType: Token,
      responseStream: false,
      options: {},
    },
    updateBotToken: {
      name: "UpdateBotToken",
      requestType: UpdateToken,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    updateUserToken: {
      name: "UpdateUserToken",
      requestType: UpdateToken,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface TokensServiceImplementation<CallContextExt = {}> {
  requestAppToken(request: Empty, context: CallContext & CallContextExt): Promise<DeepPartial<Token>>;
  requestUserToken(request: GetUserTokenRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Token>>;
  requestBotToken(request: GetBotTokenRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Token>>;
  updateBotToken(request: UpdateToken, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  updateUserToken(request: UpdateToken, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
}

export interface TokensClient<CallOptionsExt = {}> {
  requestAppToken(request: DeepPartial<Empty>, options?: CallOptions & CallOptionsExt): Promise<Token>;
  requestUserToken(request: DeepPartial<GetUserTokenRequest>, options?: CallOptions & CallOptionsExt): Promise<Token>;
  requestBotToken(request: DeepPartial<GetBotTokenRequest>, options?: CallOptions & CallOptionsExt): Promise<Token>;
  updateBotToken(request: DeepPartial<UpdateToken>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  updateUserToken(request: DeepPartial<UpdateToken>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
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
