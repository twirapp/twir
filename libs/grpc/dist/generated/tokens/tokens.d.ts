import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";
export declare const protobufPackage = "tokens";
export interface Token {
    accessToken: string;
    scopes: string[];
}
export interface GetUserTokenRequest {
    userId: string;
}
export interface GetBotTokenRequest {
    botId: string;
}
export interface UpdateToken {
    accessToken: string;
    refreshToken: string;
    expiresIn: number;
    scopes: string[];
}
export declare const Token: {
    encode(message: Token, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Token;
    fromJSON(object: any): Token;
    toJSON(message: Token): unknown;
    create(base?: DeepPartial<Token>): Token;
    fromPartial(object: DeepPartial<Token>): Token;
};
export declare const GetUserTokenRequest: {
    encode(message: GetUserTokenRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetUserTokenRequest;
    fromJSON(object: any): GetUserTokenRequest;
    toJSON(message: GetUserTokenRequest): unknown;
    create(base?: DeepPartial<GetUserTokenRequest>): GetUserTokenRequest;
    fromPartial(object: DeepPartial<GetUserTokenRequest>): GetUserTokenRequest;
};
export declare const GetBotTokenRequest: {
    encode(message: GetBotTokenRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetBotTokenRequest;
    fromJSON(object: any): GetBotTokenRequest;
    toJSON(message: GetBotTokenRequest): unknown;
    create(base?: DeepPartial<GetBotTokenRequest>): GetBotTokenRequest;
    fromPartial(object: DeepPartial<GetBotTokenRequest>): GetBotTokenRequest;
};
export declare const UpdateToken: {
    encode(message: UpdateToken, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): UpdateToken;
    fromJSON(object: any): UpdateToken;
    toJSON(message: UpdateToken): unknown;
    create(base?: DeepPartial<UpdateToken>): UpdateToken;
    fromPartial(object: DeepPartial<UpdateToken>): UpdateToken;
};
export type TokensDefinition = typeof TokensDefinition;
export declare const TokensDefinition: {
    readonly name: "Tokens";
    readonly fullName: "tokens.Tokens";
    readonly methods: {
        readonly requestAppToken: {
            readonly name: "RequestAppToken";
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
                encode(message: Token, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): Token;
                fromJSON(object: any): Token;
                toJSON(message: Token): unknown;
                create(base?: DeepPartial<Token>): Token;
                fromPartial(object: DeepPartial<Token>): Token;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly requestUserToken: {
            readonly name: "RequestUserToken";
            readonly requestType: {
                encode(message: GetUserTokenRequest, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): GetUserTokenRequest;
                fromJSON(object: any): GetUserTokenRequest;
                toJSON(message: GetUserTokenRequest): unknown;
                create(base?: DeepPartial<GetUserTokenRequest>): GetUserTokenRequest;
                fromPartial(object: DeepPartial<GetUserTokenRequest>): GetUserTokenRequest;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(message: Token, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): Token;
                fromJSON(object: any): Token;
                toJSON(message: Token): unknown;
                create(base?: DeepPartial<Token>): Token;
                fromPartial(object: DeepPartial<Token>): Token;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly requestBotToken: {
            readonly name: "RequestBotToken";
            readonly requestType: {
                encode(message: GetBotTokenRequest, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): GetBotTokenRequest;
                fromJSON(object: any): GetBotTokenRequest;
                toJSON(message: GetBotTokenRequest): unknown;
                create(base?: DeepPartial<GetBotTokenRequest>): GetBotTokenRequest;
                fromPartial(object: DeepPartial<GetBotTokenRequest>): GetBotTokenRequest;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(message: Token, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): Token;
                fromJSON(object: any): Token;
                toJSON(message: Token): unknown;
                create(base?: DeepPartial<Token>): Token;
                fromPartial(object: DeepPartial<Token>): Token;
            };
            readonly responseStream: false;
            readonly options: {};
        };
        readonly updateBotToken: {
            readonly name: "UpdateBotToken";
            readonly requestType: {
                encode(message: UpdateToken, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): UpdateToken;
                fromJSON(object: any): UpdateToken;
                toJSON(message: UpdateToken): unknown;
                create(base?: DeepPartial<UpdateToken>): UpdateToken;
                fromPartial(object: DeepPartial<UpdateToken>): UpdateToken;
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
        readonly updateUserToken: {
            readonly name: "UpdateUserToken";
            readonly requestType: {
                encode(message: UpdateToken, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): UpdateToken;
                fromJSON(object: any): UpdateToken;
                toJSON(message: UpdateToken): unknown;
                create(base?: DeepPartial<UpdateToken>): UpdateToken;
                fromPartial(object: DeepPartial<UpdateToken>): UpdateToken;
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
export interface TokensServiceImplementation<CallContextExt = {}> {
    requestAppToken(request: Empty, context: CallContext & CallContextExt): Promise<DeepPartial<Token>>;
    requestUserToken(request: GetUserTokenRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Token>>;
    requestBotToken(request: GetBotTokenRequest, context: CallContext & CallContextExt): Promise<DeepPartial<Token>>;
    updateBotToken(request: UpdateToken, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    updateUserToken(request: UpdateToken, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
}
export interface TokensClient<CallOptionsExt = {}> {
    requestAppToken(request: DeepPartial<Empty>, options?: CallOptions & CallOptionsExt): Promise<Token>;
    requestUserToken(request: DeepPartial<GetUserTokenRequest>, options?: CallOptions & CallOptionsExt): Promise<Token>;
    requestBotToken(request: DeepPartial<GetBotTokenRequest>, options?: CallOptions & CallOptionsExt): Promise<Token>;
    updateBotToken(request: DeepPartial<UpdateToken>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    updateUserToken(request: DeepPartial<UpdateToken>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
}
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
//# sourceMappingURL=tokens.d.ts.map