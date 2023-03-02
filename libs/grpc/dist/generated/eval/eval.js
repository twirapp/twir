"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.EvalDefinition = exports.EvaluateResult = exports.Evaluate = exports.protobufPackage = void 0;
const minimal_1 = __importDefault(require("protobufjs/minimal"));
exports.protobufPackage = "eval";
function createBaseEvaluate() {
    return { script: "" };
}
exports.Evaluate = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.script !== "") {
            writer.uint32(10).string(message.script);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseEvaluate();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.script = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        return { script: isSet(object.script) ? String(object.script) : "" };
    },
    toJSON(message) {
        const obj = {};
        message.script !== undefined && (obj.script = message.script);
        return obj;
    },
    create(base) {
        return exports.Evaluate.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a;
        const message = createBaseEvaluate();
        message.script = (_a = object.script) !== null && _a !== void 0 ? _a : "";
        return message;
    },
};
function createBaseEvaluateResult() {
    return { result: "" };
}
exports.EvaluateResult = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.result !== "") {
            writer.uint32(10).string(message.result);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseEvaluateResult();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.result = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        return { result: isSet(object.result) ? String(object.result) : "" };
    },
    toJSON(message) {
        const obj = {};
        message.result !== undefined && (obj.result = message.result);
        return obj;
    },
    create(base) {
        return exports.EvaluateResult.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a;
        const message = createBaseEvaluateResult();
        message.result = (_a = object.result) !== null && _a !== void 0 ? _a : "";
        return message;
    },
};
exports.EvalDefinition = {
    name: "Eval",
    fullName: "eval.Eval",
    methods: {
        process: {
            name: "Process",
            requestType: exports.Evaluate,
            requestStream: false,
            responseType: exports.EvaluateResult,
            responseStream: false,
            options: {},
        },
    },
};
function isSet(value) {
    return value !== null && value !== undefined;
}
//# sourceMappingURL=eval.js.map