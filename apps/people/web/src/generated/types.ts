export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Time: { input: string; output: string; }
};

export type CreateOrganizationInput = {
  name: Scalars['String']['input'];
  ownerEmail: Scalars['String']['input'];
  ownerFirstName: Scalars['String']['input'];
  ownerLastName: Scalars['String']['input'];
};

export type CreateOrganizationPayload = {
  __typename?: 'CreateOrganizationPayload';
  organizationId: Scalars['ID']['output'];
  owner: User;
};

export type Mutation = {
  __typename?: 'Mutation';
  createOrganization: CreateOrganizationPayload;
};


export type MutationCreateOrganizationArgs = {
  input: CreateOrganizationInput;
};

export type PageInfo = {
  __typename?: 'PageInfo';
  endCursor?: Maybe<Scalars['String']['output']>;
  hasNextPage: Scalars['Boolean']['output'];
};

export type Query = {
  __typename?: 'Query';
  user?: Maybe<User>;
  users: UserConnection;
};


export type QueryUserArgs = {
  id: Scalars['ID']['input'];
};


export type QueryUsersArgs = {
  after?: InputMaybe<Scalars['String']['input']>;
  first?: InputMaybe<Scalars['Int']['input']>;
  organizationId: Scalars['ID']['input'];
};

export type User = {
  __typename?: 'User';
  createdAt: Scalars['Time']['output'];
  email: Scalars['String']['output'];
  firstName: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  lastName: Scalars['String']['output'];
  organizationId: Scalars['ID']['output'];
  roleId: Scalars['ID']['output'];
};

export type UserConnection = {
  __typename?: 'UserConnection';
  nodes: Array<User>;
  pageInfo: PageInfo;
};
