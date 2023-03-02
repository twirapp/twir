"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.TokensDefinition = exports.UpdateToken = exports.GetBotTokenRequest = exports.GetUserTokenRequest = exports.Token = exports.protobufPackage = void 0;
const long_1 = __importDefault(require("long"));
const minimal_1 = __importDefault(require("protobufjs/minimal"));
const empty_1 = require("./google/protobuf/empty");
exports.protobufPackage = "tokens";
function createBaseToken() {
    return { accessToken: "", scopes: [] };
}
exports.Token = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.accessToken !== "") {
            writer.uint32(10).string(message.accessToken);
        }
        for (const v of message.scopes) {
            writer.uint32(18).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseToken();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.accessToken = reader.string();
                    break;
                case 2:
                    message.scopes.push(reader.string());
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
            accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
            scopes: Array.isArray(object === null || object === void 0 ? void 0 : object.scopes) ? object.scopes.map((e) => String(e)) : [],
        };
    },
    toJSON(message) {
        const obj = {};
        message.accessToken !== undefined && (obj.accessToken = message.accessToken);
        if (message.scopes) {
            obj.scopes = message.scopes.map((e) => e);
        }
        else {
            obj.scopes = [];
        }
        return obj;
    },
    create(base) {
        return exports.Token.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b;
        const message = createBaseToken();
        message.accessToken = (_a = object.accessToken) !== null && _a !== void 0 ? _a : "";
        message.scopes = ((_b = object.scopes) === null || _b === void 0 ? void 0 : _b.map((e) => e)) || [];
        return message;
    },
};
function createBaseGetUserTokenRequest() {
    return { userId: "" };
}
exports.GetUserTokenRequest = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.userId !== "") {
            writer.uint32(10).string(message.userId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetUserTokenRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
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
        return { userId: isSet(object.userId) ? String(object.userId) : "" };
    },
    toJSON(message) {
        const obj = {};
        message.userId !== undefined && (obj.userId = message.userId);
        return obj;
    },
    create(base) {
        return exports.GetUserTokenRequest.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a;
        const message = createBaseGetUserTokenRequest();
        message.userId = (_a = object.userId) !== null && _a !== void 0 ? _a : "";
        return message;
    },
};
function createBaseGetBotTokenRequest() {
    return { botId: "" };
}
exports.GetBotTokenRequest = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.botId !== "") {
            writer.uint32(10).string(message.botId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetBotTokenRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.botId = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        return { botId: isSet(object.botId) ? String(object.botId) : "" };
    },
    toJSON(message) {
        const obj = {};
        message.botId !== undefined && (obj.botId = message.botId);
        return obj;
    },
    create(base) {
        return exports.GetBotTokenRequest.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a;
        const message = createBaseGetBotTokenRequest();
        message.botId = (_a = object.botId) !== null && _a !== void 0 ? _a : "";
        return message;
    },
};
function createBaseUpdateToken() {
    return { accessToken: "", refreshToken: "", expiresIn: 0, scopes: [] };
}
exports.UpdateToken = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.accessToken !== "") {
            writer.uint32(10).string(message.accessToken);
        }
        if (message.refreshToken !== "") {
            writer.uint32(18).string(message.refreshToken);
        }
        if (message.expiresIn !== 0) {
            writer.uint32(24).int64(message.expiresIn);
        }
        for (const v of message.scopes) {
            writer.uint32(34).string(v);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseUpdateToken();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.accessToken = reader.string();
                    break;
                case 2:
                    message.refreshToken = reader.string();
                    break;
                case 3:
                    message.expiresIn = longToNumber(reader.int64());
                    break;
                case 4:
                    message.scopes.push(reader.string());
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
            accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
            refreshToken: isSet(object.refreshToken) ? String(object.refreshToken) : "",
            expiresIn: isSet(object.expiresIn) ? Number(object.expiresIn) : 0,
            scopes: Array.isArray(object === null || object === void 0 ? void 0 : object.scopes) ? object.scopes.map((e) => String(e)) : [],
        };
    },
    toJSON(message) {
        const obj = {};
        message.accessToken !== undefined && (obj.accessToken = message.accessToken);
        message.refreshToken !== undefined && (obj.refreshToken = message.refreshToken);
        message.expiresIn !== undefined && (obj.expiresIn = Math.round(message.expiresIn));
        if (message.scopes) {
            obj.scopes = message.scopes.map((e) => e);
        }
        else {
            obj.scopes = [];
        }
        return obj;
    },
    create(base) {
        return exports.UpdateToken.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d;
        const message = createBaseUpdateToken();
        message.accessToken = (_a = object.accessToken) !== null && _a !== void 0 ? _a : "";
        message.refreshToken = (_b = object.refreshToken) !== null && _b !== void 0 ? _b : "";
        message.expiresIn = (_c = object.expiresIn) !== null && _c !== void 0 ? _c : 0;
        message.scopes = ((_d = object.scopes) === null || _d === void 0 ? void 0 : _d.map((e) => e)) || [];
        return message;
    },
};
exports.TokensDefinition = {
    name: "Tokens",
    fullName: "tokens.Tokens",
    methods: {
        requestAppToken: {
            name: "RequestAppToken",
            requestType: empty_1.Empty,
            requestStream: false,
            responseType: exports.Token,
            responseStream: false,
            options: {},
        },
        requestUserToken: {
            name: "RequestUserToken",
            requestType: exports.GetUserTokenRequest,
            requestStream: false,
            responseType: exports.Token,
            responseStream: false,
            options: {},
        },
        requestBotToken: {
            name: "RequestBotToken",
            requestType: exports.GetBotTokenRequest,
            requestStream: false,
            responseType: exports.Token,
            responseStream: false,
            options: {},
        },
        updateBotToken: {
            name: "UpdateBotToken",
            requestType: exports.UpdateToken,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        updateUserToken: {
            name: "UpdateUserToken",
            requestType: exports.UpdateToken,
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
//# sourceMappingURL=tokens.js.map