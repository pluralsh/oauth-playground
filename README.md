# Instructions for running locally

*Note: the `keto` and `hydrator` directories can be ignored*

Go into the `api-server` directory and run `make setup-tilt-cluster` to create a kind cluster and a registry for tilt to use.

Next, run `tilt up` which will proceed to install Kratos, Hydra, OAuthkeeper, Keto, Grafana, the API server and the frontend to the cluster.

The frontend can be accessed at `https://localhost.pluraldev.sh:4455/`.
The GraphiQL interface can be accessed at `https://localhost.pluraldev.sh:4455/graphiql`.
Grafana can be accessed at `https://grafana.localhost.pluraldev.sh:4455/`.
Hydra is hosted on `https://hydra.localhost.pluraldev.sh:4455/`.

The helm values and other configs used for the deployment (like the oathkeeper access rules and keto namespace configs) can be found in `api-server/dev`.

Whenever a change is made to the deployment configs tilt will automatically redeploy the component.
If changes are made to the api server a new image will be built and tilt will update the deployment automatically.
For the frontend, the local filesystem is kept in-sync with the container on the cluster and it is run using `yarn start`. This way a new container for the frontend only needs to be built when there are changes to `package.sjon` or `yarn.lock`.
