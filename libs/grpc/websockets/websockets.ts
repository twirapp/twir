/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal.js";
import { Empty } from "../google/protobuf/empty.js";

export const protobufPackage = "websockets";

export enum RefreshOverlaySettingsName {
  CUSTOM = 0,
  KAPPAGEN = 1,
  BRB = 2,
  DUDES = 3,
  CHAT = 4,
  NOW_PLAYING = 5,
  UNRECOGNIZED = -1,
}

export function refreshOverlaySettingsNameFromJSON(object: any): RefreshOverlaySettingsName {
  switch (object) {
    case 0:
    case "CUSTOM":
      return RefreshOverlaySettingsName.CUSTOM;
    case 1:
    case "KAPPAGEN":
      return RefreshOverlaySettingsName.KAPPAGEN;
    case 2:
    case "BRB":
      return RefreshOverlaySettingsName.BRB;
    case 3:
    case "DUDES":
      return RefreshOverlaySettingsName.DUDES;
    case 4:
    case "CHAT":
      return RefreshOverlaySettingsName.CHAT;
    case 5:
    case "NOW_PLAYING":
      return RefreshOverlaySettingsName.NOW_PLAYING;
    case -1:
    case "UNRECOGNIZED":
    default:
      return RefreshOverlaySettingsName.UNRECOGNIZED;
  }
}

export function refreshOverlaySettingsNameToJSON(object: RefreshOverlaySettingsName): string {
  switch (object) {
    case RefreshOverlaySettingsName.CUSTOM:
      return "CUSTOM";
    case RefreshOverlaySettingsName.KAPPAGEN:
      return "KAPPAGEN";
    case RefreshOverlaySettingsName.BRB:
      return "BRB";
    case RefreshOverlaySettingsName.DUDES:
      return "DUDES";
    case RefreshOverlaySettingsName.CHAT:
      return "CHAT";
    case RefreshOverlaySettingsName.NOW_PLAYING:
      return "NOW_PLAYING";
    case RefreshOverlaySettingsName.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface YoutubeAddSongToQueueRequest {
  channelId: string;
  entityId: string;
}

export interface YoutubeRemoveSongFromQueueRequest {
  channelId: string;
  entityId: string;
}

export interface ObsSetSceneMessage {
  channelId: string;
  sceneName: string;
}

export interface ObsToggleSourceMessage {
  channelId: string;
  sourceName: string;
}

export interface ObsToggleAudioMessage {
  channelId: string;
  audioSourceName: string;
}

export interface ObsAudioSetVolumeMessage {
  channelId: string;
  audioSourceName: string;
  volume: number;
}

export interface ObsAudioIncreaseVolumeMessage {
  channelId: string;
  audioSourceName: string;
  step: number;
}

export interface ObsAudioDecreaseVolumeMessage {
  channelId: string;
  audioSourceName: string;
  step: number;
}

export interface ObsAudioDisableOrEnableMessage {
  channelId: string;
  audioSourceName: string;
}

export interface ObsStopOrStartStream {
  channelId: string;
}

export interface TTSMessage {
  channelId: string;
  text: string;
  voice: string;
  rate: string;
  pitch: string;
  volume: string;
}

export interface TTSSkipMessage {
  channelId: string;
}

export interface ObsCheckUserConnectedRequest {
  userId: string;
}

export interface ObsCheckUserConnectedResponse {
  state: boolean;
}

export interface TriggerAlertRequest {
  channelId: string;
  alertId: string;
}

export interface RefreshOverlaysRequest {
  channelId: string;
  overlayName: RefreshOverlaySettingsName;
  overlayId?: string | undefined;
}

export interface TriggerKappagenRequest {
  channelId: string;
  text: string;
  emotes: TriggerKappagenRequest_Emote[];
}

export interface TriggerKappagenRequest_Emote {
  id: string;
  positions: string[];
}

export interface TriggerKappagenByEventRequest {
  channelId: string;
  event: number;
}

export interface TriggerShowBrbRequest {
  channelId: string;
  minutes: number;
  text?: string | undefined;
}

export interface TriggerHideBrbRequest {
  channelId: string;
}

export interface DudesJumpRequest {
  channelId: string;
  userId: string;
  userDisplayName: string;
  userName: string;
  userColor: string;
}

export interface DudesUserPunishedRequest {
  channelId: string;
  userId: string;
  userDisplayName: string;
  userName: string;
}

function createBaseYoutubeAddSongToQueueRequest(): YoutubeAddSongToQueueRequest {
  return { channelId: "", entityId: "" };
}

export const YoutubeAddSongToQueueRequest = {
  encode(message: YoutubeAddSongToQueueRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.entityId !== "") {
      writer.uint32(18).string(message.entityId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): YoutubeAddSongToQueueRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseYoutubeAddSongToQueueRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.entityId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): YoutubeAddSongToQueueRequest {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      entityId: isSet(object.entityId) ? globalThis.String(object.entityId) : "",
    };
  },

  toJSON(message: YoutubeAddSongToQueueRequest): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.entityId !== "") {
      obj.entityId = message.entityId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<YoutubeAddSongToQueueRequest>, I>>(base?: I): YoutubeAddSongToQueueRequest {
    return YoutubeAddSongToQueueRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<YoutubeAddSongToQueueRequest>, I>>(object: I): YoutubeAddSongToQueueRequest {
    const message = createBaseYoutubeAddSongToQueueRequest();
    message.channelId = object.channelId ?? "";
    message.entityId = object.entityId ?? "";
    return message;
  },
};

function createBaseYoutubeRemoveSongFromQueueRequest(): YoutubeRemoveSongFromQueueRequest {
  return { channelId: "", entityId: "" };
}

export const YoutubeRemoveSongFromQueueRequest = {
  encode(message: YoutubeRemoveSongFromQueueRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.entityId !== "") {
      writer.uint32(18).string(message.entityId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): YoutubeRemoveSongFromQueueRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseYoutubeRemoveSongFromQueueRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.entityId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): YoutubeRemoveSongFromQueueRequest {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      entityId: isSet(object.entityId) ? globalThis.String(object.entityId) : "",
    };
  },

  toJSON(message: YoutubeRemoveSongFromQueueRequest): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.entityId !== "") {
      obj.entityId = message.entityId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<YoutubeRemoveSongFromQueueRequest>, I>>(
    base?: I,
  ): YoutubeRemoveSongFromQueueRequest {
    return YoutubeRemoveSongFromQueueRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<YoutubeRemoveSongFromQueueRequest>, I>>(
    object: I,
  ): YoutubeRemoveSongFromQueueRequest {
    const message = createBaseYoutubeRemoveSongFromQueueRequest();
    message.channelId = object.channelId ?? "";
    message.entityId = object.entityId ?? "";
    return message;
  },
};

function createBaseObsSetSceneMessage(): ObsSetSceneMessage {
  return { channelId: "", sceneName: "" };
}

export const ObsSetSceneMessage = {
  encode(message: ObsSetSceneMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.sceneName !== "") {
      writer.uint32(18).string(message.sceneName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ObsSetSceneMessage {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsSetSceneMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.sceneName = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ObsSetSceneMessage {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      sceneName: isSet(object.sceneName) ? globalThis.String(object.sceneName) : "",
    };
  },

  toJSON(message: ObsSetSceneMessage): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.sceneName !== "") {
      obj.sceneName = message.sceneName;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ObsSetSceneMessage>, I>>(base?: I): ObsSetSceneMessage {
    return ObsSetSceneMessage.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ObsSetSceneMessage>, I>>(object: I): ObsSetSceneMessage {
    const message = createBaseObsSetSceneMessage();
    message.channelId = object.channelId ?? "";
    message.sceneName = object.sceneName ?? "";
    return message;
  },
};

function createBaseObsToggleSourceMessage(): ObsToggleSourceMessage {
  return { channelId: "", sourceName: "" };
}

export const ObsToggleSourceMessage = {
  encode(message: ObsToggleSourceMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.sourceName !== "") {
      writer.uint32(18).string(message.sourceName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ObsToggleSourceMessage {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsToggleSourceMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.sourceName = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ObsToggleSourceMessage {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      sourceName: isSet(object.sourceName) ? globalThis.String(object.sourceName) : "",
    };
  },

  toJSON(message: ObsToggleSourceMessage): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.sourceName !== "") {
      obj.sourceName = message.sourceName;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ObsToggleSourceMessage>, I>>(base?: I): ObsToggleSourceMessage {
    return ObsToggleSourceMessage.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ObsToggleSourceMessage>, I>>(object: I): ObsToggleSourceMessage {
    const message = createBaseObsToggleSourceMessage();
    message.channelId = object.channelId ?? "";
    message.sourceName = object.sourceName ?? "";
    return message;
  },
};

function createBaseObsToggleAudioMessage(): ObsToggleAudioMessage {
  return { channelId: "", audioSourceName: "" };
}

export const ObsToggleAudioMessage = {
  encode(message: ObsToggleAudioMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.audioSourceName !== "") {
      writer.uint32(18).string(message.audioSourceName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ObsToggleAudioMessage {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsToggleAudioMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.audioSourceName = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ObsToggleAudioMessage {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      audioSourceName: isSet(object.audioSourceName) ? globalThis.String(object.audioSourceName) : "",
    };
  },

  toJSON(message: ObsToggleAudioMessage): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.audioSourceName !== "") {
      obj.audioSourceName = message.audioSourceName;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ObsToggleAudioMessage>, I>>(base?: I): ObsToggleAudioMessage {
    return ObsToggleAudioMessage.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ObsToggleAudioMessage>, I>>(object: I): ObsToggleAudioMessage {
    const message = createBaseObsToggleAudioMessage();
    message.channelId = object.channelId ?? "";
    message.audioSourceName = object.audioSourceName ?? "";
    return message;
  },
};

function createBaseObsAudioSetVolumeMessage(): ObsAudioSetVolumeMessage {
  return { channelId: "", audioSourceName: "", volume: 0 };
}

export const ObsAudioSetVolumeMessage = {
  encode(message: ObsAudioSetVolumeMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.audioSourceName !== "") {
      writer.uint32(18).string(message.audioSourceName);
    }
    if (message.volume !== 0) {
      writer.uint32(24).uint32(message.volume);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ObsAudioSetVolumeMessage {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsAudioSetVolumeMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.audioSourceName = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.volume = reader.uint32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ObsAudioSetVolumeMessage {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      audioSourceName: isSet(object.audioSourceName) ? globalThis.String(object.audioSourceName) : "",
      volume: isSet(object.volume) ? globalThis.Number(object.volume) : 0,
    };
  },

  toJSON(message: ObsAudioSetVolumeMessage): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.audioSourceName !== "") {
      obj.audioSourceName = message.audioSourceName;
    }
    if (message.volume !== 0) {
      obj.volume = Math.round(message.volume);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ObsAudioSetVolumeMessage>, I>>(base?: I): ObsAudioSetVolumeMessage {
    return ObsAudioSetVolumeMessage.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ObsAudioSetVolumeMessage>, I>>(object: I): ObsAudioSetVolumeMessage {
    const message = createBaseObsAudioSetVolumeMessage();
    message.channelId = object.channelId ?? "";
    message.audioSourceName = object.audioSourceName ?? "";
    message.volume = object.volume ?? 0;
    return message;
  },
};

function createBaseObsAudioIncreaseVolumeMessage(): ObsAudioIncreaseVolumeMessage {
  return { channelId: "", audioSourceName: "", step: 0 };
}

export const ObsAudioIncreaseVolumeMessage = {
  encode(message: ObsAudioIncreaseVolumeMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.audioSourceName !== "") {
      writer.uint32(18).string(message.audioSourceName);
    }
    if (message.step !== 0) {
      writer.uint32(24).uint32(message.step);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ObsAudioIncreaseVolumeMessage {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsAudioIncreaseVolumeMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.audioSourceName = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.step = reader.uint32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ObsAudioIncreaseVolumeMessage {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      audioSourceName: isSet(object.audioSourceName) ? globalThis.String(object.audioSourceName) : "",
      step: isSet(object.step) ? globalThis.Number(object.step) : 0,
    };
  },

  toJSON(message: ObsAudioIncreaseVolumeMessage): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.audioSourceName !== "") {
      obj.audioSourceName = message.audioSourceName;
    }
    if (message.step !== 0) {
      obj.step = Math.round(message.step);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ObsAudioIncreaseVolumeMessage>, I>>(base?: I): ObsAudioIncreaseVolumeMessage {
    return ObsAudioIncreaseVolumeMessage.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ObsAudioIncreaseVolumeMessage>, I>>(
    object: I,
  ): ObsAudioIncreaseVolumeMessage {
    const message = createBaseObsAudioIncreaseVolumeMessage();
    message.channelId = object.channelId ?? "";
    message.audioSourceName = object.audioSourceName ?? "";
    message.step = object.step ?? 0;
    return message;
  },
};

function createBaseObsAudioDecreaseVolumeMessage(): ObsAudioDecreaseVolumeMessage {
  return { channelId: "", audioSourceName: "", step: 0 };
}

export const ObsAudioDecreaseVolumeMessage = {
  encode(message: ObsAudioDecreaseVolumeMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.audioSourceName !== "") {
      writer.uint32(18).string(message.audioSourceName);
    }
    if (message.step !== 0) {
      writer.uint32(24).uint32(message.step);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ObsAudioDecreaseVolumeMessage {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsAudioDecreaseVolumeMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.audioSourceName = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.step = reader.uint32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ObsAudioDecreaseVolumeMessage {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      audioSourceName: isSet(object.audioSourceName) ? globalThis.String(object.audioSourceName) : "",
      step: isSet(object.step) ? globalThis.Number(object.step) : 0,
    };
  },

  toJSON(message: ObsAudioDecreaseVolumeMessage): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.audioSourceName !== "") {
      obj.audioSourceName = message.audioSourceName;
    }
    if (message.step !== 0) {
      obj.step = Math.round(message.step);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ObsAudioDecreaseVolumeMessage>, I>>(base?: I): ObsAudioDecreaseVolumeMessage {
    return ObsAudioDecreaseVolumeMessage.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ObsAudioDecreaseVolumeMessage>, I>>(
    object: I,
  ): ObsAudioDecreaseVolumeMessage {
    const message = createBaseObsAudioDecreaseVolumeMessage();
    message.channelId = object.channelId ?? "";
    message.audioSourceName = object.audioSourceName ?? "";
    message.step = object.step ?? 0;
    return message;
  },
};

function createBaseObsAudioDisableOrEnableMessage(): ObsAudioDisableOrEnableMessage {
  return { channelId: "", audioSourceName: "" };
}

export const ObsAudioDisableOrEnableMessage = {
  encode(message: ObsAudioDisableOrEnableMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.audioSourceName !== "") {
      writer.uint32(18).string(message.audioSourceName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ObsAudioDisableOrEnableMessage {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsAudioDisableOrEnableMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.audioSourceName = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ObsAudioDisableOrEnableMessage {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      audioSourceName: isSet(object.audioSourceName) ? globalThis.String(object.audioSourceName) : "",
    };
  },

  toJSON(message: ObsAudioDisableOrEnableMessage): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.audioSourceName !== "") {
      obj.audioSourceName = message.audioSourceName;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ObsAudioDisableOrEnableMessage>, I>>(base?: I): ObsAudioDisableOrEnableMessage {
    return ObsAudioDisableOrEnableMessage.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ObsAudioDisableOrEnableMessage>, I>>(
    object: I,
  ): ObsAudioDisableOrEnableMessage {
    const message = createBaseObsAudioDisableOrEnableMessage();
    message.channelId = object.channelId ?? "";
    message.audioSourceName = object.audioSourceName ?? "";
    return message;
  },
};

function createBaseObsStopOrStartStream(): ObsStopOrStartStream {
  return { channelId: "" };
}

export const ObsStopOrStartStream = {
  encode(message: ObsStopOrStartStream, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ObsStopOrStartStream {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsStopOrStartStream();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ObsStopOrStartStream {
    return { channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "" };
  },

  toJSON(message: ObsStopOrStartStream): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ObsStopOrStartStream>, I>>(base?: I): ObsStopOrStartStream {
    return ObsStopOrStartStream.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ObsStopOrStartStream>, I>>(object: I): ObsStopOrStartStream {
    const message = createBaseObsStopOrStartStream();
    message.channelId = object.channelId ?? "";
    return message;
  },
};

function createBaseTTSMessage(): TTSMessage {
  return { channelId: "", text: "", voice: "", rate: "", pitch: "", volume: "" };
}

export const TTSMessage = {
  encode(message: TTSMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.text !== "") {
      writer.uint32(18).string(message.text);
    }
    if (message.voice !== "") {
      writer.uint32(26).string(message.voice);
    }
    if (message.rate !== "") {
      writer.uint32(34).string(message.rate);
    }
    if (message.pitch !== "") {
      writer.uint32(42).string(message.pitch);
    }
    if (message.volume !== "") {
      writer.uint32(50).string(message.volume);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TTSMessage {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTTSMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
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

          message.voice = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.rate = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.pitch = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.volume = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TTSMessage {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      text: isSet(object.text) ? globalThis.String(object.text) : "",
      voice: isSet(object.voice) ? globalThis.String(object.voice) : "",
      rate: isSet(object.rate) ? globalThis.String(object.rate) : "",
      pitch: isSet(object.pitch) ? globalThis.String(object.pitch) : "",
      volume: isSet(object.volume) ? globalThis.String(object.volume) : "",
    };
  },

  toJSON(message: TTSMessage): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.text !== "") {
      obj.text = message.text;
    }
    if (message.voice !== "") {
      obj.voice = message.voice;
    }
    if (message.rate !== "") {
      obj.rate = message.rate;
    }
    if (message.pitch !== "") {
      obj.pitch = message.pitch;
    }
    if (message.volume !== "") {
      obj.volume = message.volume;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TTSMessage>, I>>(base?: I): TTSMessage {
    return TTSMessage.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TTSMessage>, I>>(object: I): TTSMessage {
    const message = createBaseTTSMessage();
    message.channelId = object.channelId ?? "";
    message.text = object.text ?? "";
    message.voice = object.voice ?? "";
    message.rate = object.rate ?? "";
    message.pitch = object.pitch ?? "";
    message.volume = object.volume ?? "";
    return message;
  },
};

function createBaseTTSSkipMessage(): TTSSkipMessage {
  return { channelId: "" };
}

export const TTSSkipMessage = {
  encode(message: TTSSkipMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TTSSkipMessage {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTTSSkipMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TTSSkipMessage {
    return { channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "" };
  },

  toJSON(message: TTSSkipMessage): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TTSSkipMessage>, I>>(base?: I): TTSSkipMessage {
    return TTSSkipMessage.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TTSSkipMessage>, I>>(object: I): TTSSkipMessage {
    const message = createBaseTTSSkipMessage();
    message.channelId = object.channelId ?? "";
    return message;
  },
};

function createBaseObsCheckUserConnectedRequest(): ObsCheckUserConnectedRequest {
  return { userId: "" };
}

export const ObsCheckUserConnectedRequest = {
  encode(message: ObsCheckUserConnectedRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ObsCheckUserConnectedRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsCheckUserConnectedRequest();
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

  fromJSON(object: any): ObsCheckUserConnectedRequest {
    return { userId: isSet(object.userId) ? globalThis.String(object.userId) : "" };
  },

  toJSON(message: ObsCheckUserConnectedRequest): unknown {
    const obj: any = {};
    if (message.userId !== "") {
      obj.userId = message.userId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ObsCheckUserConnectedRequest>, I>>(base?: I): ObsCheckUserConnectedRequest {
    return ObsCheckUserConnectedRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ObsCheckUserConnectedRequest>, I>>(object: I): ObsCheckUserConnectedRequest {
    const message = createBaseObsCheckUserConnectedRequest();
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseObsCheckUserConnectedResponse(): ObsCheckUserConnectedResponse {
  return { state: false };
}

export const ObsCheckUserConnectedResponse = {
  encode(message: ObsCheckUserConnectedResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.state !== false) {
      writer.uint32(8).bool(message.state);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ObsCheckUserConnectedResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsCheckUserConnectedResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.state = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ObsCheckUserConnectedResponse {
    return { state: isSet(object.state) ? globalThis.Boolean(object.state) : false };
  },

  toJSON(message: ObsCheckUserConnectedResponse): unknown {
    const obj: any = {};
    if (message.state !== false) {
      obj.state = message.state;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ObsCheckUserConnectedResponse>, I>>(base?: I): ObsCheckUserConnectedResponse {
    return ObsCheckUserConnectedResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ObsCheckUserConnectedResponse>, I>>(
    object: I,
  ): ObsCheckUserConnectedResponse {
    const message = createBaseObsCheckUserConnectedResponse();
    message.state = object.state ?? false;
    return message;
  },
};

function createBaseTriggerAlertRequest(): TriggerAlertRequest {
  return { channelId: "", alertId: "" };
}

export const TriggerAlertRequest = {
  encode(message: TriggerAlertRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.alertId !== "") {
      writer.uint32(18).string(message.alertId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TriggerAlertRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTriggerAlertRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.alertId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TriggerAlertRequest {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      alertId: isSet(object.alertId) ? globalThis.String(object.alertId) : "",
    };
  },

  toJSON(message: TriggerAlertRequest): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.alertId !== "") {
      obj.alertId = message.alertId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TriggerAlertRequest>, I>>(base?: I): TriggerAlertRequest {
    return TriggerAlertRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TriggerAlertRequest>, I>>(object: I): TriggerAlertRequest {
    const message = createBaseTriggerAlertRequest();
    message.channelId = object.channelId ?? "";
    message.alertId = object.alertId ?? "";
    return message;
  },
};

function createBaseRefreshOverlaysRequest(): RefreshOverlaysRequest {
  return { channelId: "", overlayName: 0, overlayId: undefined };
}

export const RefreshOverlaysRequest = {
  encode(message: RefreshOverlaysRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.overlayName !== 0) {
      writer.uint32(16).int32(message.overlayName);
    }
    if (message.overlayId !== undefined) {
      writer.uint32(26).string(message.overlayId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RefreshOverlaysRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRefreshOverlaysRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.overlayName = reader.int32() as any;
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.overlayId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RefreshOverlaysRequest {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      overlayName: isSet(object.overlayName) ? refreshOverlaySettingsNameFromJSON(object.overlayName) : 0,
      overlayId: isSet(object.overlayId) ? globalThis.String(object.overlayId) : undefined,
    };
  },

  toJSON(message: RefreshOverlaysRequest): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.overlayName !== 0) {
      obj.overlayName = refreshOverlaySettingsNameToJSON(message.overlayName);
    }
    if (message.overlayId !== undefined) {
      obj.overlayId = message.overlayId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RefreshOverlaysRequest>, I>>(base?: I): RefreshOverlaysRequest {
    return RefreshOverlaysRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RefreshOverlaysRequest>, I>>(object: I): RefreshOverlaysRequest {
    const message = createBaseRefreshOverlaysRequest();
    message.channelId = object.channelId ?? "";
    message.overlayName = object.overlayName ?? 0;
    message.overlayId = object.overlayId ?? undefined;
    return message;
  },
};

function createBaseTriggerKappagenRequest(): TriggerKappagenRequest {
  return { channelId: "", text: "", emotes: [] };
}

export const TriggerKappagenRequest = {
  encode(message: TriggerKappagenRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.text !== "") {
      writer.uint32(18).string(message.text);
    }
    for (const v of message.emotes) {
      TriggerKappagenRequest_Emote.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TriggerKappagenRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTriggerKappagenRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
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

          message.emotes.push(TriggerKappagenRequest_Emote.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TriggerKappagenRequest {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      text: isSet(object.text) ? globalThis.String(object.text) : "",
      emotes: globalThis.Array.isArray(object?.emotes)
        ? object.emotes.map((e: any) => TriggerKappagenRequest_Emote.fromJSON(e))
        : [],
    };
  },

  toJSON(message: TriggerKappagenRequest): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.text !== "") {
      obj.text = message.text;
    }
    if (message.emotes?.length) {
      obj.emotes = message.emotes.map((e) => TriggerKappagenRequest_Emote.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TriggerKappagenRequest>, I>>(base?: I): TriggerKappagenRequest {
    return TriggerKappagenRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TriggerKappagenRequest>, I>>(object: I): TriggerKappagenRequest {
    const message = createBaseTriggerKappagenRequest();
    message.channelId = object.channelId ?? "";
    message.text = object.text ?? "";
    message.emotes = object.emotes?.map((e) => TriggerKappagenRequest_Emote.fromPartial(e)) || [];
    return message;
  },
};

function createBaseTriggerKappagenRequest_Emote(): TriggerKappagenRequest_Emote {
  return { id: "", positions: [] };
}

export const TriggerKappagenRequest_Emote = {
  encode(message: TriggerKappagenRequest_Emote, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    for (const v of message.positions) {
      writer.uint32(18).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TriggerKappagenRequest_Emote {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTriggerKappagenRequest_Emote();
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

          message.positions.push(reader.string());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TriggerKappagenRequest_Emote {
    return {
      id: isSet(object.id) ? globalThis.String(object.id) : "",
      positions: globalThis.Array.isArray(object?.positions)
        ? object.positions.map((e: any) => globalThis.String(e))
        : [],
    };
  },

  toJSON(message: TriggerKappagenRequest_Emote): unknown {
    const obj: any = {};
    if (message.id !== "") {
      obj.id = message.id;
    }
    if (message.positions?.length) {
      obj.positions = message.positions;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TriggerKappagenRequest_Emote>, I>>(base?: I): TriggerKappagenRequest_Emote {
    return TriggerKappagenRequest_Emote.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TriggerKappagenRequest_Emote>, I>>(object: I): TriggerKappagenRequest_Emote {
    const message = createBaseTriggerKappagenRequest_Emote();
    message.id = object.id ?? "";
    message.positions = object.positions?.map((e) => e) || [];
    return message;
  },
};

function createBaseTriggerKappagenByEventRequest(): TriggerKappagenByEventRequest {
  return { channelId: "", event: 0 };
}

export const TriggerKappagenByEventRequest = {
  encode(message: TriggerKappagenByEventRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.event !== 0) {
      writer.uint32(16).int32(message.event);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TriggerKappagenByEventRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTriggerKappagenByEventRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.event = reader.int32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TriggerKappagenByEventRequest {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      event: isSet(object.event) ? globalThis.Number(object.event) : 0,
    };
  },

  toJSON(message: TriggerKappagenByEventRequest): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.event !== 0) {
      obj.event = Math.round(message.event);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TriggerKappagenByEventRequest>, I>>(base?: I): TriggerKappagenByEventRequest {
    return TriggerKappagenByEventRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TriggerKappagenByEventRequest>, I>>(
    object: I,
  ): TriggerKappagenByEventRequest {
    const message = createBaseTriggerKappagenByEventRequest();
    message.channelId = object.channelId ?? "";
    message.event = object.event ?? 0;
    return message;
  },
};

function createBaseTriggerShowBrbRequest(): TriggerShowBrbRequest {
  return { channelId: "", minutes: 0, text: undefined };
}

export const TriggerShowBrbRequest = {
  encode(message: TriggerShowBrbRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.minutes !== 0) {
      writer.uint32(16).int32(message.minutes);
    }
    if (message.text !== undefined) {
      writer.uint32(26).string(message.text);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TriggerShowBrbRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTriggerShowBrbRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.minutes = reader.int32();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.text = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TriggerShowBrbRequest {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      minutes: isSet(object.minutes) ? globalThis.Number(object.minutes) : 0,
      text: isSet(object.text) ? globalThis.String(object.text) : undefined,
    };
  },

  toJSON(message: TriggerShowBrbRequest): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.minutes !== 0) {
      obj.minutes = Math.round(message.minutes);
    }
    if (message.text !== undefined) {
      obj.text = message.text;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TriggerShowBrbRequest>, I>>(base?: I): TriggerShowBrbRequest {
    return TriggerShowBrbRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TriggerShowBrbRequest>, I>>(object: I): TriggerShowBrbRequest {
    const message = createBaseTriggerShowBrbRequest();
    message.channelId = object.channelId ?? "";
    message.minutes = object.minutes ?? 0;
    message.text = object.text ?? undefined;
    return message;
  },
};

function createBaseTriggerHideBrbRequest(): TriggerHideBrbRequest {
  return { channelId: "" };
}

export const TriggerHideBrbRequest = {
  encode(message: TriggerHideBrbRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TriggerHideBrbRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTriggerHideBrbRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TriggerHideBrbRequest {
    return { channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "" };
  },

  toJSON(message: TriggerHideBrbRequest): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TriggerHideBrbRequest>, I>>(base?: I): TriggerHideBrbRequest {
    return TriggerHideBrbRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TriggerHideBrbRequest>, I>>(object: I): TriggerHideBrbRequest {
    const message = createBaseTriggerHideBrbRequest();
    message.channelId = object.channelId ?? "";
    return message;
  },
};

function createBaseDudesJumpRequest(): DudesJumpRequest {
  return { channelId: "", userId: "", userDisplayName: "", userName: "", userColor: "" };
}

export const DudesJumpRequest = {
  encode(message: DudesJumpRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.userId !== "") {
      writer.uint32(18).string(message.userId);
    }
    if (message.userDisplayName !== "") {
      writer.uint32(26).string(message.userDisplayName);
    }
    if (message.userName !== "") {
      writer.uint32(34).string(message.userName);
    }
    if (message.userColor !== "") {
      writer.uint32(42).string(message.userColor);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DudesJumpRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDudesJumpRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.userDisplayName = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.userName = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.userColor = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DudesJumpRequest {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      userId: isSet(object.userId) ? globalThis.String(object.userId) : "",
      userDisplayName: isSet(object.userDisplayName) ? globalThis.String(object.userDisplayName) : "",
      userName: isSet(object.userName) ? globalThis.String(object.userName) : "",
      userColor: isSet(object.userColor) ? globalThis.String(object.userColor) : "",
    };
  },

  toJSON(message: DudesJumpRequest): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.userId !== "") {
      obj.userId = message.userId;
    }
    if (message.userDisplayName !== "") {
      obj.userDisplayName = message.userDisplayName;
    }
    if (message.userName !== "") {
      obj.userName = message.userName;
    }
    if (message.userColor !== "") {
      obj.userColor = message.userColor;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DudesJumpRequest>, I>>(base?: I): DudesJumpRequest {
    return DudesJumpRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DudesJumpRequest>, I>>(object: I): DudesJumpRequest {
    const message = createBaseDudesJumpRequest();
    message.channelId = object.channelId ?? "";
    message.userId = object.userId ?? "";
    message.userDisplayName = object.userDisplayName ?? "";
    message.userName = object.userName ?? "";
    message.userColor = object.userColor ?? "";
    return message;
  },
};

function createBaseDudesUserPunishedRequest(): DudesUserPunishedRequest {
  return { channelId: "", userId: "", userDisplayName: "", userName: "" };
}

export const DudesUserPunishedRequest = {
  encode(message: DudesUserPunishedRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    if (message.userId !== "") {
      writer.uint32(18).string(message.userId);
    }
    if (message.userDisplayName !== "") {
      writer.uint32(26).string(message.userDisplayName);
    }
    if (message.userName !== "") {
      writer.uint32(34).string(message.userName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DudesUserPunishedRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDudesUserPunishedRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.channelId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.userDisplayName = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.userName = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DudesUserPunishedRequest {
    return {
      channelId: isSet(object.channelId) ? globalThis.String(object.channelId) : "",
      userId: isSet(object.userId) ? globalThis.String(object.userId) : "",
      userDisplayName: isSet(object.userDisplayName) ? globalThis.String(object.userDisplayName) : "",
      userName: isSet(object.userName) ? globalThis.String(object.userName) : "",
    };
  },

  toJSON(message: DudesUserPunishedRequest): unknown {
    const obj: any = {};
    if (message.channelId !== "") {
      obj.channelId = message.channelId;
    }
    if (message.userId !== "") {
      obj.userId = message.userId;
    }
    if (message.userDisplayName !== "") {
      obj.userDisplayName = message.userDisplayName;
    }
    if (message.userName !== "") {
      obj.userName = message.userName;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DudesUserPunishedRequest>, I>>(base?: I): DudesUserPunishedRequest {
    return DudesUserPunishedRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DudesUserPunishedRequest>, I>>(object: I): DudesUserPunishedRequest {
    const message = createBaseDudesUserPunishedRequest();
    message.channelId = object.channelId ?? "";
    message.userId = object.userId ?? "";
    message.userDisplayName = object.userDisplayName ?? "";
    message.userName = object.userName ?? "";
    return message;
  },
};

export type WebsocketDefinition = typeof WebsocketDefinition;
export const WebsocketDefinition = {
  name: "Websocket",
  fullName: "websockets.Websocket",
  methods: {
    youtubeAddSongToQueue: {
      name: "YoutubeAddSongToQueue",
      requestType: YoutubeAddSongToQueueRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    youtubeRemoveSongToQueue: {
      name: "YoutubeRemoveSongToQueue",
      requestType: YoutubeRemoveSongFromQueueRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    obsCheckIsUserConnected: {
      name: "ObsCheckIsUserConnected",
      requestType: ObsCheckUserConnectedRequest,
      requestStream: false,
      responseType: ObsCheckUserConnectedResponse,
      responseStream: false,
      options: {},
    },
    obsSetScene: {
      name: "ObsSetScene",
      requestType: ObsSetSceneMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    obsToggleSource: {
      name: "ObsToggleSource",
      requestType: ObsToggleSourceMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    obsToggleAudio: {
      name: "ObsToggleAudio",
      requestType: ObsToggleAudioMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    obsAudioSetVolume: {
      name: "ObsAudioSetVolume",
      requestType: ObsAudioSetVolumeMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    obsAudioIncreaseVolume: {
      name: "ObsAudioIncreaseVolume",
      requestType: ObsAudioIncreaseVolumeMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    obsAudioDecreaseVolume: {
      name: "ObsAudioDecreaseVolume",
      requestType: ObsAudioDecreaseVolumeMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    obsAudioEnable: {
      name: "ObsAudioEnable",
      requestType: ObsAudioDisableOrEnableMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    obsAudioDisable: {
      name: "ObsAudioDisable",
      requestType: ObsAudioDisableOrEnableMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    obsStopStream: {
      name: "ObsStopStream",
      requestType: ObsStopOrStartStream,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    obsStartStream: {
      name: "ObsStartStream",
      requestType: ObsStopOrStartStream,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    textToSpeechSay: {
      name: "TextToSpeechSay",
      requestType: TTSMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    textToSpeechSkip: {
      name: "TextToSpeechSkip",
      requestType: TTSSkipMessage,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    triggerAlert: {
      name: "TriggerAlert",
      requestType: TriggerAlertRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    triggerKappagen: {
      name: "TriggerKappagen",
      requestType: TriggerKappagenRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    triggerKappagenByEvent: {
      name: "TriggerKappagenByEvent",
      requestType: TriggerKappagenByEventRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    triggerShowBrb: {
      name: "TriggerShowBrb",
      requestType: TriggerShowBrbRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    triggerHideBrb: {
      name: "TriggerHideBrb",
      requestType: TriggerHideBrbRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    refreshOverlaySettings: {
      name: "RefreshOverlaySettings",
      requestType: RefreshOverlaysRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    dudesJump: {
      name: "DudesJump",
      requestType: DudesJumpRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    dudesUserPunished: {
      name: "DudesUserPunished",
      requestType: DudesUserPunishedRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface WebsocketServiceImplementation<CallContextExt = {}> {
  youtubeAddSongToQueue(
    request: YoutubeAddSongToQueueRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<Empty>>;
  youtubeRemoveSongToQueue(
    request: YoutubeRemoveSongFromQueueRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<Empty>>;
  obsCheckIsUserConnected(
    request: ObsCheckUserConnectedRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ObsCheckUserConnectedResponse>>;
  obsSetScene(request: ObsSetSceneMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  obsToggleSource(request: ObsToggleSourceMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  obsToggleAudio(request: ObsToggleAudioMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  obsAudioSetVolume(
    request: ObsAudioSetVolumeMessage,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<Empty>>;
  obsAudioIncreaseVolume(
    request: ObsAudioIncreaseVolumeMessage,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<Empty>>;
  obsAudioDecreaseVolume(
    request: ObsAudioDecreaseVolumeMessage,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<Empty>>;
  obsAudioEnable(
    request: ObsAudioDisableOrEnableMessage,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<Empty>>;
  obsAudioDisable(
    request: ObsAudioDisableOrEnableMessage,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<Empty>>;
  obsStopStream(request: ObsStopOrStartStream, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  obsStartStream(request: ObsStopOrStartStream, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  textToSpeechSay(request: TTSMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  textToSpeechSkip(request: TTSSkipMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  triggerAlert(request: TriggerAlertRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  triggerKappagen(request: TriggerKappagenRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  triggerKappagenByEvent(
    request: TriggerKappagenByEventRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<Empty>>;
  triggerShowBrb(request: TriggerShowBrbRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  triggerHideBrb(request: TriggerHideBrbRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  refreshOverlaySettings(
    request: RefreshOverlaysRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<Empty>>;
  dudesJump(request: DudesJumpRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  dudesUserPunished(
    request: DudesUserPunishedRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<Empty>>;
}

export interface WebsocketClient<CallOptionsExt = {}> {
  youtubeAddSongToQueue(
    request: DeepPartial<YoutubeAddSongToQueueRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<Empty>;
  youtubeRemoveSongToQueue(
    request: DeepPartial<YoutubeRemoveSongFromQueueRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<Empty>;
  obsCheckIsUserConnected(
    request: DeepPartial<ObsCheckUserConnectedRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ObsCheckUserConnectedResponse>;
  obsSetScene(request: DeepPartial<ObsSetSceneMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  obsToggleSource(request: DeepPartial<ObsToggleSourceMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  obsToggleAudio(request: DeepPartial<ObsToggleAudioMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  obsAudioSetVolume(
    request: DeepPartial<ObsAudioSetVolumeMessage>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<Empty>;
  obsAudioIncreaseVolume(
    request: DeepPartial<ObsAudioIncreaseVolumeMessage>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<Empty>;
  obsAudioDecreaseVolume(
    request: DeepPartial<ObsAudioDecreaseVolumeMessage>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<Empty>;
  obsAudioEnable(
    request: DeepPartial<ObsAudioDisableOrEnableMessage>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<Empty>;
  obsAudioDisable(
    request: DeepPartial<ObsAudioDisableOrEnableMessage>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<Empty>;
  obsStopStream(request: DeepPartial<ObsStopOrStartStream>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  obsStartStream(request: DeepPartial<ObsStopOrStartStream>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  textToSpeechSay(request: DeepPartial<TTSMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  textToSpeechSkip(request: DeepPartial<TTSSkipMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  triggerAlert(request: DeepPartial<TriggerAlertRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  triggerKappagen(request: DeepPartial<TriggerKappagenRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  triggerKappagenByEvent(
    request: DeepPartial<TriggerKappagenByEventRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<Empty>;
  triggerShowBrb(request: DeepPartial<TriggerShowBrbRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  triggerHideBrb(request: DeepPartial<TriggerHideBrbRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  refreshOverlaySettings(
    request: DeepPartial<RefreshOverlaysRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<Empty>;
  dudesJump(request: DeepPartial<DudesJumpRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  dudesUserPunished(
    request: DeepPartial<DudesUserPunishedRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<Empty>;
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
