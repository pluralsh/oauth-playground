import { Checkbox } from '@mui/material';
import { NodeInputProps } from './helpers';

export function NodeInputCheckbox<T>({
  node,
  attributes,
  setValue,
  disabled
}: NodeInputProps) {
  // Render a checkbox.s
  return (
    <>
      <Checkbox
        name={attributes.name}
        defaultChecked={attributes.value === true}
        onChange={(e: any) => setValue(e.target.checked)}
        disabled={attributes.disabled || disabled}
        // label={attributes.label}
        // state={
        //   node.messages.find(({ type }) => type === "error")
        //     ? "error"
        //     : undefined
        // }
        // subtitle={node.messages.map(({ text }) => text).join("\n")}
      />
      {attributes.label}
    </>
  );
}
