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
    "\n\t\t\tquery NotificationsGetAll {\n\t\t\t\tnotificationsByUser {\n\t\t\t\t\tid\n\t\t\t\t\ttext\n\t\t\t\t\tcreatedAt\n\t\t\t\t}\n\t\t\t}\n\t\t": types.NotificationsGetAllDocument,
    "\n\t\t\tsubscription NotificationsSubscription {\n\t\t\t\tnewNotification {\n\t\t\t\t\tid\n\t\t\t\t\ttext\n\t\t\t\t\tcreatedAt\n\t\t\t\t}\n\t\t\t}\n\t\t": types.NotificationsSubscriptionDocument,
    "\n\t\t\tquery notificationsByAdmin($opts: AdminNotificationsParams!) {\n\t\t\t\tnotificationsByAdmin(opts: $opts) {\n\t\t\t\t\ttotal\n\t\t\t\t\tnotifications {\n\t\t\t\t\t\tid\n\t\t\t\t\t\ttext\n\t\t\t\t\t\tuserId\n\t\t\t\t\t\ttwitchProfile {\n\t\t\t\t\t\t\tdisplayName\n\t\t\t\t\t\t\tprofileImageUrl\n\t\t\t\t\t\t}\n\t\t\t\t\t\tcreatedAt\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t": types.NotificationsByAdminDocument,
    "\n\t\tmutation CreateNotification($text: String!, $userId: String) {\n      notificationsCreate(text: $text, userId: $userId) {\n\t\t\t\ttext\n\t\t\t\tuserId\n\t\t\t}\n    }\n\t": types.CreateNotificationDocument,
    "\n\t\tmutation DeleteNotification($id: ID!) {\n\t\t\tnotificationsDelete(id: $id)\n\t\t}\n\t": types.DeleteNotificationDocument,
    "\n\t\tmutation UpdateNotifications($id: ID!, $opts: NotificationUpdateOpts!) {\n\t\t\tnotificationsUpdate(id: $id, opts: $opts) {\n\t\t\t\tid\n\t\t\t\ttext\n\t\t\t}\n\t\t}\n\t": types.UpdateNotificationsDocument,
    "\n\t\t\tsubscription dashboardStats {\n\t\t\t\tdashboardStats {\n\t\t\t\t\tcategoryId\n\t\t\t\t\tcategoryName\n\t\t\t\t\tviewers\n\t\t\t\t\tstartedAt\n\t\t\t\t\ttitle\n\t\t\t\t\tchatMessages\n\t\t\t\t\tfollowers\n\t\t\t\t\tusedEmotes\n\t\t\t\t\trequestedSongs\n\t\t\t\t\tsubs\n\t\t\t\t}\n\t\t\t}\n\t\t": types.DashboardStatsDocument,
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
export function graphql(source: "\n\t\t\tquery NotificationsGetAll {\n\t\t\t\tnotificationsByUser {\n\t\t\t\t\tid\n\t\t\t\t\ttext\n\t\t\t\t\tcreatedAt\n\t\t\t\t}\n\t\t\t}\n\t\t"): (typeof documents)["\n\t\t\tquery NotificationsGetAll {\n\t\t\t\tnotificationsByUser {\n\t\t\t\t\tid\n\t\t\t\t\ttext\n\t\t\t\t\tcreatedAt\n\t\t\t\t}\n\t\t\t}\n\t\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\t\tsubscription NotificationsSubscription {\n\t\t\t\tnewNotification {\n\t\t\t\t\tid\n\t\t\t\t\ttext\n\t\t\t\t\tcreatedAt\n\t\t\t\t}\n\t\t\t}\n\t\t"): (typeof documents)["\n\t\t\tsubscription NotificationsSubscription {\n\t\t\t\tnewNotification {\n\t\t\t\t\tid\n\t\t\t\t\ttext\n\t\t\t\t\tcreatedAt\n\t\t\t\t}\n\t\t\t}\n\t\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\t\tquery notificationsByAdmin($opts: AdminNotificationsParams!) {\n\t\t\t\tnotificationsByAdmin(opts: $opts) {\n\t\t\t\t\ttotal\n\t\t\t\t\tnotifications {\n\t\t\t\t\t\tid\n\t\t\t\t\t\ttext\n\t\t\t\t\t\tuserId\n\t\t\t\t\t\ttwitchProfile {\n\t\t\t\t\t\t\tdisplayName\n\t\t\t\t\t\t\tprofileImageUrl\n\t\t\t\t\t\t}\n\t\t\t\t\t\tcreatedAt\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t"): (typeof documents)["\n\t\t\tquery notificationsByAdmin($opts: AdminNotificationsParams!) {\n\t\t\t\tnotificationsByAdmin(opts: $opts) {\n\t\t\t\t\ttotal\n\t\t\t\t\tnotifications {\n\t\t\t\t\t\tid\n\t\t\t\t\t\ttext\n\t\t\t\t\t\tuserId\n\t\t\t\t\t\ttwitchProfile {\n\t\t\t\t\t\t\tdisplayName\n\t\t\t\t\t\t\tprofileImageUrl\n\t\t\t\t\t\t}\n\t\t\t\t\t\tcreatedAt\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\t\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\tmutation CreateNotification($text: String!, $userId: String) {\n      notificationsCreate(text: $text, userId: $userId) {\n\t\t\t\ttext\n\t\t\t\tuserId\n\t\t\t}\n    }\n\t"): (typeof documents)["\n\t\tmutation CreateNotification($text: String!, $userId: String) {\n      notificationsCreate(text: $text, userId: $userId) {\n\t\t\t\ttext\n\t\t\t\tuserId\n\t\t\t}\n    }\n\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\tmutation DeleteNotification($id: ID!) {\n\t\t\tnotificationsDelete(id: $id)\n\t\t}\n\t"): (typeof documents)["\n\t\tmutation DeleteNotification($id: ID!) {\n\t\t\tnotificationsDelete(id: $id)\n\t\t}\n\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\tmutation UpdateNotifications($id: ID!, $opts: NotificationUpdateOpts!) {\n\t\t\tnotificationsUpdate(id: $id, opts: $opts) {\n\t\t\t\tid\n\t\t\t\ttext\n\t\t\t}\n\t\t}\n\t"): (typeof documents)["\n\t\tmutation UpdateNotifications($id: ID!, $opts: NotificationUpdateOpts!) {\n\t\t\tnotificationsUpdate(id: $id, opts: $opts) {\n\t\t\t\tid\n\t\t\t\ttext\n\t\t\t}\n\t\t}\n\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\t\tsubscription dashboardStats {\n\t\t\t\tdashboardStats {\n\t\t\t\t\tcategoryId\n\t\t\t\t\tcategoryName\n\t\t\t\t\tviewers\n\t\t\t\t\tstartedAt\n\t\t\t\t\ttitle\n\t\t\t\t\tchatMessages\n\t\t\t\t\tfollowers\n\t\t\t\t\tusedEmotes\n\t\t\t\t\trequestedSongs\n\t\t\t\t\tsubs\n\t\t\t\t}\n\t\t\t}\n\t\t"): (typeof documents)["\n\t\t\tsubscription dashboardStats {\n\t\t\t\tdashboardStats {\n\t\t\t\t\tcategoryId\n\t\t\t\t\tcategoryName\n\t\t\t\t\tviewers\n\t\t\t\t\tstartedAt\n\t\t\t\t\ttitle\n\t\t\t\t\tchatMessages\n\t\t\t\t\tfollowers\n\t\t\t\t\tusedEmotes\n\t\t\t\t\trequestedSongs\n\t\t\t\t\tsubs\n\t\t\t\t}\n\t\t\t}\n\t\t"];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;