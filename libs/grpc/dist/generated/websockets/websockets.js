"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.WebsocketDefinition = exports.TTSSkipMessage = exports.TTSMessage = exports.ObsStopOrStartStream = exports.ObsAudioDisableOrEnableMessage = exports.ObsAudioDecreaseVolumeMessage = exports.ObsAudioIncreaseVolumeMessage = exports.ObsAudioSetVolumeMessage = exports.ObsToggleAudioMessage = exports.ObsToggleSourceMessage = exports.ObsSetSceneMessage = exports.YoutubeRemoveSongFromQueueRequest = exports.YoutubeAddSongToQueueRequest = exports.protobufPackage = void 0;
const minimal_1 = __importDefault(require("protobufjs/minimal"));
const empty_1 = require("./google/protobuf/empty");
exports.protobufPackage = "websockets";
function createBaseYoutubeAddSongToQueueRequest() {
    return { channelId: "", entityId: "" };
}
exports.YoutubeAddSongToQueueRequest = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.channelId !== "") {
            writer.uint32(10).string(message.channelId);
        }
        if (message.entityId !== "") {
            writer.uint32(18).string(message.entityId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            channelId: isSet(object.channelId) ? String(object.channelId) : "",
            entityId: isSet(object.entityId) ? String(object.entityId) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.entityId !== undefined && (obj.entityId = message.entityId);
        return obj;
    },
    create(base) {
        return exports.YoutubeAddSongToQueueRequest.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b;
        const message = createBaseYoutubeAddSongToQueueRequest();
        message.channelId = (_a = object.channelId) !== null && _a !== void 0 ? _a : "";
        message.entityId = (_b = object.entityId) !== null && _b !== void 0 ? _b : "";
        return message;
    },
};
function createBaseYoutubeRemoveSongFromQueueRequest() {
    return { channelId: "", entityId: "" };
}
exports.YoutubeRemoveSongFromQueueRequest = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.channelId !== "") {
            writer.uint32(10).string(message.channelId);
        }
        if (message.entityId !== "") {
            writer.uint32(18).string(message.entityId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            channelId: isSet(object.channelId) ? String(object.channelId) : "",
            entityId: isSet(object.entityId) ? String(object.entityId) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.entityId !== undefined && (obj.entityId = message.entityId);
        return obj;
    },
    create(base) {
        return exports.YoutubeRemoveSongFromQueueRequest.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b;
        const message = createBaseYoutubeRemoveSongFromQueueRequest();
        message.channelId = (_a = object.channelId) !== null && _a !== void 0 ? _a : "";
        message.entityId = (_b = object.entityId) !== null && _b !== void 0 ? _b : "";
        return message;
    },
};
function createBaseObsSetSceneMessage() {
    return { channelId: "", sceneName: "" };
}
exports.ObsSetSceneMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.channelId !== "") {
            writer.uint32(10).string(message.channelId);
        }
        if (message.sceneName !== "") {
            writer.uint32(18).string(message.sceneName);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            channelId: isSet(object.channelId) ? String(object.channelId) : "",
            sceneName: isSet(object.sceneName) ? String(object.sceneName) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.sceneName !== undefined && (obj.sceneName = message.sceneName);
        return obj;
    },
    create(base) {
        return exports.ObsSetSceneMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b;
        const message = createBaseObsSetSceneMessage();
        message.channelId = (_a = object.channelId) !== null && _a !== void 0 ? _a : "";
        message.sceneName = (_b = object.sceneName) !== null && _b !== void 0 ? _b : "";
        return message;
    },
};
function createBaseObsToggleSourceMessage() {
    return { channelId: "", sourceName: "" };
}
exports.ObsToggleSourceMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.channelId !== "") {
            writer.uint32(10).string(message.channelId);
        }
        if (message.sourceName !== "") {
            writer.uint32(18).string(message.sourceName);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            channelId: isSet(object.channelId) ? String(object.channelId) : "",
            sourceName: isSet(object.sourceName) ? String(object.sourceName) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.sourceName !== undefined && (obj.sourceName = message.sourceName);
        return obj;
    },
    create(base) {
        return exports.ObsToggleSourceMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b;
        const message = createBaseObsToggleSourceMessage();
        message.channelId = (_a = object.channelId) !== null && _a !== void 0 ? _a : "";
        message.sourceName = (_b = object.sourceName) !== null && _b !== void 0 ? _b : "";
        return message;
    },
};
function createBaseObsToggleAudioMessage() {
    return { channelId: "", audioSourceName: "" };
}
exports.ObsToggleAudioMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.channelId !== "") {
            writer.uint32(10).string(message.channelId);
        }
        if (message.audioSourceName !== "") {
            writer.uint32(18).string(message.audioSourceName);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            channelId: isSet(object.channelId) ? String(object.channelId) : "",
            audioSourceName: isSet(object.audioSourceName) ? String(object.audioSourceName) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.audioSourceName !== undefined && (obj.audioSourceName = message.audioSourceName);
        return obj;
    },
    create(base) {
        return exports.ObsToggleAudioMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b;
        const message = createBaseObsToggleAudioMessage();
        message.channelId = (_a = object.channelId) !== null && _a !== void 0 ? _a : "";
        message.audioSourceName = (_b = object.audioSourceName) !== null && _b !== void 0 ? _b : "";
        return message;
    },
};
function createBaseObsAudioSetVolumeMessage() {
    return { channelId: "", audioSourceName: "", volume: 0 };
}
exports.ObsAudioSetVolumeMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            channelId: isSet(object.channelId) ? String(object.channelId) : "",
            audioSourceName: isSet(object.audioSourceName) ? String(object.audioSourceName) : "",
            volume: isSet(object.volume) ? Number(object.volume) : 0,
        };
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.audioSourceName !== undefined && (obj.audioSourceName = message.audioSourceName);
        message.volume !== undefined && (obj.volume = Math.round(message.volume));
        return obj;
    },
    create(base) {
        return exports.ObsAudioSetVolumeMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c;
        const message = createBaseObsAudioSetVolumeMessage();
        message.channelId = (_a = object.channelId) !== null && _a !== void 0 ? _a : "";
        message.audioSourceName = (_b = object.audioSourceName) !== null && _b !== void 0 ? _b : "";
        message.volume = (_c = object.volume) !== null && _c !== void 0 ? _c : 0;
        return message;
    },
};
function createBaseObsAudioIncreaseVolumeMessage() {
    return { channelId: "", audioSourceName: "", step: 0 };
}
exports.ObsAudioIncreaseVolumeMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            channelId: isSet(object.channelId) ? String(object.channelId) : "",
            audioSourceName: isSet(object.audioSourceName) ? String(object.audioSourceName) : "",
            step: isSet(object.step) ? Number(object.step) : 0,
        };
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.audioSourceName !== undefined && (obj.audioSourceName = message.audioSourceName);
        message.step !== undefined && (obj.step = Math.round(message.step));
        return obj;
    },
    create(base) {
        return exports.ObsAudioIncreaseVolumeMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c;
        const message = createBaseObsAudioIncreaseVolumeMessage();
        message.channelId = (_a = object.channelId) !== null && _a !== void 0 ? _a : "";
        message.audioSourceName = (_b = object.audioSourceName) !== null && _b !== void 0 ? _b : "";
        message.step = (_c = object.step) !== null && _c !== void 0 ? _c : 0;
        return message;
    },
};
function createBaseObsAudioDecreaseVolumeMessage() {
    return { channelId: "", audioSourceName: "", step: 0 };
}
exports.ObsAudioDecreaseVolumeMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            channelId: isSet(object.channelId) ? String(object.channelId) : "",
            audioSourceName: isSet(object.audioSourceName) ? String(object.audioSourceName) : "",
            step: isSet(object.step) ? Number(object.step) : 0,
        };
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.audioSourceName !== undefined && (obj.audioSourceName = message.audioSourceName);
        message.step !== undefined && (obj.step = Math.round(message.step));
        return obj;
    },
    create(base) {
        return exports.ObsAudioDecreaseVolumeMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c;
        const message = createBaseObsAudioDecreaseVolumeMessage();
        message.channelId = (_a = object.channelId) !== null && _a !== void 0 ? _a : "";
        message.audioSourceName = (_b = object.audioSourceName) !== null && _b !== void 0 ? _b : "";
        message.step = (_c = object.step) !== null && _c !== void 0 ? _c : 0;
        return message;
    },
};
function createBaseObsAudioDisableOrEnableMessage() {
    return { channelId: "", audioSourceName: "" };
}
exports.ObsAudioDisableOrEnableMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.channelId !== "") {
            writer.uint32(10).string(message.channelId);
        }
        if (message.audioSourceName !== "") {
            writer.uint32(18).string(message.audioSourceName);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            channelId: isSet(object.channelId) ? String(object.channelId) : "",
            audioSourceName: isSet(object.audioSourceName) ? String(object.audioSourceName) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.audioSourceName !== undefined && (obj.audioSourceName = message.audioSourceName);
        return obj;
    },
    create(base) {
        return exports.ObsAudioDisableOrEnableMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b;
        const message = createBaseObsAudioDisableOrEnableMessage();
        message.channelId = (_a = object.channelId) !== null && _a !== void 0 ? _a : "";
        message.audioSourceName = (_b = object.audioSourceName) !== null && _b !== void 0 ? _b : "";
        return message;
    },
};
function createBaseObsStopOrStartStream() {
    return { channelId: "" };
}
exports.ObsStopOrStartStream = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.channelId !== "") {
            writer.uint32(10).string(message.channelId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return { channelId: isSet(object.channelId) ? String(object.channelId) : "" };
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        return obj;
    },
    create(base) {
        return exports.ObsStopOrStartStream.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a;
        const message = createBaseObsStopOrStartStream();
        message.channelId = (_a = object.channelId) !== null && _a !== void 0 ? _a : "";
        return message;
    },
};
function createBaseTTSMessage() {
    return { channelId: "", text: "", voice: "", rate: "", pitch: "", volume: "" };
}
exports.TTSMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            channelId: isSet(object.channelId) ? String(object.channelId) : "",
            text: isSet(object.text) ? String(object.text) : "",
            voice: isSet(object.voice) ? String(object.voice) : "",
            rate: isSet(object.rate) ? String(object.rate) : "",
            pitch: isSet(object.pitch) ? String(object.pitch) : "",
            volume: isSet(object.volume) ? String(object.volume) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.text !== undefined && (obj.text = message.text);
        message.voice !== undefined && (obj.voice = message.voice);
        message.rate !== undefined && (obj.rate = message.rate);
        message.pitch !== undefined && (obj.pitch = message.pitch);
        message.volume !== undefined && (obj.volume = message.volume);
        return obj;
    },
    create(base) {
        return exports.TTSMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d, _e, _f;
        const message = createBaseTTSMessage();
        message.channelId = (_a = object.channelId) !== null && _a !== void 0 ? _a : "";
        message.text = (_b = object.text) !== null && _b !== void 0 ? _b : "";
        message.voice = (_c = object.voice) !== null && _c !== void 0 ? _c : "";
        message.rate = (_d = object.rate) !== null && _d !== void 0 ? _d : "";
        message.pitch = (_e = object.pitch) !== null && _e !== void 0 ? _e : "";
        message.volume = (_f = object.volume) !== null && _f !== void 0 ? _f : "";
        return message;
    },
};
function createBaseTTSSkipMessage() {
    return { channelId: "" };
}
exports.TTSSkipMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.channelId !== "") {
            writer.uint32(10).string(message.channelId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return { channelId: isSet(object.channelId) ? String(object.channelId) : "" };
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        return obj;
    },
    create(base) {
        return exports.TTSSkipMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a;
        const message = createBaseTTSSkipMessage();
        message.channelId = (_a = object.channelId) !== null && _a !== void 0 ? _a : "";
        return message;
    },
};
exports.WebsocketDefinition = {
    name: "Websocket",
    fullName: "websockets.Websocket",
    methods: {
        youtubeAddSongToQueue: {
            name: "YoutubeAddSongToQueue",
            requestType: exports.YoutubeAddSongToQueueRequest,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        youtubeRemoveSongToQueue: {
            name: "YoutubeRemoveSongToQueue",
            requestType: exports.YoutubeRemoveSongFromQueueRequest,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        obsSetScene: {
            name: "ObsSetScene",
            requestType: exports.ObsSetSceneMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        obsToggleSource: {
            name: "ObsToggleSource",
            requestType: exports.ObsToggleSourceMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        obsToggleAudio: {
            name: "ObsToggleAudio",
            requestType: exports.ObsToggleAudioMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        obsAudioSetVolume: {
            name: "ObsAudioSetVolume",
            requestType: exports.ObsAudioSetVolumeMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        obsAudioIncreaseVolume: {
            name: "ObsAudioIncreaseVolume",
            requestType: exports.ObsAudioIncreaseVolumeMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        obsAudioDecreaseVolume: {
            name: "ObsAudioDecreaseVolume",
            requestType: exports.ObsAudioDecreaseVolumeMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        obsAudioEnable: {
            name: "ObsAudioEnable",
            requestType: exports.ObsAudioDisableOrEnableMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        obsAudioDisable: {
            name: "ObsAudioDisable",
            requestType: exports.ObsAudioDisableOrEnableMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        obsStopStream: {
            name: "ObsStopStream",
            requestType: exports.ObsStopOrStartStream,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        obsStartStream: {
            name: "ObsStartStream",
            requestType: exports.ObsStopOrStartStream,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        textToSpeechSay: {
            name: "TextToSpeechSay",
            requestType: exports.TTSMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        textToSpeechSkip: {
            name: "TextToSpeechSkip",
            requestType: exports.TTSSkipMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
    },
};
function isSet(value) {
    return value !== null && value !== undefined;
}
//# sourceMappingURL=websockets.js.map