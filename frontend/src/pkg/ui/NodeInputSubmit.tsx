import { getNodeLabel } from '@ory/integrations/ui';
import { Button } from '@mui/material';
import { NodeInputProps } from './helpers';

export function NodeInputSubmit<T>({
  node,
  attributes,
  setValue,
  disabled,
  dispatchSubmit
}: NodeInputProps) {
  return (
    <Button
      name={attributes.name}
      fullWidth
      variant="contained"
      sx={{ mt: 3, mb: 2 }}
      onClick={(e: any) => {
        // On click, we set this value, and once set, dispatch the submission!
        setValue(attributes.value).then(() => dispatchSubmit(e));
      }}
      value={attributes.value || ''}
      disabled={attributes.disabled || disabled}
    >
      {getNodeLabel(node)}
    </Button>
  );
}
