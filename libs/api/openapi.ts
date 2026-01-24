/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface AuthBody {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  /** @minLength 20 */
  code: string;
  state: string;
}

export interface AuthResponseDto {
  redirect_to: string;
}

export interface BadgeWithUsers {
  /** @format int64 */
  ffzSlot: number;
  id: string;
  name: string;
  url: string;
  users: string[];
}

export interface BaseOutputBodyJsonAuthResponseDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  data: AuthResponseDto;
}

export interface BaseOutputBodyJsonCustomDomainOutputDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  data: CustomDomainOutputDto;
}

export interface BaseOutputBodyJsonIntegrationsValorantStatsOutput {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  data: IntegrationsValorantStatsOutput;
}

export interface BaseOutputBodyJsonInterface {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  data: any;
}

export interface BaseOutputBodyJsonLinkOutputDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  data: LinkOutputDto;
}

export interface BaseOutputBodyJsonLinksProfileOutputDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  data: LinksProfileOutputDto;
}

export interface BaseOutputBodyJsonListCommandResponseDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  data: CommandResponseDto[];
}

export interface BaseOutputBodyJsonListCountryStatsDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  data: CountryStatsDto[];
}

export interface BaseOutputBodyJsonListStatisticsPointDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  data: StatisticsPointDto[];
}

export interface BaseOutputBodyJsonPasteBinOutputDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  data: PasteBinOutputDto;
}

export interface BaseOutputBodyJsonProfileResponseDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  data: ProfileResponseDto;
}

export interface BaseOutputBodyJsonScheduledVipOutputDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  data: ScheduledVipOutputDto;
}

export interface BaseOutputBodyJsonStartResponseDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  data: StartResponseDto;
}

export interface BaseOutputBodyJsonStopResponseDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  data: StopResponseDto;
}

export interface BaseOutputBodyJsonStream {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  data: Stream;
}

export interface CommandDto {
  description: string | null;
  group: string | null;
  module: string;
  name: string;
  responses: CommandDtoResponse[];
}

export interface CommandDtoResponse {
  text: string;
}

export interface CommandGroupResponseDto {
  color: string;
  /** @format uuid */
  id: string;
  name: string;
}

export interface CommandResponseDto {
  aliases: string[];
  allowed_users_ids: string[];
  /** @format int64 */
  cooldown: number | null;
  cooldown_type: CommandResponseDtoCooldownTypeEnum;
  default_name: string | null;
  denied_users_ids: string[];
  description: string | null;
  enabled: boolean;
  enabled_categories: string[];
  expire: Expire;
  group: CommandGroupResponseDto;
  id: string;
  is_default: boolean;
  is_reply: boolean;
  keep_responses_order: boolean;
  module: string;
  name: string;
  offline_only: boolean;
  online_only: boolean;
  /** @format int64 */
  required_messages: number;
  /** @format int64 */
  required_used_channel_points: number;
  /** @format int64 */
  required_watch_time: number;
  responses: CommandResponsesResponseDto[];
  roles_cooldowns: CommandRoleCooldownResponseDto[];
  /** @format uuid */
  roles_ids: string[];
  visible: boolean;
}

export interface CommandResponsesResponseDto {
  /** @format uuid */
  id: string;
  offline_only: boolean;
  online_only: boolean;
  /** @format int64 */
  order: number;
  text: string;
  twitch_category_id: string[];
}

export interface CommandRoleCooldownResponseDto {
  /** @format int64 */
  cooldown: number;
  role_id: string;
}

export interface CountryStatsDto {
  /**
   * @format int64
   * @min 0
   */
  count: number;
  country: string;
}

export interface CreateCustomDomainInputBody {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  /**
   * @minLength 3
   * @maxLength 255
   */
  domain: string;
}

export interface CreateLinkInputDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  /**
   * @minLength 3
   * @maxLength 30
   * @pattern ^[a-zA-Z0-9]+$
   */
  alias?: string;
  /**
   * @format uri
   * @minLength 1
   * @maxLength 2000
   * @pattern ^https?://.*
   */
  url: string;
}

export interface CreateRequestDtoBody {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  /**
   * When to remove VIP (for time-based removal)
   * @format date-time
   */
  remove_at?: string | null;
  /** Type of removal: 'time' or 'stream_end' */
  remove_type: CreateRequestDtoBodyRemoveTypeEnum;
  /**
   * Twitch user ID
   * @minLength 1
   * @maxLength 100
   */
  user_id: string;
}

export interface CustomDomainOutputDto {
  /** @format date-time */
  created_at: string;
  domain: string;
  id: string;
  verification_target: string;
  verification_token: string;
  verified: boolean;
}

export interface EndTierStruct {
  /** @format int64 */
  id: number;
  name: string;
}

export interface ErrorDetail {
  /** Where the error occurred, e.g. 'body.items[3].tags' or 'path.thing-id' */
  location?: string;
  /** Error message text */
  message?: string;
  /** The value at the given location */
  value?: any;
}

export interface ErrorModel {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  /** A human-readable explanation specific to this occurrence of the problem. */
  detail?: string;
  /** Optional list of individual error details */
  errors?: ErrorDetail[];
  /**
   * A URI reference that identifies the specific occurrence of the problem.
   * @format uri
   */
  instance?: string;
  /**
   * HTTP status code
   * @format int64
   */
  status?: number;
  /** A short, human-readable summary of the problem type. This value should not change between occurrences of the error. */
  title?: string;
  /**
   * A URI reference to human-readable documentation for the error.
   * @format uri
   * @default "about:blank"
   */
  type?: string;
}

export interface Expire {
  /** @format date-time */
  expires_at: string;
  expires_type: ExpireExpiresTypeEnum;
}

export interface IntegrationsValorantStatsOutput {
  matches: StoredMatchesResponseMatch[];
  mmr: MmrResponseData;
}

export interface Item {
  act_wins: Item[];
  end_tier: EndTierStruct;
  /** @format int64 */
  games: number;
  leaderboard_placement: LeaderboardPlacementStruct;
  ranking_schema: string;
  season: SeasonStruct;
  /** @format int64 */
  wins: number;
}

export interface LeaderboardPlacementStruct {
  /** @format int64 */
  rank: number;
  /** @format date-time */
  updated_at: string;
}

export interface LinkOutputDto {
  /** @format date-time */
  created_at: string;
  id: string;
  short_url: string;
  url: string;
  /** @format int64 */
  views: number;
}

export interface LinksProfileOutputDto {
  items: LinkOutputDto[];
  /** @format int64 */
  total: number;
}

export interface ListResponseDtoBody {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  data: ScheduledVipOutputDto[];
}

export interface MapStruct {
  id: string;
  name: string;
}

export interface MmrResponseData {
  account: MmrResponseDataAccountStruct;
  current: MmrResponseDataCurrentStruct;
  peak: MmrResponseDataPeakStruct;
  seasonal: Item[];
}

export interface MmrResponseDataAccountStruct {
  name: string;
  puuid: string;
  tag: string;
}

export interface MmrResponseDataCurrentStruct {
  /** @format int64 */
  elo: number;
  /** @format int64 */
  games_needed_for_rating: number;
  /** @format int64 */
  last_change: number;
  leaderboard_placement: LeaderboardPlacementStruct;
  /** @format int64 */
  rr: number;
  tier: TierStruct;
}

export interface MmrResponseDataPeakStruct {
  ranking_schema: string;
  season: SeasonStruct;
  tier: TierStruct;
}

export interface PasteBinCreateRequestDtoBody {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  /**
   * @minLength 1
   * @maxLength 1000000
   */
  content: string;
  /** @format date-time */
  expire_at?: string | null;
}

export interface PasteBinOutputDto {
  content: string;
  /** @format date-time */
  created_at: string;
  /** @format date-time */
  expire_at: string | null;
  id: string;
  owner_user_id: string | null;
}

export interface ProfileResponseDto {
  items: PasteBinOutputDto[];
  /** @format int64 */
  total: number;
}

export interface ScheduledVipOutputDto {
  channel_id: string;
  /** @format date-time */
  created_at: string;
  id: string;
  /** @format date-time */
  remove_at?: string | null;
  remove_type?: ScheduledVipOutputDtoRemoveTypeEnum;
  user_id: string;
}

export interface SeasonStruct {
  id: string;
  short: string;
}

export interface StartRequestDtoBody {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  text: string | null;
  /** @format int32 */
  time: number;
}

export interface StartResponseDto {
  success: boolean;
}

export interface StatisticsPointDto {
  /** @format int64 */
  count: number;
  /** @format int64 */
  timestamp: number;
}

export interface StopResponseDto {
  success: boolean;
}

export interface StoredMatchesResponseMatch {
  meta: StoredMatchesResponseMatchMetaStruct;
  stats: StoredMatchesResponseMatchStats;
  teams: StoredMatchesResponseMatchTeamsStruct;
}

export interface StoredMatchesResponseMatchMetaStruct {
  cluster: string;
  id: string;
  map: MapStruct;
  mode: string;
  region: string;
  season: SeasonStruct;
  /** @format date-time */
  started_at: string;
  version: string;
}

export interface StoredMatchesResponseMatchStats {
  /** @format int64 */
  assists: number;
  character: StoredMatchesResponseMatchStatsCharacterStruct;
  damage: StoredMatchesResponseMatchStatsDamageStruct;
  /** @format int64 */
  deaths: number;
  /** @format int64 */
  kills: number;
  /** @format int64 */
  level: number;
  puuid: string;
  /** @format int64 */
  score: number;
  shots: StoredMatchesResponseMatchStatsShotsStruct;
  team: string;
  /** @format int64 */
  tier: number;
}

export interface StoredMatchesResponseMatchStatsCharacterStruct {
  id: string;
  name: string;
}

export interface StoredMatchesResponseMatchStatsDamageStruct {
  /** @format int64 */
  dealt: number;
  /** @format int64 */
  received: number;
}

export interface StoredMatchesResponseMatchStatsShotsStruct {
  /** @format int64 */
  body: number;
  /** @format int64 */
  head: number;
  /** @format int64 */
  leg: number;
}

export interface StoredMatchesResponseMatchTeamsStruct {
  /** @format int64 */
  blue: number;
  /** @format int64 */
  red: number;
}

export interface Stream {
  CommunityIds: string[];
  GameId: string;
  GameName: string;
  ID: string;
  IsMature: boolean;
  Language: string;
  /** @format date-time */
  StartedAt: string;
  TagIds: string[];
  Tags: string[];
  ThumbnailUrl: string;
  Title: string;
  Type: string;
  UserId: string;
  UserLogin: string;
  UserName: string;
  /** @format int64 */
  ViewerCount: number;
}

export interface TierStruct {
  /** @format int64 */
  id: number;
  name: string;
}

export interface TwirStatsResponseBody {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  /** @format int64 */
  channels: number;
  /** @format int64 */
  created_commands: number;
  /** @format int64 */
  haste_bins: number;
  /** @format int64 */
  live_channels: number;
  /** @format int64 */
  messages: number;
  /** @format int64 */
  short_urls: number;
  /** @format int64 */
  used_commands: number;
  /** @format int64 */
  used_emotes: number;
  /** @format int64 */
  viewers: number;
}

export interface UpdateRequestDtoBody {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  /**
   * @minLength 3
   * @maxLength 50
   * @pattern ^[a-zA-Z0-9]+$
   */
  new_short_id?: string;
  /**
   * @format uri
   * @minLength 1
   * @maxLength 2048
   */
  url?: string;
}

export enum CommandResponseDtoCooldownTypeEnum {
  GLOBAL = "GLOBAL",
  PER_USER = "PER_USER",
}

/** Type of removal: 'time' or 'stream_end' */
export enum CreateRequestDtoBodyRemoveTypeEnum {
  Time = "time",
  StreamEnd = "stream_end",
}

export enum ExpireExpiresTypeEnum {
  DISABLE = "DISABLE",
  DELETE = "DELETE",
}

export enum ScheduledVipOutputDtoRemoveTypeEnum {
  Time = "time",
  StreamEnd = "stream_end",
}

/** @default "views" */
export enum ShortUrlProfileParamsSortByEnum {
  Views = "views",
  CreatedAt = "created_at",
}

/** @default "day" */
export enum ShortUrlGetStatisticsParamsIntervalEnum {
  Hour = "hour",
  Day = "day",
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;

export interface FullRequestParams extends Omit<RequestInit, "body"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseFormat;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (securityData: SecurityDataType | null) => Promise<RequestParams | void> | RequestParams | void;
  customFetch?: typeof fetch;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown> extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "https://twir.localhost/api";
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private abortControllers = new Map<CancelToken, AbortController>();
  private customFetch = (...fetchParams: Parameters<typeof fetch>) => fetch(...fetchParams);

  private baseApiParams: RequestParams = {
    credentials: "same-origin",
    headers: {},
    redirect: "follow",
    referrerPolicy: "no-referrer",
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected encodeQueryParam(key: string, value: any) {
    const encodedKey = encodeURIComponent(key);
    return `${encodedKey}=${encodeURIComponent(typeof value === "number" ? value : `${value}`)}`;
  }

  protected addQueryParam(query: QueryParamsType, key: string) {
    return this.encodeQueryParam(key, query[key]);
  }

  protected addArrayQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];
    return value.map((v: any) => this.encodeQueryParam(key, v)).join("&");
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter((key) => "undefined" !== typeof query[key]);
    return keys
      .map((key) => (Array.isArray(query[key]) ? this.addArrayQueryParam(query, key) : this.addQueryParam(query, key)))
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string") ? JSON.stringify(input) : input,
    [ContentType.Text]: (input: any) => (input !== null && typeof input !== "string" ? JSON.stringify(input) : input),
    [ContentType.FormData]: (input: any) =>
      Object.keys(input || {}).reduce((formData, key) => {
        const property = input[key];
        formData.append(
          key,
          property instanceof Blob
            ? property
            : typeof property === "object" && property !== null
              ? JSON.stringify(property)
              : `${property}`,
        );
        return formData;
      }, new FormData()),
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  protected mergeRequestParams(params1: RequestParams, params2?: RequestParams): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected createAbortSignal = (cancelToken: CancelToken): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = async <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format,
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.baseApiParams.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];
    const responseFormat = format || requestParams.format;

    return this.customFetch(`${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`, {
      ...requestParams,
      headers: {
        ...(requestParams.headers || {}),
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
      },
      signal: (cancelToken ? this.createAbortSignal(cancelToken) : requestParams.signal) || null,
      body: typeof body === "undefined" || body === null ? null : payloadFormatter(body),
    }).then(async (response) => {
      const r = response.clone() as HttpResponse<T, E>;
      r.data = null as unknown as T;
      r.error = null as unknown as E;

      const data = !responseFormat
        ? r
        : await response[responseFormat]()
            .then((data) => {
              if (r.ok) {
                r.data = data;
              } else {
                r.error = data;
              }
              return r;
            })
            .catch((e) => {
              r.error = e;
              return r;
            });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
    });
  };
}

/**
 * @title Twir Api
 * @version 1.0.0
 * @baseUrl https://twir.localhost/api
 */
export class Api<SecurityDataType extends unknown> {
  http: HttpClient<SecurityDataType>;

  constructor(http: HttpClient<SecurityDataType>) {
    this.http = http;
  }

  auth = {
    /**
     * No description
     *
     * @tags Auth
     * @name AuthPostCode
     * @summary Auth post code
     * @request POST:/auth
     * @response `200` `BaseOutputBodyJsonAuthResponseDto` OK
     * @response `default` `ErrorModel` Error
     */
    authPostCode: (data: AuthBody, params: RequestParams = {}) =>
      this.http.request<BaseOutputBodyJsonAuthResponseDto, any>({
        path: `/auth`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  v1 = {
    /**
     * No description
     *
     * @tags Overlays/BRB
     * @name OverlaysBrbStart
     * @summary Start BRB overlay
     * @request PUT:/v1/channels/overlays/brb/start
     * @secure
     * @response `200` `BaseOutputBodyJsonStartResponseDto` OK
     * @response `default` `ErrorModel` Error
     */
    overlaysBrbStart: (data: StartRequestDtoBody, params: RequestParams = {}) =>
      this.http.request<BaseOutputBodyJsonStartResponseDto, any>({
        path: `/v1/channels/overlays/brb/start`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Overlays/BRB
     * @name OverlaysBrbStop
     * @summary Stop BRB overlay
     * @request PUT:/v1/channels/overlays/brb/stop
     * @secure
     * @response `200` `BaseOutputBodyJsonStopResponseDto` OK
     * @response `default` `ErrorModel` Error
     */
    overlaysBrbStop: (params: RequestParams = {}) =>
      this.http.request<BaseOutputBodyJsonStopResponseDto, any>({
        path: `/v1/channels/overlays/brb/stop`,
        method: "PUT",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Get current stream
     *
     * @tags Streams
     * @name ChannelsStreamsCurrent
     * @summary Get current stream
     * @request GET:/v1/channels/streams/current
     * @secure
     * @response `200` `BaseOutputBodyJsonStream` OK
     * @response `404` `void` No current stream
     */
    channelsStreamsCurrent: (params: RequestParams = {}) =>
      this.http.request<BaseOutputBodyJsonStream, void>({
        path: `/v1/channels/streams/current`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Commands
     * @name CommandsGetListByChannelId
     * @summary Get commands list by channel id
     * @request GET:/v1/channels/{channel_id}/commands
     * @response `200` `BaseOutputBodyJsonListCommandResponseDto` OK
     * @response `default` `ErrorModel` Error
     */
    commandsGetListByChannelId: (
      channelId: string,
      query?: {
        /** @example "CUSTOM" */
        module?: string;
      },
      params: RequestParams = {},
    ) =>
      this.http.request<BaseOutputBodyJsonListCommandResponseDto, any>({
        path: `/v1/channels/${channelId}/commands`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * @description Get file content by id
     *
     * @tags Files
     * @name ChannelsFilesContentDetail
     * @summary Get file content
     * @request GET:/v1/channels/{channel_id}/files/content/{file_id}
     * @response `200` `File` File content
     * @response `default` `ErrorModel` Error
     */
    channelsFilesContentDetail: (channelId: string, fileId: string, params: RequestParams = {}) =>
      this.http.request<File, any>({
        path: `/v1/channels/${channelId}/files/content/${fileId}`,
        method: "GET",
        ...params,
      }),

    /**
     * @description Requires api-key header.
     *
     * @tags Valorant
     * @name IntegrationsValorantStats
     * @summary Get valorant stats data
     * @request GET:/v1/integrations/valorant/stats
     * @secure
     * @response `200` `BaseOutputBodyJsonIntegrationsValorantStatsOutput` OK
     * @response `default` `ErrorModel` Error
     */
    integrationsValorantStats: (params: RequestParams = {}) =>
      this.http.request<BaseOutputBodyJsonIntegrationsValorantStatsOutput, any>({
        path: `/v1/integrations/valorant/stats`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Requires api-key header.
     *
     * @tags Pastebin
     * @name PastebinGetUserList
     * @summary Get authenticated user pastebins
     * @request GET:/v1/pastebin
     * @secure
     * @response `200` `BaseOutputBodyJsonProfileResponseDto` OK
     * @response `default` `ErrorModel` Error
     */
    pastebinGetUserList: (
      query?: {
        /**
         * @format int64
         * @min 1
         * @default 1
         * @example 1
         */
        page?: number;
        /**
         * @format int64
         * @min 0
         * @default 20
         * @example 20
         */
        perPage?: number;
      },
      params: RequestParams = {},
    ) =>
      this.http.request<BaseOutputBodyJsonProfileResponseDto, any>({
        path: `/v1/pastebin`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Pastebin
     * @name PastebinCreate
     * @summary Create pastebin
     * @request POST:/v1/pastebin
     * @response `200` `BaseOutputBodyJsonPasteBinOutputDto` OK
     * @response `default` `ErrorModel` Error
     */
    pastebinCreate: (data: PasteBinCreateRequestDtoBody, params: RequestParams = {}) =>
      this.http.request<BaseOutputBodyJsonPasteBinOutputDto, any>({
        path: `/v1/pastebin`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Pastebin
     * @name PastebinDelete
     * @summary Delete pastebin
     * @request DELETE:/v1/pastebin/{id}
     * @secure
     * @response `204` `void` No Content
     * @response `default` `ErrorModel` Error
     */
    pastebinDelete: (id: string, params: RequestParams = {}) =>
      this.http.request<void, any>({
        path: `/v1/pastebin/${id}`,
        method: "DELETE",
        secure: true,
        ...params,
      }),

    /**
     * No description
     *
     * @tags Pastebin
     * @name PastebinGetById
     * @summary Get pastebin by id
     * @request GET:/v1/pastebin/{id}
     * @response `200` `BaseOutputBodyJsonPasteBinOutputDto` OK
     * @response `default` `ErrorModel` Error
     */
    pastebinGetById: (id: string, params: RequestParams = {}) =>
      this.http.request<BaseOutputBodyJsonPasteBinOutputDto, any>({
        path: `/v1/pastebin/${id}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * @description Get created badges for twitch chat
     *
     * @tags Public
     * @name PublicTwirBadges
     * @summary Get badges
     * @request GET:/v1/public/badges
     * @response `200` `(BadgeWithUsers)[]` OK
     * @response `default` `ErrorModel` Error
     */
    publicTwirBadges: (params: RequestParams = {}) =>
      this.http.request<BadgeWithUsers[], any>({
        path: `/v1/public/badges`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * @description Get channel commands filtered by enabled and visible
     *
     * @tags Public
     * @name PublicChannelPublicCommands
     * @summary Get channel commands
     * @request GET:/v1/public/channels/{channelId}/commands
     * @response `200` `(CommandDto)[]` OK
     * @response `default` `ErrorModel` Error
     */
    publicChannelPublicCommands: (channelId: string, params: RequestParams = {}) =>
      this.http.request<CommandDto[], any>({
        path: `/v1/public/channels/${channelId}/commands`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * @description Get all scheduled VIPs for the selected dashboard
     *
     * @tags Scheduled VIPs
     * @name ScheduledVipsList
     * @summary List scheduled VIPs
     * @request GET:/v1/scheduled-vips
     * @secure
     * @response `200` `ListResponseDtoBody` OK
     * @response `default` `ErrorModel` Error
     */
    scheduledVipsList: (params: RequestParams = {}) =>
      this.http.request<ListResponseDtoBody, any>({
        path: `/v1/scheduled-vips`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Add a user as VIP on Twitch and schedule their removal
     *
     * @tags Scheduled VIPs
     * @name ScheduledVipsCreate
     * @summary Create scheduled VIP
     * @request POST:/v1/scheduled-vips
     * @secure
     * @response `200` `BaseOutputBodyJsonScheduledVipOutputDto` OK
     * @response `default` `ErrorModel` Error
     */
    scheduledVipsCreate: (data: CreateRequestDtoBody, params: RequestParams = {}) =>
      this.http.request<BaseOutputBodyJsonScheduledVipOutputDto, any>({
        path: `/v1/scheduled-vips`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Remove a scheduled VIP. Note: This only removes the schedule, not the VIP status on Twitch.
     *
     * @tags Scheduled VIPs
     * @name ScheduledVipsDelete
     * @summary Delete scheduled VIP
     * @request DELETE:/v1/scheduled-vips/{id}
     * @secure
     * @response `204` `void` No Content
     * @response `default` `ErrorModel` Error
     */
    scheduledVipsDelete: (id: string, params: RequestParams = {}) =>
      this.http.request<void, any>({
        path: `/v1/scheduled-vips/${id}`,
        method: "DELETE",
        secure: true,
        ...params,
      }),

    /**
     * No description
     *
     * @tags Short links
     * @name ShortUrlProfile
     * @summary Get user's short links from authenticated user and/or from browser session
     * @request GET:/v1/short-links
     * @response `200` `BaseOutputBodyJsonLinksProfileOutputDto` OK
     * @response `default` `ErrorModel` Error
     */
    shortUrlProfile: (
      query?: {
        /**
         * @format int64
         * @min 0
         * @default 0
         */
        page?: number;
        /**
         * @format int64
         * @min 1
         * @max 100
         * @default 20
         */
        perPage?: number;
        /** @default "views" */
        sortBy?: ShortUrlProfileParamsSortByEnum;
      },
      params: RequestParams = {},
    ) =>
      this.http.request<BaseOutputBodyJsonLinksProfileOutputDto, any>({
        path: `/v1/short-links`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Short links
     * @name ShortUrlCreate
     * @summary Create short url
     * @request POST:/v1/short-links
     * @response `200` `BaseOutputBodyJsonLinkOutputDto` OK
     * @response `default` `ErrorModel` Error
     */
    shortUrlCreate: (data: CreateLinkInputDto, params: RequestParams = {}) =>
      this.http.request<BaseOutputBodyJsonLinkOutputDto, any>({
        path: `/v1/short-links`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Short links
     * @name ShortLinksDeleteCustomDomain
     * @summary Delete custom domain configuration
     * @request DELETE:/v1/short-links/custom-domain
     * @secure
     * @response `200` `BaseOutputBodyJsonInterface` OK
     * @response `default` `ErrorModel` Error
     */
    shortLinksDeleteCustomDomain: (params: RequestParams = {}) =>
      this.http.request<BaseOutputBodyJsonInterface, any>({
        path: `/v1/short-links/custom-domain`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Short links
     * @name ShortLinksGetCustomDomain
     * @summary Get custom domain configuration
     * @request GET:/v1/short-links/custom-domain
     * @secure
     * @response `200` `BaseOutputBodyJsonCustomDomainOutputDto` OK
     * @response `default` `ErrorModel` Error
     */
    shortLinksGetCustomDomain: (params: RequestParams = {}) =>
      this.http.request<BaseOutputBodyJsonCustomDomainOutputDto, any>({
        path: `/v1/short-links/custom-domain`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Short links
     * @name ShortLinksCreateCustomDomain
     * @summary Configure custom domain
     * @request POST:/v1/short-links/custom-domain
     * @secure
     * @response `200` `BaseOutputBodyJsonCustomDomainOutputDto` OK
     * @response `default` `ErrorModel` Error
     */
    shortLinksCreateCustomDomain: (data: CreateCustomDomainInputBody, params: RequestParams = {}) =>
      this.http.request<BaseOutputBodyJsonCustomDomainOutputDto, any>({
        path: `/v1/short-links/custom-domain`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Short links
     * @name ShortLinksVerifyCustomDomain
     * @summary Verify custom domain DNS configuration
     * @request POST:/v1/short-links/custom-domain/verify
     * @secure
     * @response `200` `BaseOutputBodyJsonCustomDomainOutputDto` OK
     * @response `default` `ErrorModel` Error
     */
    shortLinksVerifyCustomDomain: (params: RequestParams = {}) =>
      this.http.request<BaseOutputBodyJsonCustomDomainOutputDto, any>({
        path: `/v1/short-links/custom-domain/verify`,
        method: "POST",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Short links
     * @name ShortUrlDelete
     * @summary Delete short url
     * @request DELETE:/v1/short-links/{shortId}
     * @secure
     * @response `200` `BaseOutputBodyJsonInterface` OK
     * @response `default` `ErrorModel` Error
     */
    shortUrlDelete: (shortId: string, params: RequestParams = {}) =>
      this.http.request<BaseOutputBodyJsonInterface, any>({
        path: `/v1/short-links/${shortId}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Short links
     * @name ShortUrlRedirect
     * @summary Redirect to url
     * @request GET:/v1/short-links/{shortId}
     * @response `301` `void` Moved Permanently
     * @response `default` `ErrorModel` Error
     */
    shortUrlRedirect: (shortId: string, params: RequestParams = {}) =>
      this.http.request<ErrorModel, void>({
        path: `/v1/short-links/${shortId}`,
        method: "GET",
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Short links
     * @name ShortUrlUpdate
     * @summary Update short url
     * @request PATCH:/v1/short-links/{shortId}
     * @secure
     * @response `200` `BaseOutputBodyJsonLinkOutputDto` OK
     * @response `default` `ErrorModel` Error
     */
    shortUrlUpdate: (shortId: string, data: UpdateRequestDtoBody, params: RequestParams = {}) =>
      this.http.request<BaseOutputBodyJsonLinkOutputDto, any>({
        path: `/v1/short-links/${shortId}`,
        method: "PATCH",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Short links
     * @name ShortUrlGetInfo
     * @summary Get short url data
     * @request GET:/v1/short-links/{shortId}/info
     * @response `200` `BaseOutputBodyJsonLinkOutputDto` OK
     * @response `default` `ErrorModel` Error
     */
    shortUrlGetInfo: (
      shortId: string,
      query: {
        /**
         * @minLength 1
         * @pattern ^[a-zA-Z0-9]+$
         */
        shortId: string;
      },
      params: RequestParams = {},
    ) =>
      this.http.request<BaseOutputBodyJsonLinkOutputDto, any>({
        path: `/v1/short-links/${shortId}/info`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Short links
     * @name ShortUrlGetStatistics
     * @summary Get short url statistics
     * @request GET:/v1/short-links/{shortId}/statistics
     * @response `200` `BaseOutputBodyJsonListStatisticsPointDto` OK
     * @response `default` `ErrorModel` Error
     */
    shortUrlGetStatistics: (
      shortId: string,
      query: {
        /** @format int64 */
        from: number;
        /** @format int64 */
        to: number;
        /** @default "day" */
        interval?: ShortUrlGetStatisticsParamsIntervalEnum;
      },
      params: RequestParams = {},
    ) =>
      this.http.request<BaseOutputBodyJsonListStatisticsPointDto, any>({
        path: `/v1/short-links/${shortId}/statistics`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Short links
     * @name ShortUrlGetTopCountries
     * @summary Get top countries by views for short url
     * @request GET:/v1/short-links/{shortId}/top-countries
     * @response `200` `BaseOutputBodyJsonListCountryStatsDto` OK
     * @response `default` `ErrorModel` Error
     */
    shortUrlGetTopCountries: (
      shortId: string,
      query?: {
        /**
         * @format int64
         * @min 1
         * @max 50
         * @default 10
         */
        limit?: number;
      },
      params: RequestParams = {},
    ) =>
      this.http.request<BaseOutputBodyJsonListCountryStatsDto, any>({
        path: `/v1/short-links/${shortId}/top-countries`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * @description Convert text to speech using the TTS service. Returns an audio file.
     *
     * @tags TTS
     * @name TtsSay
     * @summary Text-to-Speech Say
     * @request GET:/v1/tts/say
     * @response `200` `File` Successful TTS conversion
     * @response `default` `ErrorModel` Error
     */
    ttsSay: (
      query?: {
        /**
         * Voice name to use for TTS
         * @minLength 1
         * @maxLength 100
         * @example "alan"
         */
        voice?: string;
        /**
         * Text to convert to speech
         * @minLength 1
         * @maxLength 5000
         * @example "Hello world"
         */
        text?: string;
        /**
         * Voice pitch (0-100)
         * @format int64
         * @min 0
         * @max 100
         * @default 50
         * @example 50
         */
        pitch?: number;
        /**
         * Speech rate (0-100)
         * @format int64
         * @min 0
         * @max 100
         * @default 50
         * @example 50
         */
        rate?: number;
        /**
         * Volume level (0-100)
         * @format int64
         * @min 0
         * @max 100
         * @default 50
         * @example 50
         */
        volume?: number;
      },
      params: RequestParams = {},
    ) =>
      this.http.request<File, any>({
        path: `/v1/tts/say`,
        method: "GET",
        query: query,
        ...params,
      }),

    /**
     * @description Get Twir application statistics
     *
     * @tags Twir
     * @name TwirStats
     * @summary Twir Stats
     * @request GET:/v1/twir/stats
     * @response `200` `TwirStatsResponseBody` OK
     * @response `default` `ErrorModel` Error
     */
    twirStats: (params: RequestParams = {}) =>
      this.http.request<TwirStatsResponseBody, any>({
        path: `/v1/twir/stats`,
        method: "GET",
        format: "json",
        ...params,
      }),
  };
}
