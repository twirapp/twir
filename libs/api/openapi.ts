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

export interface CreateLinkInputDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  /**
   * @format uri
   * @minLength 1
   * @maxLength 2000
   * @pattern ^https?://.*
   */
  url: string;
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

export interface GetManyOutputDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  items: PasteBinOutputDto[];
  /** @format int64 */
  total: number;
}

export interface LinkOutputDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  id: string;
  short_url: string;
  url: string;
  /** @format int64 */
  views: number;
}

export interface PasteBinCreateDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  /**
   * @minLength 1
   * @maxLength 100000
   */
  content: string;
  /** @format date-time */
  expire_at?: string | null;
}

export interface PasteBinOutputDto {
  /**
   * A URL to the JSON Schema for this object.
   * @format uri
   */
  $schema?: string;
  content: string;
  /** @format date-time */
  created_at: string;
  /** @format date-time */
  expire_at: string | null;
  id: string;
  owner_user_id: string | null;
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
     * @description Get file content by id
     *
     * @tags Files
     * @name ChannelsFilesContentDetail
     * @summary Get file content
     * @request GET:/v1/channels/{channelId}/files/content/{fileId}
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
     * @tags Pastebin
     * @name PastebinGetUserList
     * @summary Get authenticated user pastebins
     * @request GET:/v1/pastebin
     * @secure
     * @response `200` `GetManyOutputDto` OK
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
      this.http.request<GetManyOutputDto, any>({
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
     * @response `200` `PasteBinOutputDto` OK
     * @response `default` `ErrorModel` Error
     */
    pastebinCreate: (data: PasteBinCreateDto, params: RequestParams = {}) =>
      this.http.request<PasteBinOutputDto, any>({
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
     * @response `200` `PasteBinOutputDto` OK
     * @response `default` `ErrorModel` Error
     */
    pastebinDelete: (id: string, params: RequestParams = {}) =>
      this.http.request<PasteBinOutputDto, any>({
        path: `/v1/pastebin/${id}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * No description
     *
     * @tags Pastebin
     * @name PastebinGetById
     * @summary Get pastebin by id
     * @request GET:/v1/pastebin/{id}
     * @response `200` `PasteBinOutputDto` OK
     * @response `default` `ErrorModel` Error
     */
    pastebinGetById: (id: string, params: RequestParams = {}) =>
      this.http.request<PasteBinOutputDto, any>({
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
     * No description
     *
     * @tags Short links
     * @name ShortUrlGetInfo
     * @summary Get short url data
     * @request GET:/v1/short-links
     * @response `200` `LinkOutputDto` OK
     * @response `default` `ErrorModel` Error
     */
    shortUrlGetInfo: (
      query: {
        /**
         * @minLength 1
         * @maxLength 5
         * @pattern ^[a-zA-Z0-9]+$
         */
        shortId: string;
      },
      params: RequestParams = {},
    ) =>
      this.http.request<LinkOutputDto, any>({
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
     * @response `200` `LinkOutputDto` OK
     * @response `default` `ErrorModel` Error
     */
    shortUrlCreate: (data: CreateLinkInputDto, params: RequestParams = {}) =>
      this.http.request<LinkOutputDto, any>({
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
  };
}
