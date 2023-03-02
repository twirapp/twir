/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "timers";

export interface Request {
  timerId: string;
}

function createBaseRequest(): Request {
  return { timerId: "" };
}

export const Request = {
  encode(message: Request, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.timerId !== "") {
      writer.uint32(10).string(message.timerId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Request {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromJSON(object: any): Request {
    return { timerId: isSet(object.timerId) ? String(object.timerId) : "" };
  },

  toJSON(message: Request): unknown {
    const obj: any = {};
    message.timerId !== undefined && (obj.timerId = message.timerId);
    return obj;
  },

  create(base?: DeepPartial<Request>): Request {
    return Request.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Request>): Request {
    const message = createBaseRequest();
    message.timerId = object.timerId ?? "";
    return message;
  },
};

export type TimersDefinition = typeof TimersDefinition;
export const TimersDefinition = {
  name: "Timers",
  fullName: "timers.Timers",
  methods: {
    addTimerToQueue: {
      name: "AddTimerToQueue",
      requestType: Request,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    removeTimerFromQueue: {
      name: "RemoveTimerFromQueue",
      requestType: Request,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface TimersServiceImplementation<CallContextExt = {}> {
  addTimerToQueue(request: Request, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  removeTimerFromQueue(request: Request, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
}

export interface TimersClient<CallOptionsExt = {}> {
  addTimerToQueue(request: DeepPartial<Request>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  removeTimerFromQueue(request: DeepPartial<Request>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
