"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.DotaDefinition = exports.GetPlayerCardResponse = exports.GetPlayerCardRequest = exports.protobufPackage = void 0;
const long_1 = __importDefault(require("long"));
const minimal_1 = __importDefault(require("protobufjs/minimal"));
exports.protobufPackage = "dota";
function createBaseGetPlayerCardRequest() {
    return { accountId: 0 };
}
exports.GetPlayerCardRequest = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.accountId !== 0) {
            writer.uint32(8).int64(message.accountId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetPlayerCardRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.accountId = longToNumber(reader.int64());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        return { accountId: isSet(object.accountId) ? Number(object.accountId) : 0 };
    },
    toJSON(message) {
        const obj = {};
        message.accountId !== undefined && (obj.accountId = Math.round(message.accountId));
        return obj;
    },
    create(base) {
        return exports.GetPlayerCardRequest.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a;
        const message = createBaseGetPlayerCardRequest();
        message.accountId = (_a = object.accountId) !== null && _a !== void 0 ? _a : 0;
        return message;
    },
};
function createBaseGetPlayerCardResponse() {
    return { accountId: "", rankTier: undefined, leaderboardRank: undefined };
}
exports.GetPlayerCardResponse = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.accountId !== "") {
            writer.uint32(10).string(message.accountId);
        }
        if (message.rankTier !== undefined) {
            writer.uint32(16).int64(message.rankTier);
        }
        if (message.leaderboardRank !== undefined) {
            writer.uint32(24).int64(message.leaderboardRank);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseGetPlayerCardResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.accountId = reader.string();
                    break;
                case 2:
                    message.rankTier = longToNumber(reader.int64());
                    break;
                case 3:
                    message.leaderboardRank = longToNumber(reader.int64());
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
            accountId: isSet(object.accountId) ? String(object.accountId) : "",
            rankTier: isSet(object.rankTier) ? Number(object.rankTier) : undefined,
            leaderboardRank: isSet(object.leaderboardRank) ? Number(object.leaderboardRank) : undefined,
        };
    },
    toJSON(message) {
        const obj = {};
        message.accountId !== undefined && (obj.accountId = message.accountId);
        message.rankTier !== undefined && (obj.rankTier = Math.round(message.rankTier));
        message.leaderboardRank !== undefined && (obj.leaderboardRank = Math.round(message.leaderboardRank));
        return obj;
    },
    create(base) {
        return exports.GetPlayerCardResponse.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c;
        const message = createBaseGetPlayerCardResponse();
        message.accountId = (_a = object.accountId) !== null && _a !== void 0 ? _a : "";
        message.rankTier = (_b = object.rankTier) !== null && _b !== void 0 ? _b : undefined;
        message.leaderboardRank = (_c = object.leaderboardRank) !== null && _c !== void 0 ? _c : undefined;
        return message;
    },
};
exports.DotaDefinition = {
    name: "Dota",
    fullName: "dota.Dota",
    methods: {
        getPlayerCard: {
            name: "GetPlayerCard",
            requestType: exports.GetPlayerCardRequest,
            requestStream: false,
            responseType: exports.GetPlayerCardResponse,
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
//# sourceMappingURL=dota.js.map