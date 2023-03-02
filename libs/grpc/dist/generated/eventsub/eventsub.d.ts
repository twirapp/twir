import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";
export declare const protobufPackage = "eventsub";
export interface SubscribeToEventsRequest {
    channelId: string;
}
export declare const SubscribeToEventsRequest: {
    encode(message: SubscribeToEventsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SubscribeToEventsRequest;
    fromJSON(object: any): SubscribeToEventsRequest;
    toJSON(message: SubscribeToEventsRequest): unknown;
    create(base?: DeepPartial<SubscribeToEventsRequest>): SubscribeToEventsRequest;
    fromPartial(object: DeepPartial<SubscribeToEventsRequest>): SubscribeToEventsRequest;
};
export type EventSubDefinition = typeof EventSubDefinition;
export declare const EventSubDefinition: {
    readonly name: "EventSub";
    readonly fullName: "eventsub.EventSub";
    readonly methods: {
        readonly subscribeToEvents: {
            readonly name: "SubscribeToEvents";
            readonly requestType: {
                encode(message: SubscribeToEventsRequest, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): SubscribeToEventsRequest;
                fromJSON(object: any): SubscribeToEventsRequest;
                toJSON(message: SubscribeToEventsRequest): unknown;
                create(base?: DeepPartial<SubscribeToEventsRequest>): SubscribeToEventsRequest;
                fromPartial(object: DeepPartial<SubscribeToEventsRequest>): SubscribeToEventsRequest;
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
export interface EventSubServiceImplementation<CallContextExt = {}> {
    subscribeToEvents(request: SubscribeToEventsRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
}
export interface EventSubClient<CallOptionsExt = {}> {
    subscribeToEvents(request: DeepPartial<SubscribeToEventsRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
}
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
//# sourceMappingURL=eventsub.d.ts.map