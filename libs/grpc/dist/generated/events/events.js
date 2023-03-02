"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.EventsDefinition = exports.GreetingSendedMessage = exports.KeywordMatchedMessage = exports.DonateMessage = exports.ChatClearMessage = exports.StreamOfflineMessage = exports.StreamOnlineMessage = exports.TitleOrCategoryChangedMessage = exports.RaidedMessage = exports.FirstUserMessageMessage = exports.CommandUsedMessage = exports.RedemptionCreatedMessage = exports.ReSubscribeMessage = exports.SubGiftMessage = exports.SubscribeMessage = exports.FollowMessage = exports.BaseInfo = exports.protobufPackage = void 0;
const long_1 = __importDefault(require("long"));
const minimal_1 = __importDefault(require("protobufjs/minimal"));
const empty_1 = require("./google/protobuf/empty");
exports.protobufPackage = "events";
function createBaseBaseInfo() {
    return { channelId: "" };
}
exports.BaseInfo = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.channelId !== "") {
            writer.uint32(10).string(message.channelId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
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
    fromJSON(object) {
        return { channelId: isSet(object.channelId) ? String(object.channelId) : "" };
    },
    toJSON(message) {
        const obj = {};
        message.channelId !== undefined && (obj.channelId = message.channelId);
        return obj;
    },
    create(base) {
        return exports.BaseInfo.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a;
        const message = createBaseBaseInfo();
        message.channelId = (_a = object.channelId) !== null && _a !== void 0 ? _a : "";
        return message;
    },
};
function createBaseFollowMessage() {
    return { baseInfo: undefined, userName: "", userDisplayName: "", userId: "" };
}
exports.FollowMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.baseInfo !== undefined) {
            exports.BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseFollowMessage();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseInfo = exports.BaseInfo.decode(reader, reader.uint32());
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
    fromJSON(object) {
        return {
            baseInfo: isSet(object.baseInfo) ? exports.BaseInfo.fromJSON(object.baseInfo) : undefined,
            userName: isSet(object.userName) ? String(object.userName) : "",
            userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
            userId: isSet(object.userId) ? String(object.userId) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? exports.BaseInfo.toJSON(message.baseInfo) : undefined);
        message.userName !== undefined && (obj.userName = message.userName);
        message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
        message.userId !== undefined && (obj.userId = message.userId);
        return obj;
    },
    create(base) {
        return exports.FollowMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c;
        const message = createBaseFollowMessage();
        message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
            ? exports.BaseInfo.fromPartial(object.baseInfo)
            : undefined;
        message.userName = (_a = object.userName) !== null && _a !== void 0 ? _a : "";
        message.userDisplayName = (_b = object.userDisplayName) !== null && _b !== void 0 ? _b : "";
        message.userId = (_c = object.userId) !== null && _c !== void 0 ? _c : "";
        return message;
    },
};
function createBaseSubscribeMessage() {
    return { baseInfo: undefined, userName: "", userDisplayName: "", level: "", userId: "" };
}
exports.SubscribeMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.baseInfo !== undefined) {
            exports.BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSubscribeMessage();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseInfo = exports.BaseInfo.decode(reader, reader.uint32());
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
    fromJSON(object) {
        return {
            baseInfo: isSet(object.baseInfo) ? exports.BaseInfo.fromJSON(object.baseInfo) : undefined,
            userName: isSet(object.userName) ? String(object.userName) : "",
            userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
            level: isSet(object.level) ? String(object.level) : "",
            userId: isSet(object.userId) ? String(object.userId) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? exports.BaseInfo.toJSON(message.baseInfo) : undefined);
        message.userName !== undefined && (obj.userName = message.userName);
        message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
        message.level !== undefined && (obj.level = message.level);
        message.userId !== undefined && (obj.userId = message.userId);
        return obj;
    },
    create(base) {
        return exports.SubscribeMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d;
        const message = createBaseSubscribeMessage();
        message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
            ? exports.BaseInfo.fromPartial(object.baseInfo)
            : undefined;
        message.userName = (_a = object.userName) !== null && _a !== void 0 ? _a : "";
        message.userDisplayName = (_b = object.userDisplayName) !== null && _b !== void 0 ? _b : "";
        message.level = (_c = object.level) !== null && _c !== void 0 ? _c : "";
        message.userId = (_d = object.userId) !== null && _d !== void 0 ? _d : "";
        return message;
    },
};
function createBaseSubGiftMessage() {
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
exports.SubGiftMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.baseInfo !== undefined) {
            exports.BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSubGiftMessage();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseInfo = exports.BaseInfo.decode(reader, reader.uint32());
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
    fromJSON(object) {
        return {
            baseInfo: isSet(object.baseInfo) ? exports.BaseInfo.fromJSON(object.baseInfo) : undefined,
            senderUserName: isSet(object.senderUserName) ? String(object.senderUserName) : "",
            senderDisplayName: isSet(object.senderDisplayName) ? String(object.senderDisplayName) : "",
            targetUserName: isSet(object.targetUserName) ? String(object.targetUserName) : "",
            targetDisplayName: isSet(object.targetDisplayName) ? String(object.targetDisplayName) : "",
            level: isSet(object.level) ? String(object.level) : "",
            senderUserId: isSet(object.senderUserId) ? String(object.senderUserId) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? exports.BaseInfo.toJSON(message.baseInfo) : undefined);
        message.senderUserName !== undefined && (obj.senderUserName = message.senderUserName);
        message.senderDisplayName !== undefined && (obj.senderDisplayName = message.senderDisplayName);
        message.targetUserName !== undefined && (obj.targetUserName = message.targetUserName);
        message.targetDisplayName !== undefined && (obj.targetDisplayName = message.targetDisplayName);
        message.level !== undefined && (obj.level = message.level);
        message.senderUserId !== undefined && (obj.senderUserId = message.senderUserId);
        return obj;
    },
    create(base) {
        return exports.SubGiftMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d, _e, _f;
        const message = createBaseSubGiftMessage();
        message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
            ? exports.BaseInfo.fromPartial(object.baseInfo)
            : undefined;
        message.senderUserName = (_a = object.senderUserName) !== null && _a !== void 0 ? _a : "";
        message.senderDisplayName = (_b = object.senderDisplayName) !== null && _b !== void 0 ? _b : "";
        message.targetUserName = (_c = object.targetUserName) !== null && _c !== void 0 ? _c : "";
        message.targetDisplayName = (_d = object.targetDisplayName) !== null && _d !== void 0 ? _d : "";
        message.level = (_e = object.level) !== null && _e !== void 0 ? _e : "";
        message.senderUserId = (_f = object.senderUserId) !== null && _f !== void 0 ? _f : "";
        return message;
    },
};
function createBaseReSubscribeMessage() {
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
exports.ReSubscribeMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.baseInfo !== undefined) {
            exports.BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseReSubscribeMessage();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseInfo = exports.BaseInfo.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.userName = reader.string();
                    break;
                case 3:
                    message.userDisplayName = reader.string();
                    break;
                case 4:
                    message.months = longToNumber(reader.int64());
                    break;
                case 5:
                    message.streak = longToNumber(reader.int64());
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
    fromJSON(object) {
        return {
            baseInfo: isSet(object.baseInfo) ? exports.BaseInfo.fromJSON(object.baseInfo) : undefined,
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
    toJSON(message) {
        const obj = {};
        message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? exports.BaseInfo.toJSON(message.baseInfo) : undefined);
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
    create(base) {
        return exports.ReSubscribeMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d, _e, _f, _g, _h;
        const message = createBaseReSubscribeMessage();
        message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
            ? exports.BaseInfo.fromPartial(object.baseInfo)
            : undefined;
        message.userName = (_a = object.userName) !== null && _a !== void 0 ? _a : "";
        message.userDisplayName = (_b = object.userDisplayName) !== null && _b !== void 0 ? _b : "";
        message.months = (_c = object.months) !== null && _c !== void 0 ? _c : 0;
        message.streak = (_d = object.streak) !== null && _d !== void 0 ? _d : 0;
        message.isPrime = (_e = object.isPrime) !== null && _e !== void 0 ? _e : false;
        message.message = (_f = object.message) !== null && _f !== void 0 ? _f : "";
        message.level = (_g = object.level) !== null && _g !== void 0 ? _g : "";
        message.userId = (_h = object.userId) !== null && _h !== void 0 ? _h : "";
        return message;
    },
};
function createBaseRedemptionCreatedMessage() {
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
exports.RedemptionCreatedMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.baseInfo !== undefined) {
            exports.BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRedemptionCreatedMessage();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseInfo = exports.BaseInfo.decode(reader, reader.uint32());
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
    fromJSON(object) {
        return {
            baseInfo: isSet(object.baseInfo) ? exports.BaseInfo.fromJSON(object.baseInfo) : undefined,
            userName: isSet(object.userName) ? String(object.userName) : "",
            userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
            id: isSet(object.id) ? String(object.id) : "",
            rewardName: isSet(object.rewardName) ? String(object.rewardName) : "",
            rewardCost: isSet(object.rewardCost) ? String(object.rewardCost) : "",
            input: isSet(object.input) ? String(object.input) : undefined,
            userId: isSet(object.userId) ? String(object.userId) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? exports.BaseInfo.toJSON(message.baseInfo) : undefined);
        message.userName !== undefined && (obj.userName = message.userName);
        message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
        message.id !== undefined && (obj.id = message.id);
        message.rewardName !== undefined && (obj.rewardName = message.rewardName);
        message.rewardCost !== undefined && (obj.rewardCost = message.rewardCost);
        message.input !== undefined && (obj.input = message.input);
        message.userId !== undefined && (obj.userId = message.userId);
        return obj;
    },
    create(base) {
        return exports.RedemptionCreatedMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d, _e, _f, _g;
        const message = createBaseRedemptionCreatedMessage();
        message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
            ? exports.BaseInfo.fromPartial(object.baseInfo)
            : undefined;
        message.userName = (_a = object.userName) !== null && _a !== void 0 ? _a : "";
        message.userDisplayName = (_b = object.userDisplayName) !== null && _b !== void 0 ? _b : "";
        message.id = (_c = object.id) !== null && _c !== void 0 ? _c : "";
        message.rewardName = (_d = object.rewardName) !== null && _d !== void 0 ? _d : "";
        message.rewardCost = (_e = object.rewardCost) !== null && _e !== void 0 ? _e : "";
        message.input = (_f = object.input) !== null && _f !== void 0 ? _f : undefined;
        message.userId = (_g = object.userId) !== null && _g !== void 0 ? _g : "";
        return message;
    },
};
function createBaseCommandUsedMessage() {
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
exports.CommandUsedMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.baseInfo !== undefined) {
            exports.BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseCommandUsedMessage();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseInfo = exports.BaseInfo.decode(reader, reader.uint32());
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
    fromJSON(object) {
        return {
            baseInfo: isSet(object.baseInfo) ? exports.BaseInfo.fromJSON(object.baseInfo) : undefined,
            commandId: isSet(object.commandId) ? String(object.commandId) : "",
            commandName: isSet(object.commandName) ? String(object.commandName) : "",
            userName: isSet(object.userName) ? String(object.userName) : "",
            userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
            commandInput: isSet(object.commandInput) ? String(object.commandInput) : "",
            userId: isSet(object.userId) ? String(object.userId) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? exports.BaseInfo.toJSON(message.baseInfo) : undefined);
        message.commandId !== undefined && (obj.commandId = message.commandId);
        message.commandName !== undefined && (obj.commandName = message.commandName);
        message.userName !== undefined && (obj.userName = message.userName);
        message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
        message.commandInput !== undefined && (obj.commandInput = message.commandInput);
        message.userId !== undefined && (obj.userId = message.userId);
        return obj;
    },
    create(base) {
        return exports.CommandUsedMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d, _e, _f;
        const message = createBaseCommandUsedMessage();
        message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
            ? exports.BaseInfo.fromPartial(object.baseInfo)
            : undefined;
        message.commandId = (_a = object.commandId) !== null && _a !== void 0 ? _a : "";
        message.commandName = (_b = object.commandName) !== null && _b !== void 0 ? _b : "";
        message.userName = (_c = object.userName) !== null && _c !== void 0 ? _c : "";
        message.userDisplayName = (_d = object.userDisplayName) !== null && _d !== void 0 ? _d : "";
        message.commandInput = (_e = object.commandInput) !== null && _e !== void 0 ? _e : "";
        message.userId = (_f = object.userId) !== null && _f !== void 0 ? _f : "";
        return message;
    },
};
function createBaseFirstUserMessageMessage() {
    return { baseInfo: undefined, userId: "", userName: "", userDisplayName: "" };
}
exports.FirstUserMessageMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.baseInfo !== undefined) {
            exports.BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseFirstUserMessageMessage();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseInfo = exports.BaseInfo.decode(reader, reader.uint32());
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
    fromJSON(object) {
        return {
            baseInfo: isSet(object.baseInfo) ? exports.BaseInfo.fromJSON(object.baseInfo) : undefined,
            userId: isSet(object.userId) ? String(object.userId) : "",
            userName: isSet(object.userName) ? String(object.userName) : "",
            userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? exports.BaseInfo.toJSON(message.baseInfo) : undefined);
        message.userId !== undefined && (obj.userId = message.userId);
        message.userName !== undefined && (obj.userName = message.userName);
        message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
        return obj;
    },
    create(base) {
        return exports.FirstUserMessageMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c;
        const message = createBaseFirstUserMessageMessage();
        message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
            ? exports.BaseInfo.fromPartial(object.baseInfo)
            : undefined;
        message.userId = (_a = object.userId) !== null && _a !== void 0 ? _a : "";
        message.userName = (_b = object.userName) !== null && _b !== void 0 ? _b : "";
        message.userDisplayName = (_c = object.userDisplayName) !== null && _c !== void 0 ? _c : "";
        return message;
    },
};
function createBaseRaidedMessage() {
    return { baseInfo: undefined, userName: "", userDisplayName: "", viewers: 0, userId: "" };
}
exports.RaidedMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.baseInfo !== undefined) {
            exports.BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRaidedMessage();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseInfo = exports.BaseInfo.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.userName = reader.string();
                    break;
                case 3:
                    message.userDisplayName = reader.string();
                    break;
                case 4:
                    message.viewers = longToNumber(reader.int64());
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
    fromJSON(object) {
        return {
            baseInfo: isSet(object.baseInfo) ? exports.BaseInfo.fromJSON(object.baseInfo) : undefined,
            userName: isSet(object.userName) ? String(object.userName) : "",
            userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
            viewers: isSet(object.viewers) ? Number(object.viewers) : 0,
            userId: isSet(object.userId) ? String(object.userId) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? exports.BaseInfo.toJSON(message.baseInfo) : undefined);
        message.userName !== undefined && (obj.userName = message.userName);
        message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
        message.viewers !== undefined && (obj.viewers = Math.round(message.viewers));
        message.userId !== undefined && (obj.userId = message.userId);
        return obj;
    },
    create(base) {
        return exports.RaidedMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d;
        const message = createBaseRaidedMessage();
        message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
            ? exports.BaseInfo.fromPartial(object.baseInfo)
            : undefined;
        message.userName = (_a = object.userName) !== null && _a !== void 0 ? _a : "";
        message.userDisplayName = (_b = object.userDisplayName) !== null && _b !== void 0 ? _b : "";
        message.viewers = (_c = object.viewers) !== null && _c !== void 0 ? _c : 0;
        message.userId = (_d = object.userId) !== null && _d !== void 0 ? _d : "";
        return message;
    },
};
function createBaseTitleOrCategoryChangedMessage() {
    return { baseInfo: undefined, oldTitle: "", newTitle: "", oldCategory: "", newCategory: "" };
}
exports.TitleOrCategoryChangedMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.baseInfo !== undefined) {
            exports.BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseTitleOrCategoryChangedMessage();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseInfo = exports.BaseInfo.decode(reader, reader.uint32());
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
    fromJSON(object) {
        return {
            baseInfo: isSet(object.baseInfo) ? exports.BaseInfo.fromJSON(object.baseInfo) : undefined,
            oldTitle: isSet(object.oldTitle) ? String(object.oldTitle) : "",
            newTitle: isSet(object.newTitle) ? String(object.newTitle) : "",
            oldCategory: isSet(object.oldCategory) ? String(object.oldCategory) : "",
            newCategory: isSet(object.newCategory) ? String(object.newCategory) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? exports.BaseInfo.toJSON(message.baseInfo) : undefined);
        message.oldTitle !== undefined && (obj.oldTitle = message.oldTitle);
        message.newTitle !== undefined && (obj.newTitle = message.newTitle);
        message.oldCategory !== undefined && (obj.oldCategory = message.oldCategory);
        message.newCategory !== undefined && (obj.newCategory = message.newCategory);
        return obj;
    },
    create(base) {
        return exports.TitleOrCategoryChangedMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d;
        const message = createBaseTitleOrCategoryChangedMessage();
        message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
            ? exports.BaseInfo.fromPartial(object.baseInfo)
            : undefined;
        message.oldTitle = (_a = object.oldTitle) !== null && _a !== void 0 ? _a : "";
        message.newTitle = (_b = object.newTitle) !== null && _b !== void 0 ? _b : "";
        message.oldCategory = (_c = object.oldCategory) !== null && _c !== void 0 ? _c : "";
        message.newCategory = (_d = object.newCategory) !== null && _d !== void 0 ? _d : "";
        return message;
    },
};
function createBaseStreamOnlineMessage() {
    return { baseInfo: undefined, title: "", category: "" };
}
exports.StreamOnlineMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.baseInfo !== undefined) {
            exports.BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
        }
        if (message.title !== "") {
            writer.uint32(18).string(message.title);
        }
        if (message.category !== "") {
            writer.uint32(26).string(message.category);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseStreamOnlineMessage();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseInfo = exports.BaseInfo.decode(reader, reader.uint32());
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
    fromJSON(object) {
        return {
            baseInfo: isSet(object.baseInfo) ? exports.BaseInfo.fromJSON(object.baseInfo) : undefined,
            title: isSet(object.title) ? String(object.title) : "",
            category: isSet(object.category) ? String(object.category) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? exports.BaseInfo.toJSON(message.baseInfo) : undefined);
        message.title !== undefined && (obj.title = message.title);
        message.category !== undefined && (obj.category = message.category);
        return obj;
    },
    create(base) {
        return exports.StreamOnlineMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b;
        const message = createBaseStreamOnlineMessage();
        message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
            ? exports.BaseInfo.fromPartial(object.baseInfo)
            : undefined;
        message.title = (_a = object.title) !== null && _a !== void 0 ? _a : "";
        message.category = (_b = object.category) !== null && _b !== void 0 ? _b : "";
        return message;
    },
};
function createBaseStreamOfflineMessage() {
    return { baseInfo: undefined };
}
exports.StreamOfflineMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.baseInfo !== undefined) {
            exports.BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseStreamOfflineMessage();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseInfo = exports.BaseInfo.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        return { baseInfo: isSet(object.baseInfo) ? exports.BaseInfo.fromJSON(object.baseInfo) : undefined };
    },
    toJSON(message) {
        const obj = {};
        message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? exports.BaseInfo.toJSON(message.baseInfo) : undefined);
        return obj;
    },
    create(base) {
        return exports.StreamOfflineMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        const message = createBaseStreamOfflineMessage();
        message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
            ? exports.BaseInfo.fromPartial(object.baseInfo)
            : undefined;
        return message;
    },
};
function createBaseChatClearMessage() {
    return { baseInfo: undefined };
}
exports.ChatClearMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.baseInfo !== undefined) {
            exports.BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseChatClearMessage();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseInfo = exports.BaseInfo.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        return { baseInfo: isSet(object.baseInfo) ? exports.BaseInfo.fromJSON(object.baseInfo) : undefined };
    },
    toJSON(message) {
        const obj = {};
        message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? exports.BaseInfo.toJSON(message.baseInfo) : undefined);
        return obj;
    },
    create(base) {
        return exports.ChatClearMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        const message = createBaseChatClearMessage();
        message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
            ? exports.BaseInfo.fromPartial(object.baseInfo)
            : undefined;
        return message;
    },
};
function createBaseDonateMessage() {
    return { baseInfo: undefined, userName: "", amount: "", currency: "", message: "" };
}
exports.DonateMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.baseInfo !== undefined) {
            exports.BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseDonateMessage();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseInfo = exports.BaseInfo.decode(reader, reader.uint32());
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
    fromJSON(object) {
        return {
            baseInfo: isSet(object.baseInfo) ? exports.BaseInfo.fromJSON(object.baseInfo) : undefined,
            userName: isSet(object.userName) ? String(object.userName) : "",
            amount: isSet(object.amount) ? String(object.amount) : "",
            currency: isSet(object.currency) ? String(object.currency) : "",
            message: isSet(object.message) ? String(object.message) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? exports.BaseInfo.toJSON(message.baseInfo) : undefined);
        message.userName !== undefined && (obj.userName = message.userName);
        message.amount !== undefined && (obj.amount = message.amount);
        message.currency !== undefined && (obj.currency = message.currency);
        message.message !== undefined && (obj.message = message.message);
        return obj;
    },
    create(base) {
        return exports.DonateMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d;
        const message = createBaseDonateMessage();
        message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
            ? exports.BaseInfo.fromPartial(object.baseInfo)
            : undefined;
        message.userName = (_a = object.userName) !== null && _a !== void 0 ? _a : "";
        message.amount = (_b = object.amount) !== null && _b !== void 0 ? _b : "";
        message.currency = (_c = object.currency) !== null && _c !== void 0 ? _c : "";
        message.message = (_d = object.message) !== null && _d !== void 0 ? _d : "";
        return message;
    },
};
function createBaseKeywordMatchedMessage() {
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
exports.KeywordMatchedMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.baseInfo !== undefined) {
            exports.BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseKeywordMatchedMessage();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseInfo = exports.BaseInfo.decode(reader, reader.uint32());
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
    fromJSON(object) {
        return {
            baseInfo: isSet(object.baseInfo) ? exports.BaseInfo.fromJSON(object.baseInfo) : undefined,
            keywordId: isSet(object.keywordId) ? String(object.keywordId) : "",
            keywordName: isSet(object.keywordName) ? String(object.keywordName) : "",
            keywordResponse: isSet(object.keywordResponse) ? String(object.keywordResponse) : "",
            userId: isSet(object.userId) ? String(object.userId) : "",
            userName: isSet(object.userName) ? String(object.userName) : "",
            userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? exports.BaseInfo.toJSON(message.baseInfo) : undefined);
        message.keywordId !== undefined && (obj.keywordId = message.keywordId);
        message.keywordName !== undefined && (obj.keywordName = message.keywordName);
        message.keywordResponse !== undefined && (obj.keywordResponse = message.keywordResponse);
        message.userId !== undefined && (obj.userId = message.userId);
        message.userName !== undefined && (obj.userName = message.userName);
        message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
        return obj;
    },
    create(base) {
        return exports.KeywordMatchedMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d, _e, _f;
        const message = createBaseKeywordMatchedMessage();
        message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
            ? exports.BaseInfo.fromPartial(object.baseInfo)
            : undefined;
        message.keywordId = (_a = object.keywordId) !== null && _a !== void 0 ? _a : "";
        message.keywordName = (_b = object.keywordName) !== null && _b !== void 0 ? _b : "";
        message.keywordResponse = (_c = object.keywordResponse) !== null && _c !== void 0 ? _c : "";
        message.userId = (_d = object.userId) !== null && _d !== void 0 ? _d : "";
        message.userName = (_e = object.userName) !== null && _e !== void 0 ? _e : "";
        message.userDisplayName = (_f = object.userDisplayName) !== null && _f !== void 0 ? _f : "";
        return message;
    },
};
function createBaseGreetingSendedMessage() {
    return { baseInfo: undefined, userId: "", userName: "", userDisplayName: "", greetingText: "" };
}
exports.GreetingSendedMessage = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.baseInfo !== undefined) {
            exports.BaseInfo.encode(message.baseInfo, writer.uint32(10).fork()).ldelim();
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
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGreetingSendedMessage();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.baseInfo = exports.BaseInfo.decode(reader, reader.uint32());
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
    fromJSON(object) {
        return {
            baseInfo: isSet(object.baseInfo) ? exports.BaseInfo.fromJSON(object.baseInfo) : undefined,
            userId: isSet(object.userId) ? String(object.userId) : "",
            userName: isSet(object.userName) ? String(object.userName) : "",
            userDisplayName: isSet(object.userDisplayName) ? String(object.userDisplayName) : "",
            greetingText: isSet(object.greetingText) ? String(object.greetingText) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        message.baseInfo !== undefined && (obj.baseInfo = message.baseInfo ? exports.BaseInfo.toJSON(message.baseInfo) : undefined);
        message.userId !== undefined && (obj.userId = message.userId);
        message.userName !== undefined && (obj.userName = message.userName);
        message.userDisplayName !== undefined && (obj.userDisplayName = message.userDisplayName);
        message.greetingText !== undefined && (obj.greetingText = message.greetingText);
        return obj;
    },
    create(base) {
        return exports.GreetingSendedMessage.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d;
        const message = createBaseGreetingSendedMessage();
        message.baseInfo = (object.baseInfo !== undefined && object.baseInfo !== null)
            ? exports.BaseInfo.fromPartial(object.baseInfo)
            : undefined;
        message.userId = (_a = object.userId) !== null && _a !== void 0 ? _a : "";
        message.userName = (_b = object.userName) !== null && _b !== void 0 ? _b : "";
        message.userDisplayName = (_c = object.userDisplayName) !== null && _c !== void 0 ? _c : "";
        message.greetingText = (_d = object.greetingText) !== null && _d !== void 0 ? _d : "";
        return message;
    },
};
exports.EventsDefinition = {
    name: "Events",
    fullName: "events.Events",
    methods: {
        follow: {
            name: "Follow",
            requestType: exports.FollowMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        subscribe: {
            name: "Subscribe",
            requestType: exports.SubscribeMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        subGift: {
            name: "SubGift",
            requestType: exports.SubGiftMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        reSubscribe: {
            name: "ReSubscribe",
            requestType: exports.ReSubscribeMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        redemptionCreated: {
            name: "RedemptionCreated",
            requestType: exports.RedemptionCreatedMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        commandUsed: {
            name: "CommandUsed",
            requestType: exports.CommandUsedMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        firstUserMessage: {
            name: "FirstUserMessage",
            requestType: exports.FirstUserMessageMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        raided: {
            name: "Raided",
            requestType: exports.RaidedMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        titleOrCategoryChanged: {
            name: "TitleOrCategoryChanged",
            requestType: exports.TitleOrCategoryChangedMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        streamOnline: {
            name: "StreamOnline",
            requestType: exports.StreamOnlineMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        streamOffline: {
            name: "StreamOffline",
            requestType: exports.StreamOfflineMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        chatClear: {
            name: "ChatClear",
            requestType: exports.ChatClearMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        donate: {
            name: "Donate",
            requestType: exports.DonateMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        keywordMatched: {
            name: "KeywordMatched",
            requestType: exports.KeywordMatchedMessage,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        greetingSended: {
            name: "GreetingSended",
            requestType: exports.GreetingSendedMessage,
            requestStream: false,
            responseType: empty_1.Empty,
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
//# sourceMappingURL=events.js.map