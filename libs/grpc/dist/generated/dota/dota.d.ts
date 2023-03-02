import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
export declare const protobufPackage = "dota";
export interface GetPlayerCardRequest {
    accountId: number;
}
export interface GetPlayerCardResponse {
    accountId: string;
    rankTier?: number | undefined;
    leaderboardRank?: number | undefined;
}
export declare const GetPlayerCardRequest: {
    encode(message: GetPlayerCardRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetPlayerCardRequest;
    fromJSON(object: any): GetPlayerCardRequest;
    toJSON(message: GetPlayerCardRequest): unknown;
    create(base?: DeepPartial<GetPlayerCardRequest>): GetPlayerCardRequest;
    fromPartial(object: DeepPartial<GetPlayerCardRequest>): GetPlayerCardRequest;
};
export declare const GetPlayerCardResponse: {
    encode(message: GetPlayerCardResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetPlayerCardResponse;
    fromJSON(object: any): GetPlayerCardResponse;
    toJSON(message: GetPlayerCardResponse): unknown;
    create(base?: DeepPartial<GetPlayerCardResponse>): GetPlayerCardResponse;
    fromPartial(object: DeepPartial<GetPlayerCardResponse>): GetPlayerCardResponse;
};
export type DotaDefinition = typeof DotaDefinition;
export declare const DotaDefinition: {
    readonly name: "Dota";
    readonly fullName: "dota.Dota";
    readonly methods: {
        readonly getPlayerCard: {
            readonly name: "GetPlayerCard";
            readonly requestType: {
                encode(message: GetPlayerCardRequest, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): GetPlayerCardRequest;
                fromJSON(object: any): GetPlayerCardRequest;
                toJSON(message: GetPlayerCardRequest): unknown;
                create(base?: DeepPartial<GetPlayerCardRequest>): GetPlayerCardRequest;
                fromPartial(object: DeepPartial<GetPlayerCardRequest>): GetPlayerCardRequest;
            };
            readonly requestStream: false;
            readonly responseType: {
                encode(message: GetPlayerCardResponse, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): GetPlayerCardResponse;
                fromJSON(object: any): GetPlayerCardResponse;
                toJSON(message: GetPlayerCardResponse): unknown;
                create(base?: DeepPartial<GetPlayerCardResponse>): GetPlayerCardResponse;
                fromPartial(object: DeepPartial<GetPlayerCardResponse>): GetPlayerCardResponse;
            };
            readonly responseStream: false;
            readonly options: {};
        };
    };
};
export interface DotaServiceImplementation<CallContextExt = {}> {
    getPlayerCard(request: GetPlayerCardRequest, context: CallContext & CallContextExt): Promise<DeepPartial<GetPlayerCardResponse>>;
}
export interface DotaClient<CallOptionsExt = {}> {
    getPlayerCard(request: DeepPartial<GetPlayerCardRequest>, options?: CallOptions & CallOptionsExt): Promise<GetPlayerCardResponse>;
}
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
//# sourceMappingURL=dota.d.ts.map