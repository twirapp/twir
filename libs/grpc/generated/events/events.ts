/* eslint-disable */
import Long from "long";
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "events";

export interface BaseInfo {
  channelId: string;
}

export interface FollowMessage {
  baseInfo: BaseInfo | undefined;
  userName: string;
  userDisplayName: string;
  userId: string;
}

export interface SubscribeMessage {
  baseInfo: BaseInfo | undefined;
  userName: string;
  userDisplayName: string;
  level: string;
  userId: string;
}

export interface SubGiftMessage {
  baseInfo: BaseInfo | undefined;
  senderUserName: string;
  senderDisplayName: string;
  targetUserName: string;
  targetDisplayName: string;
  level: string;
  senderUserId: string;
}

export interface ReSubscribeMessage {
  baseInfo: BaseInfo | undefined;
  userName: string;
  userDisplayName: string;
  months: number;
  streak: number;
  isPrime: boolean;
  message: string;
  level: string;
  userId: string;
}

export interface RedemptionCreatedMessage {
  baseInfo: BaseInfo | undefined;
  userName: string;
  userDisplayName: string;
  id: string;
  rewardName: string;
  rewardCost: string;
  input?: string | undefined;
  userId: string;
}

export interface CommandUsedMessage {
  baseInfo: BaseInfo | undefined;
  commandId: string;
  commandName: string;
  userName: string;
  userDisplayName: string;
  commandInput: string;
  userId: string;
}

export interface FirstUserMessageMessage {
  baseInfo: BaseInfo | undefined;
  userId: string;
  userName: string;
  userDisplayName: string;
}

export interface RaidedMessage {
  baseInfo: BaseInfo | undefined;
  userName: string;
  userDisplayName: string;
  viewers: number;
  userId: string;
}

export interface TitleOrCategoryChangedMessage {
  baseInfo: BaseInfo | undefined;
  oldTitle: string;
  newTitle: string;
  oldCategory: string;
  newCategory: string;
}

export interface StreamOnlineMessage {
  baseInfo: BaseInfo | undefined;
  title: string;
  category: string;
}

export interface StreamOfflineMessage {
  baseInfo: BaseInfo | undefined;
}

export interface ChatClearMessage {
  baseInfo: BaseInfo | undefined;
}

export interface DonateMessage {
  baseInfo: BaseInfo | undefined;
  userName: string;
  amount: string;
  currency: string;
  message: string;
}

export interface KeywordMatchedMessage {
  baseInfo: BaseInfo | undefined;
  keywordId: string;
  keywordName: string;
  keywordResponse: string;
  userId: string;
  userName: string;
  userDisplayName: string;
}

export interface GreetingSendedMessage {
  baseInfo: BaseInfo | undefined;
  userId: string;
  userName: string;
  userDisplayName: string;
  greetingText: string;
}

function createBaseBaseInfo(): BaseInfo {
  return { channelId: "" };
}

export const BaseInfo = {
  encode(message: BaseInfo, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BaseInfo {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBaseInfo();
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

  fromJSON(object: any): BaseInfo {
    return { channelId: isSet(object.channelId) ? String(object.channelId) : "" };
  },

  toJSON(message: BaseInfo): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    return obj;
  },

  create(base?: DeepPartial<BaseInfo>): BaseInfo {
    return BaseInfo.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<BaseInfo>): BaseInfo {
    const message = createBaseBaseInfo();
    message.channelId = object.channelId ?? "";
    return message;
  },
};

function createBaseFollowMessage(): FollowMessage {
  return { baseInfo: undefined, userName: "", userDisplayName: "", userId: "" };
}

export const FollowMessage = {
  encode(message: FollowMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.baseInfo !== undefined) {
      BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.userName !== "") {
      writer.uint32(18).string(message.userName);
    }
    if (message.userDisplayName !== "") {
      writer.uint32(26).string(message.userDisplayName);
    }
    if (message.userId !== "") {
      writer.uint32(34).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FollowMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFollowMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseInfo = BaseInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.userName = reader.string();
          break;
        case 3:
          message.userDisplayName = reader.string();
          break;
        case 4:
          message.userId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FollowMessage {
    return {
      baseInfo: isSet(object.baseInfo) ? BaseInfo.fromJSON(object.baseInfo) : undefined,
      userName: isSet(object.userName) ? String(object.userName) : "",
      userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
      userId: isSet(object.userId) ? String(object.userId) : "",
    };
  },

  toJSON(message: FollowMessage): unknown {
    const obj: any = {};
    message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? BaseInfo.toJSON(message.baseInfo) : undefined);
    message.userName !== undefined && (obj.userName = message.userName);
    message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<FollowMessage>): FollowMessage {
    return FollowMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<FollowMessage>): FollowMessage {
    const message = createBaseFollowMessage();
    message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
      ? BaseInfo.fromPartial(object.baseInfo)
      : undefined;
    message.userName = object.userName ?? "";
    message.userDisplayName = object.userDisplayName ?? "";
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseSubscribeMessage(): SubscribeMessage {
  return { baseInfo: undefined, userName: "", userDisplayName: "", level: "", userId: "" };
}

export const SubscribeMessage = {
  encode(message: SubscribeMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.baseInfo !== undefined) {
      BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.userName !== "") {
      writer.uint32(18).string(message.userName);
    }
    if (message.userDisplayName !== "") {
      writer.uint32(26).string(message.userDisplayName);
    }
    if (message.level !== "") {
      writer.uint32(34).string(message.level);
    }
    if (message.userId !== "") {
      writer.uint32(42).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubscribeMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubscribeMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseInfo = BaseInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.userName = reader.string();
          break;
        case 3:
          message.userDisplayName = reader.string();
          break;
        case 4:
          message.level = reader.string();
          break;
        case 5:
          message.userId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SubscribeMessage {
    return {
      baseInfo: isSet(object.baseInfo) ? BaseInfo.fromJSON(object.baseInfo) : undefined,
      userName: isSet(object.userName) ? String(object.userName) : "",
      userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
      level: isSet(object.level) ? String(object.level) : "",
      userId: isSet(object.userId) ? String(object.userId) : "",
    };
  },

  toJSON(message: SubscribeMessage): unknown {
    const obj: any = {};
    message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? BaseInfo.toJSON(message.baseInfo) : undefined);
    message.userName !== undefined && (obj.userName = message.userName);
    message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
    message.level !== undefined && (obj.level = message.level);
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<SubscribeMessage>): SubscribeMessage {
    return SubscribeMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SubscribeMessage>): SubscribeMessage {
    const message = createBaseSubscribeMessage();
    message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
      ? BaseInfo.fromPartial(object.baseInfo)
      : undefined;
    message.userName = object.userName ?? "";
    message.userDisplayName = object.userDisplayName ?? "";
    message.level = object.level ?? "";
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseSubGiftMessage(): SubGiftMessage {
  return {
    baseInfo: undefined,
    senderUserName: "",
    senderDisplayName: "",
    targetUserName: "",
    targetDisplayName: "",
    level: "",
    senderUserId: "",
  };
}

export const SubGiftMessage = {
  encode(message: SubGiftMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.baseInfo !== undefined) {
      BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.senderUserName !== "") {
      writer.uint32(18).string(message.senderUserName);
    }
    if (message.senderDisplayName !== "") {
      writer.uint32(26).string(message.senderDisplayName);
    }
    if (message.targetUserName !== "") {
      writer.uint32(34).string(message.targetUserName);
    }
    if (message.targetDisplayName !== "") {
      writer.uint32(42).string(message.targetDisplayName);
    }
    if (message.level !== "") {
      writer.uint32(50).string(message.level);
    }
    if (message.senderUserId !== "") {
      writer.uint32(58).string(message.senderUserId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubGiftMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubGiftMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseInfo = BaseInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.senderUserName = reader.string();
          break;
        case 3:
          message.senderDisplayName = reader.string();
          break;
        case 4:
          message.targetUserName = reader.string();
          break;
        case 5:
          message.targetDisplayName = reader.string();
          break;
        case 6:
          message.level = reader.string();
          break;
        case 7:
          message.senderUserId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SubGiftMessage {
    return {
      baseInfo: isSet(object.baseInfo) ? BaseInfo.fromJSON(object.baseInfo) : undefined,
      senderUserName: isSet(object.senderUserName) ? String(object.senderUserName) : "",
      senderDisplayName: isSet(object.senderDisplayName) ? String(object.senderDisplayName) : "",
      targetUserName: isSet(object.targetUserName) ? String(object.targetUserName) : "",
      targetDisplayName: isSet(object.targetDisplayName) ? String(object.targetDisplayName) : "",
      level: isSet(object.level) ? String(object.level) : "",
      senderUserId: isSet(object.senderUserId) ? String(object.senderUserId) : "",
    };
  },

  toJSON(message: SubGiftMessage): unknown {
    const obj: any = {};
    message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? BaseInfo.toJSON(message.baseInfo) : undefined);
    message.senderUserName !== undefined && (obj.senderUserName = message.senderUserName);
    message.senderDisplayName !== undefined && (obj.senderDisplayName = message.senderDisplayName);
    message.targetUserName !== undefined && (obj.targetUserName = message.targetUserName);
    message.targetDisplayName !== undefined && (obj.targetDisplayName = message.targetDisplayName);
    message.level !== undefined && (obj.level = message.level);
    message.senderUserId !== undefined && (obj.senderUserId = message.senderUserId);
    return obj;
  },

  create(base?: DeepPartial<SubGiftMessage>): SubGiftMessage {
    return SubGiftMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SubGiftMessage>): SubGiftMessage {
    const message = createBaseSubGiftMessage();
    message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
      ? BaseInfo.fromPartial(object.baseInfo)
      : undefined;
    message.senderUserName = object.senderUserName ?? "";
    message.senderDisplayName = object.senderDisplayName ?? "";
    message.targetUserName = object.targetUserName ?? "";
    message.targetDisplayName = object.targetDisplayName ?? "";
    message.level = object.level ?? "";
    message.senderUserId = object.senderUserId ?? "";
    return message;
  },
};

function createBaseReSubscribeMessage(): ReSubscribeMessage {
  return {
    baseInfo: undefined,
    userName: "",
    userDisplayName: "",
    months: 0,
    streak: 0,
    isPrime: false,
    message: "",
    level: "",
    userId: "",
  };
}

export const ReSubscribeMessage = {
  encode(message: ReSubscribeMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.baseInfo !== undefined) {
      BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.userName !== "") {
      writer.uint32(18).string(message.userName);
    }
    if (message.userDisplayName !== "") {
      writer.uint32(26).string(message.userDisplayName);
    }
    if (message.months !== 0) {
      writer.uint32(32).int64(message.months);
    }
    if (message.streak !== 0) {
      writer.uint32(40).int64(message.streak);
    }
    if (message.isPrime === true) {
      writer.uint32(48).bool(message.isPrime);
    }
    if (message.message !== "") {
      writer.uint32(58).string(message.message);
    }
    if (message.level !== "") {
      writer.uint32(66).string(message.level);
    }
    if (message.userId !== "") {
      writer.uint32(74).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReSubscribeMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReSubscribeMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseInfo = BaseInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.userName = reader.string();
          break;
        case 3:
          message.userDisplayName = reader.string();
          break;
        case 4:
          message.months = longToNumber(reader.int64() as Long);
          break;
        case 5:
          message.streak = longToNumber(reader.int64() as Long);
          break;
        case 6:
          message.isPrime = reader.bool();
          break;
        case 7:
          message.message = reader.string();
          break;
        case 8:
          message.level = reader.string();
          break;
        case 9:
          message.userId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ReSubscribeMessage {
    return {
      baseInfo: isSet(object.baseInfo) ? BaseInfo.fromJSON(object.baseInfo) : undefined,
      userName: isSet(object.userName) ? String(object.userName) : "",
      userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
      months: isSet(object.months) ? Number(object.months) : 0,
      streak: isSet(object.streak) ? Number(object.streak) : 0,
      isPrime: isSet(object.isPrime) ? Boolean(object.isPrime) : false,
      message: isSet(object.message) ? String(object.message) : "",
      level: isSet(object.level) ? String(object.level) : "",
      userId: isSet(object.userId) ? String(object.userId) : "",
    };
  },

  toJSON(message: ReSubscribeMessage): unknown {
    const obj: any = {};
    message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? BaseInfo.toJSON(message.baseInfo) : undefined);
    message.userName !== undefined && (obj.userName = message.userName);
    message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
    message.months !== undefined && (obj.months = Math.round(message.months));
    message.streak !== undefined && (obj.streak = Math.round(message.streak));
    message.isPrime !== undefined && (obj.isPrime = message.isPrime);
    message.message !== undefined && (obj.message = message.message);
    message.level !== undefined && (obj.level = message.level);
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<ReSubscribeMessage>): ReSubscribeMessage {
    return ReSubscribeMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ReSubscribeMessage>): ReSubscribeMessage {
    const message = createBaseReSubscribeMessage();
    message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
      ? BaseInfo.fromPartial(object.baseInfo)
      : undefined;
    message.userName = object.userName ?? "";
    message.userDisplayName = object.userDisplayName ?? "";
    message.months = object.months ?? 0;
    message.streak = object.streak ?? 0;
    message.isPrime = object.isPrime ?? false;
    message.message = object.message ?? "";
    message.level = object.level ?? "";
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseRedemptionCreatedMessage(): RedemptionCreatedMessage {
  return {
    baseInfo: undefined,
    userName: "",
    userDisplayName: "",
    id: "",
    rewardName: "",
    rewardCost: "",
    input: undefined,
    userId: "",
  };
}

export const RedemptionCreatedMessage = {
  encode(message: RedemptionCreatedMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.baseInfo !== undefined) {
      BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.userName !== "") {
      writer.uint32(18).string(message.userName);
    }
    if (message.userDisplayName !== "") {
      writer.uint32(26).string(message.userDisplayName);
    }
    if (message.id !== "") {
      writer.uint32(34).string(message.id);
    }
    if (message.rewardName !== "") {
      writer.uint32(42).string(message.rewardName);
    }
    if (message.rewardCost !== "") {
      writer.uint32(50).string(message.rewardCost);
    }
    if (message.input !== undefined) {
      writer.uint32(58).string(message.input);
    }
    if (message.userId !== "") {
      writer.uint32(66).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RedemptionCreatedMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRedemptionCreatedMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseInfo = BaseInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.userName = reader.string();
          break;
        case 3:
          message.userDisplayName = reader.string();
          break;
        case 4:
          message.id = reader.string();
          break;
        case 5:
          message.rewardName = reader.string();
          break;
        case 6:
          message.rewardCost = reader.string();
          break;
        case 7:
          message.input = reader.string();
          break;
        case 8:
          message.userId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RedemptionCreatedMessage {
    return {
      baseInfo: isSet(object.baseInfo) ? BaseInfo.fromJSON(object.baseInfo) : undefined,
      userName: isSet(object.userName) ? String(object.userName) : "",
      userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
      id: isSet(object.id) ? String(object.id) : "",
      rewardName: isSet(object.rewardName) ? String(object.rewardName) : "",
      rewardCost: isSet(object.rewardCost) ? String(object.rewardCost) : "",
      input: isSet(object.input) ? String(object.input) : undefined,
      userId: isSet(object.userId) ? String(object.userId) : "",
    };
  },

  toJSON(message: RedemptionCreatedMessage): unknown {
    const obj: any = {};
    message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? BaseInfo.toJSON(message.baseInfo) : undefined);
    message.userName !== undefined && (obj.userName = message.userName);
    message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
    message.id !== undefined && (obj.id = message.id);
    message.rewardName !== undefined && (obj.rewardName = message.rewardName);
    message.rewardCost !== undefined && (obj.rewardCost = message.rewardCost);
    message.input !== undefined && (obj.input = message.input);
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<RedemptionCreatedMessage>): RedemptionCreatedMessage {
    return RedemptionCreatedMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RedemptionCreatedMessage>): RedemptionCreatedMessage {
    const message = createBaseRedemptionCreatedMessage();
    message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
      ? BaseInfo.fromPartial(object.baseInfo)
      : undefined;
    message.userName = object.userName ?? "";
    message.userDisplayName = object.userDisplayName ?? "";
    message.id = object.id ?? "";
    message.rewardName = object.rewardName ?? "";
    message.rewardCost = object.rewardCost ?? "";
    message.input = object.input ?? undefined;
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseCommandUsedMessage(): CommandUsedMessage {
  return {
    baseInfo: undefined,
    commandId: "",
    commandName: "",
    userName: "",
    userDisplayName: "",
    commandInput: "",
    userId: "",
  };
}

export const CommandUsedMessage = {
  encode(message: CommandUsedMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.baseInfo !== undefined) {
      BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.commandId !== "") {
      writer.uint32(18).string(message.commandId);
    }
    if (message.commandName !== "") {
      writer.uint32(26).string(message.commandName);
    }
    if (message.userName !== "") {
      writer.uint32(34).string(message.userName);
    }
    if (message.userDisplayName !== "") {
      writer.uint32(42).string(message.userDisplayName);
    }
    if (message.commandInput !== "") {
      writer.uint32(50).string(message.commandInput);
    }
    if (message.userId !== "") {
      writer.uint32(58).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CommandUsedMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCommandUsedMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseInfo = BaseInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.commandId = reader.string();
          break;
        case 3:
          message.commandName = reader.string();
          break;
        case 4:
          message.userName = reader.string();
          break;
        case 5:
          message.userDisplayName = reader.string();
          break;
        case 6:
          message.commandInput = reader.string();
          break;
        case 7:
          message.userId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CommandUsedMessage {
    return {
      baseInfo: isSet(object.baseInfo) ? BaseInfo.fromJSON(object.baseInfo) : undefined,
      commandId: isSet(object.commandId) ? String(object.commandId) : "",
      commandName: isSet(object.commandName) ? String(object.commandName) : "",
      userName: isSet(object.userName) ? String(object.userName) : "",
      userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
      commandInput: isSet(object.commandInput) ? String(object.commandInput) : "",
      userId: isSet(object.userId) ? String(object.userId) : "",
    };
  },

  toJSON(message: CommandUsedMessage): unknown {
    const obj: any = {};
    message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? BaseInfo.toJSON(message.baseInfo) : undefined);
    message.commandId !== undefined && (obj.commandId = message.commandId);
    message.commandName !== undefined && (obj.commandName = message.commandName);
    message.userName !== undefined && (obj.userName = message.userName);
    message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
    message.commandInput !== undefined && (obj.commandInput = message.commandInput);
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<CommandUsedMessage>): CommandUsedMessage {
    return CommandUsedMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CommandUsedMessage>): CommandUsedMessage {
    const message = createBaseCommandUsedMessage();
    message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
      ? BaseInfo.fromPartial(object.baseInfo)
      : undefined;
    message.commandId = object.commandId ?? "";
    message.commandName = object.commandName ?? "";
    message.userName = object.userName ?? "";
    message.userDisplayName = object.userDisplayName ?? "";
    message.commandInput = object.commandInput ?? "";
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseFirstUserMessageMessage(): FirstUserMessageMessage {
  return { baseInfo: undefined, userId: "", userName: "", userDisplayName: "" };
}

export const FirstUserMessageMessage = {
  encode(message: FirstUserMessageMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.baseInfo !== undefined) {
      BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.userId !== "") {
      writer.uint32(18).string(message.userId);
    }
    if (message.userName !== "") {
      writer.uint32(26).string(message.userName);
    }
    if (message.userDisplayName !== "") {
      writer.uint32(34).string(message.userDisplayName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FirstUserMessageMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFirstUserMessageMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseInfo = BaseInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.userId = reader.string();
          break;
        case 3:
          message.userName = reader.string();
          break;
        case 4:
          message.userDisplayName = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FirstUserMessageMessage {
    return {
      baseInfo: isSet(object.baseInfo) ? BaseInfo.fromJSON(object.baseInfo) : undefined,
      userId: isSet(object.userId) ? String(object.userId) : "",
      userName: isSet(object.userName) ? String(object.userName) : "",
      userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
    };
  },

  toJSON(message: FirstUserMessageMessage): unknown {
    const obj: any = {};
    message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? BaseInfo.toJSON(message.baseInfo) : undefined);
    message.userId !== undefined && (obj.userId = message.userId);
    message.userName !== undefined && (obj.userName = message.userName);
    message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
    return obj;
  },

  create(base?: DeepPartial<FirstUserMessageMessage>): FirstUserMessageMessage {
    return FirstUserMessageMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<FirstUserMessageMessage>): FirstUserMessageMessage {
    const message = createBaseFirstUserMessageMessage();
    message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
      ? BaseInfo.fromPartial(object.baseInfo)
      : undefined;
    message.userId = object.userId ?? "";
    message.userName = object.userName ?? "";
    message.userDisplayName = object.userDisplayName ?? "";
    return message;
  },
};

function createBaseRaidedMessage(): RaidedMessage {
  return { baseInfo: undefined, userName: "", userDisplayName: "", viewers: 0, userId: "" };
}

export const RaidedMessage = {
  encode(message: RaidedMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.baseInfo !== undefined) {
      BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.userName !== "") {
      writer.uint32(18).string(message.userName);
    }
    if (message.userDisplayName !== "") {
      writer.uint32(26).string(message.userDisplayName);
    }
    if (message.viewers !== 0) {
      writer.uint32(32).int64(message.viewers);
    }
    if (message.userId !== "") {
      writer.uint32(42).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RaidedMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRaidedMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseInfo = BaseInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.userName = reader.string();
          break;
        case 3:
          message.userDisplayName = reader.string();
          break;
        case 4:
          message.viewers = longToNumber(reader.int64() as Long);
          break;
        case 5:
          message.userId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RaidedMessage {
    return {
      baseInfo: isSet(object.baseInfo) ? BaseInfo.fromJSON(object.baseInfo) : undefined,
      userName: isSet(object.userName) ? String(object.userName) : "",
      userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
      viewers: isSet(object.viewers) ? Number(object.viewers) : 0,
      userId: isSet(object.userId) ? String(object.userId) : "",
    };
  },

  toJSON(message: RaidedMessage): unknown {
    const obj: any = {};
    message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? BaseInfo.toJSON(message.baseInfo) : undefined);
    message.userName !== undefined && (obj.userName = message.userName);
    message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
    message.viewers !== undefined && (obj.viewers = Math.round(message.viewers));
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<RaidedMessage>): RaidedMessage {
    return RaidedMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RaidedMessage>): RaidedMessage {
    const message = createBaseRaidedMessage();
    message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
      ? BaseInfo.fromPartial(object.baseInfo)
      : undefined;
    message.userName = object.userName ?? "";
    message.userDisplayName = object.userDisplayName ?? "";
    message.viewers = object.viewers ?? 0;
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseTitleOrCategoryChangedMessage(): TitleOrCategoryChangedMessage {
  return { baseInfo: undefined, oldTitle: "", newTitle: "", oldCategory: "", newCategory: "" };
}

export const TitleOrCategoryChangedMessage = {
  encode(message: TitleOrCategoryChangedMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.baseInfo !== undefined) {
      BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.oldTitle !== "") {
      writer.uint32(18).string(message.oldTitle);
    }
    if (message.newTitle !== "") {
      writer.uint32(26).string(message.newTitle);
    }
    if (message.oldCategory !== "") {
      writer.uint32(34).string(message.oldCategory);
    }
    if (message.newCategory !== "") {
      writer.uint32(42).string(message.newCategory);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TitleOrCategoryChangedMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTitleOrCategoryChangedMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseInfo = BaseInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.oldTitle = reader.string();
          break;
        case 3:
          message.newTitle = reader.string();
          break;
        case 4:
          message.oldCategory = reader.string();
          break;
        case 5:
          message.newCategory = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): TitleOrCategoryChangedMessage {
    return {
      baseInfo: isSet(object.baseInfo) ? BaseInfo.fromJSON(object.baseInfo) : undefined,
      oldTitle: isSet(object.oldTitle) ? String(object.oldTitle) : "",
      newTitle: isSet(object.newTitle) ? String(object.newTitle) : "",
      oldCategory: isSet(object.oldCategory) ? String(object.oldCategory) : "",
      newCategory: isSet(object.newCategory) ? String(object.newCategory) : "",
    };
  },

  toJSON(message: TitleOrCategoryChangedMessage): unknown {
    const obj: any = {};
    message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? BaseInfo.toJSON(message.baseInfo) : undefined);
    message.oldTitle !== undefined && (obj.oldTitle = message.oldTitle);
    message.newTitle !== undefined && (obj.newTitle = message.newTitle);
    message.oldCategory !== undefined && (obj.oldCategory = message.oldCategory);
    message.newCategory !== undefined && (obj.newCategory = message.newCategory);
    return obj;
  },

  create(base?: DeepPartial<TitleOrCategoryChangedMessage>): TitleOrCategoryChangedMessage {
    return TitleOrCategoryChangedMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<TitleOrCategoryChangedMessage>): TitleOrCategoryChangedMessage {
    const message = createBaseTitleOrCategoryChangedMessage();
    message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
      ? BaseInfo.fromPartial(object.baseInfo)
      : undefined;
    message.oldTitle = object.oldTitle ?? "";
    message.newTitle = object.newTitle ?? "";
    message.oldCategory = object.oldCategory ?? "";
    message.newCategory = object.newCategory ?? "";
    return message;
  },
};

function createBaseStreamOnlineMessage(): StreamOnlineMessage {
  return { baseInfo: undefined, title: "", category: "" };
}

export const StreamOnlineMessage = {
  encode(message: StreamOnlineMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.baseInfo !== undefined) {
      BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.title !== "") {
      writer.uint32(18).string(message.title);
    }
    if (message.category !== "") {
      writer.uint32(26).string(message.category);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StreamOnlineMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStreamOnlineMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseInfo = BaseInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.title = reader.string();
          break;
        case 3:
          message.category = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): StreamOnlineMessage {
    return {
      baseInfo: isSet(object.baseInfo) ? BaseInfo.fromJSON(object.baseInfo) : undefined,
      title: isSet(object.title) ? String(object.title) : "",
      category: isSet(object.category) ? String(object.category) : "",
    };
  },

  toJSON(message: StreamOnlineMessage): unknown {
    const obj: any = {};
    message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? BaseInfo.toJSON(message.baseInfo) : undefined);
    message.title !== undefined && (obj.title = message.title);
    message.category !== undefined && (obj.category = message.category);
    return obj;
  },

  create(base?: DeepPartial<StreamOnlineMessage>): StreamOnlineMessage {
    return StreamOnlineMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<StreamOnlineMessage>): StreamOnlineMessage {
    const message = createBaseStreamOnlineMessage();
    message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
      ? BaseInfo.fromPartial(object.baseInfo)
      : undefined;
    message.title = object.title ?? "";
    message.category = object.category ?? "";
    return message;
  },
};

function createBaseStreamOfflineMessage(): StreamOfflineMessage {
  return { baseInfo: undefined };
}

export const StreamOfflineMessage = {
  encode(message: StreamOfflineMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.baseInfo !== undefined) {
      BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StreamOfflineMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStreamOfflineMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseInfo = BaseInfo.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): StreamOfflineMessage {
    return { baseInfo: isSet(object.baseInfo) ? BaseInfo.fromJSON(object.baseInfo) : undefined };
  },

  toJSON(message: StreamOfflineMessage): unknown {
    const obj: any = {};
    message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? BaseInfo.toJSON(message.baseInfo) : undefined);
    return obj;
  },

  create(base?: DeepPartial<StreamOfflineMessage>): StreamOfflineMessage {
    return StreamOfflineMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<StreamOfflineMessage>): StreamOfflineMessage {
    const message = createBaseStreamOfflineMessage();
    message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
      ? BaseInfo.fromPartial(object.baseInfo)
      : undefined;
    return message;
  },
};

function createBaseChatClearMessage(): ChatClearMessage {
  return { baseInfo: undefined };
}

export const ChatClearMessage = {
  encode(message: ChatClearMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.baseInfo !== undefined) {
      BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ChatClearMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChatClearMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseInfo = BaseInfo.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ChatClearMessage {
    return { baseInfo: isSet(object.baseInfo) ? BaseInfo.fromJSON(object.baseInfo) : undefined };
  },

  toJSON(message: ChatClearMessage): unknown {
    const obj: any = {};
    message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? BaseInfo.toJSON(message.baseInfo) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ChatClearMessage>): ChatClearMessage {
    return ChatClearMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ChatClearMessage>): ChatClearMessage {
    const message = createBaseChatClearMessage();
    message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
      ? BaseInfo.fromPartial(object.baseInfo)
      : undefined;
    return message;
  },
};

function createBaseDonateMessage(): DonateMessage {
  return { baseInfo: undefined, userName: "", amount: "", currency: "", message: "" };
}

export const DonateMessage = {
  encode(message: DonateMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.baseInfo !== undefined) {
      BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.userName !== "") {
      writer.uint32(18).string(message.userName);
    }
    if (message.amount !== "") {
      writer.uint32(26).string(message.amount);
    }
    if (message.currency !== "") {
      writer.uint32(34).string(message.currency);
    }
    if (message.message !== "") {
      writer.uint32(42).string(message.message);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DonateMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDonateMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseInfo = BaseInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.userName = reader.string();
          break;
        case 3:
          message.amount = reader.string();
          break;
        case 4:
          message.currency = reader.string();
          break;
        case 5:
          message.message = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DonateMessage {
    return {
      baseInfo: isSet(object.baseInfo) ? BaseInfo.fromJSON(object.baseInfo) : undefined,
      userName: isSet(object.userName) ? String(object.userName) : "",
      amount: isSet(object.amount) ? String(object.amount) : "",
      currency: isSet(object.currency) ? String(object.currency) : "",
      message: isSet(object.message) ? String(object.message) : "",
    };
  },

  toJSON(message: DonateMessage): unknown {
    const obj: any = {};
    message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? BaseInfo.toJSON(message.baseInfo) : undefined);
    message.userName !== undefined && (obj.userName = message.userName);
    message.amount !== undefined && (obj.amount = message.amount);
    message.currency !== undefined && (obj.currency = message.currency);
    message.message !== undefined && (obj.message = message.message);
    return obj;
  },

  create(base?: DeepPartial<DonateMessage>): DonateMessage {
    return DonateMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DonateMessage>): DonateMessage {
    const message = createBaseDonateMessage();
    message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
      ? BaseInfo.fromPartial(object.baseInfo)
      : undefined;
    message.userName = object.userName ?? "";
    message.amount = object.amount ?? "";
    message.currency = object.currency ?? "";
    message.message = object.message ?? "";
    return message;
  },
};

function createBaseKeywordMatchedMessage(): KeywordMatchedMessage {
  return {
    baseInfo: undefined,
    keywordId: "",
    keywordName: "",
    keywordResponse: "",
    userId: "",
    userName: "",
    userDisplayName: "",
  };
}

export const KeywordMatchedMessage = {
  encode(message: KeywordMatchedMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.baseInfo !== undefined) {
      BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.keywordId !== "") {
      writer.uint32(18).string(message.keywordId);
    }
    if (message.keywordName !== "") {
      writer.uint32(26).string(message.keywordName);
    }
    if (message.keywordResponse !== "") {
      writer.uint32(34).string(message.keywordResponse);
    }
    if (message.userId !== "") {
      writer.uint32(42).string(message.userId);
    }
    if (message.userName !== "") {
      writer.uint32(50).string(message.userName);
    }
    if (message.userDisplayName !== "") {
      writer.uint32(58).string(message.userDisplayName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): KeywordMatchedMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseKeywordMatchedMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseInfo = BaseInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.keywordId = reader.string();
          break;
        case 3:
          message.keywordName = reader.string();
          break;
        case 4:
          message.keywordResponse = reader.string();
          break;
        case 5:
          message.userId = reader.string();
          break;
        case 6:
          message.userName = reader.string();
          break;
        case 7:
          message.userDisplayName = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): KeywordMatchedMessage {
    return {
      baseInfo: isSet(object.baseInfo) ? BaseInfo.fromJSON(object.baseInfo) : undefined,
      keywordId: isSet(object.keywordId) ? String(object.keywordId) : "",
      keywordName: isSet(object.keywordName) ? String(object.keywordName) : "",
      keywordResponse: isSet(object.keywordResponse) ? String(object.keywordResponse) : "",
      userId: isSet(object.userId) ? String(object.userId) : "",
      userName: isSet(object.userName) ? String(object.userName) : "",
      userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
    };
  },

  toJSON(message: KeywordMatchedMessage): unknown {
    const obj: any = {};
    message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? BaseInfo.toJSON(message.baseInfo) : undefined);
    message.keywordId !== undefined && (obj.keywordId = message.keywordId);
    message.keywordName !== undefined && (obj.keywordName = message.keywordName);
    message.keywordResponse !== undefined && (obj.keywordResponse = message.keywordResponse);
    message.userId !== undefined && (obj.userId = message.userId);
    message.userName !== undefined && (obj.userName = message.userName);
    message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
    return obj;
  },

  create(base?: DeepPartial<KeywordMatchedMessage>): KeywordMatchedMessage {
    return KeywordMatchedMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<KeywordMatchedMessage>): KeywordMatchedMessage {
    const message = createBaseKeywordMatchedMessage();
    message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
      ? BaseInfo.fromPartial(object.baseInfo)
      : undefined;
    message.keywordId = object.keywordId ?? "";
    message.keywordName = object.keywordName ?? "";
    message.keywordResponse = object.keywordResponse ?? "";
    message.userId = object.userId ?? "";
    message.userName = object.userName ?? "";
    message.userDisplayName = object.userDisplayName ?? "";
    return message;
  },
};

function createBaseGreetingSendedMessage(): GreetingSendedMessage {
  return { baseInfo: undefined, userId: "", userName: "", userDisplayName: "", greetingText: "" };
}

export const GreetingSendedMessage = {
  encode(message: GreetingSendedMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.baseInfo !== undefined) {
      BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
    }
    if (message.userId !== "") {
      writer.uint32(18).string(message.userId);
    }
    if (message.userName !== "") {
      writer.uint32(26).string(message.userName);
    }
    if (message.userDisplayName !== "") {
      writer.uint32(34).string(message.userDisplayName);
    }
    if (message.greetingText !== "") {
      writer.uint32(42).string(message.greetingText);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GreetingSendedMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGreetingSendedMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.baseInfo = BaseInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.userId = reader.string();
          break;
        case 3:
          message.userName = reader.string();
          break;
        case 4:
          message.userDisplayName = reader.string();
          break;
        case 5:
          message.greetingText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GreetingSendedMessage {
    return {
      baseInfo: isSet(object.baseInfo) ? BaseInfo.fromJSON(object.baseInfo) : undefined,
      userId: isSet(object.userId) ? String(object.userId) : "",
      userName: isSet(object.userName) ? String(object.userName) : "",
      userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
      greetingText: isSet(object.greetingText) ? String(object.greetingText) : "",
    };
  },

  toJSON(message: GreetingSendedMessage): unknown {
    const obj: any = {};
    message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? BaseInfo.toJSON(message.baseInfo) : undefined);
    message.userId !== undefined && (obj.userId = message.userId);
    message.userName !== undefined && (obj.userName = message.userName);
    message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
    message.greetingText !== undefined && (obj.greetingText = message.greetingText);
    return obj;
  },

  create(base?: DeepPartial<GreetingSendedMessage>): GreetingSendedMessage {
    return GreetingSendedMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GreetingSendedMessage>): GreetingSendedMessage {
    const message = createBaseGreetingSendedMessage();
    message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
      ? BaseInfo.fromPartial(object.baseInfo)
      : undefined;
    message.userId = object.userId ?? "";
    message.userName = object.userName ?? "";
    message.userDisplayName = object.userDisplayName ?? "";
    message.greetingText = object.greetingText ?? "";
    return message;
  },
};

export type EventsDefinition = typeof EventsDefinition;
export const EventsDefinition = {
  name: "Events",
  fullName: "events.Events",
  methods: {
    follow: {
      name: "Follow",
      requestType: FollowMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    subscribe: {
      name: "Subscribe",
      requestType: SubscribeMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    subGift: {
      name: "SubGift",
      requestType: SubGiftMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    reSubscribe: {
      name: "ReSubscribe",
      requestType: ReSubscribeMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    redemptionCreated: {
      name: "RedemptionCreated",
      requestType: RedemptionCreatedMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    commandUsed: {
      name: "CommandUsed",
      requestType: CommandUsedMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    firstUserMessage: {
      name: "FirstUserMessage",
      requestType: FirstUserMessageMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    raided: {
      name: "Raided",
      requestType: RaidedMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    titleOrCategoryChanged: {
      name: "TitleOrCategoryChanged",
      requestType: TitleOrCategoryChangedMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    streamOnline: {
      name: "StreamOnline",
      requestType: StreamOnlineMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    streamOffline: {
      name: "StreamOffline",
      requestType: StreamOfflineMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    chatClear: {
      name: "ChatClear",
      requestType: ChatClearMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    donate: {
      name: "Donate",
      requestType: DonateMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    keywordMatched: {
      name: "KeywordMatched",
      requestType: KeywordMatchedMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    greetingSended: {
      name: "GreetingSended",
      requestType: GreetingSendedMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface EventsServiceImplementation<CallContextExt = {}> {
  follow(request: FollowMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  subscribe(request: SubscribeMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  subGift(request: SubGiftMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  reSubscribe(request: ReSubscribeMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  redemptionCreated(
    request: RedemptionCreatedMessage,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<Empty>>;
  commandUsed(request: CommandUsedMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  firstUserMessage(
    request: FirstUserMessageMessage,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<Empty>>;
  raided(request: RaidedMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  titleOrCategoryChanged(
    request: TitleOrCategoryChangedMessage,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<Empty>>;
  streamOnline(request: StreamOnlineMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  streamOffline(request: StreamOfflineMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  chatClear(request: ChatClearMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  donate(request: DonateMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  keywordMatched(request: KeywordMatchedMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  greetingSended(request: GreetingSendedMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
}

export interface EventsClient<CallOptionsExt = {}> {
  follow(request: DeepPartial<FollowMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  subscribe(request: DeepPartial<SubscribeMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  subGift(request: DeepPartial<SubGiftMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  reSubscribe(request: DeepPartial<ReSubscribeMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  redemptionCreated(
    request: DeepPartial<RedemptionCreatedMessage>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<Empty>;
  commandUsed(request: DeepPartial<CommandUsedMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  firstUserMessage(
    request: DeepPartial<FirstUserMessageMessage>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<Empty>;
  raided(request: DeepPartial<RaidedMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  titleOrCategoryChanged(
    request: DeepPartial<TitleOrCategoryChangedMessage>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<Empty>;
  streamOnline(request: DeepPartial<StreamOnlineMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  streamOffline(request: DeepPartial<StreamOfflineMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  chatClear(request: DeepPartial<ChatClearMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  donate(request: DeepPartial<DonateMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  keywordMatched(request: DeepPartial<KeywordMatchedMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  greetingSended(request: DeepPartial<GreetingSendedMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
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
