"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.YtsrDefinition = exports.SearchResponse = exports.Song = exports.SongAuthor = exports.SearchRequest = exports.protobufPackage = void 0;
const long_1 = __importDefault(require("long"));
const minimal_1 = __importDefault(require("protobufjs/minimal"));
exports.protobufPackage = "ytsr";
function createBaseSearchRequest() {
    return { search: "" };
}
exports.SearchRequest = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.search !== "") {
            writer.uint32(10).string(message.search);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSearchRequest();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.search = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        return { search: isSet(object.search) ? String(object.search) : "" };
    },
    toJSON(message) {
        const obj = {};
        message.search !== undefined && (obj.search = message.search);
        return obj;
    },
    create(base) {
        return exports.SearchRequest.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a;
        const message = createBaseSearchRequest();
        message.search = (_a = object.search) !== null && _a !== void 0 ? _a : "";
        return message;
    },
};
function createBaseSongAuthor() {
    return { name: "", channelId: "", avatarUrl: undefined };
}
exports.SongAuthor = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.name !== "") {
            writer.uint32(10).string(message.name);
        }
        if (message.channelId !== "") {
            writer.uint32(18).string(message.channelId);
        }
        if (message.avatarUrl !== undefined) {
            writer.uint32(26).string(message.avatarUrl);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSongAuthor();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.name = reader.string();
                    break;
                case 2:
                    message.channelId = reader.string();
                    break;
                case 3:
                    message.avatarUrl = reader.string();
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
            channelId: isSet(object.channelId) ? String(object.channelId) : "",
            avatarUrl: isSet(object.avatarUrl) ? String(object.avatarUrl) : undefined,
        };
    },
    toJSON(message) {
        const obj = {};
        message.name !== undefined && (obj.name = message.name);
        message.channelId !== undefined && (obj.channelId = message.channelId);
        message.avatarUrl !== undefined && (obj.avatarUrl = message.avatarUrl);
        return obj;
    },
    create(base) {
        return exports.SongAuthor.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c;
        const message = createBaseSongAuthor();
        message.name = (_a = object.name) !== null && _a !== void 0 ? _a : "";
        message.channelId = (_b = object.channelId) !== null && _b !== void 0 ? _b : "";
        message.avatarUrl = (_c = object.avatarUrl) !== null && _c !== void 0 ? _c : undefined;
        return message;
    },
};
function createBaseSong() {
    return { title: "", id: "", views: 0, duration: 0, thumbnailUrl: undefined, isLive: false, author: undefined };
}
exports.Song = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        if (message.title !== "") {
            writer.uint32(10).string(message.title);
        }
        if (message.id !== "") {
            writer.uint32(18).string(message.id);
        }
        if (message.views !== 0) {
            writer.uint32(24).uint64(message.views);
        }
        if (message.duration !== 0) {
            writer.uint32(32).uint64(message.duration);
        }
        if (message.thumbnailUrl !== undefined) {
            writer.uint32(42).string(message.thumbnailUrl);
        }
        if (message.isLive === true) {
            writer.uint32(48).bool(message.isLive);
        }
        if (message.author !== undefined) {
            exports.SongAuthor.encode(message.author, writer.uint32(58).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSong();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.title = reader.string();
                    break;
                case 2:
                    message.id = reader.string();
                    break;
                case 3:
                    message.views = longToNumber(reader.uint64());
                    break;
                case 4:
                    message.duration = longToNumber(reader.uint64());
                    break;
                case 5:
                    message.thumbnailUrl = reader.string();
                    break;
                case 6:
                    message.isLive = reader.bool();
                    break;
                case 7:
                    message.author = exports.SongAuthor.decode(reader, reader.uint32());
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
            title: isSet(object.title) ? String(object.title) : "",
            id: isSet(object.id) ? String(object.id) : "",
            views: isSet(object.views) ? Number(object.views) : 0,
            duration: isSet(object.duration) ? Number(object.duration) : 0,
            thumbnailUrl: isSet(object.thumbnailUrl) ? String(object.thumbnailUrl) : undefined,
            isLive: isSet(object.isLive) ? Boolean(object.isLive) : false,
            author: isSet(object.author) ? exports.SongAuthor.fromJSON(object.author) : undefined,
        };
    },
    toJSON(message) {
        const obj = {};
        message.title !== undefined && (obj.title = message.title);
        message.id !== undefined && (obj.id = message.id);
        message.views !== undefined && (obj.views = Math.round(message.views));
        message.duration !== undefined && (obj.duration = Math.round(message.duration));
        message.thumbnailUrl !== undefined && (obj.thumbnailUrl = message.thumbnailUrl);
        message.isLive !== undefined && (obj.isLive = message.isLive);
        message.author !== undefined && (obj.author = message.author ? exports.SongAuthor.toJSON(message.author) : undefined);
        return obj;
    },
    create(base) {
        return exports.Song.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a, _b, _c, _d, _e, _f;
        const message = createBaseSong();
        message.title = (_a = object.title) !== null && _a !== void 0 ? _a : "";
        message.id = (_b = object.id) !== null && _b !== void 0 ? _b : "";
        message.views = (_c = object.views) !== null && _c !== void 0 ? _c : 0;
        message.duration = (_d = object.duration) !== null && _d !== void 0 ? _d : 0;
        message.thumbnailUrl = (_e = object.thumbnailUrl) !== null && _e !== void 0 ? _e : undefined;
        message.isLive = (_f = object.isLive) !== null && _f !== void 0 ? _f : false;
        message.author = (object.author !== undefined && object.author !== null)
            ? exports.SongAuthor.fromPartial(object.author)
            : undefined;
        return message;
    },
};
function createBaseSearchResponse() {
    return { songs: [] };
}
exports.SearchResponse = {
    encode(message, writer = minimal_1.default.Writer.create()) {
        for (const v of message.songs) {
            exports.Song.encode(v, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof minimal_1.default.Reader ? input : new minimal_1.default.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseSearchResponse();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.songs.push(exports.Song.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        return { songs: Array.isArray(object === null || object === void 0 ? void 0 : object.songs) ? object.songs.map((e) => exports.Song.fromJSON(e)) : [] };
    },
    toJSON(message) {
        const obj = {};
        if (message.songs) {
            obj.songs = message.songs.map((e) => e ? exports.Song.toJSON(e) : undefined);
        }
        else {
            obj.songs = [];
        }
        return obj;
    },
    create(base) {
        return exports.SearchResponse.fromPartial(base !== null && base !== void 0 ? base : {});
    },
    fromPartial(object) {
        var _a;
        const message = createBaseSearchResponse();
        message.songs = ((_a = object.songs) === null || _a === void 0 ? void 0 : _a.map((e) => exports.Song.fromPartial(e))) || [];
        return message;
    },
};
exports.YtsrDefinition = {
    name: "Ytsr",
    fullName: "ytsr.Ytsr",
    methods: {
        search: {
            name: "Search",
            requestType: exports.SearchRequest,
            requestStream: false,
            responseType: exports.SearchResponse,
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
//# sourceMappingURL=ytsr.js.map