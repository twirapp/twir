/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal.js";
import { Empty } from "../google/protobuf/empty.js";

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
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.id = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Request {
    return { id: isSet(object.id) ? globalThis.String(object.id) : "" };
  },

  toJSON(message: Request): unknown {
    const obj: any = {};
    if (message.id !== "") {
      obj.id = message.id;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Request>, I>>(base?: I): Request {
    return Request.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Request>, I>>(object: I): Request {
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
  : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
