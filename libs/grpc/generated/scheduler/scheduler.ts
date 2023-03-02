/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "scheduler";

export interface CreateDefaultCommandsRequest {
  userId: string;
}

function createBaseCreateDefaultCommandsRequest(): CreateDefaultCommandsRequest {
  return { userId: "" };
}

export const CreateDefaultCommandsRequest = {
  encode(message: CreateDefaultCommandsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateDefaultCommandsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromJSON(object: any): CreateDefaultCommandsRequest {
    return { userId: isSet(object.userId) ? String(object.userId) : "" };
  },

  toJSON(message: CreateDefaultCommandsRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<CreateDefaultCommandsRequest>): CreateDefaultCommandsRequest {
    return CreateDefaultCommandsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CreateDefaultCommandsRequest>): CreateDefaultCommandsRequest {
    const message = createBaseCreateDefaultCommandsRequest();
    message.userId = object.userId ?? "";
    return message;
  },
};

export type SchedulerDefinition = typeof SchedulerDefinition;
export const SchedulerDefinition = {
  name: "Scheduler",
  fullName: "scheduler.Scheduler",
  methods: {
    createDefaultCommands: {
      name: "CreateDefaultCommands",
      requestType: CreateDefaultCommandsRequest,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface SchedulerServiceImplementation<CallContextExt = {}> {
  createDefaultCommands(
    request: CreateDefaultCommandsRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<Empty>>;
}

export interface SchedulerClient<CallOptionsExt = {}> {
  createDefaultCommands(
    request: DeepPartial<CreateDefaultCommandsRequest>,
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
