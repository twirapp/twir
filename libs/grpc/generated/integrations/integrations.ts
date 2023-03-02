/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";

export const protobufPackage = "integrations";

export interface Request {
  id: string;
}

function createBaseRequest(): Request {
  return { id: "" };
}

export const Request = {
  encode(message: Request, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
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
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Request {
    return { id: isSet(object.id) ? String(object.id) : "" };
  },

  toJSON(message: Request): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  create(base?: DeepPartial<Request>): Request {
    return Request.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Request>): Request {
    const message = createBaseRequest();
    message.id = object.id ?? "";
    return message;
  },
};

export type IntegrationsDefinition = typeof IntegrationsDefinition;
export const IntegrationsDefinition = {
  name: "Integrations",
  fullName: "integrations.Integrations",
  methods: {
    addIntegration: {
      name: "AddIntegration",
      requestType: Request,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
    removeIntegration: {
      name: "RemoveIntegration",
      requestType: Request,
      requestStream: false,
      responseType: Empty,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface IntegrationsServiceImplementation<CallContextExt = {}> {
  addIntegration(request: Request, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
  removeIntegration(request: Request, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
}

export interface IntegrationsClient<CallOptionsExt = {}> {
  addIntegration(request: DeepPartial<Request>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
  removeIntegration(request: DeepPartial<Request>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
