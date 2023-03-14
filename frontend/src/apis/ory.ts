import { V0alpha2Api, Configuration } from '@ory/client';

export const basePath = process.env.REACT_APP_KRATOS_PUBLIC_URL || '/.ory/kratos/public';

const ory = new V0alpha2Api(
  new Configuration({
    basePath,
    baseOptions: {
      withCredentials: true
    }
  })
);

export default ory;
