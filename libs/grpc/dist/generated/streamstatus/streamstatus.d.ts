import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
export declare const protobufPackage = "streamstatus";
export interface Evaluate {
    script: string;
}
export interface EvaluateResult {
    result: string;
}
export declare const Evaluate: {
    encode(message: Evaluate, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Evaluate;
    fromJSON(object: any): Evaluate;
    toJSON(message: Evaluate): unknown;
    create(base?: DeepPartial<Evaluate>): Evaluate;
    fromPartial(object: DeepPartial<Evaluate>): Evaluate;
};
export declare const EvaluateResult: {
    encode(message: EvaluateResult, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): EvaluateResult;
    fromJSON(object: any): EvaluateResult;
    toJSON(message: EvaluateResult): unknown;
    create(base?: DeepPartial<EvaluateResult>): EvaluateResult;
    fromPartial(object: DeepPartial<EvaluateResult>): EvaluateResult;
};
export type EvalDefinition = typeof EvalDefinition;
export declare const EvalDefinition: {
    readonly name: "Eval";
    readonly fullName: "streamstatus.Eval";
    readonly methods: {
        readonly process: {
            readonly name: "Process";
            readonly requestType: {
                encode(message: Evaluate, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): Evaluate;
                fromJSON(object: any): Evaluate;
                toJSON(message: Evaluate): unknown;
                create(base?: DeepPartial<Evaluate>): Evaluate;
                fromPartial(object: DeepPartial<Evaluate>): Evaluate;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(message: EvaluateResult, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): EvaluateResult;
                fromJSON(object: any): EvaluateResult;
                toJSON(message: EvaluateResult): unknown;
                create(base?: DeepPartial<EvaluateResult>): EvaluateResult;
                fromPartial(object: DeepPartial<EvaluateResult>): EvaluateResult;
            };
            readonly responseStream: false;
            readonly options: {};
        };
    };
};
export interface EvalServiceImplementation<CallContextExt = {}> {
    process(request: Evaluate, context: CallContext & CallContextExt): Promise<DeepPartial<EvaluateResult>>;
}
export interface EvalClient<CallOptionsExt = {}> {
    process(request: DeepPartial<Evaluate>, options?: CallOptions & CallOptionsExt): Promise<EvaluateResult>;
}
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
//# sourceMappingURL=streamstatus.d.ts.map