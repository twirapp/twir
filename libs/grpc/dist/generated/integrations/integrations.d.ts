import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";
export declare const protobufPackage = "integrations";
export interface Request {
    id: string;
}
export declare const Request: {
    encode(message: Request, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Request;
    fromJSON(object: any): Request;
    toJSON(message: Request): unknown;
    create(base?: DeepPartial<Request>): Request;
    fromPartial(object: DeepPartial<Request>): Request;
};
export type IntegrationsDefinition = typeof IntegrationsDefinition;
export declare const IntegrationsDefinition: {
    readonly name: "Integrations";
    readonly fullName: "integrations.Integrations";
    readonly methods: {
        readonly addIntegration: {
            readonly name: "AddIntegration";
            readonly requestType: {
                encode(message: Request, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): Request;
                fromJSON(object: any): Request;
                toJSON(message: Request): unknown;
                create(base?: DeepPartial<Request>): Request;
                fromPartial(object: DeepPartial<Request>): Request;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly removeIntegration: {
            readonly name: "RemoveIntegration";
            readonly requestType: {
                encode(message: Request, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): Request;
                fromJSON(object: any): Request;
                toJSON(message: Request): unknown;
                create(base?: DeepPartial<Request>): Request;
                fromPartial(object: DeepPartial<Request>): Request;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly responseStream: false;
            readonly options: {};
        };
    };
};
export interface IntegrationsServiceImplementation<CallContextExt = {}> {
    addIntegration(request: Request, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    removeIntegration(request: Request, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
}
export interface IntegrationsClient<CallOptionsExt = {}> {
    addIntegration(request: DeepPartial<Request>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    removeIntegration(request: DeepPartial<Request>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
}
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
//# sourceMappingURL=integrations.d.ts.map