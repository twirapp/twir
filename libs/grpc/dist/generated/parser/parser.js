"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.ParserDefinition = exports.ParseTextResponseData = exports.ParseTextRequestData = exports.GetDefaultCommandsResponse_DefaultCommand = exports.GetDefaultCommandsResponse = exports.GetVariablesResponse_Variable = exports.GetVariablesResponse = exports.ProcessCommandResponse = exports.ProcessCommandRequest = exports.Message_Emote = exports.Message_EmotePosition = exports.Message = exports.Channel = exports.Sender = exports.protobufPackage = void 0;
const long_1 = __importDefault(require("long"));
const minimal_1 = __importDefault(require("protobufjs/minimal"));
const empty_1 = require("./google/protobuf/empty");
exports.protobufPackage = "parser";
function createBaseSender() {
    return { id: "", name: "", displayName: "", badges: [] };
}
exports.Sender = {
    encode(message, writer = minimal_1.default.Writer.create()) {
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
            writer.uint32(34).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            id: isSet(object.id) ? String(object.id) : "",
            name: isSet(object.name) ? String(object.name) : "",
            displayName: isSet(object.displayName) ? String(object.displayName) : "",
            badges: Array.isArray(object === null || object === void 0 ? void 0 : object.badges) ? object.badges.map((e) => String(e)) : [],
        };
    },
    toJSON(message) {
        const obj = {};
        message.id !== undefined && (obj.id = message.id);
        message.name !== undefined && (obj.name = message.name);
        message.displayName !== undefined && (obj.displayName = message.displayName);
        if (message.badges) {
            obj.badges = message.badges.map((e) => e);
        }
        else {
            obj.badges = [];
        }
        return obj;
    },
    create(base) {
        return exports.Sender.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d;
        const message = createBaseSender();
        message.id = (_a = object.id) !== null && _a !== void 0 ? _a : "";
        message.name = (_b = object.name) !== null && _b !== void 0 ? _b : "";
        message.displayName = (_c = object.displayName) !== null && _c !== void 0 ? _c : "";
        message.badges = ((_d = object.badges) === null || _d === void 0 ? void 0 : _d.map((e) => e)) || [];
        return message;
    },
};
function createBaseChannel() {
    return { id: "", name: "" };
}
exports.Channel = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.id !== "") {
            writer.uint32(10).string(message.id);
        }
        if (message.name !== "") {
            writer.uint32(18).string(message.name);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return { id: isSet(object.id) ? String(object.id) : "", name: isSet(object.name) ? String(object.name) : "" };
    },
    toJSON(message) {
        const obj = {};
        message.id !== undefined && (obj.id = message.id);
        message.name !== undefined && (obj.name = message.name);
        return obj;
    },
    create(base) {
        return exports.Channel.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b;
        const message = createBaseChannel();
        message.id = (_a = object.id) !== null && _a !== void 0 ? _a : "";
        message.name = (_b = object.name) !== null && _b !== void 0 ? _b : "";
        return message;
    },
};
function createBaseMessage() {
    return { text: "", id: "", emotes: [] };
}
exports.Message = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.text !== "") {
            writer.uint32(10).string(message.text);
        }
        if (message.id !== "") {
            writer.uint32(18).string(message.id);
        }
        for (const v of message.emotes) {
            exports.Message_Emote.encode(v, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
                    message.emotes.push(exports.Message_Emote.decode(reader, reader.uint32()));
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
            text: isSet(object.text) ? String(object.text) : "",
            id: isSet(object.id) ? String(object.id) : "",
            emotes: Array.isArray(object === null || object === void 0 ? void 0 : object.emotes) ? object.emotes.map((e) => exports.Message_Emote.fromJSON(e)) : [],
        };
    },
    toJSON(message) {
        const obj = {};
        message.text !== undefined && (obj.text = message.text);
        message.id !== undefined && (obj.id = message.id);
        if (message.emotes) {
            obj.emotes = message.emotes.map((e) => e ? exports.Message_Emote.toJSON(e) : undefined);
        }
        else {
            obj.emotes = [];
        }
        return obj;
    },
    create(base) {
        return exports.Message.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c;
        const message = createBaseMessage();
        message.text = (_a = object.text) !== null && _a !== void 0 ? _a : "";
        message.id = (_b = object.id) !== null && _b !== void 0 ? _b : "";
        message.emotes = ((_c = object.emotes) === null || _c === void 0 ? void 0 : _c.map((e) => exports.Message_Emote.fromPartial(e))) || [];
        return message;
    },
};
function createBaseMessage_EmotePosition() {
    return { start: 0, end: 0 };
}
exports.Message_EmotePosition = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.start !== 0) {
            writer.uint32(8).int64(message.start);
        }
        if (message.end !== 0) {
            writer.uint32(16).int64(message.end);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseMessage_EmotePosition();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.start = longToNumber(reader.int64());
                    break;
                case 2:
                    message.end = longToNumber(reader.int64());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        return { start: isSet(object.start) ? Number(object.start) : 0, end: isSet(object.end) ? Number(object.end) : 0 };
    },
    toJSON(message) {
        const obj = {};
        message.start !== undefined && (obj.start = Math.round(message.start));
        message.end !== undefined && (obj.end = Math.round(message.end));
        return obj;
    },
    create(base) {
        return exports.Message_EmotePosition.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b;
        const message = createBaseMessage_EmotePosition();
        message.start = (_a = object.start) !== null && _a !== void 0 ? _a : 0;
        message.end = (_b = object.end) !== null && _b !== void 0 ? _b : 0;
        return message;
    },
};
function createBaseMessage_Emote() {
    return { name: "", id: "", count: 0, positions: [] };
}
exports.Message_Emote = {
    encode(message, writer = minimal_1.default.Writer.create()) {
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
            exports.Message_EmotePosition.encode(v, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
                    message.count = longToNumber(reader.int64());
                    break;
                case 4:
                    message.positions.push(exports.Message_EmotePosition.decode(reader, reader.uint32()));
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
            name: isSet(object.name) ? String(object.name) : "",
            id: isSet(object.id) ? String(object.id) : "",
            count: isSet(object.count) ? Number(object.count) : 0,
            positions: Array.isArray(object === null || object === void 0 ? void 0 : object.positions)
                ? object.positions.map((e) => exports.Message_EmotePosition.fromJSON(e))
                : [],
        };
    },
    toJSON(message) {
        const obj = {};
        message.name !== undefined && (obj.name = message.name);
        message.id !== undefined && (obj.id = message.id);
        message.count !== undefined && (obj.count = Math.round(message.count));
        if (message.positions) {
            obj.positions = message.positions.map((e) => e ? exports.Message_EmotePosition.toJSON(e) : undefined);
        }
        else {
            obj.positions = [];
        }
        return obj;
    },
    create(base) {
        return exports.Message_Emote.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d;
        const message = createBaseMessage_Emote();
        message.name = (_a = object.name) !== null && _a !== void 0 ? _a : "";
        message.id = (_b = object.id) !== null && _b !== void 0 ? _b : "";
        message.count = (_c = object.count) !== null && _c !== void 0 ? _c : 0;
        message.positions = ((_d = object.positions) === null || _d === void 0 ? void 0 : _d.map((e) => exports.Message_EmotePosition.fromPartial(e))) || [];
        return message;
    },
};
function createBaseProcessCommandRequest() {
    return { sender: undefined, channel: undefined, message: undefined };
}
exports.ProcessCommandRequest = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.sender !== undefined) {
            exports.Sender.encode(message.sender, writer.uint32(10).fork()).ldelim();
        }
        if (message.channel !== undefined) {
            exports.Channel.encode(message.channel, writer.uint32(18).fork()).ldelim();
        }
        if (message.message !== undefined) {
            exports.Message.encode(message.message, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseProcessCommandRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = exports.Sender.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.channel = exports.Channel.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.message = exports.Message.decode(reader, reader.uint32());
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
            sender: isSet(object.sender) ? exports.Sender.fromJSON(object.sender) : undefined,
            channel: isSet(object.channel) ? exports.Channel.fromJSON(object.channel) : undefined,
            message: isSet(object.message) ? exports.Message.fromJSON(object.message) : undefined,
        };
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined && (obj.sender = message.sender ? exports.Sender.toJSON(message.sender) : undefined);
        message.channel !== undefined && (obj.channel = message.channel ? exports.Channel.toJSON(message.channel) : undefined);
        message.message !== undefined && (obj.message = message.message ? exports.Message.toJSON(message.message) : undefined);
        return obj;
    },
    create(base) {
        return exports.ProcessCommandRequest.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        const message = createBaseProcessCommandRequest();
        message.sender = (object.sender !== undefined && object.sender !== null)
            ? exports.Sender.fromPartial(object.sender)
            : undefined;
        message.channel = (object.channel !== undefined && object.channel !== null)
            ? exports.Channel.fromPartial(object.channel)
            : undefined;
        message.message = (object.message !== undefined && object.message !== null)
            ? exports.Message.fromPartial(object.message)
            : undefined;
        return message;
    },
};
function createBaseProcessCommandResponse() {
    return { responses: [], isReply: false, keepOrder: undefined };
}
exports.ProcessCommandResponse = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        for (const v of message.responses) {
            writer.uint32(10).string(v);
        }
        if (message.isReply === true) {
            writer.uint32(16).bool(message.isReply);
        }
        if (message.keepOrder !== undefined) {
            writer.uint32(24).bool(message.keepOrder);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            responses: Array.isArray(object === null || object === void 0 ? void 0 : object.responses) ? object.responses.map((e) => String(e)) : [],
            isReply: isSet(object.isReply) ? Boolean(object.isReply) : false,
            keepOrder: isSet(object.keepOrder) ? Boolean(object.keepOrder) : undefined,
        };
    },
    toJSON(message) {
        const obj = {};
        if (message.responses) {
            obj.responses = message.responses.map((e) => e);
        }
        else {
            obj.responses = [];
        }
        message.isReply !== undefined && (obj.isReply = message.isReply);
        message.keepOrder !== undefined && (obj.keepOrder = message.keepOrder);
        return obj;
    },
    create(base) {
        return exports.ProcessCommandResponse.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c;
        const message = createBaseProcessCommandResponse();
        message.responses = ((_a = object.responses) === null || _a === void 0 ? void 0 : _a.map((e) => e)) || [];
        message.isReply = (_b = object.isReply) !== null && _b !== void 0 ? _b : false;
        message.keepOrder = (_c = object.keepOrder) !== null && _c !== void 0 ? _c : undefined;
        return message;
    },
};
function createBaseGetVariablesResponse() {
    return { list: [] };
}
exports.GetVariablesResponse = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        for (const v of message.list) {
            exports.GetVariablesResponse_Variable.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetVariablesResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.list.push(exports.GetVariablesResponse_Variable.decode(reader, reader.uint32()));
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
            list: Array.isArray(object === null || object === void 0 ? void 0 : object.list) ? object.list.map((e) => exports.GetVariablesResponse_Variable.fromJSON(e)) : [],
        };
    },
    toJSON(message) {
        const obj = {};
        if (message.list) {
            obj.list = message.list.map((e) => e ? exports.GetVariablesResponse_Variable.toJSON(e) : undefined);
        }
        else {
            obj.list = [];
        }
        return obj;
    },
    create(base) {
        return exports.GetVariablesResponse.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a;
        const message = createBaseGetVariablesResponse();
        message.list = ((_a = object.list) === null || _a === void 0 ? void 0 : _a.map((e) => exports.GetVariablesResponse_Variable.fromPartial(e))) || [];
        return message;
    },
};
function createBaseGetVariablesResponse_Variable() {
    return { name: "", example: "", description: "", visible: false };
}
exports.GetVariablesResponse_Variable = {
    encode(message, writer = minimal_1.default.Writer.create()) {
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            name: isSet(object.name) ? String(object.name) : "",
            example: isSet(object.example) ? String(object.example) : "",
            description: isSet(object.description) ? String(object.description) : "",
            visible: isSet(object.visible) ? Boolean(object.visible) : false,
        };
    },
    toJSON(message) {
        const obj = {};
        message.name !== undefined && (obj.name = message.name);
        message.example !== undefined && (obj.example = message.example);
        message.description !== undefined && (obj.description = message.description);
        message.visible !== undefined && (obj.visible = message.visible);
        return obj;
    },
    create(base) {
        return exports.GetVariablesResponse_Variable.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d;
        const message = createBaseGetVariablesResponse_Variable();
        message.name = (_a = object.name) !== null && _a !== void 0 ? _a : "";
        message.example = (_b = object.example) !== null && _b !== void 0 ? _b : "";
        message.description = (_c = object.description) !== null && _c !== void 0 ? _c : "";
        message.visible = (_d = object.visible) !== null && _d !== void 0 ? _d : false;
        return message;
    },
};
function createBaseGetDefaultCommandsResponse() {
    return { list: [] };
}
exports.GetDefaultCommandsResponse = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        for (const v of message.list) {
            exports.GetDefaultCommandsResponse_DefaultCommand.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetDefaultCommandsResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.list.push(exports.GetDefaultCommandsResponse_DefaultCommand.decode(reader, reader.uint32()));
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
            list: Array.isArray(object === null || object === void 0 ? void 0 : object.list)
                ? object.list.map((e) => exports.GetDefaultCommandsResponse_DefaultCommand.fromJSON(e))
                : [],
        };
    },
    toJSON(message) {
        const obj = {};
        if (message.list) {
            obj.list = message.list.map((e) => e ? exports.GetDefaultCommandsResponse_DefaultCommand.toJSON(e) : undefined);
        }
        else {
            obj.list = [];
        }
        return obj;
    },
    create(base) {
        return exports.GetDefaultCommandsResponse.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a;
        const message = createBaseGetDefaultCommandsResponse();
        message.list = ((_a = object.list) === null || _a === void 0 ? void 0 : _a.map((e) => exports.GetDefaultCommandsResponse_DefaultCommand.fromPartial(e))) || [];
        return message;
    },
};
function createBaseGetDefaultCommandsResponse_DefaultCommand() {
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
exports.GetDefaultCommandsResponse_DefaultCommand = {
    encode(message, writer = minimal_1.default.Writer.create()) {
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
            writer.uint32(34).string(v);
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
            writer.uint32(66).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            name: isSet(object.name) ? String(object.name) : "",
            description: isSet(object.description) ? String(object.description) : "",
            visible: isSet(object.visible) ? Boolean(object.visible) : false,
            rolesNames: Array.isArray(object === null || object === void 0 ? void 0 : object.rolesNames) ? object.rolesNames.map((e) => String(e)) : [],
            module: isSet(object.module) ? String(object.module) : "",
            isReply: isSet(object.isReply) ? Boolean(object.isReply) : false,
            keepResponsesOrder: isSet(object.keepResponsesOrder) ? Boolean(object.keepResponsesOrder) : false,
            aliases: Array.isArray(object === null || object === void 0 ? void 0 : object.aliases) ? object.aliases.map((e) => String(e)) : [],
        };
    },
    toJSON(message) {
        const obj = {};
        message.name !== undefined && (obj.name = message.name);
        message.description !== undefined && (obj.description = message.description);
        message.visible !== undefined && (obj.visible = message.visible);
        if (message.rolesNames) {
            obj.rolesNames = message.rolesNames.map((e) => e);
        }
        else {
            obj.rolesNames = [];
        }
        message.module !== undefined && (obj.module = message.module);
        message.isReply !== undefined && (obj.isReply = message.isReply);
        message.keepResponsesOrder !== undefined && (obj.keepResponsesOrder = message.keepResponsesOrder);
        if (message.aliases) {
            obj.aliases = message.aliases.map((e) => e);
        }
        else {
            obj.aliases = [];
        }
        return obj;
    },
    create(base) {
        return exports.GetDefaultCommandsResponse_DefaultCommand.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d, _e, _f, _g, _h;
        const message = createBaseGetDefaultCommandsResponse_DefaultCommand();
        message.name = (_a = object.name) !== null && _a !== void 0 ? _a : "";
        message.description = (_b = object.description) !== null && _b !== void 0 ? _b : "";
        message.visible = (_c = object.visible) !== null && _c !== void 0 ? _c : false;
        message.rolesNames = ((_d = object.rolesNames) === null || _d === void 0 ? void 0 : _d.map((e) => e)) || [];
        message.module = (_e = object.module) !== null && _e !== void 0 ? _e : "";
        message.isReply = (_f = object.isReply) !== null && _f !== void 0 ? _f : false;
        message.keepResponsesOrder = (_g = object.keepResponsesOrder) !== null && _g !== void 0 ? _g : false;
        message.aliases = ((_h = object.aliases) === null || _h === void 0 ? void 0 : _h.map((e) => e)) || [];
        return message;
    },
};
function createBaseParseTextRequestData() {
    return { sender: undefined, channel: undefined, message: undefined, parseVariables: undefined };
}
exports.ParseTextRequestData = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.sender !== undefined) {
            exports.Sender.encode(message.sender, writer.uint32(10).fork()).ldelim();
        }
        if (message.channel !== undefined) {
            exports.Channel.encode(message.channel, writer.uint32(18).fork()).ldelim();
        }
        if (message.message !== undefined) {
            exports.Message.encode(message.message, writer.uint32(26).fork()).ldelim();
        }
        if (message.parseVariables !== undefined) {
            writer.uint32(32).bool(message.parseVariables);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseParseTextRequestData();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = exports.Sender.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.channel = exports.Channel.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.message = exports.Message.decode(reader, reader.uint32());
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
    fromJSON(object) {
        return {
            sender: isSet(object.sender) ? exports.Sender.fromJSON(object.sender) : undefined,
            channel: isSet(object.channel) ? exports.Channel.fromJSON(object.channel) : undefined,
            message: isSet(object.message) ? exports.Message.fromJSON(object.message) : undefined,
            parseVariables: isSet(object.parseVariables) ? Boolean(object.parseVariables) : undefined,
        };
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined && (obj.sender = message.sender ? exports.Sender.toJSON(message.sender) : undefined);
        message.channel !== undefined && (obj.channel = message.channel ? exports.Channel.toJSON(message.channel) : undefined);
        message.message !== undefined && (obj.message = message.message ? exports.Message.toJSON(message.message) : undefined);
        message.parseVariables !== undefined && (obj.parseVariables = message.parseVariables);
        return obj;
    },
    create(base) {
        return exports.ParseTextRequestData.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a;
        const message = createBaseParseTextRequestData();
        message.sender = (object.sender !== undefined && object.sender !== null)
            ? exports.Sender.fromPartial(object.sender)
            : undefined;
        message.channel = (object.channel !== undefined && object.channel !== null)
            ? exports.Channel.fromPartial(object.channel)
            : undefined;
        message.message = (object.message !== undefined && object.message !== null)
            ? exports.Message.fromPartial(object.message)
            : undefined;
        message.parseVariables = (_a = object.parseVariables) !== null && _a !== void 0 ? _a : undefined;
        return message;
    },
};
function createBaseParseTextResponseData() {
    return { responses: [] };
}
exports.ParseTextResponseData = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        for (const v of message.responses) {
            writer.uint32(10).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return { responses: Array.isArray(object === null || object === void 0 ? void 0 : object.responses) ? object.responses.map((e) => String(e)) : [] };
    },
    toJSON(message) {
        const obj = {};
        if (message.responses) {
            obj.responses = message.responses.map((e) => e);
        }
        else {
            obj.responses = [];
        }
        return obj;
    },
    create(base) {
        return exports.ParseTextResponseData.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a;
        const message = createBaseParseTextResponseData();
        message.responses = ((_a = object.responses) === null || _a === void 0 ? void 0 : _a.map((e) => e)) || [];
        return message;
    },
};
exports.ParserDefinition = {
    name: "Parser",
    fullName: "parser.Parser",
    methods: {
        processCommand: {
            name: "ProcessCommand",
            requestType: exports.ProcessCommandRequest,
            requestStream: false,
            responseType: exports.ProcessCommandResponse,
            responseStream: false,
            options: {},
        },
        parseTextResponse: {
            name: "ParseTextResponse",
            requestType: exports.ParseTextRequestData,
            requestStream: false,
            responseType: exports.ParseTextResponseData,
            responseStream: false,
            options: {},
        },
        getDefaultCommands: {
            name: "GetDefaultCommands",
            requestType: empty_1.Empty,
            requestStream: false,
            responseType: exports.GetDefaultCommandsResponse,
            responseStream: false,
            options: {},
        },
        getDefaultVariables: {
            name: "GetDefaultVariables",
            requestType: empty_1.Empty,
            requestStream: false,
            responseType: exports.GetVariablesResponse,
            responseStream: false,
            options: {},
        },
    },
};
var tsProtoGlobalThis = (() => {
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
function longToNumber(long) {
    if (long.gt(Number.MAX_SAFE_INTEGER)) {
        throw new tsProtoGlobalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
    }
    return long.toNumber();
}
if (minimal_1.default.util.Long !== long_1.default) {
    minimal_1.default.util.Long = long_1.default;
    minimal_1.default.configure();
}
function isSet(value) {
    return value !== null && value !== undefined;
}
//# sourceMappingURL=parser.js.map