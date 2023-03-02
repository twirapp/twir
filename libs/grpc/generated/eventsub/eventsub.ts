/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "eventsub";

export interface SubscribeToEventsRequest {
  channelId: string;
}

function createBaseSubscribeToEventsRequest(): SubscribeToEventsRequest {
  return { channelId: "" };
}

export const SubscribeToEventsRequest = {
  encode(message: SubscribeToEventsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.channelId !== "") {
      writer.uint32(10).string(message.channelId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubscribeToEventsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubscribeToEventsRequest();
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

  fromJSON(object: any): SubscribeToEventsRequest {
    return { channelId: isSet(object.channelId) ? String(object.channelId) : "" };
  },

  toJSON(message: SubscribeToEventsRequest): unknown {
    const obj: any = {};
    message.channelId !== undefined && (obj.channelId = message.channelId);
    return obj;
  },

  create(base?: DeepPartial<SubscribeToEventsRequest>): SubscribeToEventsRequest {
    return SubscribeToEventsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SubscribeToEventsRequest>): SubscribeToEventsRequest {
    const message = createBaseSubscribeToEventsRequest();
    message.channelId = object.channelId ?? "";
    return message;
  },
};

export type EventSubDefinition = typeof EventSubDefinition;
export const EventSubDefinition = {
  name: "EventSub",
  fullName: "eventsub.EventSub",
  methods: {
    subscribeToEvents: {
      name: "SubscribeToEvents",
      requestType: SubscribeToEventsRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface EventSubServiceImplementation<CallContextExt = {}> {
  subscribeToEvents(
    request: SubscribeToEventsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<Empty>>;
}

export interface EventSubClient<CallOptionsExt = {}> {
  subscribeToEvents(
    request: DeepPartial<SubscribeToEventsRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<Empty>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
