"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.SchedulerDefinition = exports.CreateDefaultCommandsRequest = exports.protobufPackage = void 0;
const minimal_1 = __importDefault(require("protobufjs/minimal"));
const empty_1 = require("./google/protobuf/empty");
exports.protobufPackage = "scheduler";
function createBaseCreateDefaultCommandsRequest() {
    return { userId: "" };
}
exports.CreateDefaultCommandsRequest = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.userId !== "") {
            writer.uint32(10).string(message.userId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseCreateDefaultCommandsRequest();
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
        return exports.CreateDefaultCommandsRequest.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a;
        const message = createBaseCreateDefaultCommandsRequest();
        message.userId = (_a = object.userId) !== null && _a !== void 0 ? _a : "";
        return message;
    },
};
exports.SchedulerDefinition = {
    name: "Scheduler",
    fullName: "scheduler.Scheduler",
    methods: {
        createDefaultCommands: {
            name: "CreateDefaultCommands",
            requestType: exports.CreateDefaultCommandsRequest,
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
//# sourceMappingURL=scheduler.js.map