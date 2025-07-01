/* eslint-disable */
import Long from "long";
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal.js";
import { Empty } from "../google/protobuf/empty.js";

export const protobufPackage = "tokens";

export interface Token {
  accessToken: string;
  scopes: string[];
  expiresIn: number;
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
  return { accessToken: "", scopes: [], expiresIn: 0 };
}

export const Token = {
  encode(message: Token, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.accessToken !== "") {
      writer.uint32(10).string(message.accessToken);
    }
    for (const v of message.scopes) {
      writer.uint32(18).string(v!);
    }
    if (message.expiresIn !== 0) {
      writer.uint32(24).int32(message.expiresIn);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Token {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseToken();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.accessToken = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.scopes.push(reader.string());
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.expiresIn = reader.int32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Token {
    return {
      accessToken: isSet(object.accessToken) ? globalThis.String(object.accessToken) : "",
      scopes: globalThis.Array.isArray(object?.scopes) ? object.scopes.map((e: any) => globalThis.String(e)) : [],
      expiresIn: isSet(object.expiresIn) ? globalThis.Number(object.expiresIn) : 0,
    };
  },

  toJSON(message: Token): unknown {
    const obj: any = {};
    if (message.accessToken !== "") {
      obj.accessToken = message.accessToken;
    }
    if (message.scopes?.length) {
      obj.scopes = message.scopes;
    }
    if (message.expiresIn !== 0) {
      obj.expiresIn = Math.round(message.expiresIn);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Token>, I>>(base?: I): Token {
    return Token.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Token>, I>>(object: I): Token {
    const message = createBaseToken();
    message.accessToken = object.accessToken ?? "";
    message.scopes = object.scopes?.map((e) => e) || [];
    message.expiresIn = object.expiresIn ?? 0;
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
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetUserTokenRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.userId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetUserTokenRequest {
    return { userId: isSet(object.userId) ? globalThis.String(object.userId) : "" };
  },

  toJSON(message: GetUserTokenRequest): unknown {
    const obj: any = {};
    if (message.userId !== "") {
      obj.userId = message.userId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetUserTokenRequest>, I>>(base?: I): GetUserTokenRequest {
    return GetUserTokenRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetUserTokenRequest>, I>>(object: I): GetUserTokenRequest {
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
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetBotTokenRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.botId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetBotTokenRequest {
    return { botId: isSet(object.botId) ? globalThis.String(object.botId) : "" };
  },

  toJSON(message: GetBotTokenRequest): unknown {
    const obj: any = {};
    if (message.botId !== "") {
      obj.botId = message.botId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetBotTokenRequest>, I>>(base?: I): GetBotTokenRequest {
    return GetBotTokenRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetBotTokenRequest>, I>>(object: I): GetBotTokenRequest {
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
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateToken();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.accessToken = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.refreshToken = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.expiresIn = longToNumber(reader.int64() as Long);
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.scopes.push(reader.string());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateToken {
    return {
      accessToken: isSet(object.accessToken) ? globalThis.String(object.accessToken) : "",
      refreshToken: isSet(object.refreshToken) ? globalThis.String(object.refreshToken) : "",
      expiresIn: isSet(object.expiresIn) ? globalThis.Number(object.expiresIn) : 0,
      scopes: globalThis.Array.isArray(object?.scopes) ? object.scopes.map((e: any) => globalThis.String(e)) : [],
    };
  },

  toJSON(message: UpdateToken): unknown {
    const obj: any = {};
    if (message.accessToken !== "") {
      obj.accessToken = message.accessToken;
    }
    if (message.refreshToken !== "") {
      obj.refreshToken = message.refreshToken;
    }
    if (message.expiresIn !== 0) {
      obj.expiresIn = Math.round(message.expiresIn);
    }
    if (message.scopes?.length) {
      obj.scopes = message.scopes;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateToken>, I>>(base?: I): UpdateToken {
    return UpdateToken.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateToken>, I>>(object: I): UpdateToken {
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

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToNumber(long: Long): number {
  if (long.gt(globalThis.Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
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
