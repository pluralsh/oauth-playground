import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import {
  ApolloClient,
  createHttpLink,
  InMemoryCache,
  ApolloProvider
} from '@apollo/client';
import { RecoilRoot } from 'recoil';

import { ThemeProvider } from "@ory/elements"
import "@ory/elements/style.css"
import { BrowserRouter } from 'react-router-dom';

// // Ory Elements
// // optional fontawesome icons
// import "@ory/elements/assets/fa-brands.min.css"
// import "@ory/elements/assets/fa-solid.min.css"
// import "@ory/elements/assets/fontawesome.min.css"

// // optional fonts
// import "@ory/elements/assets/inter-font.css"
// import "@ory/elements/assets/jetbrains-mono-font.css"

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

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <ApolloProvider client={client}>
      <RecoilRoot>
        <ThemeProvider themeOverrides={{
          fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif",
          }}>
          <App />
        </ThemeProvider>
      </RecoilRoot>
    </ApolloProvider>
  </React.StrictMode>,
);
