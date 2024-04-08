/* eslint-disable */
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  DateTime: { input: any; output: any; }
};

export type Command = {
  __typename?: 'Command';
  aliases?: Maybe<Array<Scalars['String']['output']>>;
  createdAt: Scalars['DateTime']['output'];
  description?: Maybe<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  responses?: Maybe<Array<CommandResponse>>;
  updatedAt: Scalars['DateTime']['output'];
};

export type CommandResponse = {
  __typename?: 'CommandResponse';
  commandId: Scalars['ID']['output'];
  id: Scalars['ID']['output'];
  order: Scalars['Int']['output'];
  text: Scalars['String']['output'];
};

export type CreateCommandInput = {
  aliases?: InputMaybe<Array<Scalars['String']['input']>>;
  description?: InputMaybe<Scalars['String']['input']>;
  name: Scalars['String']['input'];
  responses?: InputMaybe<Array<CreateCommandResponseInput>>;
};

export type CreateCommandResponseInput = {
  order: Scalars['Int']['input'];
  text: Scalars['String']['input'];
};

export type Empty = {
  __typename?: 'Empty';
  empty?: Maybe<Scalars['Boolean']['output']>;
};

export type Mutation = {
  __typename?: 'Mutation';
  createCommand: Command;
  createNotification: Notification;
  empty?: Maybe<Empty>;
  updateCommand: Command;
};


export type MutationCreateCommandArgs = {
  opts: CreateCommandInput;
};


export type MutationCreateNotificationArgs = {
  text: Scalars['String']['input'];
  userId?: InputMaybe<Scalars['String']['input']>;
};


export type MutationUpdateCommandArgs = {
  id: Scalars['String']['input'];
  opts: UpdateCommandOpts;
};

export type Notification = {
  __typename?: 'Notification';
  id: Scalars['ID']['output'];
  text: Scalars['String']['output'];
  userId: Scalars['ID']['output'];
};

export type Query = {
  __typename?: 'Query';
  authedUser: User;
  commands: Array<Command>;
  empty?: Maybe<Empty>;
  notifications: Array<Notification>;
};


export type QueryNotificationsArgs = {
  userId: Scalars['String']['input'];
};

export type UpdateCommandOpts = {
  aliases?: InputMaybe<Array<Scalars['String']['input']>>;
  description?: InputMaybe<Scalars['String']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
};

export type User = {
  __typename?: 'User';
  apiKey: Scalars['String']['output'];
  channel: UserChannel;
  hideOnLandingPage: Scalars['Boolean']['output'];
  id: Scalars['ID']['output'];
  isBanned: Scalars['Boolean']['output'];
  isBotAdmin: Scalars['Boolean']['output'];
};

export type UserChannel = {
  __typename?: 'UserChannel';
  botId: Scalars['ID']['output'];
  isBotModerator: Scalars['Boolean']['output'];
  isEnabled: Scalars['Boolean']['output'];
};

export type NewCommandMutationVariables = Exact<{
  name: Scalars['String']['input'];
  aliases?: InputMaybe<Array<Scalars['String']['input']> | Scalars['String']['input']>;
  description?: InputMaybe<Scalars['String']['input']>;
  responses?: InputMaybe<Array<CreateCommandResponseInput> | CreateCommandResponseInput>;
}>;


export type NewCommandMutation = { __typename?: 'Mutation', createCommand: { __typename?: 'Command', id: string } };

export type GetCommandsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetCommandsQuery = { __typename?: 'Query', commands: Array<{ __typename?: 'Command', name: string, id: string, aliases?: Array<string> | null, description?: string | null, createdAt: any, updatedAt: any, responses?: Array<{ __typename?: 'CommandResponse', text: string }> | null }> };

export type GetUserQueryVariables = Exact<{ [key: string]: never; }>;


export type GetUserQuery = { __typename?: 'Query', authedUser: { __typename?: 'User', id: string, apiKey: string, hideOnLandingPage: boolean, isBanned: boolean, isBotAdmin: boolean, channel: { __typename?: 'UserChannel', botId: string, isBotModerator: boolean, isEnabled: boolean } } };


export const NewCommandDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"newCommand"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"name"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"aliases"}},"type":{"kind":"ListType","type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"description"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"responses"}},"type":{"kind":"ListType","type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateCommandResponseInput"}}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createCommand"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"opts"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"name"},"value":{"kind":"Variable","name":{"kind":"Name","value":"name"}}},{"kind":"ObjectField","name":{"kind":"Name","value":"description"},"value":{"kind":"Variable","name":{"kind":"Name","value":"description"}}},{"kind":"ObjectField","name":{"kind":"Name","value":"aliases"},"value":{"kind":"Variable","name":{"kind":"Name","value":"aliases"}}},{"kind":"ObjectField","name":{"kind":"Name","value":"responses"},"value":{"kind":"Variable","name":{"kind":"Name","value":"responses"}}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]} as unknown as DocumentNode<NewCommandMutation, NewCommandMutationVariables>;
export const GetCommandsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"getCommands"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"commands"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"aliases"}},{"kind":"Field","name":{"kind":"Name","value":"responses"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"text"}}]}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}}]} as unknown as DocumentNode<GetCommandsQuery, GetCommandsQueryVariables>;
export const GetUserDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"getUser"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"authedUser"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"apiKey"}},{"kind":"Field","name":{"kind":"Name","value":"channel"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"botId"}},{"kind":"Field","name":{"kind":"Name","value":"isBotModerator"}},{"kind":"Field","name":{"kind":"Name","value":"isEnabled"}}]}},{"kind":"Field","name":{"kind":"Name","value":"hideOnLandingPage"}},{"kind":"Field","name":{"kind":"Name","value":"isBanned"}},{"kind":"Field","name":{"kind":"Name","value":"isBotAdmin"}}]}}]}}]} as unknown as DocumentNode<GetUserQuery, GetUserQueryVariables>;