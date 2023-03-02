import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Empty } from "./google/protobuf/empty";
export declare const protobufPackage = "events";
export interface BaseInfo {
    channelId: string;
}
export interface FollowMessage {
    baseInfo: BaseInfo | undefined;
    userName: string;
    userDisplayName: string;
    userId: string;
}
export interface SubscribeMessage {
    baseInfo: BaseInfo | undefined;
    userName: string;
    userDisplayName: string;
    level: string;
    userId: string;
}
export interface SubGiftMessage {
    baseInfo: BaseInfo | undefined;
    senderUserName: string;
    senderDisplayName: string;
    targetUserName: string;
    targetDisplayName: string;
    level: string;
    senderUserId: string;
}
export interface ReSubscribeMessage {
    baseInfo: BaseInfo | undefined;
    userName: string;
    userDisplayName: string;
    months: number;
    streak: number;
    isPrime: boolean;
    message: string;
    level: string;
    userId: string;
}
export interface RedemptionCreatedMessage {
    baseInfo: BaseInfo | undefined;
    userName: string;
    userDisplayName: string;
    id: string;
    rewardName: string;
    rewardCost: string;
    input?: string | undefined;
    userId: string;
}
export interface CommandUsedMessage {
    baseInfo: BaseInfo | undefined;
    commandId: string;
    commandName: string;
    userName: string;
    userDisplayName: string;
    commandInput: string;
    userId: string;
}
export interface FirstUserMessageMessage {
    baseInfo: BaseInfo | undefined;
    userId: string;
    userName: string;
    userDisplayName: string;
}
export interface RaidedMessage {
    baseInfo: BaseInfo | undefined;
    userName: string;
    userDisplayName: string;
    viewers: number;
    userId: string;
}
export interface TitleOrCategoryChangedMessage {
    baseInfo: BaseInfo | undefined;
    oldTitle: string;
    newTitle: string;
    oldCategory: string;
    newCategory: string;
}
export interface StreamOnlineMessage {
    baseInfo: BaseInfo | undefined;
    title: string;
    category: string;
}
export interface StreamOfflineMessage {
    baseInfo: BaseInfo | undefined;
}
export interface ChatClearMessage {
    baseInfo: BaseInfo | undefined;
}
export interface DonateMessage {
    baseInfo: BaseInfo | undefined;
    userName: string;
    amount: string;
    currency: string;
    message: string;
}
export interface KeywordMatchedMessage {
    baseInfo: BaseInfo | undefined;
    keywordId: string;
    keywordName: string;
    keywordResponse: string;
    userId: string;
    userName: string;
    userDisplayName: string;
}
export interface GreetingSendedMessage {
    baseInfo: BaseInfo | undefined;
    userId: string;
    userName: string;
    userDisplayName: string;
    greetingText: string;
}
export declare const BaseInfo: {
    encode(message: BaseInfo, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): BaseInfo;
    fromJSON(object: any): BaseInfo;
    toJSON(message: BaseInfo): unknown;
    create(base?: DeepPartial<BaseInfo>): BaseInfo;
    fromPartial(object: DeepPartial<BaseInfo>): BaseInfo;
};
export declare const FollowMessage: {
    encode(message: FollowMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): FollowMessage;
    fromJSON(object: any): FollowMessage;
    toJSON(message: FollowMessage): unknown;
    create(base?: DeepPartial<FollowMessage>): FollowMessage;
    fromPartial(object: DeepPartial<FollowMessage>): FollowMessage;
};
export declare const SubscribeMessage: {
    encode(message: SubscribeMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SubscribeMessage;
    fromJSON(object: any): SubscribeMessage;
    toJSON(message: SubscribeMessage): unknown;
    create(base?: DeepPartial<SubscribeMessage>): SubscribeMessage;
    fromPartial(object: DeepPartial<SubscribeMessage>): SubscribeMessage;
};
export declare const SubGiftMessage: {
    encode(message: SubGiftMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SubGiftMessage;
    fromJSON(object: any): SubGiftMessage;
    toJSON(message: SubGiftMessage): unknown;
    create(base?: DeepPartial<SubGiftMessage>): SubGiftMessage;
    fromPartial(object: DeepPartial<SubGiftMessage>): SubGiftMessage;
};
export declare const ReSubscribeMessage: {
    encode(message: ReSubscribeMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ReSubscribeMessage;
    fromJSON(object: any): ReSubscribeMessage;
    toJSON(message: ReSubscribeMessage): unknown;
    create(base?: DeepPartial<ReSubscribeMessage>): ReSubscribeMessage;
    fromPartial(object: DeepPartial<ReSubscribeMessage>): ReSubscribeMessage;
};
export declare const RedemptionCreatedMessage: {
    encode(message: RedemptionCreatedMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): RedemptionCreatedMessage;
    fromJSON(object: any): RedemptionCreatedMessage;
    toJSON(message: RedemptionCreatedMessage): unknown;
    create(base?: DeepPartial<RedemptionCreatedMessage>): RedemptionCreatedMessage;
    fromPartial(object: DeepPartial<RedemptionCreatedMessage>): RedemptionCreatedMessage;
};
export declare const CommandUsedMessage: {
    encode(message: CommandUsedMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): CommandUsedMessage;
    fromJSON(object: any): CommandUsedMessage;
    toJSON(message: CommandUsedMessage): unknown;
    create(base?: DeepPartial<CommandUsedMessage>): CommandUsedMessage;
    fromPartial(object: DeepPartial<CommandUsedMessage>): CommandUsedMessage;
};
export declare const FirstUserMessageMessage: {
    encode(message: FirstUserMessageMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): FirstUserMessageMessage;
    fromJSON(object: any): FirstUserMessageMessage;
    toJSON(message: FirstUserMessageMessage): unknown;
    create(base?: DeepPartial<FirstUserMessageMessage>): FirstUserMessageMessage;
    fromPartial(object: DeepPartial<FirstUserMessageMessage>): FirstUserMessageMessage;
};
export declare const RaidedMessage: {
    encode(message: RaidedMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): RaidedMessage;
    fromJSON(object: any): RaidedMessage;
    toJSON(message: RaidedMessage): unknown;
    create(base?: DeepPartial<RaidedMessage>): RaidedMessage;
    fromPartial(object: DeepPartial<RaidedMessage>): RaidedMessage;
};
export declare const TitleOrCategoryChangedMessage: {
    encode(message: TitleOrCategoryChangedMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TitleOrCategoryChangedMessage;
    fromJSON(object: any): TitleOrCategoryChangedMessage;
    toJSON(message: TitleOrCategoryChangedMessage): unknown;
    create(base?: DeepPartial<TitleOrCategoryChangedMessage>): TitleOrCategoryChangedMessage;
    fromPartial(object: DeepPartial<TitleOrCategoryChangedMessage>): TitleOrCategoryChangedMessage;
};
export declare const StreamOnlineMessage: {
    encode(message: StreamOnlineMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): StreamOnlineMessage;
    fromJSON(object: any): StreamOnlineMessage;
    toJSON(message: StreamOnlineMessage): unknown;
    create(base?: DeepPartial<StreamOnlineMessage>): StreamOnlineMessage;
    fromPartial(object: DeepPartial<StreamOnlineMessage>): StreamOnlineMessage;
};
export declare const StreamOfflineMessage: {
    encode(message: StreamOfflineMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): StreamOfflineMessage;
    fromJSON(object: any): StreamOfflineMessage;
    toJSON(message: StreamOfflineMessage): unknown;
    create(base?: DeepPartial<StreamOfflineMessage>): StreamOfflineMessage;
    fromPartial(object: DeepPartial<StreamOfflineMessage>): StreamOfflineMessage;
};
export declare const ChatClearMessage: {
    encode(message: ChatClearMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ChatClearMessage;
    fromJSON(object: any): ChatClearMessage;
    toJSON(message: ChatClearMessage): unknown;
    create(base?: DeepPartial<ChatClearMessage>): ChatClearMessage;
    fromPartial(object: DeepPartial<ChatClearMessage>): ChatClearMessage;
};
export declare const DonateMessage: {
    encode(message: DonateMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DonateMessage;
    fromJSON(object: any): DonateMessage;
    toJSON(message: DonateMessage): unknown;
    create(base?: DeepPartial<DonateMessage>): DonateMessage;
    fromPartial(object: DeepPartial<DonateMessage>): DonateMessage;
};
export declare const KeywordMatchedMessage: {
    encode(message: KeywordMatchedMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): KeywordMatchedMessage;
    fromJSON(object: any): KeywordMatchedMessage;
    toJSON(message: KeywordMatchedMessage): unknown;
    create(base?: DeepPartial<KeywordMatchedMessage>): KeywordMatchedMessage;
    fromPartial(object: DeepPartial<KeywordMatchedMessage>): KeywordMatchedMessage;
};
export declare const GreetingSendedMessage: {
    encode(message: GreetingSendedMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GreetingSendedMessage;
    fromJSON(object: any): GreetingSendedMessage;
    toJSON(message: GreetingSendedMessage): unknown;
    create(base?: DeepPartial<GreetingSendedMessage>): GreetingSendedMessage;
    fromPartial(object: DeepPartial<GreetingSendedMessage>): GreetingSendedMessage;
};
export type EventsDefinition = typeof EventsDefinition;
export declare const EventsDefinition: {
    readonly name: "Events";
    readonly fullName: "events.Events";
    readonly methods: {
        readonly follow: {
            readonly name: "Follow";
            readonly requestType: {
                encode(message: FollowMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): FollowMessage;
                fromJSON(object: any): FollowMessage;
                toJSON(message: FollowMessage): unknown;
                create(base?: DeepPartial<FollowMessage>): FollowMessage;
                fromPartial(object: DeepPartial<FollowMessage>): FollowMessage;
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
        readonly subscribe: {
            readonly name: "Subscribe";
            readonly requestType: {
                encode(message: SubscribeMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): SubscribeMessage;
                fromJSON(object: any): SubscribeMessage;
                toJSON(message: SubscribeMessage): unknown;
                create(base?: DeepPartial<SubscribeMessage>): SubscribeMessage;
                fromPartial(object: DeepPartial<SubscribeMessage>): SubscribeMessage;
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
        readonly subGift: {
            readonly name: "SubGift";
            readonly requestType: {
                encode(message: SubGiftMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): SubGiftMessage;
                fromJSON(object: any): SubGiftMessage;
                toJSON(message: SubGiftMessage): unknown;
                create(base?: DeepPartial<SubGiftMessage>): SubGiftMessage;
                fromPartial(object: DeepPartial<SubGiftMessage>): SubGiftMessage;
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
        readonly reSubscribe: {
            readonly name: "ReSubscribe";
            readonly requestType: {
                encode(message: ReSubscribeMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): ReSubscribeMessage;
                fromJSON(object: any): ReSubscribeMessage;
                toJSON(message: ReSubscribeMessage): unknown;
                create(base?: DeepPartial<ReSubscribeMessage>): ReSubscribeMessage;
                fromPartial(object: DeepPartial<ReSubscribeMessage>): ReSubscribeMessage;
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
        readonly redemptionCreated: {
            readonly name: "RedemptionCreated";
            readonly requestType: {
                encode(message: RedemptionCreatedMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): RedemptionCreatedMessage;
                fromJSON(object: any): RedemptionCreatedMessage;
                toJSON(message: RedemptionCreatedMessage): unknown;
                create(base?: DeepPartial<RedemptionCreatedMessage>): RedemptionCreatedMessage;
                fromPartial(object: DeepPartial<RedemptionCreatedMessage>): RedemptionCreatedMessage;
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
        readonly commandUsed: {
            readonly name: "CommandUsed";
            readonly requestType: {
                encode(message: CommandUsedMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): CommandUsedMessage;
                fromJSON(object: any): CommandUsedMessage;
                toJSON(message: CommandUsedMessage): unknown;
                create(base?: DeepPartial<CommandUsedMessage>): CommandUsedMessage;
                fromPartial(object: DeepPartial<CommandUsedMessage>): CommandUsedMessage;
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
        readonly firstUserMessage: {
            readonly name: "FirstUserMessage";
            readonly requestType: {
                encode(message: FirstUserMessageMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): FirstUserMessageMessage;
                fromJSON(object: any): FirstUserMessageMessage;
                toJSON(message: FirstUserMessageMessage): unknown;
                create(base?: DeepPartial<FirstUserMessageMessage>): FirstUserMessageMessage;
                fromPartial(object: DeepPartial<FirstUserMessageMessage>): FirstUserMessageMessage;
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
        readonly raided: {
            readonly name: "Raided";
            readonly requestType: {
                encode(message: RaidedMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): RaidedMessage;
                fromJSON(object: any): RaidedMessage;
                toJSON(message: RaidedMessage): unknown;
                create(base?: DeepPartial<RaidedMessage>): RaidedMessage;
                fromPartial(object: DeepPartial<RaidedMessage>): RaidedMessage;
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
        readonly titleOrCategoryChanged: {
            readonly name: "TitleOrCategoryChanged";
            readonly requestType: {
                encode(message: TitleOrCategoryChangedMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): TitleOrCategoryChangedMessage;
                fromJSON(object: any): TitleOrCategoryChangedMessage;
                toJSON(message: TitleOrCategoryChangedMessage): unknown;
                create(base?: DeepPartial<TitleOrCategoryChangedMessage>): TitleOrCategoryChangedMessage;
                fromPartial(object: DeepPartial<TitleOrCategoryChangedMessage>): TitleOrCategoryChangedMessage;
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
        readonly streamOnline: {
            readonly name: "StreamOnline";
            readonly requestType: {
                encode(message: StreamOnlineMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): StreamOnlineMessage;
                fromJSON(object: any): StreamOnlineMessage;
                toJSON(message: StreamOnlineMessage): unknown;
                create(base?: DeepPartial<StreamOnlineMessage>): StreamOnlineMessage;
                fromPartial(object: DeepPartial<StreamOnlineMessage>): StreamOnlineMessage;
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
        readonly streamOffline: {
            readonly name: "StreamOffline";
            readonly requestType: {
                encode(message: StreamOfflineMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): StreamOfflineMessage;
                fromJSON(object: any): StreamOfflineMessage;
                toJSON(message: StreamOfflineMessage): unknown;
                create(base?: DeepPartial<StreamOfflineMessage>): StreamOfflineMessage;
                fromPartial(object: DeepPartial<StreamOfflineMessage>): StreamOfflineMessage;
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
        readonly chatClear: {
            readonly name: "ChatClear";
            readonly requestType: {
                encode(message: ChatClearMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): ChatClearMessage;
                fromJSON(object: any): ChatClearMessage;
                toJSON(message: ChatClearMessage): unknown;
                create(base?: DeepPartial<ChatClearMessage>): ChatClearMessage;
                fromPartial(object: DeepPartial<ChatClearMessage>): ChatClearMessage;
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
        readonly donate: {
            readonly name: "Donate";
            readonly requestType: {
                encode(message: DonateMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): DonateMessage;
                fromJSON(object: any): DonateMessage;
                toJSON(message: DonateMessage): unknown;
                create(base?: DeepPartial<DonateMessage>): DonateMessage;
                fromPartial(object: DeepPartial<DonateMessage>): DonateMessage;
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
        readonly keywordMatched: {
            readonly name: "KeywordMatched";
            readonly requestType: {
                encode(message: KeywordMatchedMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): KeywordMatchedMessage;
                fromJSON(object: any): KeywordMatchedMessage;
                toJSON(message: KeywordMatchedMessage): unknown;
                create(base?: DeepPartial<KeywordMatchedMessage>): KeywordMatchedMessage;
                fromPartial(object: DeepPartial<KeywordMatchedMessage>): KeywordMatchedMessage;
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
        readonly greetingSended: {
            readonly name: "GreetingSended";
            readonly requestType: {
                encode(message: GreetingSendedMessage, writer?: _m0.Writer): _m0.Writer;
                decode(input: _m0.Reader | Uint8Array, length?: number): GreetingSendedMessage;
                fromJSON(object: any): GreetingSendedMessage;
                toJSON(message: GreetingSendedMessage): unknown;
                create(base?: DeepPartial<GreetingSendedMessage>): GreetingSendedMessage;
                fromPartial(object: DeepPartial<GreetingSendedMessage>): GreetingSendedMessage;
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
export interface EventsServiceImplementation<CallContextExt = {}> {
    follow(request: FollowMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    subscribe(request: SubscribeMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    subGift(request: SubGiftMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    reSubscribe(request: ReSubscribeMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    redemptionCreated(request: RedemptionCreatedMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    commandUsed(request: CommandUsedMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    firstUserMessage(request: FirstUserMessageMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    raided(request: RaidedMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    titleOrCategoryChanged(request: TitleOrCategoryChangedMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    streamOnline(request: StreamOnlineMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    streamOffline(request: StreamOfflineMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    chatClear(request: ChatClearMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    donate(request: DonateMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    keywordMatched(request: KeywordMatchedMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
    greetingSended(request: GreetingSendedMessage, context: CallContext & CallContextExt): Promise<DeepPartial<Empty>>;
}
export interface EventsClient<CallOptionsExt = {}> {
    follow(request: DeepPartial<FollowMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    subscribe(request: DeepPartial<SubscribeMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    subGift(request: DeepPartial<SubGiftMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    reSubscribe(request: DeepPartial<ReSubscribeMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    redemptionCreated(request: DeepPartial<RedemptionCreatedMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    commandUsed(request: DeepPartial<CommandUsedMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    firstUserMessage(request: DeepPartial<FirstUserMessageMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    raided(request: DeepPartial<RaidedMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    titleOrCategoryChanged(request: DeepPartial<TitleOrCategoryChangedMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    streamOnline(request: DeepPartial<StreamOnlineMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    streamOffline(request: DeepPartial<StreamOfflineMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    chatClear(request: DeepPartial<ChatClearMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    donate(request: DeepPartial<DonateMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    keywordMatched(request: DeepPartial<KeywordMatchedMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
    greetingSended(request: DeepPartial<GreetingSendedMessage>, options?: CallOptions & CallOptionsExt): Promise<Empty>;
}
type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;
export type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
//# sourceMappingURL=events.d.ts.map