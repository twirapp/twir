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
    "\n\tmutation newCommand($name: String!, $aliases: [String!], $description: String, $responses: [CreateCommandResponseInput!]) {\n\t\tcreateCommand(\n    opts: {name: $name, description: $description, aliases: $aliases, responses: $responses}\n  ) {\n    id\n  }\n\t}\n": types.NewCommandDocument,
    "\n\t\t\tquery getCommands {\n\t\t\t\tcommands {\n\t\t\t\t\tname\n\t\t\t\t\tid\n\t\t\t\t\taliases\n\t\t\t\t\tresponses {\n\t\t\t\t\t\ttext\n\t\t\t\t\t}\n\t\t\t\t\tdescription\n\t\t\t\t\tcreatedAt\n\t\t\t\t\tupdatedAt\n\t\t\t\t}\n\t\t\t}\n  ": types.GetCommandsDocument,
    "\n\t\tquery getUser {\n\t\t\tauthedUser {\n\t\t\tid\n\t\t\tapiKey\n\t\t\tchannel {\n\t\t\t\tbotId\n\t\t\t\tisBotModerator\n\t\t\t\tisEnabled\n\t\t\t}\n\t\t\thideOnLandingPage\n\t\t\tisBanned\n\t\t\tisBotAdmin\n\t\t\t}\n\t\t}\n\t": types.GetUserDocument,
    "\n\t\tsubscription newC {\n\t\t\tnewCommand {\n\t\t\t\tid\n\t\t\t}\n\t\t}\n\t": types.NewCDocument,
    "\n\t\tsubscription newN {\n\t\t\tnewNotification {\n\t\t\t\tid\n\t\t\t}\n\t\t}\n\t": types.NewNDocument,
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
export function graphql(source: "\n\tmutation newCommand($name: String!, $aliases: [String!], $description: String, $responses: [CreateCommandResponseInput!]) {\n\t\tcreateCommand(\n    opts: {name: $name, description: $description, aliases: $aliases, responses: $responses}\n  ) {\n    id\n  }\n\t}\n"): (typeof documents)["\n\tmutation newCommand($name: String!, $aliases: [String!], $description: String, $responses: [CreateCommandResponseInput!]) {\n\t\tcreateCommand(\n    opts: {name: $name, description: $description, aliases: $aliases, responses: $responses}\n  ) {\n    id\n  }\n\t}\n"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\t\tquery getCommands {\n\t\t\t\tcommands {\n\t\t\t\t\tname\n\t\t\t\t\tid\n\t\t\t\t\taliases\n\t\t\t\t\tresponses {\n\t\t\t\t\t\ttext\n\t\t\t\t\t}\n\t\t\t\t\tdescription\n\t\t\t\t\tcreatedAt\n\t\t\t\t\tupdatedAt\n\t\t\t\t}\n\t\t\t}\n  "): (typeof documents)["\n\t\t\tquery getCommands {\n\t\t\t\tcommands {\n\t\t\t\t\tname\n\t\t\t\t\tid\n\t\t\t\t\taliases\n\t\t\t\t\tresponses {\n\t\t\t\t\t\ttext\n\t\t\t\t\t}\n\t\t\t\t\tdescription\n\t\t\t\t\tcreatedAt\n\t\t\t\t\tupdatedAt\n\t\t\t\t}\n\t\t\t}\n  "];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\tquery getUser {\n\t\t\tauthedUser {\n\t\t\tid\n\t\t\tapiKey\n\t\t\tchannel {\n\t\t\t\tbotId\n\t\t\t\tisBotModerator\n\t\t\t\tisEnabled\n\t\t\t}\n\t\t\thideOnLandingPage\n\t\t\tisBanned\n\t\t\tisBotAdmin\n\t\t\t}\n\t\t}\n\t"): (typeof documents)["\n\t\tquery getUser {\n\t\t\tauthedUser {\n\t\t\tid\n\t\t\tapiKey\n\t\t\tchannel {\n\t\t\t\tbotId\n\t\t\t\tisBotModerator\n\t\t\t\tisEnabled\n\t\t\t}\n\t\t\thideOnLandingPage\n\t\t\tisBanned\n\t\t\tisBotAdmin\n\t\t\t}\n\t\t}\n\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\tsubscription newC {\n\t\t\tnewCommand {\n\t\t\t\tid\n\t\t\t}\n\t\t}\n\t"): (typeof documents)["\n\t\tsubscription newC {\n\t\t\tnewCommand {\n\t\t\t\tid\n\t\t\t}\n\t\t}\n\t"];
/**
 * The graphql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function graphql(source: "\n\t\tsubscription newN {\n\t\t\tnewNotification {\n\t\t\t\tid\n\t\t\t}\n\t\t}\n\t"): (typeof documents)["\n\t\tsubscription newN {\n\t\t\tnewNotification {\n\t\t\t\tid\n\t\t\t}\n\t\t}\n\t"];

export function graphql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;