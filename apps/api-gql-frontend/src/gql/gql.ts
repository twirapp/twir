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
    "\n\tmutation newCommand($name: String!, $description: String) {\n\t\tcreateCommand(name: $name, description: $description) {\n\t\t\tid\n\t\t\tname\n\t\t}\n\t}\n": types.NewCommandDocument,
    "\n\tmutation newNotification($text: String!, $userId: String) {\n\t\tcreateNotification(text: $text, userId: $userId) {\n\t\t\tid\n\t\t\tuserId\n\t\t\ttext\n\t\t}\n\t}\n": types.NewNotificationDocument,
    "\n\t\t\tquery getCommandsAndNotifications($userId: String!) {\n\t\t\t\tcommands {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\taliases\n\t\t\t\t}\n\t\t\t\tnotifications(userId: $userId) {\n\t\t\t\t\tid\n\t\t\t\t\tuserId\n\t\t\t\t\ttext\n\t\t\t\t}\n\t\t\t}\n  ": types.GetCommandsAndNotificationsDocument,
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
export function graphql(source: "\n\tmutation newCommand($name: String!, $description: String) {\n\t\tcreateCommand(name: $name, description: $description) {\n\t\t\tid\n\t\t\tname\n\t\t}\n\t}\n"): (typeof documents)["\n\tmutation newCommand($name: String!, $description: String) {\n\t\tcreateCommand(name: $name, description: $description) {\n\t\t\tid\n\t\t\tname\n\t\t}\n\t}\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\tmutation newNotification($text: String!, $userId: String) {\n\t\tcreateNotification(text: $text, userId: $userId) {\n\t\t\tid\n\t\t\tuserId\n\t\t\ttext\n\t\t}\n\t}\n"): (typeof documents)["\n\tmutation newNotification($text: String!, $userId: String) {\n\t\tcreateNotification(text: $text, userId: $userId) {\n\t\t\tid\n\t\t\tuserId\n\t\t\ttext\n\t\t}\n\t}\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\t\tquery getCommandsAndNotifications($userId: String!) {\n\t\t\t\tcommands {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\taliases\n\t\t\t\t}\n\t\t\t\tnotifications(userId: $userId) {\n\t\t\t\t\tid\n\t\t\t\t\tuserId\n\t\t\t\t\ttext\n\t\t\t\t}\n\t\t\t}\n  "): (typeof documents)["\n\t\t\tquery getCommandsAndNotifications($userId: String!) {\n\t\t\t\tcommands {\n\t\t\t\t\tid\n\t\t\t\t\tname\n\t\t\t\t\taliases\n\t\t\t\t}\n\t\t\t\tnotifications(userId: $userId) {\n\t\t\t\t\tid\n\t\t\t\t\tuserId\n\t\t\t\t\ttext\n\t\t\t\t}\n\t\t\t}\n  "];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;