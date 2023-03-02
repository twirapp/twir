/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "bots";

export interface DeleteMessagesRequest {
  channelId: string;
  channelName: string;
  messageIds: string[];
}

export interface SendMessageRequest {
  channelId: string;
  channelName?: string | undefined;
  message: string;
  isAnnounce?: boolean | undefined;
}

export interface JoinOrLeaveRequest {
  botId: string;
  userName: string;
}

function createBaseDeleteMessagesRequest(): DeleteMessagesRequest {
  return { channelId: "", channelName: "", messageIds: [] };
}

export const DeleteMessagesRequest = {
  encode(message: DeleteMessagesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.channelName !== "") {
      writer.uint32(18).string(message.channelName);
    }
    for (const v of message.messageIds) {
      writer.uint32(26).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteMessagesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteMessagesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.channelId = reader.string();
          break;
        case 2:
          message.channelName = reader.string();
          break;
        case 3:
          message.messageIds.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DeleteMessagesRequest {
    return {
      channelId: isSet(object.channelId) ? String(object.channelId) : "",
      channelName: isSet(object.channelName) ? String(object.channelName) : "",
      messageIds: Array.isArray(object?.messageIds) ? object.messageIds.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: DeleteMessagesRequest): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    message.channelName !== undefined && (obj.channelName = message.channelName);
    if (message.messageIds) {
      obj.messageIds = message.messageIds.map((e) => e);
    } else {
      obj.messageIds = [];
    }
    return obj;
  },

  create(base?: DeepPartial<DeleteMessagesRequest>): DeleteMessagesRequest {
    return DeleteMessagesRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeleteMessagesRequest>): DeleteMessagesRequest {
    const message = createBaseDeleteMessagesRequest();
    message.channelId = object.channelId ?? "";
    message.channelName = object.channelName ?? "";
    message.messageIds = object.messageIds?.map((e) => e) || [];
    return message;
  },
};

function createBaseSendMessageRequest(): SendMessageRequest {
  return { channelId: "", channelName: undefined, message: "", isAnnounce: undefined };
}

export const SendMessageRequest = {
  encode(message: SendMessageRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.channelName !== undefined) {
      writer.uint32(18).string(message.channelName);
    }
    if (message.message !== "") {
      writer.uint32(26).string(message.message);
    }
    if (message.isAnnounce !== undefined) {
      writer.uint32(32).bool(message.isAnnounce);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SendMessageRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSendMessageRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.channelId = reader.string();
          break;
        case 2:
          message.channelName = reader.string();
          break;
        case 3:
          message.message = reader.string();
          break;
        case 4:
          message.isAnnounce = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SendMessageRequest {
    return {
      channelId: isSet(object.channelId) ? String(object.channelId) : "",
      channelName: isSet(object.channelName) ? String(object.channelName) : undefined,
      message: isSet(object.message) ? String(object.message) : "",
      isAnnounce: isSet(object.isAnnounce) ? Boolean(object.isAnnounce) : undefined,
    };
  },

  toJSON(message: SendMessageRequest): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    message.channelName !== undefined && (obj.channelName = message.channelName);
    message.message !== undefined && (obj.message = message.message);
    message.isAnnounce !== undefined && (obj.isAnnounce = message.isAnnounce);
    return obj;
  },

  create(base?: DeepPartial<SendMessageRequest>): SendMessageRequest {
    return SendMessageRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SendMessageRequest>): SendMessageRequest {
    const message = createBaseSendMessageRequest();
    message.channelId = object.channelId ?? "";
    message.channelName = object.channelName ?? undefined;
    message.message = object.message ?? "";
    message.isAnnounce = object.isAnnounce ?? undefined;
    return message;
  },
};

function createBaseJoinOrLeaveRequest(): JoinOrLeaveRequest {
  return { botId: "", userName: "" };
}

export const JoinOrLeaveRequest = {
  encode(message: JoinOrLeaveRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.botId !== "") {
      writer.uint32(18).string(message.botId);
    }
    if (message.userName !== "") {
      writer.uint32(26).string(message.userName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): JoinOrLeaveRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseJoinOrLeaveRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.botId = reader.string();
          break;
        case 3:
          message.userName = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): JoinOrLeaveRequest {
    return {
      botId: isSet(object.botId) ? String(object.botId) : "",
      userName: isSet(object.userName) ? String(object.userName) : "",
    };
  },

  toJSON(message: JoinOrLeaveRequest): unknown {
    const obj: any = {};
    message.botId !== undefined && (obj.botId = message.botId);
    message.userName !== undefined && (obj.userName = message.userName);
    return obj;
  },

  create(base?: DeepPartial<JoinOrLeaveRequest>): JoinOrLeaveRequest {
    return JoinOrLeaveRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<JoinOrLeaveRequest>): JoinOrLeaveRequest {
    const message = createBaseJoinOrLeaveRequest();
    message.botId = object.botId ?? "";
    message.userName = object.userName ?? "";
    return message;
  },
};

export type BotsDefinition = typeof BotsDefinition;
export const BotsDefinition = {
  name: "Bots",
  fullName: "bots.Bots",
  methods: {
    deleteMessage: {
      name: "DeleteMessage",
      requestType: DeleteMessagesRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    sendMessage: {
      name: "SendMessage",
      requestType: SendMessageRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    join: {
      name: "Join",
      requestType: JoinOrLeaveRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    leave: {
      name: "Leave",
      requestType: JoinOrLeaveRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface BotsServiceImplementation<CallContextExt = {}> {
  deleteMessage(request: DeleteMessagesRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  sendMessage(request: SendMessageRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  join(request: JoinOrLeaveRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  leave(request: JoinOrLeaveRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
}

export interface BotsClient<CallOptionsExt = {}> {
  deleteMessage(request: DeepPartial<DeleteMessagesRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  sendMessage(request: DeepPartial<SendMessageRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  join(request: DeepPartial<JoinOrLeaveRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  leave(request: DeepPartial<JoinOrLeaveRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
