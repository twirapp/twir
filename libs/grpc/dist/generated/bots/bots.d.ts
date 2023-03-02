import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";
export declare const protobufPackage = "bots";
export interface DeleteMessagesRequest {
    channelId: string;
    channelName: string;
    messageIds: string[];
}
export interface SendMessageRequest {
    channelId: string;
    channelName?: string | undefined;
    message: string;
    isAnnounce?: boolean | undefined;
}
export interface JoinOrLeaveRequest {
    botId: string;
    userName: string;
}
export declare const DeleteMessagesRequest: {
    encode(message: DeleteMessagesRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DeleteMessagesRequest;
    fromJSON(object: any): DeleteMessagesRequest;
    toJSON(message: DeleteMessagesRequest): unknown;
    create(base?: DeepPartial<DeleteMessagesRequest>): DeleteMessagesRequest;
    fromPartial(object: DeepPartial<DeleteMessagesRequest>): DeleteMessagesRequest;
};
export declare const SendMessageRequest: {
    encode(message: SendMessageRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SendMessageRequest;
    fromJSON(object: any): SendMessageRequest;
    toJSON(message: SendMessageRequest): unknown;
    create(base?: DeepPartial<SendMessageRequest>): SendMessageRequest;
    fromPartial(object: DeepPartial<SendMessageRequest>): SendMessageRequest;
};
export declare const JoinOrLeaveRequest: {
    encode(message: JoinOrLeaveRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): JoinOrLeaveRequest;
    fromJSON(object: any): JoinOrLeaveRequest;
    toJSON(message: JoinOrLeaveRequest): unknown;
    create(base?: DeepPartial<JoinOrLeaveRequest>): JoinOrLeaveRequest;
    fromPartial(object: DeepPartial<JoinOrLeaveRequest>): JoinOrLeaveRequest;
};
export type BotsDefinition = typeof BotsDefinition;
export declare const BotsDefinition: {
    readonly name: "Bots";
    readonly fullName: "bots.Bots";
    readonly methods: {
        readonly deleteMessage: {
            readonly name: "DeleteMessage";
            readonly requestType: {
                encode(message: DeleteMessagesRequest, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): DeleteMessagesRequest;
                fromJSON(object: any): DeleteMessagesRequest;
                toJSON(message: DeleteMessagesRequest): unknown;
                create(base?: DeepPartial<DeleteMessagesRequest>): DeleteMessagesRequest;
                fromPartial(object: DeepPartial<DeleteMessagesRequest>): DeleteMessagesRequest;
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
        readonly sendMessage: {
            readonly name: "SendMessage";
            readonly requestType: {
                encode(message: SendMessageRequest, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): SendMessageRequest;
                fromJSON(object: any): SendMessageRequest;
                toJSON(message: SendMessageRequest): unknown;
                create(base?: DeepPartial<SendMessageRequest>): SendMessageRequest;
                fromPartial(object: DeepPartial<SendMessageRequest>): SendMessageRequest;
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
        readonly join: {
            readonly name: "Join";
            readonly requestType: {
                encode(message: JoinOrLeaveRequest, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): JoinOrLeaveRequest;
                fromJSON(object: any): JoinOrLeaveRequest;
                toJSON(message: JoinOrLeaveRequest): unknown;
                create(base?: DeepPartial<JoinOrLeaveRequest>): JoinOrLeaveRequest;
                fromPartial(object: DeepPartial<JoinOrLeaveRequest>): JoinOrLeaveRequest;
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
        readonly leave: {
            readonly name: "Leave";
            readonly requestType: {
                encode(message: JoinOrLeaveRequest, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): JoinOrLeaveRequest;
                fromJSON(object: any): JoinOrLeaveRequest;
                toJSON(message: JoinOrLeaveRequest): unknown;
                create(base?: DeepPartial<JoinOrLeaveRequest>): JoinOrLeaveRequest;
                fromPartial(object: DeepPartial<JoinOrLeaveRequest>): JoinOrLeaveRequest;
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
export interface BotsServiceImplementation<CallContextExt = {}> {
    deleteMessage(request: DeleteMessagesRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    sendMessage(request: SendMessageRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    join(request: JoinOrLeaveRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    leave(request: JoinOrLeaveRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
}
export interface BotsClient<CallOptionsExt = {}> {
    deleteMessage(request: DeepPartial<DeleteMessagesRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    sendMessage(request: DeepPartial<SendMessageRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    join(request: DeepPartial<JoinOrLeaveRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    leave(request: DeepPartial<JoinOrLeaveRequest>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
}
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
//# sourceMappingURL=bots.d.ts.map