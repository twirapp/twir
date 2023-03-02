/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "websockets";

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
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseYoutubeAddSongToQueueRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.channelId = reader.string();
          break;
        case 2:
          message.entityId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): YoutubeAddSongToQueueRequest {
    return {
      channelId: isSet(object.channelId) ? String(object.channelId) : "",
      entityId: isSet(object.entityId) ? String(object.entityId) : "",
    };
  },

  toJSON(message: YoutubeAddSongToQueueRequest): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    message.entityId !== undefined && (obj.entityId = message.entityId);
    return obj;
  },

  create(base?: DeepPartial<YoutubeAddSongToQueueRequest>): YoutubeAddSongToQueueRequest {
    return YoutubeAddSongToQueueRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<YoutubeAddSongToQueueRequest>): YoutubeAddSongToQueueRequest {
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
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseYoutubeRemoveSongFromQueueRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.channelId = reader.string();
          break;
        case 2:
          message.entityId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): YoutubeRemoveSongFromQueueRequest {
    return {
      channelId: isSet(object.channelId) ? String(object.channelId) : "",
      entityId: isSet(object.entityId) ? String(object.entityId) : "",
    };
  },

  toJSON(message: YoutubeRemoveSongFromQueueRequest): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    message.entityId !== undefined && (obj.entityId = message.entityId);
    return obj;
  },

  create(base?: DeepPartial<YoutubeRemoveSongFromQueueRequest>): YoutubeRemoveSongFromQueueRequest {
    return YoutubeRemoveSongFromQueueRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<YoutubeRemoveSongFromQueueRequest>): YoutubeRemoveSongFromQueueRequest {
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
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsSetSceneMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.channelId = reader.string();
          break;
        case 2:
          message.sceneName = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ObsSetSceneMessage {
    return {
      channelId: isSet(object.channelId) ? String(object.channelId) : "",
      sceneName: isSet(object.sceneName) ? String(object.sceneName) : "",
    };
  },

  toJSON(message: ObsSetSceneMessage): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    message.sceneName !== undefined && (obj.sceneName = message.sceneName);
    return obj;
  },

  create(base?: DeepPartial<ObsSetSceneMessage>): ObsSetSceneMessage {
    return ObsSetSceneMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ObsSetSceneMessage>): ObsSetSceneMessage {
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
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsToggleSourceMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.channelId = reader.string();
          break;
        case 2:
          message.sourceName = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ObsToggleSourceMessage {
    return {
      channelId: isSet(object.channelId) ? String(object.channelId) : "",
      sourceName: isSet(object.sourceName) ? String(object.sourceName) : "",
    };
  },

  toJSON(message: ObsToggleSourceMessage): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    message.sourceName !== undefined && (obj.sourceName = message.sourceName);
    return obj;
  },

  create(base?: DeepPartial<ObsToggleSourceMessage>): ObsToggleSourceMessage {
    return ObsToggleSourceMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ObsToggleSourceMessage>): ObsToggleSourceMessage {
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
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsToggleAudioMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.channelId = reader.string();
          break;
        case 2:
          message.audioSourceName = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ObsToggleAudioMessage {
    return {
      channelId: isSet(object.channelId) ? String(object.channelId) : "",
      audioSourceName: isSet(object.audioSourceName) ? String(object.audioSourceName) : "",
    };
  },

  toJSON(message: ObsToggleAudioMessage): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    message.audioSourceName !== undefined && (obj.audioSourceName = message.audioSourceName);
    return obj;
  },

  create(base?: DeepPartial<ObsToggleAudioMessage>): ObsToggleAudioMessage {
    return ObsToggleAudioMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ObsToggleAudioMessage>): ObsToggleAudioMessage {
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
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsAudioSetVolumeMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.channelId = reader.string();
          break;
        case 2:
          message.audioSourceName = reader.string();
          break;
        case 3:
          message.volume = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ObsAudioSetVolumeMessage {
    return {
      channelId: isSet(object.channelId) ? String(object.channelId) : "",
      audioSourceName: isSet(object.audioSourceName) ? String(object.audioSourceName) : "",
      volume: isSet(object.volume) ? Number(object.volume) : 0,
    };
  },

  toJSON(message: ObsAudioSetVolumeMessage): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    message.audioSourceName !== undefined && (obj.audioSourceName = message.audioSourceName);
    message.volume !== undefined && (obj.volume = Math.round(message.volume));
    return obj;
  },

  create(base?: DeepPartial<ObsAudioSetVolumeMessage>): ObsAudioSetVolumeMessage {
    return ObsAudioSetVolumeMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ObsAudioSetVolumeMessage>): ObsAudioSetVolumeMessage {
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
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsAudioIncreaseVolumeMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.channelId = reader.string();
          break;
        case 2:
          message.audioSourceName = reader.string();
          break;
        case 3:
          message.step = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ObsAudioIncreaseVolumeMessage {
    return {
      channelId: isSet(object.channelId) ? String(object.channelId) : "",
      audioSourceName: isSet(object.audioSourceName) ? String(object.audioSourceName) : "",
      step: isSet(object.step) ? Number(object.step) : 0,
    };
  },

  toJSON(message: ObsAudioIncreaseVolumeMessage): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    message.audioSourceName !== undefined && (obj.audioSourceName = message.audioSourceName);
    message.step !== undefined && (obj.step = Math.round(message.step));
    return obj;
  },

  create(base?: DeepPartial<ObsAudioIncreaseVolumeMessage>): ObsAudioIncreaseVolumeMessage {
    return ObsAudioIncreaseVolumeMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ObsAudioIncreaseVolumeMessage>): ObsAudioIncreaseVolumeMessage {
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
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsAudioDecreaseVolumeMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.channelId = reader.string();
          break;
        case 2:
          message.audioSourceName = reader.string();
          break;
        case 3:
          message.step = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ObsAudioDecreaseVolumeMessage {
    return {
      channelId: isSet(object.channelId) ? String(object.channelId) : "",
      audioSourceName: isSet(object.audioSourceName) ? String(object.audioSourceName) : "",
      step: isSet(object.step) ? Number(object.step) : 0,
    };
  },

  toJSON(message: ObsAudioDecreaseVolumeMessage): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    message.audioSourceName !== undefined && (obj.audioSourceName = message.audioSourceName);
    message.step !== undefined && (obj.step = Math.round(message.step));
    return obj;
  },

  create(base?: DeepPartial<ObsAudioDecreaseVolumeMessage>): ObsAudioDecreaseVolumeMessage {
    return ObsAudioDecreaseVolumeMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ObsAudioDecreaseVolumeMessage>): ObsAudioDecreaseVolumeMessage {
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
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsAudioDisableOrEnableMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.channelId = reader.string();
          break;
        case 2:
          message.audioSourceName = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ObsAudioDisableOrEnableMessage {
    return {
      channelId: isSet(object.channelId) ? String(object.channelId) : "",
      audioSourceName: isSet(object.audioSourceName) ? String(object.audioSourceName) : "",
    };
  },

  toJSON(message: ObsAudioDisableOrEnableMessage): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    message.audioSourceName !== undefined && (obj.audioSourceName = message.audioSourceName);
    return obj;
  },

  create(base?: DeepPartial<ObsAudioDisableOrEnableMessage>): ObsAudioDisableOrEnableMessage {
    return ObsAudioDisableOrEnableMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ObsAudioDisableOrEnableMessage>): ObsAudioDisableOrEnableMessage {
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
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseObsStopOrStartStream();
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

  fromJSON(object: any): ObsStopOrStartStream {
    return { channelId: isSet(object.channelId) ? String(object.channelId) : "" };
  },

  toJSON(message: ObsStopOrStartStream): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    return obj;
  },

  create(base?: DeepPartial<ObsStopOrStartStream>): ObsStopOrStartStream {
    return ObsStopOrStartStream.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ObsStopOrStartStream>): ObsStopOrStartStream {
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
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTTSMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.channelId = reader.string();
          break;
        case 2:
          message.text = reader.string();
          break;
        case 3:
          message.voice = reader.string();
          break;
        case 4:
          message.rate = reader.string();
          break;
        case 5:
          message.pitch = reader.string();
          break;
        case 6:
          message.volume = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): TTSMessage {
    return {
      channelId: isSet(object.channelId) ? String(object.channelId) : "",
      text: isSet(object.text) ? String(object.text) : "",
      voice: isSet(object.voice) ? String(object.voice) : "",
      rate: isSet(object.rate) ? String(object.rate) : "",
      pitch: isSet(object.pitch) ? String(object.pitch) : "",
      volume: isSet(object.volume) ? String(object.volume) : "",
    };
  },

  toJSON(message: TTSMessage): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    message.text !== undefined && (obj.text = message.text);
    message.voice !== undefined && (obj.voice = message.voice);
    message.rate !== undefined && (obj.rate = message.rate);
    message.pitch !== undefined && (obj.pitch = message.pitch);
    message.volume !== undefined && (obj.volume = message.volume);
    return obj;
  },

  create(base?: DeepPartial<TTSMessage>): TTSMessage {
    return TTSMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<TTSMessage>): TTSMessage {
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
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTTSSkipMessage();
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

  fromJSON(object: any): TTSSkipMessage {
    return { channelId: isSet(object.channelId) ? String(object.channelId) : "" };
  },

  toJSON(message: TTSSkipMessage): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    return obj;
  },

  create(base?: DeepPartial<TTSSkipMessage>): TTSSkipMessage {
    return TTSSkipMessage.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<TTSSkipMessage>): TTSSkipMessage {
    const message = createBaseTTSSkipMessage();
    message.channelId = object.channelId ?? "";
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
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
