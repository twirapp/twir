/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal.js";
import { Empty } from "../google/protobuf/empty.js";

export const protobufPackage = "discord";

export enum ChannelType {
  VOICE = 0,
  TEXT = 1,
  UNRECOGNIZED = -1,
}

export function channelTypeFromJSON(object: any): ChannelType {
  switch (object) {
    case 0:
    case "VOICE":
      return ChannelType.VOICE;
    case 1:
    case "TEXT":
      return ChannelType.TEXT;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ChannelType.UNRECOGNIZED;
  }
}

export function channelTypeToJSON(object: ChannelType): string {
  switch (object) {
    case ChannelType.VOICE:
      return "VOICE";
    case ChannelType.TEXT:
      return "TEXT";
    case ChannelType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface GetGuildChannelsRequest {
  guildId: string;
}

export interface GuildChannel {
  id: string;
  name: string;
  type: ChannelType;
  canSendMessages: boolean;
}

export interface GetGuildChannelsResponse {
  channels: GuildChannel[];
}

export interface GetGuildInfoRequest {
  guildId: string;
}

export interface GetGuildInfoResponse {
  id: string;
  name: string;
  icon: string;
  channels: GuildChannel[];
  roles: Role[];
}

export interface LeaveGuildRequest {
  guildId: string;
}

export interface GetGuildRolesRequest {
  guildId: string;
}

export interface Role {
  id: string;
  name: string;
  color: string;
}

export interface GetGuildRolesResponse {
  roles: Role[];
}

function createBaseGetGuildChannelsRequest(): GetGuildChannelsRequest {
  return { guildId: "" };
}

export const GetGuildChannelsRequest = {
  encode(message: GetGuildChannelsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.guildId !== "") {
      writer.uint32(10).string(message.guildId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetGuildChannelsRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetGuildChannelsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.guildId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetGuildChannelsRequest {
    return { guildId: isSet(object.guildId) ? globalThis.String(object.guildId) : "" };
  },

  toJSON(message: GetGuildChannelsRequest): unknown {
    const obj: any = {};
    if (message.guildId !== "") {
      obj.guildId = message.guildId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetGuildChannelsRequest>, I>>(base?: I): GetGuildChannelsRequest {
    return GetGuildChannelsRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetGuildChannelsRequest>, I>>(object: I): GetGuildChannelsRequest {
    const message = createBaseGetGuildChannelsRequest();
    message.guildId = object.guildId ?? "";
    return message;
  },
};

function createBaseGuildChannel(): GuildChannel {
  return { id: "", name: "", type: 0, canSendMessages: false };
}

export const GuildChannel = {
  encode(message: GuildChannel, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.type !== 0) {
      writer.uint32(24).int32(message.type);
    }
    if (message.canSendMessages !== false) {
      writer.uint32(32).bool(message.canSendMessages);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GuildChannel {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGuildChannel();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.id = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.name = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.type = reader.int32() as any;
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.canSendMessages = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GuildChannel {
    return {
      id: isSet(object.id) ? globalThis.String(object.id) : "",
      name: isSet(object.name) ? globalThis.String(object.name) : "",
      type: isSet(object.type) ? channelTypeFromJSON(object.type) : 0,
      canSendMessages: isSet(object.canSendMessages) ? globalThis.Boolean(object.canSendMessages) : false,
    };
  },

  toJSON(message: GuildChannel): unknown {
    const obj: any = {};
    if (message.id !== "") {
      obj.id = message.id;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.type !== 0) {
      obj.type = channelTypeToJSON(message.type);
    }
    if (message.canSendMessages !== false) {
      obj.canSendMessages = message.canSendMessages;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GuildChannel>, I>>(base?: I): GuildChannel {
    return GuildChannel.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GuildChannel>, I>>(object: I): GuildChannel {
    const message = createBaseGuildChannel();
    message.id = object.id ?? "";
    message.name = object.name ?? "";
    message.type = object.type ?? 0;
    message.canSendMessages = object.canSendMessages ?? false;
    return message;
  },
};

function createBaseGetGuildChannelsResponse(): GetGuildChannelsResponse {
  return { channels: [] };
}

export const GetGuildChannelsResponse = {
  encode(message: GetGuildChannelsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.channels) {
      GuildChannel.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetGuildChannelsResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetGuildChannelsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channels.push(GuildChannel.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetGuildChannelsResponse {
    return {
      channels: globalThis.Array.isArray(object?.channels)
        ? object.channels.map((e: any) => GuildChannel.fromJSON(e))
        : [],
    };
  },

  toJSON(message: GetGuildChannelsResponse): unknown {
    const obj: any = {};
    if (message.channels?.length) {
      obj.channels = message.channels.map((e) => GuildChannel.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetGuildChannelsResponse>, I>>(base?: I): GetGuildChannelsResponse {
    return GetGuildChannelsResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetGuildChannelsResponse>, I>>(object: I): GetGuildChannelsResponse {
    const message = createBaseGetGuildChannelsResponse();
    message.channels = object.channels?.map((e) => GuildChannel.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetGuildInfoRequest(): GetGuildInfoRequest {
  return { guildId: "" };
}

export const GetGuildInfoRequest = {
  encode(message: GetGuildInfoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.guildId !== "") {
      writer.uint32(10).string(message.guildId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetGuildInfoRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetGuildInfoRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.guildId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetGuildInfoRequest {
    return { guildId: isSet(object.guildId) ? globalThis.String(object.guildId) : "" };
  },

  toJSON(message: GetGuildInfoRequest): unknown {
    const obj: any = {};
    if (message.guildId !== "") {
      obj.guildId = message.guildId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetGuildInfoRequest>, I>>(base?: I): GetGuildInfoRequest {
    return GetGuildInfoRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetGuildInfoRequest>, I>>(object: I): GetGuildInfoRequest {
    const message = createBaseGetGuildInfoRequest();
    message.guildId = object.guildId ?? "";
    return message;
  },
};

function createBaseGetGuildInfoResponse(): GetGuildInfoResponse {
  return { id: "", name: "", icon: "", channels: [], roles: [] };
}

export const GetGuildInfoResponse = {
  encode(message: GetGuildInfoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.icon !== "") {
      writer.uint32(26).string(message.icon);
    }
    for (const v of message.channels) {
      GuildChannel.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.roles) {
      Role.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetGuildInfoResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetGuildInfoResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.id = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.name = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.icon = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.channels.push(GuildChannel.decode(reader, reader.uint32()));
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.roles.push(Role.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetGuildInfoResponse {
    return {
      id: isSet(object.id) ? globalThis.String(object.id) : "",
      name: isSet(object.name) ? globalThis.String(object.name) : "",
      icon: isSet(object.icon) ? globalThis.String(object.icon) : "",
      channels: globalThis.Array.isArray(object?.channels)
        ? object.channels.map((e: any) => GuildChannel.fromJSON(e))
        : [],
      roles: globalThis.Array.isArray(object?.roles) ? object.roles.map((e: any) => Role.fromJSON(e)) : [],
    };
  },

  toJSON(message: GetGuildInfoResponse): unknown {
    const obj: any = {};
    if (message.id !== "") {
      obj.id = message.id;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.icon !== "") {
      obj.icon = message.icon;
    }
    if (message.channels?.length) {
      obj.channels = message.channels.map((e) => GuildChannel.toJSON(e));
    }
    if (message.roles?.length) {
      obj.roles = message.roles.map((e) => Role.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetGuildInfoResponse>, I>>(base?: I): GetGuildInfoResponse {
    return GetGuildInfoResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetGuildInfoResponse>, I>>(object: I): GetGuildInfoResponse {
    const message = createBaseGetGuildInfoResponse();
    message.id = object.id ?? "";
    message.name = object.name ?? "";
    message.icon = object.icon ?? "";
    message.channels = object.channels?.map((e) => GuildChannel.fromPartial(e)) || [];
    message.roles = object.roles?.map((e) => Role.fromPartial(e)) || [];
    return message;
  },
};

function createBaseLeaveGuildRequest(): LeaveGuildRequest {
  return { guildId: "" };
}

export const LeaveGuildRequest = {
  encode(message: LeaveGuildRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.guildId !== "") {
      writer.uint32(10).string(message.guildId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LeaveGuildRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLeaveGuildRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.guildId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): LeaveGuildRequest {
    return { guildId: isSet(object.guildId) ? globalThis.String(object.guildId) : "" };
  },

  toJSON(message: LeaveGuildRequest): unknown {
    const obj: any = {};
    if (message.guildId !== "") {
      obj.guildId = message.guildId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<LeaveGuildRequest>, I>>(base?: I): LeaveGuildRequest {
    return LeaveGuildRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<LeaveGuildRequest>, I>>(object: I): LeaveGuildRequest {
    const message = createBaseLeaveGuildRequest();
    message.guildId = object.guildId ?? "";
    return message;
  },
};

function createBaseGetGuildRolesRequest(): GetGuildRolesRequest {
  return { guildId: "" };
}

export const GetGuildRolesRequest = {
  encode(message: GetGuildRolesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.guildId !== "") {
      writer.uint32(10).string(message.guildId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetGuildRolesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetGuildRolesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.guildId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetGuildRolesRequest {
    return { guildId: isSet(object.guildId) ? globalThis.String(object.guildId) : "" };
  },

  toJSON(message: GetGuildRolesRequest): unknown {
    const obj: any = {};
    if (message.guildId !== "") {
      obj.guildId = message.guildId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetGuildRolesRequest>, I>>(base?: I): GetGuildRolesRequest {
    return GetGuildRolesRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetGuildRolesRequest>, I>>(object: I): GetGuildRolesRequest {
    const message = createBaseGetGuildRolesRequest();
    message.guildId = object.guildId ?? "";
    return message;
  },
};

function createBaseRole(): Role {
  return { id: "", name: "", color: "" };
}

export const Role = {
  encode(message: Role, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.color !== "") {
      writer.uint32(26).string(message.color);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Role {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRole();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.id = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.name = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.color = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Role {
    return {
      id: isSet(object.id) ? globalThis.String(object.id) : "",
      name: isSet(object.name) ? globalThis.String(object.name) : "",
      color: isSet(object.color) ? globalThis.String(object.color) : "",
    };
  },

  toJSON(message: Role): unknown {
    const obj: any = {};
    if (message.id !== "") {
      obj.id = message.id;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.color !== "") {
      obj.color = message.color;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Role>, I>>(base?: I): Role {
    return Role.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Role>, I>>(object: I): Role {
    const message = createBaseRole();
    message.id = object.id ?? "";
    message.name = object.name ?? "";
    message.color = object.color ?? "";
    return message;
  },
};

function createBaseGetGuildRolesResponse(): GetGuildRolesResponse {
  return { roles: [] };
}

export const GetGuildRolesResponse = {
  encode(message: GetGuildRolesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.roles) {
      Role.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetGuildRolesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetGuildRolesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.roles.push(Role.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetGuildRolesResponse {
    return { roles: globalThis.Array.isArray(object?.roles) ? object.roles.map((e: any) => Role.fromJSON(e)) : [] };
  },

  toJSON(message: GetGuildRolesResponse): unknown {
    const obj: any = {};
    if (message.roles?.length) {
      obj.roles = message.roles.map((e) => Role.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetGuildRolesResponse>, I>>(base?: I): GetGuildRolesResponse {
    return GetGuildRolesResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetGuildRolesResponse>, I>>(object: I): GetGuildRolesResponse {
    const message = createBaseGetGuildRolesResponse();
    message.roles = object.roles?.map((e) => Role.fromPartial(e)) || [];
    return message;
  },
};

export type DiscordDefinition = typeof DiscordDefinition;
export const DiscordDefinition = {
  name: "Discord",
  fullName: "discord.Discord",
  methods: {
    getGuildChannels: {
      name: "GetGuildChannels",
      requestType: GetGuildChannelsRequest,
      requestStream: false,
      responseType: GetGuildChannelsResponse,
      responseStream: false,
      options: {},
    },
    getGuildInfo: {
      name: "GetGuildInfo",
      requestType: GetGuildInfoRequest,
      requestStream: false,
      responseType: GetGuildInfoResponse,
      responseStream: false,
      options: {},
    },
    leaveGuild: {
      name: "LeaveGuild",
      requestType: LeaveGuildRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    getGuildRoles: {
      name: "GetGuildRoles",
      requestType: GetGuildRolesRequest,
      requestStream: false,
      responseType: GetGuildRolesResponse,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface DiscordServiceImplementation<CallContextExt = {}> {
  getGuildChannels(
    request: GetGuildChannelsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetGuildChannelsResponse>>;
  getGuildInfo(
    request: GetGuildInfoRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetGuildInfoResponse>>;
  leaveGuild(request: LeaveGuildRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  getGuildRoles(
    request: GetGuildRolesRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetGuildRolesResponse>>;
}

export interface DiscordClient<CallOptionsExt = {}> {
  getGuildChannels(
    request: DeepPartial<GetGuildChannelsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetGuildChannelsResponse>;
  getGuildInfo(
    request: DeepPartial<GetGuildInfoRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetGuildInfoResponse>;
  leaveGuild(request: DeepPartial<LeaveGuildRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  getGuildRoles(
    request: DeepPartial<GetGuildRolesRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetGuildRolesResponse>;
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

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
