import { UiText } from '@ory/client';
import { Alert } from '@mui/material';

interface MessageProps {
  message: UiText;
}

export const Message = ({ message }: MessageProps) => {
  return (
    <Alert
      sx={{ marginTop: 0.5 }}
      severity={message.type === 'error' ? 'error' : 'info'}
    >
      {message.text}
    </Alert>
  );
};

interface MessagesProps {
  messages?: Array<UiText>;
}

export const Messages = ({ messages }: MessagesProps) => {
  if (!messages) {
    return null;
  }

  return (
    <div>
      {messages.map(message => (
        <Message key={message.id} message={message} />
      ))}
    </div>
  );
};
