import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";
export declare const protobufPackage = "parser";
export interface Sender {
    id: string;
    name: string;
    displayName: string;
    badges: string[];
}
export interface Channel {
    id: string;
    name: string;
}
export interface Message {
    text: string;
    id: string;
    emotes: Message_Emote[];
}
export interface Message_EmotePosition {
    start: number;
    end: number;
}
export interface Message_Emote {
    name: string;
    id: string;
    count: number;
    positions: Message_EmotePosition[];
}
export interface ProcessCommandRequest {
    sender: Sender | undefined;
    channel: Channel | undefined;
    message: Message | undefined;
}
export interface ProcessCommandResponse {
    responses: string[];
    isReply: boolean;
    keepOrder?: boolean | undefined;
}
export interface GetVariablesResponse {
    list: GetVariablesResponse_Variable[];
}
export interface GetVariablesResponse_Variable {
    name: string;
    example: string;
    description: string;
    visible: boolean;
}
export interface GetDefaultCommandsResponse {
    list: GetDefaultCommandsResponse_DefaultCommand[];
}
export interface GetDefaultCommandsResponse_DefaultCommand {
    name: string;
    description: string;
    visible: boolean;
    rolesNames: string[];
    module: string;
    isReply: boolean;
    keepResponsesOrder: boolean;
    aliases: string[];
}
export interface ParseTextRequestData {
    sender: Sender | undefined;
    channel: Channel | undefined;
    message: Message | undefined;
    parseVariables?: boolean | undefined;
}
export interface ParseTextResponseData {
    responses: string[];
}
export declare const Sender: {
    encode(message: Sender, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Sender;
    fromJSON(object: any): Sender;
    toJSON(message: Sender): unknown;
    create(base?: DeepPartial<Sender>): Sender;
    fromPartial(object: DeepPartial<Sender>): Sender;
};
export declare const Channel: {
    encode(message: Channel, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Channel;
    fromJSON(object: any): Channel;
    toJSON(message: Channel): unknown;
    create(base?: DeepPartial<Channel>): Channel;
    fromPartial(object: DeepPartial<Channel>): Channel;
};
export declare const Message: {
    encode(message: Message, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Message;
    fromJSON(object: any): Message;
    toJSON(message: Message): unknown;
    create(base?: DeepPartial<Message>): Message;
    fromPartial(object: DeepPartial<Message>): Message;
};
export declare const Message_EmotePosition: {
    encode(message: Message_EmotePosition, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Message_EmotePosition;
    fromJSON(object: any): Message_EmotePosition;
    toJSON(message: Message_EmotePosition): unknown;
    create(base?: DeepPartial<Message_EmotePosition>): Message_EmotePosition;
    fromPartial(object: DeepPartial<Message_EmotePosition>): Message_EmotePosition;
};
export declare const Message_Emote: {
    encode(message: Message_Emote, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Message_Emote;
    fromJSON(object: any): Message_Emote;
    toJSON(message: Message_Emote): unknown;
    create(base?: DeepPartial<Message_Emote>): Message_Emote;
    fromPartial(object: DeepPartial<Message_Emote>): Message_Emote;
};
export declare const ProcessCommandRequest: {
    encode(message: ProcessCommandRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ProcessCommandRequest;
    fromJSON(object: any): ProcessCommandRequest;
    toJSON(message: ProcessCommandRequest): unknown;
    create(base?: DeepPartial<ProcessCommandRequest>): ProcessCommandRequest;
    fromPartial(object: DeepPartial<ProcessCommandRequest>): ProcessCommandRequest;
};
export declare const ProcessCommandResponse: {
    encode(message: ProcessCommandResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ProcessCommandResponse;
    fromJSON(object: any): ProcessCommandResponse;
    toJSON(message: ProcessCommandResponse): unknown;
    create(base?: DeepPartial<ProcessCommandResponse>): ProcessCommandResponse;
    fromPartial(object: DeepPartial<ProcessCommandResponse>): ProcessCommandResponse;
};
export declare const GetVariablesResponse: {
    encode(message: GetVariablesResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetVariablesResponse;
    fromJSON(object: any): GetVariablesResponse;
    toJSON(message: GetVariablesResponse): unknown;
    create(base?: DeepPartial<GetVariablesResponse>): GetVariablesResponse;
    fromPartial(object: DeepPartial<GetVariablesResponse>): GetVariablesResponse;
};
export declare const GetVariablesResponse_Variable: {
    encode(message: GetVariablesResponse_Variable, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetVariablesResponse_Variable;
    fromJSON(object: any): GetVariablesResponse_Variable;
    toJSON(message: GetVariablesResponse_Variable): unknown;
    create(base?: DeepPartial<GetVariablesResponse_Variable>): GetVariablesResponse_Variable;
    fromPartial(object: DeepPartial<GetVariablesResponse_Variable>): GetVariablesResponse_Variable;
};
export declare const GetDefaultCommandsResponse: {
    encode(message: GetDefaultCommandsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultCommandsResponse;
    fromJSON(object: any): GetDefaultCommandsResponse;
    toJSON(message: GetDefaultCommandsResponse): unknown;
    create(base?: DeepPartial<GetDefaultCommandsResponse>): GetDefaultCommandsResponse;
    fromPartial(object: DeepPartial<GetDefaultCommandsResponse>): GetDefaultCommandsResponse;
};
export declare const GetDefaultCommandsResponse_DefaultCommand: {
    encode(message: GetDefaultCommandsResponse_DefaultCommand, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultCommandsResponse_DefaultCommand;
    fromJSON(object: any): GetDefaultCommandsResponse_DefaultCommand;
    toJSON(message: GetDefaultCommandsResponse_DefaultCommand): unknown;
    create(base?: DeepPartial<GetDefaultCommandsResponse_DefaultCommand>): GetDefaultCommandsResponse_DefaultCommand;
    fromPartial(object: DeepPartial<GetDefaultCommandsResponse_DefaultCommand>): GetDefaultCommandsResponse_DefaultCommand;
};
export declare const ParseTextRequestData: {
    encode(message: ParseTextRequestData, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ParseTextRequestData;
    fromJSON(object: any): ParseTextRequestData;
    toJSON(message: ParseTextRequestData): unknown;
    create(base?: DeepPartial<ParseTextRequestData>): ParseTextRequestData;
    fromPartial(object: DeepPartial<ParseTextRequestData>): ParseTextRequestData;
};
export declare const ParseTextResponseData: {
    encode(message: ParseTextResponseData, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ParseTextResponseData;
    fromJSON(object: any): ParseTextResponseData;
    toJSON(message: ParseTextResponseData): unknown;
    create(base?: DeepPartial<ParseTextResponseData>): ParseTextResponseData;
    fromPartial(object: DeepPartial<ParseTextResponseData>): ParseTextResponseData;
};
export type ParserDefinition = typeof ParserDefinition;
export declare const ParserDefinition: {
    readonly name: "Parser";
    readonly fullName: "parser.Parser";
    readonly methods: {
        readonly processCommand: {
            readonly name: "ProcessCommand";
            readonly requestType: {
                encode(message: ProcessCommandRequest, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): ProcessCommandRequest;
                fromJSON(object: any): ProcessCommandRequest;
                toJSON(message: ProcessCommandRequest): unknown;
                create(base?: DeepPartial<ProcessCommandRequest>): ProcessCommandRequest;
                fromPartial(object: DeepPartial<ProcessCommandRequest>): ProcessCommandRequest;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(message: ProcessCommandResponse, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): ProcessCommandResponse;
                fromJSON(object: any): ProcessCommandResponse;
                toJSON(message: ProcessCommandResponse): unknown;
                create(base?: DeepPartial<ProcessCommandResponse>): ProcessCommandResponse;
                fromPartial(object: DeepPartial<ProcessCommandResponse>): ProcessCommandResponse;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly parseTextResponse: {
            readonly name: "ParseTextResponse";
            readonly requestType: {
                encode(message: ParseTextRequestData, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): ParseTextRequestData;
                fromJSON(object: any): ParseTextRequestData;
                toJSON(message: ParseTextRequestData): unknown;
                create(base?: DeepPartial<ParseTextRequestData>): ParseTextRequestData;
                fromPartial(object: DeepPartial<ParseTextRequestData>): ParseTextRequestData;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(message: ParseTextResponseData, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): ParseTextResponseData;
                fromJSON(object: any): ParseTextResponseData;
                toJSON(message: ParseTextResponseData): unknown;
                create(base?: DeepPartial<ParseTextResponseData>): ParseTextResponseData;
                fromPartial(object: DeepPartial<ParseTextResponseData>): ParseTextResponseData;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly getDefaultCommands: {
            readonly name: "GetDefaultCommands";
            readonly requestType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(message: GetDefaultCommandsResponse, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): GetDefaultCommandsResponse;
                fromJSON(object: any): GetDefaultCommandsResponse;
                toJSON(message: GetDefaultCommandsResponse): unknown;
                create(base?: DeepPartial<GetDefaultCommandsResponse>): GetDefaultCommandsResponse;
                fromPartial(object: DeepPartial<GetDefaultCommandsResponse>): GetDefaultCommandsResponse;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly getDefaultVariables: {
            readonly name: "GetDefaultVariables";
            readonly requestType: {
                encode(_: Empty, writer?: _m0.Writer): _m0.Writer;
                decode(input: Uint8Array | _m0.Reader, length?: number): Empty;
                fromJSON(_: any): Empty;
                toJSON(_: Empty): unknown;
                create(base?: {}): Empty;
                fromPartial(_: {}): Empty;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(message: GetVariablesResponse, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): GetVariablesResponse;
                fromJSON(object: any): GetVariablesResponse;
                toJSON(message: GetVariablesResponse): unknown;
                create(base?: DeepPartial<GetVariablesResponse>): GetVariablesResponse;
                fromPartial(object: DeepPartial<GetVariablesResponse>): GetVariablesResponse;
            };
            readonly responseStream: false;
            readonly options: {};
        };
    };
};
export interface ParserServiceImplementation<CallContextExt = {}> {
    processCommand(request: ProcessCommandRequest, context: CallContext & CallContextExt): Promise<DeepPartial<ProcessCommandResponse>>;
    parseTextResponse(request: ParseTextRequestData, context: CallContext & CallContextExt): Promise<DeepPartial<ParseTextResponseData>>;
    getDefaultCommands(request: Empty, context: CallContext & CallContextExt): Promise<DeepPartial<GetDefaultCommandsResponse>>;
    getDefaultVariables(request: Empty, context: CallContext & CallContextExt): Promise<DeepPartial<GetVariablesResponse>>;
}
export interface ParserClient<CallOptionsExt = {}> {
    processCommand(request: DeepPartial<ProcessCommandRequest>, options?: CallOptions & CallOptionsExt): Promise<ProcessCommandResponse>;
    parseTextResponse(request: DeepPartial<ParseTextRequestData>, options?: CallOptions & CallOptionsExt): Promise<ParseTextResponseData>;
    getDefaultCommands(request: DeepPartial<Empty>, options?: CallOptions & CallOptionsExt): Promise<GetDefaultCommandsResponse>;
    getDefaultVariables(request: DeepPartial<Empty>, options?: CallOptions & CallOptionsExt): Promise<GetVariablesResponse>;
}
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
//# sourceMappingURL=parser.d.ts.map