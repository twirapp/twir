"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.WatchedDefinition = exports.Request = exports.protobufPackage = void 0;
const minimal_1 = __importDefault(require("protobufjs/minimal"));
const empty_1 = require("./google/protobuf/empty");
exports.protobufPackage = "watched";
function createBaseRequest() {
    return { channelsId: [], botId: "" };
}
exports.Request = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        for (const v of message.channelsId) {
            writer.uint32(10).string(v);
        }
        if (message.botId !== "") {
            writer.uint32(18).string(message.botId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.channelsId.push(reader.string());
                    break;
                case 2:
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
        return {
            channelsId: Array.isArray(object === null || object === void 0 ? void 0 : object.channelsId) ? object.channelsId.map((e) => String(e)) : [],
            botId: isSet(object.botId) ? String(object.botId) : "",
        };
    },
    toJSON(message) {
        const obj = {};
        if (message.channelsId) {
            obj.channelsId = message.channelsId.map((e) => e);
        }
        else {
            obj.channelsId = [];
        }
        message.botId !== undefined && (obj.botId = message.botId);
        return obj;
    },
    create(base) {
        return exports.Request.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b;
        const message = createBaseRequest();
        message.channelsId = ((_a = object.channelsId) === null || _a === void 0 ? void 0 : _a.map((e) => e)) || [];
        message.botId = (_b = object.botId) !== null && _b !== void 0 ? _b : "";
        return message;
    },
};
exports.WatchedDefinition = {
    name: "Watched",
    fullName: "watched.Watched",
    methods: {
        incrementByChannelId: {
            name: "IncrementByChannelId",
            requestType: exports.Request,
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
//# sourceMappingURL=watched.js.map