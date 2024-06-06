/* eslint-disable */
import * as types from './graphql';
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 */
const documents = {
    "\n\t\t\tquery ChatOverlayWithAdditionalData($overlayId: String!) {\n\t\t\t\tauthenticatedUser {\n\t\t\t\t\tid\n\t\t\t\t\ttwitchProfile {\n\t\t\t\t\t\tlogin\n\t\t\t\t\t\tdisplayName\n\t\t\t\t\t\tprofileImageUrl\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t\tchatOverlaysById(id: $overlayId) {\n\t\t\t\t\tid\n\t\t\t\t\tmessageHideTimeout\n\t\t\t\t}\n\t\t\t\ttwitchGetGlobalBadges {\n\t\t\t\t\tbadges {\n\t\t\t\t\t\tset_id\n\t\t\t\t\t\tversions {\n\t\t\t\t\t\t\timage_url_1x\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t\ttwitchGetChannelBadges {\n\t\t\t\t\tbadges {\n\t\t\t\t\t\tset_id\n\t\t\t\t\t\tversions {\n\t\t\t\t\t\t\timage_url_1x\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t": types.ChatOverlayWithAdditionalDataDocument,
    "\n\t\t\tsubscription ChatOverlaySettings($id: String!, $apiKey: String!) {\n\t\t\t\tchatOverlaySettings(id: $id, apiKey: $apiKey) {\n\t\t\t\t\tid\n\t\t\t\t\tanimation\n\t\t\t\t\tfontSize\n\t\t\t\t}\n\t\t\t}\n\t\t": types.ChatOverlaySettingsDocument,
};

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = graphql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function graphql(source: string): unknown;

/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\t\tquery ChatOverlayWithAdditionalData($overlayId: String!) {\n\t\t\t\tauthenticatedUser {\n\t\t\t\t\tid\n\t\t\t\t\ttwitchProfile {\n\t\t\t\t\t\tlogin\n\t\t\t\t\t\tdisplayName\n\t\t\t\t\t\tprofileImageUrl\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t\tchatOverlaysById(id: $overlayId) {\n\t\t\t\t\tid\n\t\t\t\t\tmessageHideTimeout\n\t\t\t\t}\n\t\t\t\ttwitchGetGlobalBadges {\n\t\t\t\t\tbadges {\n\t\t\t\t\t\tset_id\n\t\t\t\t\t\tversions {\n\t\t\t\t\t\t\timage_url_1x\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t\ttwitchGetChannelBadges {\n\t\t\t\t\tbadges {\n\t\t\t\t\t\tset_id\n\t\t\t\t\t\tversions {\n\t\t\t\t\t\t\timage_url_1x\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t"): (typeof documents)["\n\t\t\tquery ChatOverlayWithAdditionalData($overlayId: String!) {\n\t\t\t\tauthenticatedUser {\n\t\t\t\t\tid\n\t\t\t\t\ttwitchProfile {\n\t\t\t\t\t\tlogin\n\t\t\t\t\t\tdisplayName\n\t\t\t\t\t\tprofileImageUrl\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t\tchatOverlaysById(id: $overlayId) {\n\t\t\t\t\tid\n\t\t\t\t\tmessageHideTimeout\n\t\t\t\t}\n\t\t\t\ttwitchGetGlobalBadges {\n\t\t\t\t\tbadges {\n\t\t\t\t\t\tset_id\n\t\t\t\t\t\tversions {\n\t\t\t\t\t\t\timage_url_1x\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t\ttwitchGetChannelBadges {\n\t\t\t\t\tbadges {\n\t\t\t\t\t\tset_id\n\t\t\t\t\t\tversions {\n\t\t\t\t\t\t\timage_url_1x\n\t\t\t\t\t\t}\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\t\tsubscription ChatOverlaySettings($id: String!, $apiKey: String!) {\n\t\t\t\tchatOverlaySettings(id: $id, apiKey: $apiKey) {\n\t\t\t\t\tid\n\t\t\t\t\tanimation\n\t\t\t\t\tfontSize\n\t\t\t\t}\n\t\t\t}\n\t\t"): (typeof documents)["\n\t\t\tsubscription ChatOverlaySettings($id: String!, $apiKey: String!) {\n\t\t\t\tchatOverlaySettings(id: $id, apiKey: $apiKey) {\n\t\t\t\t\tid\n\t\t\t\t\tanimation\n\t\t\t\t\tfontSize\n\t\t\t\t}\n\t\t\t}\n\t\t"];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;