/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal.js";

export const protobufPackage = "shared";

export enum FragmentType {
  TEXT = 0,
  CHEERMOTE = 1,
  EMOTE = 2,
  MENTION = 3,
  UNRECOGNIZED = -1,
}

export function fragmentTypeFromJSON(object: any): FragmentType {
  switch (object) {
    case 0:
    case "TEXT":
      return FragmentType.TEXT;
    case 1:
    case "CHEERMOTE":
      return FragmentType.CHEERMOTE;
    case 2:
    case "EMOTE":
      return FragmentType.EMOTE;
    case 3:
    case "MENTION":
      return FragmentType.MENTION;
    case -1:
    case "UNRECOGNIZED":
    default:
      return FragmentType.UNRECOGNIZED;
  }
}

export function fragmentTypeToJSON(object: FragmentType): string {
  switch (object) {
    case FragmentType.TEXT:
      return "TEXT";
    case FragmentType.CHEERMOTE:
      return "CHEERMOTE";
    case FragmentType.EMOTE:
      return "EMOTE";
    case FragmentType.MENTION:
      return "MENTION";
    case FragmentType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface ChatMessageMessageFragmentCheermote {
  prefix: string;
  bits: number;
  tier: number;
}

export interface ChatMessageMessageFragmentEmote {
  id: string;
  emoteSetId: string;
  ownerId: string;
  format: string[];
}

export interface ChatMessageMessageFragmentMention {
  userId: string;
  userName: string;
  userLogin: string;
}

export interface ChatMessageMessageFragment {
  type: FragmentType;
  text: string;
  cheermote: ChatMessageMessageFragmentCheermote | undefined;
  emote: ChatMessageMessageFragmentEmote | undefined;
  mention: ChatMessageMessageFragmentMention | undefined;
}

export interface ChatMessageMessage {
  text: string;
  fragments: ChatMessageMessageFragment[];
}

export interface ChatMessageBadge {
  id: string;
  setId: string;
  info: string;
}

export interface ChatMessageCheer {
  bits: number;
}

export interface ChatMessageReply {
  parentMessageId: string;
  parentMessageBody: string;
  parentUserId: string;
  parentUserName: string;
  parentUserLogin: string;
  threadMessageId: string;
  threadUserId: string;
  threadUserName: string;
  threadUserLogin: string;
}

export interface TwitchChatMessage {
  broadcasterUserId: string;
  broadcasterUserName: string;
  broadcasterUserLogin: string;
  chatterUserId: string;
  chatterUserName: string;
  chatterUserLogin: string;
  messageId: string;
  message: ChatMessageMessage | undefined;
  color: string;
  badges: ChatMessageBadge[];
  messageType: string;
  cheer: ChatMessageCheer | undefined;
  reply: ChatMessageReply | undefined;
  channelPointsCustomRewardId: string;
}

function createBaseChatMessageMessageFragmentCheermote(): ChatMessageMessageFragmentCheermote {
  return { prefix: "", bits: 0, tier: 0 };
}

export const ChatMessageMessageFragmentCheermote = {
  encode(message: ChatMessageMessageFragmentCheermote, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.prefix !== "") {
      writer.uint32(10).string(message.prefix);
    }
    if (message.bits !== 0) {
      writer.uint32(16).int64(message.bits);
    }
    if (message.tier !== 0) {
      writer.uint32(24).int64(message.tier);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ChatMessageMessageFragmentCheermote {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChatMessageMessageFragmentCheermote();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.prefix = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.bits = longToNumber(reader.int64() as Long);
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.tier = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ChatMessageMessageFragmentCheermote {
    return {
      prefix: isSet(object.prefix) ? globalThis.String(object.prefix) : "",
      bits: isSet(object.bits) ? globalThis.Number(object.bits) : 0,
      tier: isSet(object.tier) ? globalThis.Number(object.tier) : 0,
    };
  },

  toJSON(message: ChatMessageMessageFragmentCheermote): unknown {
    const obj: any = {};
    if (message.prefix !== "") {
      obj.prefix = message.prefix;
    }
    if (message.bits !== 0) {
      obj.bits = Math.round(message.bits);
    }
    if (message.tier !== 0) {
      obj.tier = Math.round(message.tier);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ChatMessageMessageFragmentCheermote>, I>>(
    base?: I,
  ): ChatMessageMessageFragmentCheermote {
    return ChatMessageMessageFragmentCheermote.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ChatMessageMessageFragmentCheermote>, I>>(
    object: I,
  ): ChatMessageMessageFragmentCheermote {
    const message = createBaseChatMessageMessageFragmentCheermote();
    message.prefix = object.prefix ?? "";
    message.bits = object.bits ?? 0;
    message.tier = object.tier ?? 0;
    return message;
  },
};

function createBaseChatMessageMessageFragmentEmote(): ChatMessageMessageFragmentEmote {
  return { id: "", emoteSetId: "", ownerId: "", format: [] };
}

export const ChatMessageMessageFragmentEmote = {
  encode(message: ChatMessageMessageFragmentEmote, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.emoteSetId !== "") {
      writer.uint32(18).string(message.emoteSetId);
    }
    if (message.ownerId !== "") {
      writer.uint32(26).string(message.ownerId);
    }
    for (const v of message.format) {
      writer.uint32(34).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ChatMessageMessageFragmentEmote {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChatMessageMessageFragmentEmote();
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

          message.emoteSetId = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.ownerId = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.format.push(reader.string());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ChatMessageMessageFragmentEmote {
    return {
      id: isSet(object.id) ? globalThis.String(object.id) : "",
      emoteSetId: isSet(object.emoteSetId) ? globalThis.String(object.emoteSetId) : "",
      ownerId: isSet(object.ownerId) ? globalThis.String(object.ownerId) : "",
      format: globalThis.Array.isArray(object?.format) ? object.format.map((e: any) => globalThis.String(e)) : [],
    };
  },

  toJSON(message: ChatMessageMessageFragmentEmote): unknown {
    const obj: any = {};
    if (message.id !== "") {
      obj.id = message.id;
    }
    if (message.emoteSetId !== "") {
      obj.emoteSetId = message.emoteSetId;
    }
    if (message.ownerId !== "") {
      obj.ownerId = message.ownerId;
    }
    if (message.format?.length) {
      obj.format = message.format;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ChatMessageMessageFragmentEmote>, I>>(base?: I): ChatMessageMessageFragmentEmote {
    return ChatMessageMessageFragmentEmote.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ChatMessageMessageFragmentEmote>, I>>(
    object: I,
  ): ChatMessageMessageFragmentEmote {
    const message = createBaseChatMessageMessageFragmentEmote();
    message.id = object.id ?? "";
    message.emoteSetId = object.emoteSetId ?? "";
    message.ownerId = object.ownerId ?? "";
    message.format = object.format?.map((e) => e) || [];
    return message;
  },
};

function createBaseChatMessageMessageFragmentMention(): ChatMessageMessageFragmentMention {
  return { userId: "", userName: "", userLogin: "" };
}

export const ChatMessageMessageFragmentMention = {
  encode(message: ChatMessageMessageFragmentMention, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.userName !== "") {
      writer.uint32(18).string(message.userName);
    }
    if (message.userLogin !== "") {
      writer.uint32(26).string(message.userLogin);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ChatMessageMessageFragmentMention {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChatMessageMessageFragmentMention();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.userName = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.userLogin = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ChatMessageMessageFragmentMention {
    return {
      userId: isSet(object.userId) ? globalThis.String(object.userId) : "",
      userName: isSet(object.userName) ? globalThis.String(object.userName) : "",
      userLogin: isSet(object.userLogin) ? globalThis.String(object.userLogin) : "",
    };
  },

  toJSON(message: ChatMessageMessageFragmentMention): unknown {
    const obj: any = {};
    if (message.userId !== "") {
      obj.userId = message.userId;
    }
    if (message.userName !== "") {
      obj.userName = message.userName;
    }
    if (message.userLogin !== "") {
      obj.userLogin = message.userLogin;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ChatMessageMessageFragmentMention>, I>>(
    base?: I,
  ): ChatMessageMessageFragmentMention {
    return ChatMessageMessageFragmentMention.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ChatMessageMessageFragmentMention>, I>>(
    object: I,
  ): ChatMessageMessageFragmentMention {
    const message = createBaseChatMessageMessageFragmentMention();
    message.userId = object.userId ?? "";
    message.userName = object.userName ?? "";
    message.userLogin = object.userLogin ?? "";
    return message;
  },
};

function createBaseChatMessageMessageFragment(): ChatMessageMessageFragment {
  return { type: 0, text: "", cheermote: undefined, emote: undefined, mention: undefined };
}

export const ChatMessageMessageFragment = {
  encode(message: ChatMessageMessageFragment, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== 0) {
      writer.uint32(8).int32(message.type);
    }
    if (message.text !== "") {
      writer.uint32(18).string(message.text);
    }
    if (message.cheermote !== undefined) {
      ChatMessageMessageFragmentCheermote.encode(message.cheermote, writer.uint32(26).fork()).ldelim();
    }
    if (message.emote !== undefined) {
      ChatMessageMessageFragmentEmote.encode(message.emote, writer.uint32(34).fork()).ldelim();
    }
    if (message.mention !== undefined) {
      ChatMessageMessageFragmentMention.encode(message.mention, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ChatMessageMessageFragment {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChatMessageMessageFragment();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.type = reader.int32() as any;
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.text = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.cheermote = ChatMessageMessageFragmentCheermote.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.emote = ChatMessageMessageFragmentEmote.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.mention = ChatMessageMessageFragmentMention.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ChatMessageMessageFragment {
    return {
      type: isSet(object.type) ? fragmentTypeFromJSON(object.type) : 0,
      text: isSet(object.text) ? globalThis.String(object.text) : "",
      cheermote: isSet(object.cheermote) ? ChatMessageMessageFragmentCheermote.fromJSON(object.cheermote) : undefined,
      emote: isSet(object.emote) ? ChatMessageMessageFragmentEmote.fromJSON(object.emote) : undefined,
      mention: isSet(object.mention) ? ChatMessageMessageFragmentMention.fromJSON(object.mention) : undefined,
    };
  },

  toJSON(message: ChatMessageMessageFragment): unknown {
    const obj: any = {};
    if (message.type !== 0) {
      obj.type = fragmentTypeToJSON(message.type);
    }
    if (message.text !== "") {
      obj.text = message.text;
    }
    if (message.cheermote !== undefined) {
      obj.cheermote = ChatMessageMessageFragmentCheermote.toJSON(message.cheermote);
    }
    if (message.emote !== undefined) {
      obj.emote = ChatMessageMessageFragmentEmote.toJSON(message.emote);
    }
    if (message.mention !== undefined) {
      obj.mention = ChatMessageMessageFragmentMention.toJSON(message.mention);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ChatMessageMessageFragment>, I>>(base?: I): ChatMessageMessageFragment {
    return ChatMessageMessageFragment.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ChatMessageMessageFragment>, I>>(object: I): ChatMessageMessageFragment {
    const message = createBaseChatMessageMessageFragment();
    message.type = object.type ?? 0;
    message.text = object.text ?? "";
    message.cheermote = (object.cheermote !== undefined && object.cheermote !== null)
      ? ChatMessageMessageFragmentCheermote.fromPartial(object.cheermote)
      : undefined;
    message.emote = (object.emote !== undefined && object.emote !== null)
      ? ChatMessageMessageFragmentEmote.fromPartial(object.emote)
      : undefined;
    message.mention = (object.mention !== undefined && object.mention !== null)
      ? ChatMessageMessageFragmentMention.fromPartial(object.mention)
      : undefined;
    return message;
  },
};

function createBaseChatMessageMessage(): ChatMessageMessage {
  return { text: "", fragments: [] };
}

export const ChatMessageMessage = {
  encode(message: ChatMessageMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.text !== "") {
      writer.uint32(10).string(message.text);
    }
    for (const v of message.fragments) {
      ChatMessageMessageFragment.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ChatMessageMessage {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChatMessageMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.text = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.fragments.push(ChatMessageMessageFragment.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ChatMessageMessage {
    return {
      text: isSet(object.text) ? globalThis.String(object.text) : "",
      fragments: globalThis.Array.isArray(object?.fragments)
        ? object.fragments.map((e: any) => ChatMessageMessageFragment.fromJSON(e))
        : [],
    };
  },

  toJSON(message: ChatMessageMessage): unknown {
    const obj: any = {};
    if (message.text !== "") {
      obj.text = message.text;
    }
    if (message.fragments?.length) {
      obj.fragments = message.fragments.map((e) => ChatMessageMessageFragment.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ChatMessageMessage>, I>>(base?: I): ChatMessageMessage {
    return ChatMessageMessage.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ChatMessageMessage>, I>>(object: I): ChatMessageMessage {
    const message = createBaseChatMessageMessage();
    message.text = object.text ?? "";
    message.fragments = object.fragments?.map((e) => ChatMessageMessageFragment.fromPartial(e)) || [];
    return message;
  },
};

function createBaseChatMessageBadge(): ChatMessageBadge {
  return { id: "", setId: "", info: "" };
}

export const ChatMessageBadge = {
  encode(message: ChatMessageBadge, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.setId !== "") {
      writer.uint32(18).string(message.setId);
    }
    if (message.info !== "") {
      writer.uint32(26).string(message.info);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ChatMessageBadge {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChatMessageBadge();
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

          message.setId = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.info = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ChatMessageBadge {
    return {
      id: isSet(object.id) ? globalThis.String(object.id) : "",
      setId: isSet(object.setId) ? globalThis.String(object.setId) : "",
      info: isSet(object.info) ? globalThis.String(object.info) : "",
    };
  },

  toJSON(message: ChatMessageBadge): unknown {
    const obj: any = {};
    if (message.id !== "") {
      obj.id = message.id;
    }
    if (message.setId !== "") {
      obj.setId = message.setId;
    }
    if (message.info !== "") {
      obj.info = message.info;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ChatMessageBadge>, I>>(base?: I): ChatMessageBadge {
    return ChatMessageBadge.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ChatMessageBadge>, I>>(object: I): ChatMessageBadge {
    const message = createBaseChatMessageBadge();
    message.id = object.id ?? "";
    message.setId = object.setId ?? "";
    message.info = object.info ?? "";
    return message;
  },
};

function createBaseChatMessageCheer(): ChatMessageCheer {
  return { bits: 0 };
}

export const ChatMessageCheer = {
  encode(message: ChatMessageCheer, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.bits !== 0) {
      writer.uint32(8).int64(message.bits);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ChatMessageCheer {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChatMessageCheer();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.bits = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ChatMessageCheer {
    return { bits: isSet(object.bits) ? globalThis.Number(object.bits) : 0 };
  },

  toJSON(message: ChatMessageCheer): unknown {
    const obj: any = {};
    if (message.bits !== 0) {
      obj.bits = Math.round(message.bits);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ChatMessageCheer>, I>>(base?: I): ChatMessageCheer {
    return ChatMessageCheer.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ChatMessageCheer>, I>>(object: I): ChatMessageCheer {
    const message = createBaseChatMessageCheer();
    message.bits = object.bits ?? 0;
    return message;
  },
};

function createBaseChatMessageReply(): ChatMessageReply {
  return {
    parentMessageId: "",
    parentMessageBody: "",
    parentUserId: "",
    parentUserName: "",
    parentUserLogin: "",
    threadMessageId: "",
    threadUserId: "",
    threadUserName: "",
    threadUserLogin: "",
  };
}

export const ChatMessageReply = {
  encode(message: ChatMessageReply, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.parentMessageId !== "") {
      writer.uint32(10).string(message.parentMessageId);
    }
    if (message.parentMessageBody !== "") {
      writer.uint32(18).string(message.parentMessageBody);
    }
    if (message.parentUserId !== "") {
      writer.uint32(26).string(message.parentUserId);
    }
    if (message.parentUserName !== "") {
      writer.uint32(34).string(message.parentUserName);
    }
    if (message.parentUserLogin !== "") {
      writer.uint32(42).string(message.parentUserLogin);
    }
    if (message.threadMessageId !== "") {
      writer.uint32(50).string(message.threadMessageId);
    }
    if (message.threadUserId !== "") {
      writer.uint32(58).string(message.threadUserId);
    }
    if (message.threadUserName !== "") {
      writer.uint32(66).string(message.threadUserName);
    }
    if (message.threadUserLogin !== "") {
      writer.uint32(74).string(message.threadUserLogin);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ChatMessageReply {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChatMessageReply();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.parentMessageId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.parentMessageBody = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.parentUserId = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.parentUserName = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.parentUserLogin = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.threadMessageId = reader.string();
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.threadUserId = reader.string();
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.threadUserName = reader.string();
          continue;
        case 9:
          if (tag !== 74) {
            break;
          }

          message.threadUserLogin = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ChatMessageReply {
    return {
      parentMessageId: isSet(object.parentMessageId) ? globalThis.String(object.parentMessageId) : "",
      parentMessageBody: isSet(object.parentMessageBody) ? globalThis.String(object.parentMessageBody) : "",
      parentUserId: isSet(object.parentUserId) ? globalThis.String(object.parentUserId) : "",
      parentUserName: isSet(object.parentUserName) ? globalThis.String(object.parentUserName) : "",
      parentUserLogin: isSet(object.parentUserLogin) ? globalThis.String(object.parentUserLogin) : "",
      threadMessageId: isSet(object.threadMessageId) ? globalThis.String(object.threadMessageId) : "",
      threadUserId: isSet(object.threadUserId) ? globalThis.String(object.threadUserId) : "",
      threadUserName: isSet(object.threadUserName) ? globalThis.String(object.threadUserName) : "",
      threadUserLogin: isSet(object.threadUserLogin) ? globalThis.String(object.threadUserLogin) : "",
    };
  },

  toJSON(message: ChatMessageReply): unknown {
    const obj: any = {};
    if (message.parentMessageId !== "") {
      obj.parentMessageId = message.parentMessageId;
    }
    if (message.parentMessageBody !== "") {
      obj.parentMessageBody = message.parentMessageBody;
    }
    if (message.parentUserId !== "") {
      obj.parentUserId = message.parentUserId;
    }
    if (message.parentUserName !== "") {
      obj.parentUserName = message.parentUserName;
    }
    if (message.parentUserLogin !== "") {
      obj.parentUserLogin = message.parentUserLogin;
    }
    if (message.threadMessageId !== "") {
      obj.threadMessageId = message.threadMessageId;
    }
    if (message.threadUserId !== "") {
      obj.threadUserId = message.threadUserId;
    }
    if (message.threadUserName !== "") {
      obj.threadUserName = message.threadUserName;
    }
    if (message.threadUserLogin !== "") {
      obj.threadUserLogin = message.threadUserLogin;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ChatMessageReply>, I>>(base?: I): ChatMessageReply {
    return ChatMessageReply.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ChatMessageReply>, I>>(object: I): ChatMessageReply {
    const message = createBaseChatMessageReply();
    message.parentMessageId = object.parentMessageId ?? "";
    message.parentMessageBody = object.parentMessageBody ?? "";
    message.parentUserId = object.parentUserId ?? "";
    message.parentUserName = object.parentUserName ?? "";
    message.parentUserLogin = object.parentUserLogin ?? "";
    message.threadMessageId = object.threadMessageId ?? "";
    message.threadUserId = object.threadUserId ?? "";
    message.threadUserName = object.threadUserName ?? "";
    message.threadUserLogin = object.threadUserLogin ?? "";
    return message;
  },
};

function createBaseTwitchChatMessage(): TwitchChatMessage {
  return {
    broadcasterUserId: "",
    broadcasterUserName: "",
    broadcasterUserLogin: "",
    chatterUserId: "",
    chatterUserName: "",
    chatterUserLogin: "",
    messageId: "",
    message: undefined,
    color: "",
    badges: [],
    messageType: "",
    cheer: undefined,
    reply: undefined,
    channelPointsCustomRewardId: "",
  };
}

export const TwitchChatMessage = {
  encode(message: TwitchChatMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.broadcasterUserId !== "") {
      writer.uint32(10).string(message.broadcasterUserId);
    }
    if (message.broadcasterUserName !== "") {
      writer.uint32(18).string(message.broadcasterUserName);
    }
    if (message.broadcasterUserLogin !== "") {
      writer.uint32(26).string(message.broadcasterUserLogin);
    }
    if (message.chatterUserId !== "") {
      writer.uint32(34).string(message.chatterUserId);
    }
    if (message.chatterUserName !== "") {
      writer.uint32(42).string(message.chatterUserName);
    }
    if (message.chatterUserLogin !== "") {
      writer.uint32(50).string(message.chatterUserLogin);
    }
    if (message.messageId !== "") {
      writer.uint32(58).string(message.messageId);
    }
    if (message.message !== undefined) {
      ChatMessageMessage.encode(message.message, writer.uint32(66).fork()).ldelim();
    }
    if (message.color !== "") {
      writer.uint32(74).string(message.color);
    }
    for (const v of message.badges) {
      ChatMessageBadge.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    if (message.messageType !== "") {
      writer.uint32(90).string(message.messageType);
    }
    if (message.cheer !== undefined) {
      ChatMessageCheer.encode(message.cheer, writer.uint32(98).fork()).ldelim();
    }
    if (message.reply !== undefined) {
      ChatMessageReply.encode(message.reply, writer.uint32(106).fork()).ldelim();
    }
    if (message.channelPointsCustomRewardId !== "") {
      writer.uint32(114).string(message.channelPointsCustomRewardId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TwitchChatMessage {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTwitchChatMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.broadcasterUserId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.broadcasterUserName = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.broadcasterUserLogin = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.chatterUserId = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.chatterUserName = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.chatterUserLogin = reader.string();
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.messageId = reader.string();
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.message = ChatMessageMessage.decode(reader, reader.uint32());
          continue;
        case 9:
          if (tag !== 74) {
            break;
          }

          message.color = reader.string();
          continue;
        case 10:
          if (tag !== 82) {
            break;
          }

          message.badges.push(ChatMessageBadge.decode(reader, reader.uint32()));
          continue;
        case 11:
          if (tag !== 90) {
            break;
          }

          message.messageType = reader.string();
          continue;
        case 12:
          if (tag !== 98) {
            break;
          }

          message.cheer = ChatMessageCheer.decode(reader, reader.uint32());
          continue;
        case 13:
          if (tag !== 106) {
            break;
          }

          message.reply = ChatMessageReply.decode(reader, reader.uint32());
          continue;
        case 14:
          if (tag !== 114) {
            break;
          }

          message.channelPointsCustomRewardId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TwitchChatMessage {
    return {
      broadcasterUserId: isSet(object.broadcasterUserId) ? globalThis.String(object.broadcasterUserId) : "",
      broadcasterUserName: isSet(object.broadcasterUserName) ? globalThis.String(object.broadcasterUserName) : "",
      broadcasterUserLogin: isSet(object.broadcasterUserLogin) ? globalThis.String(object.broadcasterUserLogin) : "",
      chatterUserId: isSet(object.chatterUserId) ? globalThis.String(object.chatterUserId) : "",
      chatterUserName: isSet(object.chatterUserName) ? globalThis.String(object.chatterUserName) : "",
      chatterUserLogin: isSet(object.chatterUserLogin) ? globalThis.String(object.chatterUserLogin) : "",
      messageId: isSet(object.messageId) ? globalThis.String(object.messageId) : "",
      message: isSet(object.message) ? ChatMessageMessage.fromJSON(object.message) : undefined,
      color: isSet(object.color) ? globalThis.String(object.color) : "",
      badges: globalThis.Array.isArray(object?.badges)
        ? object.badges.map((e: any) => ChatMessageBadge.fromJSON(e))
        : [],
      messageType: isSet(object.messageType) ? globalThis.String(object.messageType) : "",
      cheer: isSet(object.cheer) ? ChatMessageCheer.fromJSON(object.cheer) : undefined,
      reply: isSet(object.reply) ? ChatMessageReply.fromJSON(object.reply) : undefined,
      channelPointsCustomRewardId: isSet(object.channelPointsCustomRewardId)
        ? globalThis.String(object.channelPointsCustomRewardId)
        : "",
    };
  },

  toJSON(message: TwitchChatMessage): unknown {
    const obj: any = {};
    if (message.broadcasterUserId !== "") {
      obj.broadcasterUserId = message.broadcasterUserId;
    }
    if (message.broadcasterUserName !== "") {
      obj.broadcasterUserName = message.broadcasterUserName;
    }
    if (message.broadcasterUserLogin !== "") {
      obj.broadcasterUserLogin = message.broadcasterUserLogin;
    }
    if (message.chatterUserId !== "") {
      obj.chatterUserId = message.chatterUserId;
    }
    if (message.chatterUserName !== "") {
      obj.chatterUserName = message.chatterUserName;
    }
    if (message.chatterUserLogin !== "") {
      obj.chatterUserLogin = message.chatterUserLogin;
    }
    if (message.messageId !== "") {
      obj.messageId = message.messageId;
    }
    if (message.message !== undefined) {
      obj.message = ChatMessageMessage.toJSON(message.message);
    }
    if (message.color !== "") {
      obj.color = message.color;
    }
    if (message.badges?.length) {
      obj.badges = message.badges.map((e) => ChatMessageBadge.toJSON(e));
    }
    if (message.messageType !== "") {
      obj.messageType = message.messageType;
    }
    if (message.cheer !== undefined) {
      obj.cheer = ChatMessageCheer.toJSON(message.cheer);
    }
    if (message.reply !== undefined) {
      obj.reply = ChatMessageReply.toJSON(message.reply);
    }
    if (message.channelPointsCustomRewardId !== "") {
      obj.channelPointsCustomRewardId = message.channelPointsCustomRewardId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TwitchChatMessage>, I>>(base?: I): TwitchChatMessage {
    return TwitchChatMessage.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TwitchChatMessage>, I>>(object: I): TwitchChatMessage {
    const message = createBaseTwitchChatMessage();
    message.broadcasterUserId = object.broadcasterUserId ?? "";
    message.broadcasterUserName = object.broadcasterUserName ?? "";
    message.broadcasterUserLogin = object.broadcasterUserLogin ?? "";
    message.chatterUserId = object.chatterUserId ?? "";
    message.chatterUserName = object.chatterUserName ?? "";
    message.chatterUserLogin = object.chatterUserLogin ?? "";
    message.messageId = object.messageId ?? "";
    message.message = (object.message !== undefined && object.message !== null)
      ? ChatMessageMessage.fromPartial(object.message)
      : undefined;
    message.color = object.color ?? "";
    message.badges = object.badges?.map((e) => ChatMessageBadge.fromPartial(e)) || [];
    message.messageType = object.messageType ?? "";
    message.cheer = (object.cheer !== undefined && object.cheer !== null)
      ? ChatMessageCheer.fromPartial(object.cheer)
      : undefined;
    message.reply = (object.reply !== undefined && object.reply !== null)
      ? ChatMessageReply.fromPartial(object.reply)
      : undefined;
    message.channelPointsCustomRewardId = object.channelPointsCustomRewardId ?? "";
    return message;
  },
};

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
