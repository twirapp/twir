/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "eval";

export interface Evaluate {
  script: string;
}

export interface EvaluateResult {
  result: string;
}

function createBaseEvaluate(): Evaluate {
  return { script: "" };
}

export const Evaluate = {
  encode(message: Evaluate, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.script !== "") {
      writer.uint32(10).string(message.script);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Evaluate {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromJSON(object: any): Evaluate {
    return { script: isSet(object.script) ? String(object.script) : "" };
  },

  toJSON(message: Evaluate): unknown {
    const obj: any = {};
    message.script !== undefined && (obj.script = message.script);
    return obj;
  },

  create(base?: DeepPartial<Evaluate>): Evaluate {
    return Evaluate.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Evaluate>): Evaluate {
    const message = createBaseEvaluate();
    message.script = object.script ?? "";
    return message;
  },
};

function createBaseEvaluateResult(): EvaluateResult {
  return { result: "" };
}

export const EvaluateResult = {
  encode(message: EvaluateResult, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.result !== "") {
      writer.uint32(10).string(message.result);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EvaluateResult {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromJSON(object: any): EvaluateResult {
    return { result: isSet(object.result) ? String(object.result) : "" };
  },

  toJSON(message: EvaluateResult): unknown {
    const obj: any = {};
    message.result !== undefined && (obj.result = message.result);
    return obj;
  },

  create(base?: DeepPartial<EvaluateResult>): EvaluateResult {
    return EvaluateResult.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<EvaluateResult>): EvaluateResult {
    const message = createBaseEvaluateResult();
    message.result = object.result ?? "";
    return message;
  },
};

export type EvalDefinition = typeof EvalDefinition;
export const EvalDefinition = {
  name: "Eval",
  fullName: "eval.Eval",
  methods: {
    process: {
      name: "Process",
      requestType: Evaluate,
      requestStream: false,
      responseType: EvaluateResult,
      responseStream: false,
      options: {},
    },
  },
} as const;

export interface EvalServiceImplementation<CallContextExt = {}> {
  process(request: Evaluate, context: CallContext & CallContextExt): Promise<DeepPartial<EvaluateResult>>;
}

export interface EvalClient<CallOptionsExt = {}> {
  process(request: DeepPartial<Evaluate>, options?: CallOptions & CallOptionsExt): Promise<EvaluateResult>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
