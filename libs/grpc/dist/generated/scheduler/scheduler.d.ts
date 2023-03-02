import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";
export declare const protobufPackage = "scheduler";
export interface CreateDefaultCommandsRequest {
    userId: string;
}
export declare const CreateDefaultCommandsRequest: {
    encode(message: CreateDefaultCommandsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): CreateDefaultCommandsRequest;
    fromJSON(object: any): CreateDefaultCommandsRequest;
    toJSON(message: CreateDefaultCommandsRequest): unknown;
    create(base?: DeepPartial<CreateDefaultCommandsRequest>): CreateDefaultCommandsRequest;
    fromPartial(object: DeepPartial<CreateDefaultCommandsRequest>): CreateDefaultCommandsRequest;
};
export type SchedulerDefinition = typeof SchedulerDefinition;
export declare const SchedulerDefinition: {
    readonly name: "Scheduler";
    readonly fullName: "scheduler.Scheduler";
    readonly methods: {
        readonly createDefaultCommands: {
            readonly name: "CreateDefaultCommands";
            readonly requestType: {
                encode(message: CreateDefaultCommandsRequest, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): CreateDefaultCommandsRequest;
                fromJSON(object: any): CreateDefaultCommandsRequest;
                toJSON(message: CreateDefaultCommandsRequest): unknown;
                create(base?: DeepPartial<CreateDefaultCommandsRequest>): CreateDefaultCommandsRequest;
                fromPartial(object: DeepPartial<CreateDefaultCommandsRequest>): CreateDefaultCommandsRequest;
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
export interface SchedulerServiceImplementation<CallContextExt = {}> {
    createDefaultCommands(request: CreateDefaultCommandsRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
}
export interface SchedulerClient<CallOptionsExt = {}> {
    createDefaultCommands(request: DeepPartial<CreateDefaultCommandsRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
}
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
//# sourceMappingURL=scheduler.d.ts.map