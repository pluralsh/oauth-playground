import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import {
  ApolloClient,
  createHttpLink,
  InMemoryCache,
  ApolloProvider
} from '@apollo/client';
import { RecoilRoot } from 'recoil';

const httpLink = createHttpLink({
  uri: '/graphql',
  credentials: 'include'
});

const client = new ApolloClient({
  // uri: 'http://localhost:8082/query',
  link: httpLink,
  cache: new InMemoryCache(),
  defaultOptions: {
    watchQuery: {
      fetchPolicy: 'cache-and-network'
    }
  }
});

ReactDOM.render(
  <React.StrictMode>
    <ApolloProvider client={client}>
      <RecoilRoot>
        <App />
      </RecoilRoot>
    </ApolloProvider>
  </React.StrictMode>,
  document.getElementById('root')
);
