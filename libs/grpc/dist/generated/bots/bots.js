"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.BotsDefinition = exports.JoinOrLeaveRequest = exports.SendMessageRequest = exports.DeleteMessagesRequest = exports.protobufPackage = void 0;
const minimal_1 = __importDefault(require("protobufjs/minimal"));
const empty_1 = require("./google/protobuf/empty");
exports.protobufPackage = "bots";
function createBaseDeleteMessagesRequest() {
    return { channelId: "", channelName: "", messageIds: [] };
}
exports.DeleteMessagesRequest = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.channelId !== "") {
            writer.uint32(10).string(message.channelId);
        }
        if (message.channelName !== "") {
            writer.uint32(18).string(message.channelName);
        }
        for (const v of message.messageIds) {
            writer.uint32(26).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            channelId: isSet(object.channelId) ? String(object.channelId) : "",
            channelName: isSet(object.channelName) ? String(object.channelName) : "",
            messageIds: Array.isArray(object === null || object === void 0 ? void 0 : object.messageIds) ? object.messageIds.map((e) => String(e)) : [],
        };
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.channelName !== undefined && (obj.channelName = message.channelName);
        if (message.messageIds) {
            obj.messageIds = message.messageIds.map((e) => e);
        }
        else {
            obj.messageIds = [];
        }
        return obj;
    },
    create(base) {
        return exports.DeleteMessagesRequest.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c;
        const message = createBaseDeleteMessagesRequest();
        message.channelId = (_a = object.channelId) !== null && _a !== void 0 ? _a : "";
        message.channelName = (_b = object.channelName) !== null && _b !== void 0 ? _b : "";
        message.messageIds = ((_c = object.messageIds) === null || _c === void 0 ? void 0 : _c.map((e) => e)) || [];
        return message;
    },
};
function createBaseSendMessageRequest() {
    return { channelId: "", channelName: undefined, message: "", isAnnounce: undefined };
}
exports.SendMessageRequest = {
    encode(message, writer = minimal_1.default.Writer.create()) {
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            channelId: isSet(object.channelId) ? String(object.channelId) : "",
            channelName: isSet(object.channelName) ? String(object.channelName) : undefined,
            message: isSet(object.message) ? String(object.message) : "",
            isAnnounce: isSet(object.isAnnounce) ? Boolean(object.isAnnounce) : undefined,
        };
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.channelName !== undefined && (obj.channelName = message.channelName);
        message.message !== undefined && (obj.message = message.message);
        message.isAnnounce !== undefined && (obj.isAnnounce = message.isAnnounce);
        return obj;
    },
    create(base) {
        return exports.SendMessageRequest.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d;
        const message = createBaseSendMessageRequest();
        message.channelId = (_a = object.channelId) !== null && _a !== void 0 ? _a : "";
        message.channelName = (_b = object.channelName) !== null && _b !== void 0 ? _b : undefined;
        message.message = (_c = object.message) !== null && _c !== void 0 ? _c : "";
        message.isAnnounce = (_d = object.isAnnounce) !== null && _d !== void 0 ? _d : undefined;
        return message;
    },
};
function createBaseJoinOrLeaveRequest() {
    return { botId: "", userName: "" };
}
exports.JoinOrLeaveRequest = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.botId !== "") {
            writer.uint32(18).string(message.botId);
        }
        if (message.userName !== "") {
            writer.uint32(26).string(message.userName);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return {
            botId: isSet(object.botId) ? String(object.botId) : "",
            userName: isSet(object.userName) ? String(object.userName) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.botId !== undefined && (obj.botId = message.botId);
        message.userName !== undefined && (obj.userName = message.userName);
        return obj;
    },
    create(base) {
        return exports.JoinOrLeaveRequest.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b;
        const message = createBaseJoinOrLeaveRequest();
        message.botId = (_a = object.botId) !== null && _a !== void 0 ? _a : "";
        message.userName = (_b = object.userName) !== null && _b !== void 0 ? _b : "";
        return message;
    },
};
exports.BotsDefinition = {
    name: "Bots",
    fullName: "bots.Bots",
    methods: {
        deleteMessage: {
            name: "DeleteMessage",
            requestType: exports.DeleteMessagesRequest,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        sendMessage: {
            name: "SendMessage",
            requestType: exports.SendMessageRequest,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        join: {
            name: "Join",
            requestType: exports.JoinOrLeaveRequest,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        leave: {
            name: "Leave",
            requestType: exports.JoinOrLeaveRequest,
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
//# sourceMappingURL=bots.js.map