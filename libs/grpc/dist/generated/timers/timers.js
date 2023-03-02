"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.TimersDefinition = exports.Request = exports.protobufPackage = void 0;
const minimal_1 = __importDefault(require("protobufjs/minimal"));
const empty_1 = require("./google/protobuf/empty");
exports.protobufPackage = "timers";
function createBaseRequest() {
    return { timerId: "" };
}
exports.Request = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.timerId !== "") {
            writer.uint32(10).string(message.timerId);
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
                    message.timerId = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        return { timerId: isSet(object.timerId) ? String(object.timerId) : "" };
    },
    toJSON(message) {
        const obj = {};
        message.timerId !== undefined && (obj.timerId = message.timerId);
        return obj;
    },
    create(base) {
        return exports.Request.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a;
        const message = createBaseRequest();
        message.timerId = (_a = object.timerId) !== null && _a !== void 0 ? _a : "";
        return message;
    },
};
exports.TimersDefinition = {
    name: "Timers",
    fullName: "timers.Timers",
    methods: {
        addTimerToQueue: {
            name: "AddTimerToQueue",
            requestType: exports.Request,
            requestStream: false,
            responseType: empty_1.Empty,
            responseStream: false,
            options: {},
        },
        removeTimerFromQueue: {
            name: "RemoveTimerFromQueue",
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
//# sourceMappingURL=timers.js.map