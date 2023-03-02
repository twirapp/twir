/* eslint-disable */
import Long from "long";
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "parser";

export interface Sender {
  id: string;
  name: string;
  displayName: string;
  badges: string[];
}

export interface Channel {
  id: string;
  name: string;
}

export interface Message {
  text: string;
  id: string;
  emotes: Message_Emote[];
}

export interface Message_EmotePosition {
  start: number;
  end: number;
}

export interface Message_Emote {
  name: string;
  id: string;
  count: number;
  positions: Message_EmotePosition[];
}

export interface ProcessCommandRequest {
  sender: Sender | undefined;
  channel: Channel | undefined;
  message: Message | undefined;
}

export interface ProcessCommandResponse {
  responses: string[];
  isReply: boolean;
  keepOrder?: boolean | undefined;
}

export interface GetVariablesResponse {
  list: GetVariablesResponse_Variable[];
}

export interface GetVariablesResponse_Variable {
  name: string;
  example: string;
  description: string;
  visible: boolean;
}

export interface GetDefaultCommandsResponse {
  list: GetDefaultCommandsResponse_DefaultCommand[];
}

export interface GetDefaultCommandsResponse_DefaultCommand {
  name: string;
  description: string;
  visible: boolean;
  rolesNames: string[];
  module: string;
  isReply: boolean;
  keepResponsesOrder: boolean;
  aliases: string[];
}

export interface ParseTextRequestData {
  sender: Sender | undefined;
  channel: Channel | undefined;
  message: Message | undefined;
  parseVariables?: boolean | undefined;
}

export interface ParseTextResponseData {
  responses: string[];
}

function createBaseSender(): Sender {
  return { id: "", name: "", displayName: "", badges: [] };
}

export const Sender = {
  encode(message: Sender, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.displayName !== "") {
      writer.uint32(26).string(message.displayName);
    }
    for (const v of message.badges) {
      writer.uint32(34).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Sender {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSender();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
          message.displayName = reader.string();
          break;
        case 4:
          message.badges.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Sender {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      name: isSet(object.name) ? String(object.name) : "",
      displayName: isSet(object.displayName) ? String(object.displayName) : "",
      badges: Array.isArray(object?.badges) ? object.badges.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: Sender): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.name !== undefined && (obj.name = message.name);
    message.displayName !== undefined && (obj.displayName = message.displayName);
    if (message.badges) {
      obj.badges = message.badges.map((e) => e);
    } else {
      obj.badges = [];
    }
    return obj;
  },

  create(base?: DeepPartial<Sender>): Sender {
    return Sender.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Sender>): Sender {
    const message = createBaseSender();
    message.id = object.id ?? "";
    message.name = object.name ?? "";
    message.displayName = object.displayName ?? "";
    message.badges = object.badges?.map((e) => e) || [];
    return message;
  },
};

function createBaseChannel(): Channel {
  return { id: "", name: "" };
}

export const Channel = {
  encode(message: Channel, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Channel {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChannel();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Channel {
    return { id: isSet(object.id) ? String(object.id) : "", name: isSet(object.name) ? String(object.name) : "" };
  },

  toJSON(message: Channel): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.name !== undefined && (obj.name = message.name);
    return obj;
  },

  create(base?: DeepPartial<Channel>): Channel {
    return Channel.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Channel>): Channel {
    const message = createBaseChannel();
    message.id = object.id ?? "";
    message.name = object.name ?? "";
    return message;
  },
};

function createBaseMessage(): Message {
  return { text: "", id: "", emotes: [] };
}

export const Message = {
  encode(message: Message, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.text !== "") {
      writer.uint32(10).string(message.text);
    }
    if (message.id !== "") {
      writer.uint32(18).string(message.id);
    }
    for (const v of message.emotes) {
      Message_Emote.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Message {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.text = reader.string();
          break;
        case 2:
          message.id = reader.string();
          break;
        case 3:
          message.emotes.push(Message_Emote.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Message {
    return {
      text: isSet(object.text) ? String(object.text) : "",
      id: isSet(object.id) ? String(object.id) : "",
      emotes: Array.isArray(object?.emotes) ? object.emotes.map((e: any) => Message_Emote.fromJSON(e)) : [],
    };
  },

  toJSON(message: Message): unknown {
    const obj: any = {};
    message.text !== undefined && (obj.text = message.text);
    message.id !== undefined && (obj.id = message.id);
    if (message.emotes) {
      obj.emotes = message.emotes.map((e) => e ? Message_Emote.toJSON(e) : undefined);
    } else {
      obj.emotes = [];
    }
    return obj;
  },

  create(base?: DeepPartial<Message>): Message {
    return Message.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Message>): Message {
    const message = createBaseMessage();
    message.text = object.text ?? "";
    message.id = object.id ?? "";
    message.emotes = object.emotes?.map((e) => Message_Emote.fromPartial(e)) || [];
    return message;
  },
};

function createBaseMessage_EmotePosition(): Message_EmotePosition {
  return { start: 0, end: 0 };
}

export const Message_EmotePosition = {
  encode(message: Message_EmotePosition, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.start !== 0) {
      writer.uint32(8).int64(message.start);
    }
    if (message.end !== 0) {
      writer.uint32(16).int64(message.end);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Message_EmotePosition {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMessage_EmotePosition();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.start = longToNumber(reader.int64() as Long);
          break;
        case 2:
          message.end = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Message_EmotePosition {
    return { start: isSet(object.start) ? Number(object.start) : 0, end: isSet(object.end) ? Number(object.end) : 0 };
  },

  toJSON(message: Message_EmotePosition): unknown {
    const obj: any = {};
    message.start !== undefined && (obj.start = Math.round(message.start));
    message.end !== undefined && (obj.end = Math.round(message.end));
    return obj;
  },

  create(base?: DeepPartial<Message_EmotePosition>): Message_EmotePosition {
    return Message_EmotePosition.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Message_EmotePosition>): Message_EmotePosition {
    const message = createBaseMessage_EmotePosition();
    message.start = object.start ?? 0;
    message.end = object.end ?? 0;
    return message;
  },
};

function createBaseMessage_Emote(): Message_Emote {
  return { name: "", id: "", count: 0, positions: [] };
}

export const Message_Emote = {
  encode(message: Message_Emote, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.id !== "") {
      writer.uint32(18).string(message.id);
    }
    if (message.count !== 0) {
      writer.uint32(24).int64(message.count);
    }
    for (const v of message.positions) {
      Message_EmotePosition.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Message_Emote {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMessage_Emote();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.id = reader.string();
          break;
        case 3:
          message.count = longToNumber(reader.int64() as Long);
          break;
        case 4:
          message.positions.push(Message_EmotePosition.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Message_Emote {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      id: isSet(object.id) ? String(object.id) : "",
      count: isSet(object.count) ? Number(object.count) : 0,
      positions: Array.isArray(object?.positions)
        ? object.positions.map((e: any) => Message_EmotePosition.fromJSON(e))
        : [],
    };
  },

  toJSON(message: Message_Emote): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.id !== undefined && (obj.id = message.id);
    message.count !== undefined && (obj.count = Math.round(message.count));
    if (message.positions) {
      obj.positions = message.positions.map((e) => e ? Message_EmotePosition.toJSON(e) : undefined);
    } else {
      obj.positions = [];
    }
    return obj;
  },

  create(base?: DeepPartial<Message_Emote>): Message_Emote {
    return Message_Emote.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Message_Emote>): Message_Emote {
    const message = createBaseMessage_Emote();
    message.name = object.name ?? "";
    message.id = object.id ?? "";
    message.count = object.count ?? 0;
    message.positions = object.positions?.map((e) => Message_EmotePosition.fromPartial(e)) || [];
    return message;
  },
};

function createBaseProcessCommandRequest(): ProcessCommandRequest {
  return { sender: undefined, channel: undefined, message: undefined };
}

export const ProcessCommandRequest = {
  encode(message: ProcessCommandRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sender !== undefined) {
      Sender.encode(message.sender, writer.uint32(10).fork()).ldelim();
    }
    if (message.channel !== undefined) {
      Channel.encode(message.channel, writer.uint32(18).fork()).ldelim();
    }
    if (message.message !== undefined) {
      Message.encode(message.message, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ProcessCommandRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProcessCommandRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sender = Sender.decode(reader, reader.uint32());
          break;
        case 2:
          message.channel = Channel.decode(reader, reader.uint32());
          break;
        case 3:
          message.message = Message.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ProcessCommandRequest {
    return {
      sender: isSet(object.sender) ? Sender.fromJSON(object.sender) : undefined,
      channel: isSet(object.channel) ? Channel.fromJSON(object.channel) : undefined,
      message: isSet(object.message) ? Message.fromJSON(object.message) : undefined,
    };
  },

  toJSON(message: ProcessCommandRequest): unknown {
    const obj: any = {};
    message.sender !== undefined && (obj.sender = message.sender ? Sender.toJSON(message.sender) : undefined);
    message.channel !== undefined && (obj.channel = message.channel ? Channel.toJSON(message.channel) : undefined);
    message.message !== undefined && (obj.message = message.message ? Message.toJSON(message.message) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ProcessCommandRequest>): ProcessCommandRequest {
    return ProcessCommandRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ProcessCommandRequest>): ProcessCommandRequest {
    const message = createBaseProcessCommandRequest();
    message.sender = (object.sender !== undefined && object.sender !== null)
      ? Sender.fromPartial(object.sender)
      : undefined;
    message.channel = (object.channel !== undefined && object.channel !== null)
      ? Channel.fromPartial(object.channel)
      : undefined;
    message.message = (object.message !== undefined && object.message !== null)
      ? Message.fromPartial(object.message)
      : undefined;
    return message;
  },
};

function createBaseProcessCommandResponse(): ProcessCommandResponse {
  return { responses: [], isReply: false, keepOrder: undefined };
}

export const ProcessCommandResponse = {
  encode(message: ProcessCommandResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.responses) {
      writer.uint32(10).string(v!);
    }
    if (message.isReply === true) {
      writer.uint32(16).bool(message.isReply);
    }
    if (message.keepOrder !== undefined) {
      writer.uint32(24).bool(message.keepOrder);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ProcessCommandResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProcessCommandResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.responses.push(reader.string());
          break;
        case 2:
          message.isReply = reader.bool();
          break;
        case 3:
          message.keepOrder = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ProcessCommandResponse {
    return {
      responses: Array.isArray(object?.responses) ? object.responses.map((e: any) => String(e)) : [],
      isReply: isSet(object.isReply) ? Boolean(object.isReply) : false,
      keepOrder: isSet(object.keepOrder) ? Boolean(object.keepOrder) : undefined,
    };
  },

  toJSON(message: ProcessCommandResponse): unknown {
    const obj: any = {};
    if (message.responses) {
      obj.responses = message.responses.map((e) => e);
    } else {
      obj.responses = [];
    }
    message.isReply !== undefined && (obj.isReply = message.isReply);
    message.keepOrder !== undefined && (obj.keepOrder = message.keepOrder);
    return obj;
  },

  create(base?: DeepPartial<ProcessCommandResponse>): ProcessCommandResponse {
    return ProcessCommandResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ProcessCommandResponse>): ProcessCommandResponse {
    const message = createBaseProcessCommandResponse();
    message.responses = object.responses?.map((e) => e) || [];
    message.isReply = object.isReply ?? false;
    message.keepOrder = object.keepOrder ?? undefined;
    return message;
  },
};

function createBaseGetVariablesResponse(): GetVariablesResponse {
  return { list: [] };
}

export const GetVariablesResponse = {
  encode(message: GetVariablesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.list) {
      GetVariablesResponse_Variable.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetVariablesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetVariablesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.list.push(GetVariablesResponse_Variable.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetVariablesResponse {
    return {
      list: Array.isArray(object?.list) ? object.list.map((e: any) => GetVariablesResponse_Variable.fromJSON(e)) : [],
    };
  },

  toJSON(message: GetVariablesResponse): unknown {
    const obj: any = {};
    if (message.list) {
      obj.list = message.list.map((e) => e ? GetVariablesResponse_Variable.toJSON(e) : undefined);
    } else {
      obj.list = [];
    }
    return obj;
  },

  create(base?: DeepPartial<GetVariablesResponse>): GetVariablesResponse {
    return GetVariablesResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetVariablesResponse>): GetVariablesResponse {
    const message = createBaseGetVariablesResponse();
    message.list = object.list?.map((e) => GetVariablesResponse_Variable.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetVariablesResponse_Variable(): GetVariablesResponse_Variable {
  return { name: "", example: "", description: "", visible: false };
}

export const GetVariablesResponse_Variable = {
  encode(message: GetVariablesResponse_Variable, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.example !== "") {
      writer.uint32(18).string(message.example);
    }
    if (message.description !== "") {
      writer.uint32(26).string(message.description);
    }
    if (message.visible === true) {
      writer.uint32(32).bool(message.visible);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetVariablesResponse_Variable {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetVariablesResponse_Variable();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.example = reader.string();
          break;
        case 3:
          message.description = reader.string();
          break;
        case 4:
          message.visible = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetVariablesResponse_Variable {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      example: isSet(object.example) ? String(object.example) : "",
      description: isSet(object.description) ? String(object.description) : "",
      visible: isSet(object.visible) ? Boolean(object.visible) : false,
    };
  },

  toJSON(message: GetVariablesResponse_Variable): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.example !== undefined && (obj.example = message.example);
    message.description !== undefined && (obj.description = message.description);
    message.visible !== undefined && (obj.visible = message.visible);
    return obj;
  },

  create(base?: DeepPartial<GetVariablesResponse_Variable>): GetVariablesResponse_Variable {
    return GetVariablesResponse_Variable.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetVariablesResponse_Variable>): GetVariablesResponse_Variable {
    const message = createBaseGetVariablesResponse_Variable();
    message.name = object.name ?? "";
    message.example = object.example ?? "";
    message.description = object.description ?? "";
    message.visible = object.visible ?? false;
    return message;
  },
};

function createBaseGetDefaultCommandsResponse(): GetDefaultCommandsResponse {
  return { list: [] };
}

export const GetDefaultCommandsResponse = {
  encode(message: GetDefaultCommandsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.list) {
      GetDefaultCommandsResponse_DefaultCommand.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultCommandsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultCommandsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.list.push(GetDefaultCommandsResponse_DefaultCommand.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultCommandsResponse {
    return {
      list: Array.isArray(object?.list)
        ? object.list.map((e: any) => GetDefaultCommandsResponse_DefaultCommand.fromJSON(e))
        : [],
    };
  },

  toJSON(message: GetDefaultCommandsResponse): unknown {
    const obj: any = {};
    if (message.list) {
      obj.list = message.list.map((e) => e ? GetDefaultCommandsResponse_DefaultCommand.toJSON(e) : undefined);
    } else {
      obj.list = [];
    }
    return obj;
  },

  create(base?: DeepPartial<GetDefaultCommandsResponse>): GetDefaultCommandsResponse {
    return GetDefaultCommandsResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetDefaultCommandsResponse>): GetDefaultCommandsResponse {
    const message = createBaseGetDefaultCommandsResponse();
    message.list = object.list?.map((e) => GetDefaultCommandsResponse_DefaultCommand.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetDefaultCommandsResponse_DefaultCommand(): GetDefaultCommandsResponse_DefaultCommand {
  return {
    name: "",
    description: "",
    visible: false,
    rolesNames: [],
    module: "",
    isReply: false,
    keepResponsesOrder: false,
    aliases: [],
  };
}

export const GetDefaultCommandsResponse_DefaultCommand = {
  encode(message: GetDefaultCommandsResponse_DefaultCommand, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.description !== "") {
      writer.uint32(18).string(message.description);
    }
    if (message.visible === true) {
      writer.uint32(24).bool(message.visible);
    }
    for (const v of message.rolesNames) {
      writer.uint32(34).string(v!);
    }
    if (message.module !== "") {
      writer.uint32(42).string(message.module);
    }
    if (message.isReply === true) {
      writer.uint32(48).bool(message.isReply);
    }
    if (message.keepResponsesOrder === true) {
      writer.uint32(56).bool(message.keepResponsesOrder);
    }
    for (const v of message.aliases) {
      writer.uint32(66).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultCommandsResponse_DefaultCommand {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetDefaultCommandsResponse_DefaultCommand();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.description = reader.string();
          break;
        case 3:
          message.visible = reader.bool();
          break;
        case 4:
          message.rolesNames.push(reader.string());
          break;
        case 5:
          message.module = reader.string();
          break;
        case 6:
          message.isReply = reader.bool();
          break;
        case 7:
          message.keepResponsesOrder = reader.bool();
          break;
        case 8:
          message.aliases.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetDefaultCommandsResponse_DefaultCommand {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      description: isSet(object.description) ? String(object.description) : "",
      visible: isSet(object.visible) ? Boolean(object.visible) : false,
      rolesNames: Array.isArray(object?.rolesNames) ? object.rolesNames.map((e: any) => String(e)) : [],
      module: isSet(object.module) ? String(object.module) : "",
      isReply: isSet(object.isReply) ? Boolean(object.isReply) : false,
      keepResponsesOrder: isSet(object.keepResponsesOrder) ? Boolean(object.keepResponsesOrder) : false,
      aliases: Array.isArray(object?.aliases) ? object.aliases.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: GetDefaultCommandsResponse_DefaultCommand): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.description !== undefined && (obj.description = message.description);
    message.visible !== undefined && (obj.visible = message.visible);
    if (message.rolesNames) {
      obj.rolesNames = message.rolesNames.map((e) => e);
    } else {
      obj.rolesNames = [];
    }
    message.module !== undefined && (obj.module = message.module);
    message.isReply !== undefined && (obj.isReply = message.isReply);
    message.keepResponsesOrder !== undefined && (obj.keepResponsesOrder = message.keepResponsesOrder);
    if (message.aliases) {
      obj.aliases = message.aliases.map((e) => e);
    } else {
      obj.aliases = [];
    }
    return obj;
  },

  create(base?: DeepPartial<GetDefaultCommandsResponse_DefaultCommand>): GetDefaultCommandsResponse_DefaultCommand {
    return GetDefaultCommandsResponse_DefaultCommand.fromPartial(base ?? {});
  },

  fromPartial(
    object: DeepPartial<GetDefaultCommandsResponse_DefaultCommand>,
  ): GetDefaultCommandsResponse_DefaultCommand {
    const message = createBaseGetDefaultCommandsResponse_DefaultCommand();
    message.name = object.name ?? "";
    message.description = object.description ?? "";
    message.visible = object.visible ?? false;
    message.rolesNames = object.rolesNames?.map((e) => e) || [];
    message.module = object.module ?? "";
    message.isReply = object.isReply ?? false;
    message.keepResponsesOrder = object.keepResponsesOrder ?? false;
    message.aliases = object.aliases?.map((e) => e) || [];
    return message;
  },
};

function createBaseParseTextRequestData(): ParseTextRequestData {
  return { sender: undefined, channel: undefined, message: undefined, parseVariables: undefined };
}

export const ParseTextRequestData = {
  encode(message: ParseTextRequestData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sender !== undefined) {
      Sender.encode(message.sender, writer.uint32(10).fork()).ldelim();
    }
    if (message.channel !== undefined) {
      Channel.encode(message.channel, writer.uint32(18).fork()).ldelim();
    }
    if (message.message !== undefined) {
      Message.encode(message.message, writer.uint32(26).fork()).ldelim();
    }
    if (message.parseVariables !== undefined) {
      writer.uint32(32).bool(message.parseVariables);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ParseTextRequestData {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParseTextRequestData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.sender = Sender.decode(reader, reader.uint32());
          break;
        case 2:
          message.channel = Channel.decode(reader, reader.uint32());
          break;
        case 3:
          message.message = Message.decode(reader, reader.uint32());
          break;
        case 4:
          message.parseVariables = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ParseTextRequestData {
    return {
      sender: isSet(object.sender) ? Sender.fromJSON(object.sender) : undefined,
      channel: isSet(object.channel) ? Channel.fromJSON(object.channel) : undefined,
      message: isSet(object.message) ? Message.fromJSON(object.message) : undefined,
      parseVariables: isSet(object.parseVariables) ? Boolean(object.parseVariables) : undefined,
    };
  },

  toJSON(message: ParseTextRequestData): unknown {
    const obj: any = {};
    message.sender !== undefined && (obj.sender = message.sender ? Sender.toJSON(message.sender) : undefined);
    message.channel !== undefined && (obj.channel = message.channel ? Channel.toJSON(message.channel) : undefined);
    message.message !== undefined && (obj.message = message.message ? Message.toJSON(message.message) : undefined);
    message.parseVariables !== undefined && (obj.parseVariables = message.parseVariables);
    return obj;
  },

  create(base?: DeepPartial<ParseTextRequestData>): ParseTextRequestData {
    return ParseTextRequestData.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ParseTextRequestData>): ParseTextRequestData {
    const message = createBaseParseTextRequestData();
    message.sender = (object.sender !== undefined && object.sender !== null)
      ? Sender.fromPartial(object.sender)
      : undefined;
    message.channel = (object.channel !== undefined && object.channel !== null)
      ? Channel.fromPartial(object.channel)
      : undefined;
    message.message = (object.message !== undefined && object.message !== null)
      ? Message.fromPartial(object.message)
      : undefined;
    message.parseVariables = object.parseVariables ?? undefined;
    return message;
  },
};

function createBaseParseTextResponseData(): ParseTextResponseData {
  return { responses: [] };
}

export const ParseTextResponseData = {
  encode(message: ParseTextResponseData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.responses) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ParseTextResponseData {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParseTextResponseData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.responses.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ParseTextResponseData {
    return { responses: Array.isArray(object?.responses) ? object.responses.map((e: any) => String(e)) : [] };
  },

  toJSON(message: ParseTextResponseData): unknown {
    const obj: any = {};
    if (message.responses) {
      obj.responses = message.responses.map((e) => e);
    } else {
      obj.responses = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ParseTextResponseData>): ParseTextResponseData {
    return ParseTextResponseData.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ParseTextResponseData>): ParseTextResponseData {
    const message = createBaseParseTextResponseData();
    message.responses = object.responses?.map((e) => e) || [];
    return message;
  },
};

export type ParserDefinition = typeof ParserDefinition;
export const ParserDefinition = {
  name: "Parser",
  fullName: "parser.Parser",
  methods: {
    processCommand: {
      name: "ProcessCommand",
      requestType: ProcessCommandRequest,
      requestStream: false,
      responseType: ProcessCommandResponse,
      responseStream: false,
      options: {},
    },
    parseTextResponse: {
      name: "ParseTextResponse",
      requestType: ParseTextRequestData,
      requestStream: false,
      responseType: ParseTextResponseData,
      responseStream: false,
      options: {},
    },
    getDefaultCommands: {
      name: "GetDefaultCommands",
      requestType: Empty,
      requestStream: false,
      responseType: GetDefaultCommandsResponse,
      responseStream: false,
      options: {},
    },
    getDefaultVariables: {
      name: "GetDefaultVariables",
      requestType: Empty,
      requestStream: false,
      responseType: GetVariablesResponse,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface ParserServiceImplementation<CallContextExt = {}> {
  processCommand(
    request: ProcessCommandRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ProcessCommandResponse>>;
  parseTextResponse(
    request: ParseTextRequestData,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ParseTextResponseData>>;
  getDefaultCommands(
    request: Empty,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetDefaultCommandsResponse>>;
  getDefaultVariables(
    request: Empty,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetVariablesResponse>>;
}

export interface ParserClient<CallOptionsExt = {}> {
  processCommand(
    request: DeepPartial<ProcessCommandRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ProcessCommandResponse>;
  parseTextResponse(
    request: DeepPartial<ParseTextRequestData>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ParseTextResponseData>;
  getDefaultCommands(
    request: DeepPartial<Empty>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetDefaultCommandsResponse>;
  getDefaultVariables(
    request: DeepPartial<Empty>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetVariablesResponse>;
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
